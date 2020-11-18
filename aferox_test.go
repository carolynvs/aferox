package aferox

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFsWd_Getwd(t *testing.T) {
	f := NewFsFd("/home", afero.NewMemMapFs())
	pwd := f.Getwd()
	assert.Equal(t, "/home", pwd)
}

func TestFsWd_Setwd(t *testing.T) {
	f := NewFsFd("/home", afero.NewMemMapFs())
	f.Setwd("/bin")
	pwd := f.Getwd()
	assert.Equal(t, "/bin", pwd)
}

func TestFsWd_ioutil(t *testing.T) {
	backend := &afero.Afero{Fs:afero.NewMemMapFs()}
	backend.Mkdir("/home", 0755)
	backend.WriteFile("/home/user.txt", []byte("sally"), 0644)

	f := NewFsFd("/home", backend)

	contents, err := f.ReadFile("/home/user.txt")
	require.NoError(t, err)
	assert.Equal(t, "sally", contents)
}
