package mock

type Uploader struct {
	UploadFileFn        func(file []byte, course string, filename string) (error, string)
	UploadFileFnInvoked bool
}

func (u *Uploader) UploadFile(file []byte, course string, filename string) (error, string) {
	u.UploadFileFnInvoked = true
	return u.UploadFileFn(file, course, filename)
}
