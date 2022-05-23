package forum

import (
	"fmt"
	"net/http"
	"time"
)

// func ReqLimiter(f func(http.ResponseWriter, *http.Request)) http.Handler {
func RateLimiter(h http.Handler) http.Handler {
	// limiter := make(chan struct{}, 2)
	burstyLimiter := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		filler := time.NewTicker(5000 * time.Millisecond)
		for t := range filler.C {
			burstyLimiter <- t
		}
	}()

	// h := http.HandlerFunc(f)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// limiter <- struct{}{}
		// defer func() { <-limiter }()
		fmt.Printf("Please wait until the traffic is clear...\n")
		tL := <-burstyLimiter
		h.ServeHTTP(w, r)
		fmt.Printf("Serving req at %v\n", tL)
	})
}
