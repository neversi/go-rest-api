package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
)

// AuthorizationUser ...
func AuthorizationUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
		userID, ok := r.Context().Value(String("userID")).(string)
		if ok == false {
			misc.JSONWrite(w, misc.WriteResponse(true, "Not authorized"), http.StatusUnauthorized)
		}
		vars := mux.Vars(r);
		if userID != vars["id"] {
			misc.JSONWrite(w, misc.WriteResponse(true, "Not allowed"), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r);
	})
}

