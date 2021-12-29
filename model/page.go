package model

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Page struct {
	Id          int64
	User        string
	Filename    string
	Title       string
	Description string
	Content     []byte
	Comments    int
	LastUpdated time.Time
}

var extensions = []string{"", ".html", ".htm", ".png", ".jpeg", ".jpg", ".css", ".js", ".gmi"}
var binary = []string{".png", ".jpeg", ".jpg"}

func (p Page) Validate() error {
	if len(p.Filename) == 0 {
		return errors.New("filename is empty")
	}
	if !p.IsBinary() {
		if len(p.Content) == 0 {
			return errors.New("content is empty")
		}
	}
	for _, ext := range extensions {
		if filepath.Ext(p.Filename) == ext {
			return nil
		}
	}
	return errors.New("extension not supported")
}

func (p Page) SaveToFile() error {
	return ioutil.WriteFile(p.path(), p.Content, 0644)
}

func (p Page) UpdateFile(prevFilename string) error {
	prev := Page{User: p.User, Filename: prevFilename}
	err := os.Rename(prev.path(), p.path())
	if err != nil {
		return err
	}
	if p.IsBinary() {
		return nil
	}
	return p.SaveToFile()
}

func (p Page) IsBinary() bool {
	for _, ext := range binary {
		if filepath.Ext(p.Filename) == ext {
			return true
		}
	}
	return false
}

func (p *Page) ReadFile() error {
	b, err := ioutil.ReadFile(p.path())
	if err != nil {
		return err
	}
	p.Content = b
	return nil
}

func (p Page) DeleteFile() error {
	return os.Remove(p.path())
}

func (p Page) path() string {
	return path.Join(User{Name: p.User}.Folder(), p.Filename)
}

func (p Page) Date() string {
	return p.LastUpdated.Format("2006-01-02")
}
