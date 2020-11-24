package aferox

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Getwd(t *testing.T) {
	a := NewAferox("/home", afero.NewMemMapFs())
	f := NewFsx("/home", afero.NewMemMapFs())

	pwd := a.Getwd()
	assert.Equal(t, "/home", pwd)

	pwd = f.Getwd()
	assert.Equal(t, "/home", pwd)
}

func Test_Chdir(t *testing.T) {
	a := NewAferox("/home", afero.NewMemMapFs())
	f := NewFsx("/home", afero.NewMemMapFs())

	a.Chdir("/bin")
	pwd := a.Getwd()
	assert.Equal(t, "/bin", pwd)

	f.Chdir("/bin")
	pwd = f.Getwd()
	assert.Equal(t, "/bin", pwd)
}

func Test_Abs(t *testing.T) {
	a := NewAferox("/home", afero.NewMemMapFs())
	fs := NewFsx("/home", afero.NewMemMapFs())

	t.Run("empty", func(t *testing.T) {
		p := fs.Abs("")
		assert.Equal(t, "/home", p)

		p = a.Abs("")
		assert.Equal(t, "/home", p)
	})

	t.Run("relative", func(t *testing.T) {
		p := fs.Abs("me")
		assert.Equal(t, "/home/me", p)

		p = a.Abs("me")
		assert.Equal(t, "/home/me", p)
	})

	t.Run("absolute", func(t *testing.T) {
		p := fs.Abs("/tmp")
		assert.Equal(t, "/tmp", p)

		p = fs.Abs("/tmp")
		assert.Equal(t, "/tmp", p)
	})
}

func TestAferox_LookPath(t *testing.T) {
	t.Run("osfs", func(t *testing.T) {
		pwd, err := os.Getwd()
		require.NoError(t, err, "Getwd failed")

		f := NewAferox(pwd, afero.NewOsFs())
		goPath, hasGo := f.LookPath("go", os.Getenv("PATH"))
		require.True(t, hasGo)
		assert.NotEmpty(t, goPath)
	})

	t.Run("memfs", func(t *testing.T) {
		f := NewAferox("/home", afero.NewMemMapFs())

		// /usr/local/bin not executable
		_, err := f.Create("/usr/local/bin/go")
		require.NoError(t, err, "Create failed")

		// /bin/go executable
		_, err = f.Create("/bin/go")
		require.NoError(t, err, "Create failed")
		err = f.Chmod("/bin/go", 0744)
		require.NoError(t, err, "Chmod faild")

		path := strings.Join([]string{"/home/bin", "/usr/local/bin", "/bin", "/home/go/bin"}, string(os.PathListSeparator))
		goPath, hasGo := f.LookPath("go", path)
		require.True(t, hasGo)
		assert.Equal(t, "/bin/go", goPath)
	})
}

func TestAferox_ReadDir(t *testing.T) {
	a := NewAferox("/home", afero.NewMemMapFs())

	err := a.WriteFile("/home/homefile.txt", nil, 0644)
	require.NoError(t, err, "WriteFile failed for /home/homefile.txt")

	err = a.WriteFile("/home/me/mefile.txt", nil, 0644)
	require.NoError(t, err, "WriteFile failed for /home/me/mefile.txt")

	err = a.WriteFile("/tmp/tmpfile.txt", []byte("tmpfile"), 0644)
	require.NoError(t, err, "WriteFile failed for /tmp/tmpfile.txt")

	t.Run("empty", func(t *testing.T) {
		items, err := a.ReadDir("")

		require.NoError(t, err, "ReadDir failed")
		assert.Len(t, items, 2, "expected 2 children")
		assert.Equal(t, "homefile.txt", items[0].Name())
		assert.Equal(t, "me", items[1].Name())
	})

	t.Run("relative", func(t *testing.T) {
		items, err := a.ReadDir("me")

		require.NoError(t, err, "ReadDir failed")
		assert.Len(t, items, 1, "expected 1 file")
		assert.Equal(t, "mefile.txt", items[0].Name())
	})

	t.Run("absolute", func(t *testing.T) {
		items, err := a.ReadDir("/home/me/")

		require.NoError(t, err, "ReadDir failed")
		assert.Len(t, items, 1, "expected 1 file")
		assert.Equal(t, "mefile.txt", items[0].Name())
	})
}

func TestAferox_ReadFile(t *testing.T) {
	a := NewAferox("/home", afero.NewMemMapFs())

	err := a.WriteFile("/home/me/mefile.txt", []byte("mefile"), 0644)
	require.NoError(t, err, "WriteFile failed for /home/me/mefile.txt")

	err = a.WriteFile("/tmp/tmpfile.txt", []byte("tmpfile"), 0644)
	require.NoError(t, err, "WriteFile failed for /tmp/tmpfile.txt")

	t.Run("relative", func(t *testing.T) {
		file, err := a.ReadFile("me/mefile.txt")

		require.NoError(t, err, "ReadFile failed")
		assert.Equal(t, "mefile", string(file))
	})

	t.Run("absolute", func(t *testing.T) {
		file, err := a.ReadFile("/tmp/tmpfile.txt")

		require.NoError(t, err, "ReadFile failed")
		assert.Equal(t, "tmpfile", string(file))
	})
}

func TestAferox_TempDir(t *testing.T) {
	a := NewAferox("/home", afero.NewMemMapFs())

	t.Run("empty", func(t *testing.T) {
		gotTmp, err := a.TempDir("", "aferox")
		require.NoError(t, err)

		wantTmp := filepath.Join(os.TempDir(), "aferox")
		assert.Contains(t, gotTmp, wantTmp)
	})

	t.Run("relative", func(t *testing.T) {
		gotTmp, err := a.TempDir("me", "aferox")
		require.NoError(t, err)

		wantTmp := "me/aferox"
		assert.Contains(t, gotTmp, wantTmp)
	})

	t.Run("absolute", func(t *testing.T) {
		gotTmp, err := a.TempDir("/etc", "aferox")
		require.NoError(t, err)

		wantTmp := "/etc/aferox"
		assert.Contains(t, gotTmp, wantTmp)
	})
}

func TestAferox_TempFile(t *testing.T) {
	a := NewAferox("/home", afero.NewMemMapFs())

	t.Run("empty", func(t *testing.T) {
		gotTmp, err := a.TempFile("", "aferox")
		require.NoError(t, err)

		wantTmp := filepath.Join(os.TempDir(), "aferox")
		assert.Contains(t, gotTmp.Name(), wantTmp)
	})

	t.Run("relative", func(t *testing.T) {
		gotTmp, err := a.TempFile("me", "aferox")
		require.NoError(t, err)

		wantTmp := "/home/me/aferox"
		assert.Contains(t, gotTmp.Name(), wantTmp)
	})

	t.Run("absolute", func(t *testing.T) {
		gotTmp, err := a.TempFile("/etc", "aferox")
		require.NoError(t, err)

		wantTmp := "/etc/aferox"
		assert.Contains(t, gotTmp.Name(), wantTmp)
	})
}

func TestAferox_WriteFile(t *testing.T) {
	a := NewAferox("/home", afero.NewMemMapFs())

	t.Run("relative", func(t *testing.T) {
		err := a.WriteFile("homefile.txt", []byte("homefile"), 0644)
		require.NoError(t, err, "WriteFile failed")
		file, err := a.ReadFile("homefile.txt")
		require.NoError(t, err, "ReadFile failed")
		assert.Equal(t, "homefile", string(file))
	})

	t.Run("absolute", func(t *testing.T) {
		err := a.WriteFile("/tmp/tmpfile.txt", []byte("tmpfile"), 0644)
		require.NoError(t, err, "WriteFile failed")
		file, err := a.ReadFile("/tmp/tmpfile.txt")
		require.NoError(t, err, "ReadFile failed")
		assert.Equal(t, "tmpfile", string(file))
	})
}
