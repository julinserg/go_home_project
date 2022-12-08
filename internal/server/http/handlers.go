package internalhttp

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
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

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		ph.sendError(err, http.StatusInternalServerError, w)
		return
	}

	f, err := os.Create("./pic.jpg")
	if err != nil {
		ph.sendError(err, http.StatusInternalServerError, w)
		return
	}
	wf := bufio.NewWriter(f)
	wf.Write(bodyBytes)

	// decode jpeg into image.Image
	readerJpeg := bytes.NewReader(bodyBytes)
	img, err := jpeg.Decode(readerJpeg)
	if err != nil {
		log.Fatal(err)
	}

	srcX := img.Bounds().Dx()
	srcY := img.Bounds().Dy()
	srcP := srcX / srcY
	dstP := width / height
	if srcP != dstP {
		dstX := srcY * width / height
		dstY := srcX * height / width
		img, err = cutter.Crop(img, cutter.Config{
			Width:  dstX,
			Height: dstY,
			Mode:   cutter.Centered,
		})
	}
	resizeImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	out, err := os.Create("./test_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, resizeImage, nil)
	out.Write(buf.Bytes())
	w.Write(buf.Bytes())
	w.WriteHeader(res.StatusCode)
}
