package aferox

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFsWd_Getwd(t *testing.T) {
	f := NewAferox("/home", afero.NewMemMapFs())
	pwd := f.Getwd()
	assert.Equal(t, "/home", pwd)
}

func TestFsWd_Setwd(t *testing.T) {
	f := NewAferox("/home", afero.NewMemMapFs())
	f.Setwd("/bin")
	pwd := f.Getwd()
	assert.Equal(t, "/bin", pwd)
}

func TestFsWd_Abs(t *testing.T) {
	f := NewAferox("/home/me", afero.NewMemMapFs())
	path := f.Abs("../you")
	assert.Equal(t, "/home/you", path)
}

func TestFsWd_ReadFile(t *testing.T) {
	f := NewAferox("/home", afero.NewMemMapFs())
	f.Mkdir("/home", 0755)
	f.WriteFile("/home/user.txt", []byte("sally"), 0644)

	contents, err := f.ReadFile("/home/user.txt")
	require.NoError(t, err)
	assert.Equal(t, "sally", string(contents))

	contents, err = f.ReadFile("user.txt")
	require.NoError(t, err)
	assert.Equal(t, "sally", string(contents))
}

func TestFsWd_Create(t *testing.T) {
	f := NewAferox("/home", afero.NewMemMapFs())
	f.Mkdir("/home", 0755)

	_, err := f.Create("user.txt")
	require.NoError(t, err)
	exists, _ := f.Exists("/home/user.txt")
	assert.True(t, exists)
}

func TestFsWd_LookPath(t *testing.T) {
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
