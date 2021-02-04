package core

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"os"
	"path/filepath"
	"strings"
)

type RenameObject struct {
	Src     string `form:"src"`
	NewName string `form:"newName"`
}

// Upload files via multipart/form-data
func Upload(ctx iris.Context) {
	maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	form := ctx.Request().MultipartForm
	dest := "./" + form.Value["dest"][0]
	files := form.File["files"]
	failures := make([]string, 0)
	for _, file := range files {
		_, err = ctx.SaveFormFile(file, dest+"/"+file.Filename)
		if err != nil {
			failures = append(failures, file.Filename)
		}
	}

	if len(failures) > 0 {
		errMsg := fmt.Sprintf("%s upload failed.", strings.Join(failures, ", "))
		ctx.JSON(ResponseError(iris.StatusBadRequest, errMsg))
	} else {
		ctx.JSON(ResponseOK())
	}
}

// Rename file or directory
func Rename(ctx iris.Context) {
	var renameObj RenameObject
	err := ctx.ReadForm(&renameObj)
	if err != nil {
		ctx.JSON(ResponseError(iris.StatusBadRequest, err.Error()))
		return
	}

	oldPath := "./" + renameObj.Src
	err = os.Rename(oldPath, filepath.Dir(oldPath)+"/"+renameObj.NewName)
	if err != nil {
		ctx.JSON(ResponseError(iris.StatusBadRequest, err.Error()))
	} else {
		ctx.JSON(ResponseOK())
	}
}

// Remove file or directory(even if directory is not empty)
func Remove(ctx iris.Context) {
	path := "./" + ctx.Params().Get("p")
	err := os.RemoveAll(path)
	if err != nil {
		ctx.JSON(ResponseError(iris.StatusBadRequest, err.Error()))
	} else {
		ctx.JSON(ResponseOK())
	}
}
