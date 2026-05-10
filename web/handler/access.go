package handler

import "vpub/model"

// canSeeForum controls discovery (listings, search): a non-member may still see
// that a restricted forum exists when RestrictedVisibility is not "hidden".
// canAccessForum enforces actual content access and always requires membership.
func canSeeForum(forum model.Forum, user model.User) bool {
	if user.IsAdmin || forum.GroupID == 0 {
		return true
	}
	if forum.RestrictedVisibility.IsHidden() && !isMember(forum.GroupID, user) {
		return false
	}
	return true
}

func canAccessForum(forum model.Forum, user model.User) bool {
	if user.IsAdmin || forum.GroupID == 0 {
		return true
	}
	return isMember(forum.GroupID, user)
}

func isMember(groupID int64, user model.User) bool {
	for _, id := range user.GroupIDs {
		if id == groupID {
			return true
		}
	}
	return false
}
