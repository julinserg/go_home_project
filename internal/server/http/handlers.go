package internalhttp

import (
	"bytes"
	"context"
	"errors"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

var (
	ErrNotSetParams          = errors.New("not set width, height or image path in URL")
	ErrWidthNotInt           = errors.New("width in URL not integer")
	ErrHeightNotInt          = errors.New("height in URL not integer")
	ErrDimmensionIsVeryLarge = errors.New("width or height is very large")
)

type previewerHandler struct {
	logger Logger
	app    Application
}

type inputParams struct {
	width    int
	height   int
	imageURL string
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

func (ph *previewerHandler) proxyError(body []byte, code int, w http.ResponseWriter) {
	ph.logger.Error("Image Server Error Code: " + strconv.Itoa(code))
	w.WriteHeader(code)
	w.Write(body)
}

func (ph *previewerHandler) validateInputParameter(url string, w http.ResponseWriter) (bool, inputParams) {
	splitURL := strings.Split(url, "/")
	if len(splitURL) < 3 {
		ph.sendError(ErrNotSetParams, http.StatusBadRequest, w)
		return false, inputParams{}
	}
	splitURLParam := splitURL[2:]

	if len(splitURLParam) < 3 {
		ph.sendError(ErrNotSetParams, http.StatusBadRequest, w)
		return false, inputParams{}
	}
	width, err := strconv.Atoi(splitURLParam[0])
	if err != nil {
		ph.sendError(ErrWidthNotInt, http.StatusBadRequest, w)
		return false, inputParams{}
	}

	height, err := strconv.Atoi(splitURLParam[1])
	if err != nil {
		ph.sendError(ErrHeightNotInt, http.StatusBadRequest, w)
		return false, inputParams{}
	}

	if width > 3840 || height > 2160 {
		ph.sendError(ErrDimmensionIsVeryLarge, http.StatusBadRequest, w)
		return false, inputParams{}
	}

	imageURL := strings.Join(splitURLParam[2:], "/")

	return true, inputParams{width, height, imageURL}
}

func (ph *previewerHandler) getImageFromRemoteServer(imageURL string, header http.Header, w http.ResponseWriter) (bool, []byte) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://"+imageURL, nil)
	if err != nil {
		ph.sendError(err, http.StatusInternalServerError, w)
		return false, nil
	}

	req.Header = header
	res, err := client.Do(req)
	if err != nil {
		ph.sendError(err, http.StatusBadRequest, w)
		return false, nil
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		ph.sendError(err, http.StatusInternalServerError, w)
		return false, nil
	}

	if res.StatusCode != http.StatusOK {
		ph.proxyError(bodyBytes, res.StatusCode, w)
		return false, nil
	}
	return true, bodyBytes
}

func (ph *previewerHandler) cropAndResizeImage(imageRaw []byte, width int, height int, w http.ResponseWriter) (bool, []byte) {
	readerJpeg := bytes.NewReader(imageRaw)
	img, err := jpeg.Decode(readerJpeg)
	if err != nil {
		ph.sendError(err, http.StatusInternalServerError, w)
		return false, nil
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

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, resizeImage, nil)
	return true, buf.Bytes()
}

func (ph *previewerHandler) saveImage(image []byte, w http.ResponseWriter) bool {
	out, err := os.Create("./resized_image.jpg")
	if err != nil {
		ph.sendError(err, http.StatusInternalServerError, w)
		return false
	}
	defer out.Close()
	out.Write(image)
	return true
}

func (ph *previewerHandler) mainHandler(w http.ResponseWriter, r *http.Request) {
	isValid, params := ph.validateInputParameter(r.URL.Path, w)
	if !isValid {
		return
	}

	isSuccess, imagesRaw := ph.getImageFromRemoteServer(params.imageURL, r.Header, w)
	if !isSuccess {
		return
	}

	isSuccess, imagePreview := ph.cropAndResizeImage(imagesRaw, params.width, params.height, w)
	if !isSuccess {
		return
	}

	isSuccess = ph.saveImage(imagePreview, w)
	if !isSuccess {
		return
	}

	w.Write(imagePreview)
	w.WriteHeader(http.StatusOK)
}
