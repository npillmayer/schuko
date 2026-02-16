package schuko

// TODO include XDG
// for example:   https://github.com/adrg/xdg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilepathMatch(t *testing.T) {
	if !fm("abc.*", "abc.yml") {
		t.Error("no match")
	}
}

func TestLocateConfig(t *testing.T) {
	appTag, tag := "MyApp", "myapp"

	// setup
	tmpdir := os.TempDir()
	configPath := filepath.Join(tmpdir, "abc-98qewqw", ".config", tag)
	err := os.MkdirAll(configPath, os.ModeDir|0770)
	if err != nil {
		t.Fatal(err)
	}
	fname := filepath.Join(configPath, tag+".yml")
	f, err := os.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("created fake config file %q", f.Name())
	path := []string{tmpdir, "abc-98qewqw", ".config", tag}
	defer eraser(fname, f, path)
	os.Setenv("HOME", filepath.Join(tmpdir, "abc-98qewqw"))
	H, _ := os.UserHomeDir()
	t.Logf("HOME = %q", H)

	// now search
	suffixes := []string{"json", "toml", "nt", "yml"}
	fileNames := LocateConfig(appTag, "", suffixes)
	if len(fileNames) == 0 {
		t.Errorf("expected 1 config file to be found, didn't")
	} else {
		for _, fn := range fileNames {
			t.Logf("found config: %q", fn)
		}
	}
	suffixes = []string{"yml"}
	fileNames = LocateConfig(appTag, "my*.*", suffixes)
	if len(fileNames) == 0 {
		t.Errorf("expected 1 pattern file to be found, didn't")
	} else {
		for _, fn := range fileNames {
			t.Logf("found config: %q", fn)
		}
	}
	fileNames = LocateConfig(appTag, "your*.*", suffixes)
	if len(fileNames) == 0 {
		t.Errorf("expected no pattern file to be found")
	}
}

func eraser(fname string, f *os.File, path []string) {
	f.Close()
	os.Remove(fname)
	for i := len(path); i > 1; i-- {
		p := filepath.Join(path[:i]...)
		os.Remove(p)
	}
}
