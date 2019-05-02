package testutil

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Mkdir(t *testing.T, dir string) (cleanup func()) {
	err := os.Mkdir(dir, 0755)
	require.NoError(t, err)
	return func() {
		err := os.RemoveAll(dir)
		require.NoError(t, err)
	}
}

func WriteFile(t *testing.T, path string, buf []byte) (cleanup func()) {
	err := ioutil.WriteFile(path, buf, 0644)
	require.NoError(t, err)
	return func() {
		err := os.Remove(path)
		require.NoError(t, err)
	}
}

func RemoveFile(t *testing.T, path string) {
	err := os.Remove(path)
	require.NoError(t, err)
}

func ReadFile(t *testing.T, path string) []byte {
	buf, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return buf
}

func NormalizeTestName(name string) string {
	r := strings.NewReplacer(
		" ", "_",
		"'", "",
		`"`, "",
		",", "",
	)
	return r.Replace(strings.ToLower(name))
}
