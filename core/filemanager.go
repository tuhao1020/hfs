package core

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"os"
	"strings"
)

type BatchRemoveParam struct {
	Files []string `json:"files"`
}

type MoveParam struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

type BatchMoveParam struct {
	Params []MoveParam `json:"params"`
}

type FileManager struct {
	RootDir string
}

// Upload files via multipart/form-data
func (fm FileManager) Upload(ctx iris.Context) {
	maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	form := ctx.Request().MultipartForm
	destParam := form.Value["dest"]
	// if dest is not set, the default path is
	dest := fm.RootDir
	if destParam != nil {
		dest = fm.web2LocalPath(destParam[0])
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

// Remove file or directory(even if directory is not empty)
func (fm FileManager) Remove(ctx iris.Context) {
	path := fm.web2LocalPath(ctx.Params().Get("p"))
	err := os.RemoveAll(path)
	if err != nil {
		ResponseError(ctx, iris.StatusBadRequest, err.Error())
	} else {
		ResponseOK(ctx)
	}
}

// BatchRemove remove multiple files or directories once
func (fm FileManager) BatchRemove(ctx iris.Context) {
	var batchRemoveParam BatchRemoveParam
	err := ctx.ReadJSON(&batchRemoveParam)
	if err != nil {
		ResponseError(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	failures := make([]string, 0)
	for _, path := range batchRemoveParam.Files {
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

// BatchMove or rename multiple files or directories to another exist location
func (fm FileManager) BatchMove(ctx iris.Context) {
	var batchMoveParam BatchMoveParam
	err := ctx.ReadJSON(&batchMoveParam)
	if err != nil {
		ResponseError(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	failures := make([]string, 0)
	for _, param := range batchMoveParam.Params {
		err = fm.moveTo(param.Src, param.Dest)
		if err != nil {
			failures = append(failures, param.Src)
		}
	}

	if len(failures) > 0 {
		errMsg := fmt.Sprintf("%s move failed.", strings.Join(failures, ", "))
		ResponseError(ctx, iris.StatusBadRequest, errMsg)
	} else {
		ResponseOK(ctx)
	}
}

// CreateFolder for the given path
func (fm FileManager) CreateFolder(ctx iris.Context) {
	path := fm.web2LocalPath(ctx.Params().Get("p"))
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		ResponseError(ctx, iris.StatusBadRequest, err.Error())
	} else {
		ResponseOK(ctx)
	}
}

func (fm FileManager) web2LocalPath(webPath string) string {
	return fm.RootDir + "/" + webPath
}

func (fm FileManager) moveTo(src, dest string) error {
	srcLocal := fm.web2LocalPath(src)
	destLocal := fm.web2LocalPath(dest)
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
