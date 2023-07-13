package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
  update = flag.Bool("update", false, "update .golden files")
)

var binaryName = "clibgen"
var binaryPath = ""

func TestMain(m *testing.M) {
  flag.Parse()

  err := os.Chdir("..")
  if err != nil {
    fmt.Printf("could not change dir: %v", err)
    os.Exit(1)
  }

  dir, err := os.Getwd()

  if err != nil {
    fmt.Printf("could not get current dir: %v", err)
  }

  binaryPath = filepath.Join(dir, binaryName)

  fmt.Println("binary path", binaryPath)

  os.Exit(m.Run())
}


func runBinary(args []string) ([]byte, error) {
  cmd := exec.Command(binaryPath, args...)
  cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
  return cmd.CombinedOutput()
}

func runBinaryWithFileInput(args []string, bytesToWrite []byte) ([]byte, error) {
  fmt.Println("Running binary with file input", bytesToWrite)
  fmt.Printf("Executing command: %s", binaryPath)

  cmd := exec.Command(binaryPath, args...)
  cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")

  stdout, err := cmd.StdoutPipe()
  stdin, err := cmd.StdinPipe()

  cmd.Stderr = cmd.Stdout

  if err != nil {
    return nil, err
  }

  if err = cmd.Start(); err != nil {
    return nil, err
  }

  time.Sleep(1 * time.Second)

  _, err = stdin.Write(bytesToWrite)

  if err != nil {
    fmt.Println("Error writing to stdin", err)
  }

  time.Sleep(15 * time.Second)

  if err = cmd.Process.Signal(os.Interrupt); err != nil {
    fmt.Println("signal")
    return nil, err
  }

  out, _ := ioutil.ReadAll(stdout)

  if err = cmd.Wait(); err != nil {
    if strings.Contains(err.Error(), "interrupt") {
      return out, nil
    } else {
      return nil, err
    }
  }

  return nil, err
}

func TestStaticCliArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
	}{
    {"help args", []string{"--help"}, "help.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args)

			if err != nil {
				t.Fatal(err)
			}

			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := string(output)

			expected := loadFixture(t, tt.fixture)

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
    bytesToWrite []byte
	}{
    {"search test", []string{"search", "Eloquent JavaScript"}, "eloquent.golden", []byte{14, 14, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinaryWithFileInput(tt.args, tt.bytesToWrite)

			if err != nil {
        fmt.Println("error", err.Error())
				t.Fatal(err)
			}

			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := removeLastLine(string(output), 1)
			expected := strings.TrimSpace(loadFixture(t, tt.fixture))

      fmt.Println("\nactual:\n", actual)
      fmt.Println("\nexpected:\n", expected)

			if !reflect.DeepEqual(actual, expected) {
        t.Fatalf("actual: \n'%s'\n, expected: \n'%s'\n", actual, expected)
			}
		})
	}
}

func removeLastLine(str string, num int) string {
	lines := strings.Split(str, "\n")

	if len(lines) > 0 {
		lines = lines[:len(lines) - num]
	}

	return strings.Join(lines, "\n")
}

func writeFixture(t *testing.T, goldenFile string, actual []byte) {
	t.Helper()
	goldenPath := "testdata/" + goldenFile

	f, err := os.OpenFile(goldenPath, os.O_RDWR, 0644)
	defer f.Close()

  _, err = f.WriteString(string(actual))

  if err != nil {
    t.Fatalf("Error writing to file %s: %s", goldenPath, err)
  }
}

func loadFixture(t *testing.T, goldenFile string) string {
	goldenPath := "testdata/" + goldenFile

  f, err := os.OpenFile(goldenPath, os.O_RDWR, 0644)

  content, err := ioutil.ReadAll(f)
  if err != nil {
    t.Fatalf("Error opening file %s: %s", goldenPath, err)
  }

  return string(content)
}
