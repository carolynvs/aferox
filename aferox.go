package aferox

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// Aferox adjusts all relative paths based on the stored
// working directory, instead of relying on the default behavior for relative
// paths defined by the implementing Fs.
type Aferox struct {
	*afero.Afero

	// Fs is the working directory aware filesystem.
	Fs *Fsx
}

// NewAferox creates a wrapper around a filesystem representation with
// an independent working directory.
func NewAferox(dir string, fs afero.Fs) Aferox {
	wrapper := NewFsx(dir, fs)
	return Aferox{
		Afero: &afero.Afero{Fs: wrapper},
		Fs:    wrapper,
	}
}

// Getwd returns a rooted path name corresponding to the current directory.
// Use in place of os.Getwd.
func (a Aferox) Getwd() string {
	return a.Fs.Getwd()
}

// Chdir changes the current working directory to the named directory.
// Use in place of os.Chdir.
func (a Aferox) Chdir(dir string) {
	a.Fs.Chdir(dir)
}

// Abs returns an absolute representation of path. If the path is not absolute
// it will be joined with the current working directory to turn it into an
// absolute path. The absolute path name for a given file is not guaranteed to
// be unique. Abs calls Clean on the result.
// Use in place of filepath.Abs.
func (a Aferox) Abs(path string) string {
	return a.Fs.Abs(path)
}

// This is a simplified exec.LookPath that checks if command is accessible given
// a PATH environment variable.
// Use in place of exec.LookPath.
func (a Aferox) LookPath(cmd string, path string) (string, bool) {
	paths := strings.Split(path, string(os.PathListSeparator))
	for _, p := range paths {
		files, err := a.ReadDir(p)
		if err != nil {
			continue
		}

		for _, f := range files {
			if f.Name() != cmd {
				continue
			}

			// Return if the file is executable MAYBE
			// Simplified check, we aren't checking if it's executable by the current user
			executable := f.Mode()&0111 != 0
			if executable {
				return filepath.Join(p, f.Name()), true
			}
		}
	}

	return "", false
}
