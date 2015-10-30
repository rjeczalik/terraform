package path

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestPathOrContents_Path(t *testing.T) {
	isPath := true
	f, err := ioutil.TempFile("", "tf")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(f.Name())

	if _, err := io.WriteString(f, "foobar"); err != nil {
		t.Fatalf("err: %s", err)
	}
	f.Close()

	contents, wasPath, err := PathOrContents(f.Name())

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}
	if contents != "foobar" {
		t.Fatalf("expected contents %s, got %s", "foobar", contents)
	}
}

func TestPathOrContents_TildePath(t *testing.T) {
	isPath := true
	home, err := homedir.Dir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	f, err := ioutil.TempFile(home, "tf")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(f.Name())

	if _, err := io.WriteString(f, "foobar"); err != nil {
		t.Fatalf("err: %s", err)
	}
	f.Close()

	r := strings.NewReplacer(home, "~")
	homePath := r.Replace(f.Name())
	contents, wasPath, err := PathOrContents(homePath)

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}
	if contents != "foobar" {
		t.Fatalf("expected contents %s, got %s", "foobar", contents)
	}
}

func TestPathOrContents_PathNoPermission(t *testing.T) {
	isPath := true
	f, err := ioutil.TempFile("", "tf")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(f.Name())

	if _, err := io.WriteString(f, "foobar"); err != nil {
		t.Fatalf("err: %s", err)
	}
	f.Close()

	if err := os.Chmod(f.Name(), 0); err != nil {
		t.Fatalf("err: %s", err)
	}

	contents, wasPath, err := PathOrContents(f.Name())

	if err == nil {
		t.Fatal("Expected error, got none!")
	}
	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}
	if contents != "" {
		t.Fatalf("expected contents %s, got %s", "", contents)
	}
}

func TestPathOrContents_Contents(t *testing.T) {
	isPath := false
	input := "hello"

	contents, wasPath, err := PathOrContents(input)

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}
	if contents != input {
		t.Fatalf("expected contents %s, got %s", input, contents)
	}
}

func TestPathOrContents_TildeContents(t *testing.T) {
	isPath := false
	input := "~/hello/notafile"

	contents, wasPath, err := PathOrContents(input)

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if wasPath != isPath {
		t.Fatalf("expected wasPath: %t, got %t", isPath, wasPath)
	}
	if contents != input {
		t.Fatalf("expected contents %s, got %s", input, contents)
	}
}

func testTempFile(t *testing.T, baseDir *string) (*os.File, func()) {
	f, err := ioutil.TempFile("", "tf")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return f, func() {
		os.Remove(f.Name())
	}
}
