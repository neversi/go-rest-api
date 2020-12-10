package middleware

import (
	"context"
	"net/http"

	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/auth"
	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/misc"
)

// String ...
type String string;

// IsAuthenticated ...
func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		if r.Header["Authorization"] != nil {
			if err := auth.TokenValid(r); err != nil {
				misc.JSONWrite(w, misc.WriteResponse(true, "Token is not valid"), http.StatusUnauthorized)
			} else {
				accessToken, err := auth.ExtractTokenData(r)
				
				if err != nil {
					misc.JSONWrite(w, misc.WriteResponse(true, "Problem with Token"), http.StatusUnauthorized)
					return
				}
				
				r.WithContext(context.WithValue(r.Context(), String("userID"), accessToken.Userid))
				r.WithContext(context.WithValue(r.Context(), String("role"), accessToken.URole))
				next.ServeHTTP(w, r)
				return
			}
			
		} else {
			misc.JSONWrite(w, misc.WriteResponse(true, "Not Authenticated"), http.StatusUnauthorized)
		}
	})
}

