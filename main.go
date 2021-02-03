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
	flag.Parse()

	if !core.IsExist(*dir) {
		log.Fatalln("dir not found.")
	}

	app := iris.New()
	app.Favicon("./favicon.ico")
	app.HandleDir("/", iris.Dir(*dir), iris.DirOptions{
		ShowList:   true,
		ShowHidden: *showHidden,
		DirList:    core.TemplateDirList(),
	})
	_ = app.Listen(fmt.Sprintf(":%d", *port))
}
