package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris/v12"
	"hfs/core"
	"log"
)

func main() {
	dir := flag.String("dir", "./", "file directory")
	port := flag.Int("port", 8090, "server port")
	flag.Parse()

	if !core.IsExist(*dir) {
		log.Fatalln("dir not found.")
	}

	app := iris.New()
	app.Favicon("./favicon.ico")
	app.HandleDir("/", iris.Dir(*dir), iris.DirOptions{
		ShowList: true,
		DirList: core.TemplateDirList(),
	})
	_ = app.Listen(fmt.Sprintf(":%d", *port))
}
