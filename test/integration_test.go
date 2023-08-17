package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

	fmt.Println("Starting command...")
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second)

	_, err = stdin.Write(bytesToWrite)

	fmt.Println("Wrote to stdin", bytesToWrite)

	if err != nil {
		fmt.Println("Error writing to stdin", err)
	}

	time.Sleep(14 * time.Second)

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

	return out, err
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		bytesToWrite []byte
	}{
		{"search test", []string{"search", "Eloquent JavaScript", "-l", "faster"}, []byte{14, 14, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinaryWithFileInput(tt.args, tt.bytesToWrite)

			fmt.Println(string(output))
			if _, err := os.Stat("eloquent-javascript-a-modern-introduction-to-programming.epub"); err == nil {
				fmt.Printf("File exists\n")
			} else {
				t.Fatal("File wasn't found")
			}

			if err != nil {
				fmt.Println("error", err.Error())
				t.Fatal(err)
			}
		})
	}
}
