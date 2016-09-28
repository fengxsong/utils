package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
)

const (
	B float64 = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
	EB
)

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func readableBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	return fmt.Sprintf(f+"%s", val, suffix)
}

func ReadableFileSize(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return readableBytes(s, 1024, sizes)
}

func FileMTime(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.ModTime().Unix(), nil
}

func FileSize(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

func Copy(src, dest string) error {
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(target, dest)
	}
	sr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sr.Close()

	dw, err := os.Open(dest)
	if err != nil {
		return err
	}
	defer dw.Close()

	if _, err = io.Copy(dw, sr); err != nil {
		return err
	}
	if err = os.Chtimes(dest, si.ModTime(), si.ModTime()); err != nil {
		return err
	}
	return os.Chmod(dest, si.Mode())
}

func WriteFile(filename string, msg []byte) error {
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, msg, 0655)
}

func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ToBytes(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func ToString(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", nil
	}
	return string(b), nil
}
