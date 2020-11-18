package aferox

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

type Aferox struct {
	*afero.Afero
	wrapper *FsWd
}

func NewAferox(dir string, fs afero.Fs) Aferox {
	wrapper := NewFsWd(dir, fs)
	return Aferox{
		Afero:   &afero.Afero{Fs: wrapper},
		wrapper: wrapper,
	}
}

func (a Aferox) Getwd() string {
	return a.wrapper.Getwd()
}

func (a Aferox) Setwd(dir string) {
	a.wrapper.Setwd(dir)
}

func (a Aferox) Abs(path string) string {
	return a.wrapper.Abs(path)
}

// This is a simplified exec.LookPath that checks if command is accessible given
// a PATH environment variable.
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

var _ afero.Fs = &FsWd{}

// FsWd adjusts all relative paths based on the stored
// working directory, instead of relying on the default behavior for relative
// paths defined by the implementing Fs.
type FsWd struct {
	fs afero.Fs

	dir string
}

func NewFsWd(dir string, fs afero.Fs) *FsWd {
	return &FsWd{
		dir: dir,
		fs:  fs,
	}
}

func (f *FsWd) Getwd() string {
	return f.dir
}

func (f *FsWd) Setwd(dir string) {
	f.dir = dir
}

func (f *FsWd) Abs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	path = filepath.Clean(path)
	return filepath.Join(f.dir, path)
}

func (f *FsWd) Create(name string) (afero.File, error) {
	name = f.absolute(name)
	return f.fs.Create(name)
}

func (f *FsWd) Mkdir(name string, perm os.FileMode) error {
	name = f.absolute(name)
	return f.fs.Mkdir(name, perm)
}

func (f *FsWd) MkdirAll(path string, perm os.FileMode) error {
	path = f.absolute(path)
	return f.fs.MkdirAll(path, perm)
}

func (f *FsWd) Open(name string) (afero.File, error) {
	name = f.absolute(name)
	return f.fs.Open(name)
}

func (f *FsWd) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	name = f.absolute(name)
	return f.fs.OpenFile(name, flag, perm)
}

func (f *FsWd) Remove(name string) error {
	name = f.absolute(name)
	return f.fs.Remove(name)
}

func (f *FsWd) RemoveAll(path string) error {
	path = f.absolute(path)
	return f.fs.RemoveAll(path)
}

func (f *FsWd) Rename(oldname, newname string) error {
	oldname = f.absolute(oldname)
	newname = f.absolute(newname)
	return f.fs.Rename(oldname, newname)
}

func (f *FsWd) Stat(name string) (os.FileInfo, error) {
	name = f.absolute(name)
	return f.fs.Stat(name)
}

func (f *FsWd) Name() string {
	return "FsWd"
}

func (f *FsWd) Chmod(name string, mode os.FileMode) error {
	name = f.absolute(name)
	return f.fs.Chmod(name, mode)
}

func (f *FsWd) Chtimes(name string, atime time.Time, mtime time.Time) error {
	name = f.absolute(name)
	return f.fs.Chtimes(name, atime, mtime)
}

func (f *FsWd) absolute(path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(f.dir, path)
	}

	return path
}
