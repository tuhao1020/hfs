package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/markbates/pkger"
	"hfs/core"
	"io/ioutil"
	"log"
	"os"
)

var (
	Version string
)

const (
	CustomTemplatePath = "./template/"
	FaviconPath        = "./favicon.ico"
)

func embeddedFaviconHandler(ctx iris.Context) {
	ctx.ContentType("image/x-icon")
	r, _ := pkger.Open("/favicon.ico")
	rawBytes, _ := ioutil.ReadAll(r)
	_, _ = ctx.Write(rawBytes)
}

func main() {
	dir := flag.String("dir", "./", "root directory")
	port := flag.Int("port", 8090, "server port")
	showHidden := flag.Bool("showHidden", false, "show hidden files or directories")
	tls := flag.Bool("tls", false, "enable tls")
	certFile := flag.String("certFile", "", "cert file")
	keyFile := flag.String("keyFile", "", "key file")
	tmplName := flag.String("tmpl", "", "custom template file name")
	flag.Bool("v", false, "show version")
	flag.Parse()

	// show version
	if len(os.Args) == 2 && os.Args[1] == "-v" {
		fmt.Println(Version)
		return
	}

	if !core.IsExist(*dir) {
		log.Fatalln("dir not found.")
	}

	app := iris.New()

	// embed template and ico file
	pkger.Include("/template")
	pkger.Include(FaviconPath)

	// support for custom ico and template file
	if core.IsExist(FaviconPath) {
		app.Favicon(FaviconPath)
	} else {
		app.Get("/favicon.ico", embeddedFaviconHandler)
	}

	opts := core.TemplateOption{}
	if !core.IsExist(CustomTemplatePath) || *tmplName == "" {
		fmt.Println("Use default template")
	} else {
		opts.TemplatePath = CustomTemplatePath + (*tmplName)
		fmt.Println("Template file: ", opts.TemplatePath)
	}

	// http file server
	app.HandleDir("/", iris.Dir(*dir), iris.DirOptions{
		ShowList:   true,
		ShowHidden: *showHidden,
		DirList:    core.TemplateDirList(opts),
	})

	// http file management interfaces
	party := app.Party("/manager")
	fm := core.FileManager{
		RootDir: *dir,
	}
	party.Put("/upload", fm.Upload)
	party.Put("/md/{p:path}", fm.CreateFolder)
	party.Post("/bmv", fm.BatchMove)
	party.Delete("/rm/{p:path}", fm.Remove)
	party.Delete("/brm", fm.BatchRemove)

	var err error
	if *tls {
		err = app.Run(iris.TLS(fmt.Sprintf(":%d", *port), *certFile, *keyFile), iris.WithPostMaxMemory(32*iris.MB))
	} else {
		err = app.Listen(fmt.Sprintf(":%d", *port), iris.WithPostMaxMemory(32*iris.MB))
	}

	if err != nil {
		log.Fatalln(err.Error())
	}
}
