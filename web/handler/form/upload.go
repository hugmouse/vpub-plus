package form

import (
	"io"
	"net/http"
)

type UploadForm struct {
	Filename string
	File     io.Reader
}

func NewUploadForm(r *http.Request) (*UploadForm, error) {
	if err := r.ParseMultipartForm(2 << 20); err != nil { // TODO: Limit should be somewhere else?
		return nil, err
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return &UploadForm{
		Filename: handler.Filename,
		File:     file,
	}, err
}
