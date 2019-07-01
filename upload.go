// Package fileupload provides an easy way to upload files to a filehost.
package fileupload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// UploadToHost takes a url and a file as arguments and uploads the file to the provided url with HTTP POST.
// It returns the url to the uploaded file as a string and any error encountered.
func UploadToHost(file *os.File, name string) (string, error) {
	var err error

	var client http.Client
	var b bytes.Buffer

	writer := multipart.NewWriter(&b)

	var fw io.Writer

	defer file.Close()

	if fw, err = writer.CreateFormFile("file:1", file.Name()); err != nil {
		return "", err
	}

	if fw, err = writer.CreateFormFile("name:1", name); err != nil {
		return "", err
	}

	if _, err = io.Copy(fw, file); err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://ix.io", &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Submit the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return strings.Replace(bodyString, "\n", "", -1), nil
	}
	return "", fmt.Errorf("bad status: %s", resp.Status)
}
