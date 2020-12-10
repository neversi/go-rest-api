package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
)

// AuthorizationUser ...
func AuthorizationUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
		userID, ok := r.Context().Value("userID").(uint64)

		role, ok := r.Context().Value("role").(string)
		if ok == false {
			misc.JSONWrite(w, misc.WriteResponse(true, userID), http.StatusUnauthorized)
		}
		vars := mux.Vars(r);
		checkID, _ := strconv.ParseUint(vars["id"], 10, 0)
		if userID != checkID && "user" == strings.ToLower(role)  {
			misc.JSONWrite(w, misc.WriteResponse(true, "Not allowed"), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r);
	})
}

