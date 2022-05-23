package forum

import (
	"fmt"
	"net/http"
)

// func ReqLimiter(f func(http.ResponseWriter, *http.Request)) http.Handler {
func ReqLimiter(h http.Handler) http.Handler {
	limiter := make(chan struct{}, 2)

	// h := http.HandlerFunc(f)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter <- struct{}{}
		// defer func() { <-limiter }()

		h.ServeHTTP(w, r)
		fmt.Printf("Serving req\n")
	})
}
