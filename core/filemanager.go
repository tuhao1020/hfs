package core

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"os"
	"path/filepath"
	"strings"
)

type RenameObject struct {
	Src     string `json:"src"`
	NewName string `json:"newName"`
}

type BatchRemoveObject struct {
	Files []string `json:"files"`
}

type MoveObject struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
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
	destParam := form.Value["dest"]
	dest := "./"
	if destParam != nil {
		dest += destParam[0]
	}

	files := form.File["files"]
	if len(files) == 0 {
		ResponseError(ctx, iris.StatusBadRequest, "files not found")
		return
	}

	failures := make([]string, 0)
	for _, file := range files {
		_, err = ctx.SaveFormFile(file, dest+"/"+file.Filename)
		if err != nil {
			failures = append(failures, file.Filename)
		}
	}

	if len(failures) > 0 {
		errMsg := fmt.Sprintf("%s upload failed.", strings.Join(failures, ", "))
		ResponseError(ctx, iris.StatusBadRequest, errMsg)
	} else {
		ResponseOK(ctx)
	}
}

// Rename a file or directory
func Rename(ctx iris.Context) {
	var renameObj RenameObject
	err := ctx.ReadJSON(&renameObj)
	if err != nil {
		ResponseError(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	oldPath := "./" + renameObj.Src
	err = os.Rename(oldPath, filepath.Dir(oldPath)+"/"+renameObj.NewName)
	if err != nil {
		ResponseError(ctx, iris.StatusBadRequest, err.Error())
	} else {
		ResponseOK(ctx)
	}
}

// Remove file or directory(even if directory is not empty)
func Remove(ctx iris.Context) {
	path := "./" + ctx.Params().Get("p")
	err := os.RemoveAll(path)
	if err != nil {
		ResponseError(ctx, iris.StatusBadRequest, err.Error())
	} else {
		ResponseOK(ctx)
	}
}

// BatchRemove remove multiple files or directories once
func BatchRemove(ctx iris.Context) {
	var batchRemoveObj BatchRemoveObject
	err := ctx.ReadJSON(&batchRemoveObj)
	if err != nil {
		ResponseError(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	failures := make([]string, 0)
	for _, path := range batchRemoveObj.Files {
		err = os.RemoveAll(path)
		if err != nil {
			failures = append(failures, path)
		}
	}

	if len(failures) > 0 {
		errMsg := fmt.Sprintf("%s delete failed.", strings.Join(failures, ", "))
		ResponseError(ctx, iris.StatusBadRequest, errMsg)
	} else {
		ResponseOK(ctx)
	}
}

// Move a file to another exist location
func Move(ctx iris.Context) {
	var moveObj MoveObject
	err := ctx.ReadJSON(&moveObj)
	if err != nil {
		ResponseError(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	err = os.Rename("./"+moveObj.Src, "./"+moveObj.Dest)
	if err != nil {
		ResponseError(ctx, iris.StatusBadRequest, err.Error())
	} else {
		ResponseOK(ctx)
	}
}

// CreateFolder for the given path
func CreateFolder(ctx iris.Context) {
	path := "./" + ctx.Params().Get("p")
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		ResponseError(ctx, iris.StatusBadRequest, err.Error())
	} else {
		ResponseOK(ctx)
	}
}
