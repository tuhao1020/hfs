package core

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
)

type TemplateOption struct {
	TemplatePath string
}

// generate navigator urls from request path
func requestPath2Nav(path string) []NavInfo {
	navs := make([]NavInfo, 1)
	navs[0] = NavInfo{
		Label: "/",
		Url:   "/",
	}

	if path == "/" {
		return navs
	}

	pathSlice := strings.Split(path[1:], "/")
	url := new(strings.Builder)
	for idx, path := range pathSlice {
		url.WriteString("/")
		url.WriteString(path)
		label := "/" + path
		if idx ==0 {
			label = path
		}
		navs = append(navs, NavInfo{
			Label: label,
			Url: url.String(),
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

	if option.TemplatePath == "" {
		option.TemplatePath = "./template/dirlist.tmpl"
	}

	return func(ctx *context.Context, dirOptions iris.DirOptions, dirName string, dir http.File) error {
		dirs, err := dir.Readdir(-1)
		if err != nil {
			return err
		}

		pagData := PageData{
			Files: make([]FileInfo, 0),
			Navs: requestPath2Nav(ctx.RequestPath(true)),
		}

		dirSlice := make([]FileInfo, 0)
		fileSlice := make([]FileInfo, 0)

		for _, d := range dirs {
			if !dirOptions.ShowHidden && router.IsHidden(d) {
				continue
			}

			name := ToBaseName(d.Name())

			upath := path.Join(ctx.Request().RequestURI, name)
			url := url.URL{Path: upath}

			viewName := name
			if d.IsDir() {
				viewName += "/"
				dirSlice = append(dirSlice, FileInfo{d, url.String(), html.EscapeString(viewName)})
			} else {
				fileSlice = append(fileSlice, FileInfo{d, url.String(), html.EscapeString(viewName)})
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

		tmpl, _ := template.ParseFiles(option.TemplatePath)
		return tmpl.Execute(ctx, pagData)
	}
}