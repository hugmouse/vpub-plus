package model

// TODO: maybe not the best name? Because this is essentially RBAC? Maybe "Role"?
type Group struct {
	ID          int64
	Name        string
	MemberCount int64
}

type GroupRequest struct {
	Name string
}
