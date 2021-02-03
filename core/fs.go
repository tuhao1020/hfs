package core

import (
	"fmt"
	"os"
)

type FileInfo struct {
	os.FileInfo
	Url string
	ViewName string
}

type NavInfo struct {
	Label string
	Url string
}

type PageData struct {
	Files []FileInfo
	Navs []NavInfo
}

func (f FileInfo) formattedSize() (string, string) {
	if f.IsDir() {
		return "-", "-"
	}

	var sizeStr, unitStr string
	s := f.Size()
	if s < 1024 {
		sizeStr = fmt.Sprintf("%6d", s)
		unitStr = "B"
	} else if s < 1024 * 1024 {
		sizeStr = fmt.Sprintf("%6.1f", float64(s)/1024)
		unitStr = "KB"
	} else if s < 1024 * 1024 * 1024 {
		sizeStr = fmt.Sprintf("%6.1f", float64(s)/1024/1024)
		unitStr = "MB"
	} else if s < 1024 * 1024 * 1024 * 1024 {
		sizeStr = fmt.Sprintf("%6.1f", float64(s)/1024/1024/1024)
		unitStr = "T"
	}
	return sizeStr, unitStr
}

func (f FileInfo) ReadableSize() string {
	s, _ := f.formattedSize()
	return s
}

func (f FileInfo) ReadableUnit() string {
	_, unit := f.formattedSize()
	return unit
}

// Instead of path.Base(filepath.ToSlash(s))
// let's do something like that, it is faster
// (used to list directories on serve-time too):
func ToBaseName(s string) string {
	n := len(s) - 1
	for i := n; i >= 0; i-- {
		if c := s[i]; c == '/' || c == '\\' {
			if i == n {
				// "s" ends with a slash, remove it and retry.
				return ToBaseName(s[:n])
			}

			return s[i+1:] // return the rest, trimming the slash.
		}
	}
	return s
}

func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
