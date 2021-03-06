package aferox

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFsx_Create(t *testing.T) {
	f := NewAferox("/home", afero.NewMemMapFs())
	err := f.Mkdir("/home", 0755)
	require.NoError(t, err, "Mkdir failed")

	_, err = f.Create("user.txt")
	require.NoError(t, err)
	exists, _ := f.Exists("/home/user.txt")
	assert.True(t, exists)
}

func TestFsx_Chmod(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	filename := "/root/.ssh/id_pem"
	_, err := f.Create(filename)
	require.NoError(t, err, "Create failed")

	var wantMode os.FileMode = 0600
	err = f.Chmod(filename, wantMode)
	require.NoError(t, err, "Chmod failed")

	fi, err := f.Stat(filename)
	require.NoError(t, err, "Stat failed")
	assert.Equal(t, fi.Mode()&wantMode, wantMode)
}

func TestFsx_Chown(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("This test fails on other OS's without running as root")
	}

	tmp, err := ioutil.TempDir("", "aferox")
	require.NoError(t, err, "create temp directory failed")
	defer os.RemoveAll(tmp)

	// Testing against the OS because that's the only way to get at the uid/gid
	f := NewFsx("", afero.NewOsFs())

	filename := filepath.Join(tmp, "myfile.txt")
	_, err = f.Create(filename)
	require.NoError(t, err, "Create failed")

	var wantUid = 0
	var wantGid = 0
	err = f.Chown(filename, wantUid, wantGid)

	require.NoError(t, err, "Chown failed")

	fi, err := f.Stat(filename)
	require.NoError(t, err, "Stat failed")
	stat, ok := fi.Sys().(syscall.Stat_t)
	if !ok {
		t.Skip("the current OS doesn't support extended stat info with uid/gid")
	}

	gotUid := int(stat.Uid)
	gotGid := int(stat.Gid)
	assert.Equal(t, wantUid, gotUid, "incorrect uid")
	assert.Equal(t, wantGid, gotGid, "incorrect gid")
}

func TestFsx_Chtimes(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	filename := "test.sh"
	_, err := f.Create(filename)
	require.NoError(t, err, "Create failed")

	sometime := time.Now().Add(time.Hour)
	err = f.Chtimes(filename, sometime, sometime)
	require.NoError(t, err, "Chtimes failed")

	fi, err := f.Stat(filename)
	require.NoError(t, err, "Stat failed")
	assert.Equal(t, sometime, fi.ModTime())
}

func TestFsx_Mkdir(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	dir := "me"
	var wantMode os.FileMode = 0755
	err := f.MkdirAll(dir, wantMode)
	require.NoError(t, err, "Mkdir failed")

	fi, err := f.Stat(dir)
	require.NoError(t, err, "Stat failed")
	assert.Equal(t, fi.Mode()&wantMode, wantMode)
}

func TestFsx_MkdirAll(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	dir := "/tmp/aferox"
	var wantMode os.FileMode = 0755
	err := f.MkdirAll(dir, wantMode)
	require.NoError(t, err, "Mkdir failed")

	fi, err := f.Stat(dir)
	require.NoError(t, err, "Stat failed")
	assert.Equal(t, fi.Mode()&wantMode, wantMode)
}

func TestFsx_Name(t *testing.T) {
	f := Fsx{}
	assert.Equal(t, "Fsx", f.Name())
}

func TestFsx_Open(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	filename := "test.txt"
	_, err := f.Create(filename)
	require.NoError(t, err, "Create failed")

	fi, err := f.Open(filename)
	require.NoError(t, err, "Open failed")
	assert.Equal(t, xplat("/home/test.txt"), fi.Name())
}

func TestFsx_OpenFile(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	var wantMode os.FileMode = 0644
	file, err := f.OpenFile("test.txt", os.O_CREATE, wantMode)
	require.NoError(t, err, "OpenFile failed")

	fi, err := file.Stat()
	require.NoError(t, err, "Stat failed")
	assert.Equal(t, "test.txt", fi.Name())
	assert.Equal(t, wantMode, fi.Mode()&wantMode)
}

func TestFsx_Remove(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	filename := "test.txt"
	_, err := f.Create(filename)
	require.NoError(t, err, "Create failed")

	err = f.Remove(filename)
	require.NoError(t, err, "Remove failed")

	_, err = f.Stat(filename)
	assert.True(t, os.IsNotExist(err))
}

func TestFsx_RemoveAll(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	_, err := f.Create("test1.txt")
	require.NoError(t, err, "Create test1.txt failed")

	_, err = f.Create("test2.txt")
	require.NoError(t, err, "Create test2.txt failed")

	err = f.RemoveAll("/home")
	require.NoError(t, err, "Remove failed")

	_, err = f.Stat("/home")
	assert.True(t, os.IsNotExist(err))
}

func TestFsx_Rename(t *testing.T) {
	f := NewFsx("/home", afero.NewMemMapFs())

	_, err := f.Create("test1.txt")
	require.NoError(t, err, "Create test1.txt failed")

	err = f.Rename("test1.txt", "test2.txt")
	require.NoError(t, err, "Rename failed")

	fi, err := f.Stat("test2.txt")
	require.NoError(t, err, "Stat failed")
	assert.Equal(t, "test2.txt", fi.Name())
}

func TestFsx_Abs(t *testing.T) {
	f := NewFsx("/", afero.NewMemMapFs())

	path := f.Abs("test")
	wantPath, _ := filepath.Abs("/test")
	assert.Equal(t, wantPath, path)
}
