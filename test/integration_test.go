package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
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

  os.Exit(m.Run())
}


func runBinary(args []string) ([]byte, error) {
  cmd := exec.Command(binaryPath, args...)
  // Set the command's output to the standard output of the current process
  // cmd.Stdout = os.Stdout
  // cmd.Stderr = os.Stderr
  cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")

  // Start the command
  // err := cmd.Start()
  // if err != nil {
  //   fmt.Printf("Failed to start command: %s\n", err)
  //   return nil, err
  // }

  // Set up a signal channel to listen for interrupts
  sigChan := make(chan os.Signal, 1)
  signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

  go func() {
    // Wait for a signal
    <-sigChan

    // Send the interrupt signal to the command's process
    cmd.Process.Signal(os.Interrupt)
    }()

  // Wait for the command to complete
  err := cmd.Wait()
  if err != nil {
    fmt.Printf("Command failed: %s\n", err)
  }
  return cmd.CombinedOutput()
}

func TestCliArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
	}{
		// {"no arguments", []string{}, "no-args.golden"},
    {"help args", []string{"--help"}, "help.golden"},
    {"search test", []string{"search", "Eloquent JavaScript"}, "eloquent.golden"},
    
		// {"one argument", []string{"ciao"}, "one-argument.golden"},
		// {"multiple arguments", []string{"ciao", "hello"}, "multiple-arguments.golden"},
		// {"shout arg", []string{"--shout", "ciao"}, "shout-arg.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args)

			if err != nil {
				t.Fatal(err)
			}

      fmt.Println(string(output))

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

// func TestSearch(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		args    []string
// 		fixture string
// 	}{
// 		// {"no arguments", []string{}, "no-args.golden"},
//     {"search test", []string{"search \"Eloquent JavaScript\""}, "eloquent.golden"},
// 		// {"one argument", []string{"ciao"}, "one-argument.golden"},
// 		// {"multiple arguments", []string{"ciao", "hello"}, "multiple-arguments.golden"},
// 		// {"shout arg", []string{"--shout", "ciao"}, "shout-arg.golden"},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			output, err := runBinary(tt.args)
//
// 			if err != nil {
// 				t.Fatal(err)
// 			}
//
//       fmt.Println(string(output))
//
// 			if *update {
// 				writeFixture(t, tt.fixture, output)
// 			}
//
// 			actual := string(output)
//
// 			expected := loadFixture(t, tt.fixture)
//
// 			if !reflect.DeepEqual(actual, expected) {
// 				t.Fatalf("actual = %s, expected = %s", actual, expected)
// 			}
// 		})
// 	}
// }

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
