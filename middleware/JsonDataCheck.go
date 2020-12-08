package middleware

import (
	"net/http"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
)

// JSONDataCheck ...
func JSONDataCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		if r.Header.Get("Content-Type") != "application/json" {
			misc.JSONWrite(w, misc.WriteResponse(true, "Unsupported Media Type: need application/json"), http.StatusUnsupportedMediaType)
			return
		} 
		next.ServeHTTP(w, r);
	})
}