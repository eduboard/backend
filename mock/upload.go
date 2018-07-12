package mock

type Uploader struct {
	UploadFileFn        func(file []byte, course string, filename string) (string, error)
	UploadFileFnInvoked bool
}

func (u *Uploader) UploadFile(file []byte, course string, filename string) (string, error) {
	u.UploadFileFnInvoked = true
	return u.UploadFileFn(file, course, filename)
}
