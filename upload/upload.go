package upload

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Uploader struct{}

func (u *Uploader) UploadFile(file []byte, course string, filename string) (error, string) {
	// check content type
	dir := string("./files/" + course + "/")
	path := string(dir + filename)
	serverFile := string("/files/" + course + "/" + filename)

	contentType := http.DetectContentType(file)

	fmt.Println(contentType)

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
