package aferox

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

var _ afero.Fs = &FsWd{}

// FsWd adjusts all relative paths based on the stored
// working directory, instead of relying on the default behavior for relative
// paths defined by the implementing Fs.
type FsWd struct {
	dir string
	fs afero.Fs
}

func NewFsFd(dir string, fs afero.Fs) *FsWd {
	return &FsWd{
		dir: dir,
		fs: fs,
	}
}

func (o *FsWd) Getwd() string {
	return o.dir
}

func (o *FsWd) Setwd(dir string) {
	o.dir = dir
}

func (o *FsWd) Create(name string) (afero.File, error) {
	name = o.absolute(name)
	return o.fs.Create(name)
}

func (o *FsWd) Mkdir(name string, perm os.FileMode) error {
	name = o.absolute(name)
	return o.fs.Mkdir(name, perm)
}

func (o *FsWd) MkdirAll(path string, perm os.FileMode) error {
	path = o.absolute(path)
	return o.fs.MkdirAll(path, perm)
}

func (o *FsWd) Open(name string) (afero.File, error) {
	name = o.absolute(name)
	return o.fs.Open(name)
}

func (o *FsWd) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	name = o.absolute(name)
	return o.fs.OpenFile(name, flag, perm)
}

func (o *FsWd) Remove(name string) error {
	name = o.absolute(name)
	return o.fs.Remove(name)
}

func (o *FsWd) RemoveAll(path string) error {
	path = o.absolute(path)
	return o.fs.RemoveAll(path)
}

func (o *FsWd) Rename(oldname, newname string) error {
	oldname = o.absolute(oldname)
	newname = o.absolute(newname)
	return o.fs.Rename(oldname, newname)
}

func (o *FsWd) Stat(name string) (os.FileInfo, error) {
	name = o.absolute(name)
	return o.fs.Stat(name)
}

func (o *FsWd) Name() string {
	return "FsWd"
}

func (o *FsWd) Chmod(name string, mode os.FileMode) error {
	name = o.absolute(name)
	return o.fs.Chmod(name, mode)
}

func (o *FsWd) Chtimes(name string, atime time.Time, mtime time.Time) error {
	name = o.absolute(name)
	return o.fs.Chtimes(name, atime, mtime)
}

func (o *FsWd) absolute(path string) string {
	if !filepath.IsAbs(path){
		path = filepath.Join(o.dir, path)
	}

	return path
}