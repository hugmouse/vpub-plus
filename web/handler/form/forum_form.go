package form

import (
	"net/http"
	"strconv"
	"vpub/model"
)

type ForumForm struct {
	Name     string
	IsLocked bool
	Position int64
}

//func (f *ForumForm) Validate() error {
//	if len(strings.TrimSpace(f.Name)) == 0 {
//		return errors.New("Forum name can't be empty")
//	}
//	return nil
//}

func (f *ForumForm) Merge(forum model.Forum) model.Forum {
	forum.Name = f.Name
	forum.Position = f.Position
	forum.IsLocked = f.IsLocked
	return forum
}

func NewForumForm(r *http.Request) *ForumForm {
	position, _ := strconv.ParseInt(r.FormValue("position"), 10, 64)
	return &ForumForm{
		Name:     r.FormValue("name"),
		IsLocked: r.FormValue("locked") == "on",
		Position: position,
	}
}
