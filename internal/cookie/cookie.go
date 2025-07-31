package cookie

import "net/http"

const SessionKey = "bs"

func Set(rw http.ResponseWriter, r *http.Request, key, value string) {
	http.SetCookie(rw, &http.Cookie{
		Name:     key,
		Value:    value,
		HttpOnly: true,
		Domain:   r.URL.Host,
		Path:     "/",
		MaxAge:   86400 * 30,
	})
}

func Delete(rw http.ResponseWriter, r *http.Request, key string) {
	http.SetCookie(rw, &http.Cookie{
		Name:     key,
		Value:    "",
		HttpOnly: true,
		Domain:   r.URL.Host,
		Path:     "/",
		MaxAge:   0,
	})
}
