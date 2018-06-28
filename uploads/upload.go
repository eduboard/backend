package upload

import (
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
)

func UploadFile(file []byte, course string, filename string) (error, string) {
	// check content type
	dir := string("./test/"+course+"/")
	path := string(dir+filename)
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

	return nil, path
}
