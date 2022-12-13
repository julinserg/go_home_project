package app

import (
	"bytes"
	"context"
	"errors"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julinserg/go_home_project/internal/lrucache"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

type App struct {
	logger       Logger
	cache        lrucache.Cache
	tempCacheDir string
}

type Logger interface {
	Error(msg string)
}

type InputParams struct {
	Width    int
	Height   int
	ImageURL string
}

func (in *InputParams) key() lrucache.Key {
	url := strings.ReplaceAll(in.ImageURL, ".", "_")
	url = strings.ReplaceAll(url, "/", "_")
	return lrucache.Key(url + "_" + strconv.Itoa(in.Width) + "_" + strconv.Itoa(in.Height) + ".jpg")
}

var ErrFromRemoteServer = errors.New("error from remote server")

func (a *App) getImageFromRemoteServer(
	imageURL string,
	header http.Header,
) ([]byte, int, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://"+imageURL, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header = header
	res, err := client.Do(req)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if res.StatusCode != http.StatusOK {
		return bodyBytes, res.StatusCode, ErrFromRemoteServer
	}
	return bodyBytes, http.StatusOK, nil
}

func (a *App) cropAndResizeImage(imageRaw []byte, width int, height int) ([]byte, error) {
	readerJpeg := bytes.NewReader(imageRaw)
	img, err := jpeg.Decode(readerJpeg)
	if err != nil {
		return nil, err
	}

	srcX := img.Bounds().Dx()
	srcY := img.Bounds().Dy()

	if srcX/srcY != width/height {
		dstX := srcY * width / height
		dstY := srcX * height / width
		img, err = cutter.Crop(img, cutter.Config{
			Width:  dstX,
			Height: dstY,
			Mode:   cutter.Centered,
		})
		if err != nil {
			return nil, err
		}
	}
	resizeImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, resizeImage, nil)
	return buf.Bytes(), nil
}

func (a *App) saveImageOnDisk(image []byte, pathToFile string) error {
	out, err := os.Create(pathToFile)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write(image)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) readImageFromDisk(pathToFile string) ([]byte, error) {
	dat, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (a *App) GetImagePreview(params InputParams, header http.Header) ([]byte, int, bool, error) {
	_, ok := a.cache.Get(params.key())
	if ok {
		image, err := a.readImageFromDisk(filepath.Join(a.tempCacheDir, string(params.key())))
		if err != nil {
			return nil, http.StatusInternalServerError, false, err
		}
		return image, http.StatusOK, true, nil
	}
	image, code, err := a.getImageFromRemoteServer(params.ImageURL, header)
	if err != nil {
		return image, code, false, err
	}

	imagePreview, err := a.cropAndResizeImage(image, params.Width, params.Height)
	if err != nil {
		return nil, http.StatusInternalServerError, false, err
	}
	err = a.saveImageOnDisk(imagePreview, filepath.Join(a.tempCacheDir, string(params.key())))
	if err != nil {
		return nil, http.StatusInternalServerError, false, err
	}
	a.cache.Set(params.key(), nil)
	return imagePreview, http.StatusOK, false, nil
}

func (a *App) ClearCache() {
	a.cache.Clear()
}

func New(logger Logger, cacheSize int, tempCacheDir string) *App {
	return &App{logger: logger, cache: lrucache.NewCache(cacheSize), tempCacheDir: tempCacheDir}
}
