package main

import (
	"fmt"
	"forum/forum"
	"log"
	"net/http"
	"os"
)

func main() {
	forum.InitDB()
	// forum.ClearUsers()
	// forum.ClearPosts()
	// forum.ClearComments()
	// exec.Command("xdg-open", "http://localhost:8080/").Start()

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	mux.Handle("/", forum.RateLimiter(forum.HomeHandler))
	mux.Handle("/login", forum.RateLimiter(forum.LoginHandler))
	mux.Handle("/register", forum.RateLimiter(forum.RegisterHandler))
	mux.Handle("/logout", forum.RateLimiter(forum.LogoutHandler))
	mux.Handle("/postpage", forum.RateLimiter(forum.PostPageHandler))
	// http.HandleFunc("/delete", forum.DeleteHandler)
	fmt.Println("Starting server")

	// err := http.ListenAndServe(":8080", mux)
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
