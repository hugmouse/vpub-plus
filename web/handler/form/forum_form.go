package form

import (
	"net/http"
	"strconv"
	"vpub/model"
)

type ForumForm struct {
	Name                 string
	IsLocked             bool
	Position             int64
	GroupID              int64
	RestrictedVisibility model.RestrictedVisibility
	Groups               []model.Group
}

func NewForumForm(r *http.Request) *ForumForm {
	position, _ := strconv.ParseInt(r.FormValue("position"), 10, 64)
	groupID, _ := strconv.ParseInt(r.FormValue("group_id"), 10, 64)
	return &ForumForm{
		Name:                 r.FormValue("name"),
		IsLocked:             r.FormValue("locked") == "on",
		Position:             position,
		GroupID:              groupID,
		RestrictedVisibility: model.NormalizeRestrictedVisibility(r.FormValue("restricted_visibility")),
	}
}
