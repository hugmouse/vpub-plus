package handler

//
//import (
//	"github.com/gorilla/csrf"
//	"net/http"
//)
//
//func (h *Handler) showNotificationsView(w http.ResponseWriter, r *http.Request, user string) {
//	notifications, err := h.storage.NotificationsByUser(user)
//	if err != nil {
//		serverError(w, err)
//		return
//	}
//	h.renderLayout(w, "notifications", map[string]interface{}{
//		"notifications":  notifications,
//		csrf.TemplateTag: csrf.TemplateField(r),
//	}, user)
//}
//
//func (h *Handler) markRead(w http.ResponseWriter, r *http.Request, user string) {
//	id := RouteInt64Param(r, "notificationId")
//	if err := h.storage.DeleteNotification(id); err != nil {
//		serverError(w, err)
//		return
//	}
//	http.Redirect(w, r, "/notifications", http.StatusFound)
//}
//
//func (h *Handler) markAllRead(w http.ResponseWriter, r *http.Request, user string) {
//	if err := h.storage.DeleteNotificationsFromUser(user); err != nil {
//		serverError(w, err)
//		return
//	}
//	http.Redirect(w, r, "/notifications", http.StatusFound)
//}
