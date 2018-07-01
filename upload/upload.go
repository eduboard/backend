package upload

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Uploader struct{}

func (u *Uploader) UploadFile(file []byte, course string, filename string) (error, string) {
	// check content type
	dir := string("./files/" + course + "/")
	path := string(dir + filename)
	serverFile := string("/files/" + course + "/" + filename)

	contentType := http.DetectContentType(file)

	if strings.Split(contentType, "/")[0] != "image" {
		return errors.New("filetype not supported"), ""
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	err := ioutil.WriteFile(path, file, 0644)
	if err != nil {
		panic(err)
		return err, ""
	}

	return nil, serverFile
}
