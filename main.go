package main

import (
	"fmt"
	"forum/forum"
	"log"
	"net"
	"net/http"
	"os/exec"
)

func hanleConn(conn net.Conn) {

}

func main() {
	forum.InitDB()
	// forum.ClearUsers()
	// forum.ClearPosts()
	// forum.ClearComments()
	exec.Command("xdg-open", "http://localhost:8080/").Start()

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/", forum.HomeHandler)
	mux.HandleFunc("/login", forum.LoginHandler)
	mux.HandleFunc("/register", forum.RegisterHandler)
	mux.HandleFunc("/logout", forum.LogoutHandler)
	mux.HandleFunc("/postpage", forum.PostPageHandler)
	// http.HandleFunc("/delete", forum.DeleteHandler)
	fmt.Println("Starting server at port 8080")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	// limit rate here?
	http.Serve(ln, mux)

	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
