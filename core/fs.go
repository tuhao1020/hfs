package core

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type FileInfo struct {
	os.FileInfo
	Url      string
	ViewName string
}

type NavInfo struct {
	Label string
	Url   string
}

type PageData struct {
	Files []FileInfo
	Navs  []NavInfo
}

const (
	SizeFormat = "%6.1f"
)

func (f FileInfo) formattedSize() (string, string) {
	if f.IsDir() {
		return "-", "-"
	}

	var sizeStr, unitStr string
	s := float64(f.Size())
	if s < iris.KB {
		sizeStr = fmt.Sprintf(SizeFormat, s)
		unitStr = "B"
	} else if s < 1024*iris.KB {
		sizeStr = fmt.Sprintf(SizeFormat, s/iris.KB)
		unitStr = "KB"
	} else if s < 1024*iris.MB {
		sizeStr = fmt.Sprintf(SizeFormat, s/iris.MB)
		unitStr = "MB"
	} else if s < 1024*iris.GB {
		sizeStr = fmt.Sprintf(SizeFormat, s/iris.GB)
		unitStr = "GB"
	} else if s < 1024*iris.TB {
		sizeStr = fmt.Sprintf(SizeFormat, s/iris.TB)
		unitStr = "TB"
	} else if s < 1024*iris.PB {
		sizeStr = fmt.Sprintf(SizeFormat, s/iris.PB)
		unitStr = "PB"
	} else if s < 1024*iris.EB {
		sizeStr = fmt.Sprintf(SizeFormat, s/iris.EB)
		unitStr = "EB"
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

func web2LocalPath(webPath string) string {
	if strings.HasPrefix(webPath, "/") {
		return "." + webPath
	}
	return "./" + webPath
}

func MoveTo(src, dest string) error {
	srcLocal := web2LocalPath(src)
	destLocal := web2LocalPath(dest)
	if !IsExist(srcLocal) {
		return errors.New(srcLocal + " not exist")
	}

	srcInfo, err := os.Stat(srcLocal)
	if srcInfo.IsDir() {
		err = os.Rename(srcLocal, destLocal+"/"+srcInfo.Name())
	} else {
		err = os.Rename(srcLocal, destLocal)
	}

	return err
}
