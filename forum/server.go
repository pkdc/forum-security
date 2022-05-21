package forum

import (
	"net/http"
	"time"
)

func MakeServer() http.Server {
	return http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
