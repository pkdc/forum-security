package forum

import (
	"fmt"
	"time"
)

func RateLimit() {
	// regular interval rate limiter
	// request := make(chan int, 5)
	// for i := 1; i <= 5; i++ {
	// 	request <- i
	// }

	// burstyLimiter
	burstyLimiter := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		filler := time.NewTicker(1000 * time.Millisecond)
		for t := range filler.C {
			burstyLimiter <- t
			// fmt.Fprint(w, `<h1>Please wait...<h1>`)
		}
	}()

	// for req := range request {
	// 	tS := <-burstyLimiter
	// 	fmt.Printf("Burst...%d in %v\n", req, tS)
	// }

	fmt.Printf("Please wait...\n")
}
