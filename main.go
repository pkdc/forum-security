package main

import (
	"fmt"
	"forum/forum"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	forum.InitDB()
	// forum.ClearUsers()
	// forum.ClearPosts()
	// forum.ClearComments()
	exec.Command("xdg-open", "http://localhost:8080/").Start()

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", forum.HomeHandler)
	http.HandleFunc("/login", forum.LoginHandler)
	http.HandleFunc("/register", forum.RegisterHandler)
	http.HandleFunc("/logout", forum.LogoutHandler)
	http.HandleFunc("/postpage", forum.PostPageHandler)
	// http.HandleFunc("/delete", forum.DeleteHandler)
	fmt.Println("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
