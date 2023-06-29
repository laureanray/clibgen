package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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
  cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
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
