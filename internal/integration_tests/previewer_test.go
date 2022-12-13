//go:build integration

package previewerintegrationtests

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

type previewerTest struct {
	responseStatusCode int
	responseBody       []byte
	header             http.Header
}

func (test *previewerTest) iSendRequestTo(httpMethod, addr string) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodGet:
		r, err = http.Get(addr) //nolint:gosec,noctx
		defer r.Body.Close()
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}
	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	test.header = r.Header

	return
}

func (test *previewerTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *previewerTest) theResponseShouldMatchText(text string) error {
	if string(test.responseBody) != text {
		return fmt.Errorf("unexpected text: %s != %s", test.responseBody, text)
	}
	return nil
}

func (test *previewerTest) theResponseShouldMatchTextMultiLine(text string) error {
	if !strings.Contains(string(test.responseBody), text) {
		return fmt.Errorf("unexpected text: %s != %s", test.responseBody, text)
	}
	return nil
}

func (test *previewerTest) compareWithImage(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	imageFromFile, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	res := bytes.Compare(imageFromFile, test.responseBody)
	if res != 0 {
		return fmt.Errorf("response body images bounds and file images bounds %s not equivalent", filePath)
	}
	return nil
}

func (test *previewerTest) imageGetFromCache() error {
	if test.header.Get("is-image-from-cache") != "true" {
		return fmt.Errorf("is-image-from-cache: %s != true", test.header.Get("is-image-from-cache"))
	}
	return nil
}

func (test *previewerTest) imageGetFromRemoteServer() error {
	if test.header.Get("is-image-from-cache") != "false" {
		return fmt.Errorf("is-image-from-cache: %s != false", test.header.Get("is-image-from-cache"))
	}
	return nil
}

func InitializeScenario(s *godog.ScenarioContext) {
	test := new(previewerTest)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match text "([^"]*)"$`, test.theResponseShouldMatchText)
	s.Step(`^The response should contains text$`, test.theResponseShouldMatchTextMultiLine)
	s.Step(`^The response equivalent image "([^"]*)"$`, test.compareWithImage)
	s.Step(`^Image get from cache$`, test.imageGetFromCache)
	s.Step(`^Image get from remote server$`, test.imageGetFromRemoteServer)
}
