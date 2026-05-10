package handler

import (
	"testing"
	"vpub/model"
)

func TestCanSeeForum_PublicForumIsAlwaysVisible(t *testing.T) {
	forum := model.Forum{GroupID: 0}
	if !canSeeForum(forum, model.User{}) {
		t.Error("public forum must be visible to anonymous user")
	}
}

func TestCanSeeForum_AdminSeesHiddenForum(t *testing.T) {
	forum := model.Forum{GroupID: 1, RestrictedVisibility: model.RestrictedVisibilityHidden}
	if !canSeeForum(forum, model.User{IsAdmin: true}) {
		t.Error("admin must see all forums")
	}
}

func TestCanSeeForum_HiddenForumInvisibleToNonMember(t *testing.T) {
	forum := model.Forum{GroupID: 1, RestrictedVisibility: model.RestrictedVisibilityHidden}
	user := model.User{GroupIDs: []int64{2}}
	if canSeeForum(forum, user) {
		t.Error("non-member must not see hidden forum in listing")
	}
}

func TestCanSeeForum_HiddenForumInvisibleToAnonymous(t *testing.T) {
	forum := model.Forum{GroupID: 1, RestrictedVisibility: model.RestrictedVisibilityHidden}
	if canSeeForum(forum, model.User{}) {
		t.Error("anonymous user must not see hidden forum in listing")
	}
}

func TestCanSeeForum_VisibleForumListedForNonMember(t *testing.T) {
	forum := model.Forum{GroupID: 1, RestrictedVisibility: model.RestrictedVisibilityVisible}
	user := model.User{GroupIDs: []int64{2}}
	if !canSeeForum(forum, user) {
		t.Error("non-member must see visible-mode restricted forum in listing")
	}
}

func TestCanSeeForum_MemberSeesHiddenForum(t *testing.T) {
	forum := model.Forum{GroupID: 1, RestrictedVisibility: model.RestrictedVisibilityHidden}
	user := model.User{GroupIDs: []int64{1}}
	if !canSeeForum(forum, user) {
		t.Error("group member must see their restricted forum")
	}
}

func TestCanAccessForum_PublicForumAccessible(t *testing.T) {
	forum := model.Forum{GroupID: 0}
	if !canAccessForum(forum, model.User{}) {
		t.Error("public forum must be accessible to everyone")
	}
}

func TestCanAccessForum_AdminAccessesRestrictedForum(t *testing.T) {
	forum := model.Forum{GroupID: 1}
	if !canAccessForum(forum, model.User{IsAdmin: true}) {
		t.Error("admin must access all forums")
	}
}

func TestCanAccessForum_MemberAccessesRestrictedForum(t *testing.T) {
	forum := model.Forum{GroupID: 1}
	user := model.User{GroupIDs: []int64{1, 2}}
	if !canAccessForum(forum, user) {
		t.Error("group member must access restricted forum")
	}
}

func TestCanAccessForum_NonMemberBlocked(t *testing.T) {
	forum := model.Forum{GroupID: 1}
	user := model.User{GroupIDs: []int64{2}}
	if canAccessForum(forum, user) {
		t.Error("non-member must not access restricted forum")
	}
}

func TestCanAccessForum_AnonymousUserBlocked(t *testing.T) {
	forum := model.Forum{GroupID: 1}
	if canAccessForum(forum, model.User{}) {
		t.Error("anonymous user (nil GroupIDs) must not access restricted forum")
	}
}
