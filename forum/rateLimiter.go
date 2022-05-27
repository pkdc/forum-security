package forum

import (
	"fmt"
	"net/http"
	"time"
)

func RateLimiter(f func(http.ResponseWriter, *http.Request)) http.Handler {
	h := http.HandlerFunc(f)
	// func RateLimiter(h http.Handler) http.Handler {
	// limiter := make(chan struct{}, 2)
	burstyLimiter := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		filler := time.NewTicker(200 * time.Millisecond)
		for t := range filler.C {
			burstyLimiter <- t
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// limiter <- struct{}{}
		// defer func() { <-limiter }()
		fmt.Printf("Please try again when the traffic is cleared...\n")
		select {
		case tL := <-burstyLimiter:
			h.ServeHTTP(w, r)
			fmt.Printf("Serving req at %v\n", tL)
		default:
			fmt.Fprintf(w, "<h1>Please try again when the traffic is cleared...<h1>")
		}
	})
}
