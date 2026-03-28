package proxy

import "net/http"

func Middleware(next http.Handler) http.Handler {

	// here I can read the config file and do "proxy pass"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-From-Proxy", "true")
		w.Write([]byte("Hello, world!"))
	})
}
