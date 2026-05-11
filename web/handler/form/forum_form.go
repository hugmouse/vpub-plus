package form

import (
	"errors"
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

var ErrInvalidGroupID = errors.New("invalid group_id")

func NewForumForm(r *http.Request) (*ForumForm, error) {
	position, _ := strconv.ParseInt(r.FormValue("position"), 10, 64)
	raw := r.FormValue("group_id")
	var groupID int64
	if raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || parsed < 0 {
			return nil, ErrInvalidGroupID
		}
		groupID = parsed
	}
	return &ForumForm{
		Name:                 r.FormValue("name"),
		IsLocked:             r.FormValue("locked") == "on",
		Position:             position,
		GroupID:              groupID,
		RestrictedVisibility: model.NormalizeRestrictedVisibility(r.FormValue("restricted_visibility")),
	}, nil
}
