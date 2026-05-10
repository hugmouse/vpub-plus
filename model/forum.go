package model

type RestrictedVisibility string

const (
	RestrictedVisibilityHidden  RestrictedVisibility = "hidden"
	RestrictedVisibilityVisible RestrictedVisibility = "visible"
)

func NormalizeRestrictedVisibility(v string) RestrictedVisibility {
	switch RestrictedVisibility(v) {
	case RestrictedVisibilityVisible:
		return RestrictedVisibilityVisible
	default:
		return RestrictedVisibilityHidden
	}
}

func (v RestrictedVisibility) String() string  { return string(v) }
func (v RestrictedVisibility) IsHidden() bool  { return v == RestrictedVisibilityHidden }
func (v RestrictedVisibility) IsVisible() bool { return v == RestrictedVisibilityVisible }

type Forum struct {
	ID                   int64
	Name                 string
	Boards               []Board
	Position             int64
	IsLocked             bool
	GroupID              int64  // 0 = no restriction
	RestrictedVisibility RestrictedVisibility
}

type ForumRequest struct {
	Name                 string
	Position             int64
	IsLocked             bool
	GroupID              int64
	RestrictedVisibility RestrictedVisibility
}
