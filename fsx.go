package aferox

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

var _ afero.Fs = &Fsx{}

// Fsx adjusts all relative paths based on the stored
// working directory, instead of relying on the default behavior for relative
// paths defined by the implementing Fs.
type Fsx struct {
	fs afero.Fs

	dir string
}

func NewFsx(dir string, fs afero.Fs) *Fsx {
	pwd, _ := filepath.Abs(dir)
	return &Fsx{
		dir: pwd,
		fs:  fs,
	}
}

// Getwd returns a rooted path name corresponding to the current directory.
func (f *Fsx) Getwd() string {
	return f.dir
}

// Chdir changes the current working directory to the named directory.
func (f *Fsx) Chdir(dir string) {
	f.dir = f.Abs(dir)
}

// Chown changes the uid and gid of the named file.
func (f *Fsx) Chown(name string, uid, gid int) error {
	return f.fs.Chown(name, uid, gid)
}

// Abs returns an absolute representation of path. If the path is not absolute
// it will be joined with the current working directory to turn it into an
// absolute path. The absolute path name for a given file is not guaranteed to
// be unique. Abs calls Clean on the result.
func (f *Fsx) Abs(path string) string {
	var fullPath string
	if filepath.IsAbs(path) {
		fullPath = path
	} else {
		prefix := f.dir
		// On Windows /foo resolves to DRIVEPATH:\foo, so treat anything that starts with a slash as absolute that just needs cleaning up
		if strings.HasPrefix(path, `/`) || strings.HasPrefix(path, `\`) {
			prefix, _ = filepath.Abs("/")
		}
		fullPath = filepath.Join(prefix, path)
	}
	return filepath.Clean(fullPath)
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
func (f *Fsx) Create(name string) (afero.File, error) {
	name = f.Abs(name)
	return f.fs.Create(name)
}

// Mkdir creates a new directory with the specified name and permission
// bits (before umask).
func (f *Fsx) Mkdir(name string, perm os.FileMode) error {
	name = f.Abs(name)
	return f.fs.Mkdir(name, perm)
}

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func (f *Fsx) MkdirAll(path string, perm os.FileMode) error {
	path = f.Abs(path)
	return f.fs.MkdirAll(path, perm)
}

// OpenFile is the generalized open call; most users will use Open
// or Create instead. It opens the named file with specified flag
// (O_RDONLY etc.). If the file does not exist, and the O_CREATE flag
// is passed, it is created with mode perm (before umask). If successful,
// methods on the returned File can be used for I/O.
func (f *Fsx) Open(name string) (afero.File, error) {
	name = f.Abs(name)
	return f.fs.Open(name)
}

// OpenFile opens a file using the given flags and the given mode.
func (f *Fsx) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	name = f.Abs(name)
	return f.fs.OpenFile(name, flag, perm)
}

// Remove removes the named file or (empty) directory.
func (f *Fsx) Remove(name string) error {
	name = f.Abs(name)
	return f.fs.Remove(name)
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
func (f *Fsx) RemoveAll(path string) error {
	path = f.Abs(path)
	return f.fs.RemoveAll(path)
}

// Rename renames (moves) oldpath to newpath.
// If newpath already exists and is not a directory, Rename replaces it.
// OS-specific restrictions may apply when oldpath and newpath are in different directories.
func (f *Fsx) Rename(oldname, newname string) error {
	oldname = f.Abs(oldname)
	newname = f.Abs(newname)
	return f.fs.Rename(oldname, newname)
}

// Stat returns a FileInfo describing the named file.
func (f *Fsx) Stat(name string) (os.FileInfo, error) {
	name = f.Abs(name)
	return f.fs.Stat(name)
}

// The name of this FileSystem.
func (f *Fsx) Name() string {
	return "Fsx"
}

// Chmod changes the mode of the named file to mode.
// If the file is a symbolic link, it changes the mode of the link's target.
// If there is an error, it will be of type *PathError.
//
// A different subset of the mode bits are used, depending on the
// operating system.
func (f *Fsx) Chmod(name string, mode os.FileMode) error {
	name = f.Abs(name)
	return f.fs.Chmod(name, mode)
}

// Chtimes changes the access and modification times of the named
// file, similar to the Unix utime() or utimes() functions.
func (f *Fsx) Chtimes(name string, atime time.Time, mtime time.Time) error {
	name = f.Abs(name)
	return f.fs.Chtimes(name, atime, mtime)
}
