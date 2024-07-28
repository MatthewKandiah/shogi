package util

import "net/http"

func ValueFromCookie(key string, r *http.Request) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}
