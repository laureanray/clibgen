package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)


const (
  workers = 5
)

type Downloader struct {
  url string
}

type Part struct {
  Data []byte
  Index int
}

func New(url string) *Downloader {
  downloader := &Downloader{
    url: url,
  } 

  return downloader
}


func (d *Downloader) Do() {
  var size int

  results := make(chan Part, workers) 
  parts := [workers][]byte{}

  client := &http.Client{}

  req, err := http.NewRequest("HEAD", d.url, nil)

  if err != nil {
    log.Fatal(err)
  }

  resp, err := client.Do(req)

  if err != nil {
    log.Fatal(err)
  }

  if header, ok := resp.Header["Content-Length"]; ok {
    fileSize, err := strconv.Atoi(header[0])

    if err != nil {
      log.Fatal("Filesize couldn't be determined: ", err)
    }

    size = fileSize / workers
  }


  for i := 0; i < workers; i++ {
    go d.downloadPart(i, size, results)
  }

  counter := 0

  for part := range results {
    counter++

    parts[part.Index] = part.Data
    if counter == workers {
      break
    }
  }

  file := []byte{}

  for _, part := range parts {
    file = append(file, part...)
  }

  // Set permissions accordingly, 0700 may not
  // be the best choice
  err = ioutil.WriteFile("data", file, 0700)

  if err != nil {
    log.Fatal(err)
  }
}

func (d *Downloader) downloadPart(index, size int, c chan Part) {
  client := &http.Client{}

  start := index * size

  dataRange := fmt.Sprintf("bytes=%d-%d", start, start+size-1)
  if index == workers - 1 {
    dataRange = fmt.Sprintf("bytes=%d-", start)
  }

  log.Println(dataRange)

  req, err := http.NewRequest("GET", d.url, nil)

  if err != nil {
    return 
  }

  req.Header.Add("Range", dataRange)

  resp, err := client.Do(req)

  if err != nil {
    return
  }
  
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)

  if err != nil {
    return 
  }

  c <- Part{Index: index, Data: body}
}
