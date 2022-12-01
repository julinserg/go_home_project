package internalhttp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrNotSetWidthOrHeights = errors.New("not set width or height in URL")
	ErrWidthNotInt          = errors.New("width in URL not integer")
	ErrHeightNotInt         = errors.New("height in URL not integer")
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
	splitURL := strings.Split(r.URL.Path, "/")
	if len(splitURL) < 3 {
		ph.sendError(ErrNotSetWidthOrHeights, http.StatusBadRequest, w)
		return
	}
	splitURLParam := splitURL[2:]

	if len(splitURLParam) < 3 {
		ph.sendError(ErrNotSetWidthOrHeights, http.StatusBadRequest, w)
		return
	}
	width, err := strconv.Atoi(splitURLParam[0])
	if err != nil {
		ph.sendError(ErrWidthNotInt, http.StatusBadRequest, w)
		return
	}

	height, err := strconv.Atoi(splitURLParam[1])
	if err != nil {
		ph.sendError(ErrHeightNotInt, http.StatusBadRequest, w)
		return
	}

	imageURL := strings.Join(splitURLParam[2:], "/")

	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://"+imageURL, nil)
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
	defer res.Body.Close()

	_ = width
	_ = height

	io.Copy(w, res.Body)
	w.WriteHeader(res.StatusCode)
}
