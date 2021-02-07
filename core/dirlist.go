package core

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/markbates/pkger"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
)

type TemplateOption struct {
	TemplatePath string
}

// generate navigator urls from request path
func requestPath2Nav(rpath string) []NavInfo {
	navs := make([]NavInfo, 1)
	navs[0] = NavInfo{
		Label: "/",
		Url:   "/",
	}

	if rpath == "/" {
		return navs
	}

	pathSlice := strings.Split(rpath[1:], "/")
	urlBuilder := new(strings.Builder)
	for idx, fpath := range pathSlice {
		urlBuilder.WriteString("/")
		urlBuilder.WriteString(fpath)
		label := "/" + fpath
		if idx == 0 {
			label = fpath
		}

		label, _ = url.PathUnescape(label)
		navs = append(navs, NavInfo{
			Label: label,
			Url:   urlBuilder.String(),
		})
	}
	return navs
}

// show file list form a template file
func TemplateDirList(opts ...TemplateOption) router.DirListFunc {
	var option TemplateOption
	if len(opts) > 0 {
		option = opts[0]
	}

	var f io.Reader
	if option.TemplatePath == "" {
		option.TemplatePath = "/template/dirlist.tmpl"
		f, _ = pkger.Open(option.TemplatePath)
	} else {
		f, _ = os.Open(option.TemplatePath)
	}

	content, _ := ioutil.ReadAll(f)
	templateContent := string(content)

	return func(ctx *context.Context, dirOptions iris.DirOptions, dirName string, dir http.File) error {
		dirs, err := dir.Readdir(-1)
		if err != nil {
			return err
		}

		pagData := PageData{
			Files: make([]FileInfo, 0),
			Navs:  requestPath2Nav(ctx.RequestPath(true)),
		}

		dirSlice := make([]FileInfo, 0)
		fileSlice := make([]FileInfo, 0)

		for _, d := range dirs {
			if !dirOptions.ShowHidden && router.IsHidden(d) {
				continue
			}

			name := ToBaseName(d.Name())
			upath := path.Join(ctx.Request().RequestURI, url.PathEscape(name))

			viewName := name
			if d.IsDir() {
				viewName += "/"
				dirSlice = append(dirSlice, FileInfo{d, upath, html.EscapeString(viewName)})
			} else {
				fileSlice = append(fileSlice, FileInfo{d, upath, html.EscapeString(viewName)})
			}
		}

		sort.Slice(dirSlice, func(i, j int) bool {
			return dirSlice[i].Name() < dirSlice[j].Name()
		})

		sort.Slice(fileSlice, func(i, j int) bool {
			return fileSlice[i].Name() < fileSlice[j].Name()
		})
		pagData.Files = append(pagData.Files, dirSlice...)
		pagData.Files = append(pagData.Files, fileSlice...)

		tmpl, err := template.New("default").Parse(templateContent)
		if err != nil {
			return err
		}
		return tmpl.Execute(ctx, pagData)
	}
}
