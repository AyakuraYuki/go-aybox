package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func PathExist(path string) bool {
	exist, _ := PathExistE(path)
	return exist
}

func PathExistE(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SaveTextToFile(dst, content string) {
	f, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	buf := bufio.NewWriter(f)
	_, _ = fmt.Fprintln(buf, content)
	_ = buf.Flush()
}

func AppendTextToFile(dst, content string) {
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	buf := bufio.NewWriter(f)
	_, _ = fmt.Fprintln(buf, content)
	_ = buf.Flush()
}

func SaveLinesToFile(dst string, lines []string) {
	f, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	buf := bufio.NewWriter(f)
	for _, line := range lines {
		_, _ = fmt.Fprintln(buf, line)
	}
	_ = buf.Flush()
}

func ReadBytes(src string) []byte {
	f, err := os.Open(src)
	if err != nil {
		return nil
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	r, err := io.ReadAll(f)
	if err != nil {
		return nil
	}
	return r
}

func ReadTextFile(src string) string {
	bs := ReadBytes(src)
	if bs == nil {
		return ""
	}
	return string(bs)
}

func Readlines(src string) (lines []string) {
	lines = make([]string, 0)
	f, err := os.Open(src)
	if err != nil {
		return
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			lines = append(lines, line)
		}
		if err == io.EOF {
			break
		}
	}
	return
}

func ListFileNoRecursive(dir string) (files []string, err error) {
	files = make([]string, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // ignore sub dir
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}
	return
}

func ListFile(dir string) (files []string) {
	files = make([]string, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			files = append(files, ListFile(filepath.Join(dir, entry.Name()))...)
		} else {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	return
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}
