package testutil

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Mkdir(t *testing.T, dir string) (cleanup func()) {
	t.Helper()

	err := os.MkdirAll(dir, 0755)
	require.NoError(t, err)
	return func() {
		err := os.RemoveAll(dir)
		require.NoError(t, err)
	}
}

func RemoveAll(t *testing.T, dir string) {
	t.Helper()

	err := os.RemoveAll(dir)
	require.NoError(t, err)
}

func ReadDir(t *testing.T, dir string) []os.FileInfo {
	t.Helper()

	infos, err := ioutil.ReadDir(dir)
	require.NoError(t, err)
	return infos
}

func WriteFile(t *testing.T, path string, buf []byte) (cleanup func()) {
	t.Helper()

	err := ioutil.WriteFile(path, buf, 0644)
	require.NoError(t, err)
	return func() {
		err := os.Remove(path)
		require.NoError(t, err)
	}
}

func ReadFile(t *testing.T, path string) []byte {
	t.Helper()

	buf, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return buf
}

func CopyFile(t *testing.T, srcName, dstName string) (cleanup func()) {
	t.Helper()

	_, err := os.Stat(srcName)
	if err != nil {
		t.Logf("%s does not exist", srcName)
		return
	}

	src, err := os.Open(srcName)
	require.NoError(t, err)
	dst, err := os.Create(dstName)
	require.NoError(t, err)
	_, err = io.Copy(dst, src)
	require.NoError(t, err)

	return func() {
		err := os.Remove(dstName)
		require.NoError(t, err)
	}
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
