package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris/v12"
	"hfs/core"
	"log"
)

func main() {
	dir := flag.String("dir", "./", "root directory")
	port := flag.Int("port", 8090, "server port")
	showHidden := flag.Bool("showHidden", false, "show hidden files or directories")
	tls := flag.Bool("tls", false, "enable tls")
	certFile := flag.String("certFile", "", "cert file")
	keyFile := flag.String("keyFile", "", "key file")
	flag.Parse()

	if !core.IsExist(*dir) {
		log.Fatalln("dir not found.")
	}

	app := iris.New()
	app.Favicon("./favicon.ico")

	// http file server
	app.HandleDir("/", iris.Dir(*dir), iris.DirOptions{
		ShowList:   true,
		ShowHidden: *showHidden,
		DirList:    core.TemplateDirList(),
	})

	// http file management interfaces
	party := app.Party("/manager")
	party.Put("/upload", core.Upload)
	party.Put("/md/{p:path}", core.CreateFolder)
	party.Post("/bmv", core.BatchMove)
	party.Delete("/rm/{p:path}", core.Remove)
	party.Delete("/brm", core.BatchRemove)

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
