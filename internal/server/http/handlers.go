package internalhttp

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrNotSetWidthOrHeights = errors.New("Not set width or height in URL")
	ErrWidthNotInt          = errors.New("Width in URL not integer")
	ErrHeightNotInt         = errors.New("Height in URL not integer")
)

type previewerHandler struct {
	logger Logger
	app    Application
}

func (ph *previewerHandler) hellowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is my previewer!"))
}

func (ph *previewerHandler) sendError(err error, code int, w http.ResponseWriter) {
	ph.logger.Error(err.Error())
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func (ph *previewerHandler) mainHandler(w http.ResponseWriter, r *http.Request) {

	splitUrl := strings.Split(r.URL.Path, "/")
	if len(splitUrl) < 3 {
		ph.sendError(ErrNotSetWidthOrHeights, http.StatusBadRequest, w)
		return
	}
	splitUrlParam := splitUrl[2:len(splitUrl)]

	if len(splitUrlParam) < 3 {
		ph.sendError(ErrNotSetWidthOrHeights, http.StatusBadRequest, w)
		return
	}
	width, err := strconv.Atoi(splitUrlParam[0])
	if err != nil {
		ph.sendError(ErrWidthNotInt, http.StatusBadRequest, w)
		return
	}

	height, err := strconv.Atoi(splitUrlParam[1])
	if err != nil {
		ph.sendError(ErrHeightNotInt, http.StatusBadRequest, w)
		return
	}

	imageUrl := strings.Join(splitUrlParam[2:len(splitUrlParam)], "/")

	client := http.Client{}
	req, err := http.NewRequest("GET", "http://"+imageUrl, nil)
	if err != nil {
		ph.sendError(ErrHeightNotInt, http.StatusInternalServerError, w)
		return
	}

	req.Header = r.Header
	res, err := client.Do(req)
	if err != nil {
		ph.sendError(ErrHeightNotInt, http.StatusBadRequest, w)
		return
	}

	_ = width
	_ = height

	io.Copy(w, res.Body)
	w.WriteHeader(res.StatusCode)
}
