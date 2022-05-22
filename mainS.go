package main

// import (
// 	"crypto/tls"
// 	"fmt"
// 	"forum/forum"
// 	"log"
// 	"net/http"

// 	"golang.org/x/crypto/acme/autocert"
// )

// func main() {
// 	forum.InitDB()
// 	forum.ClearUsers()
// 	forum.ClearPosts()
// 	forum.ClearComments()
// 	// exec.Command("xdg-open", "http://localhost:8080/").Start()

// 	dir := "./certs"
// 	certMan := &autocert.Manager{
// 		Prompt: autocert.AcceptTOS,
// 		// HostPolicy: autocert.HostWhitelist("www.domain.com"),
// 		HostPolicy: nil,
// 		Cache:      autocert.DirCache(dir),
// 	}
// 	go func() {
// 		httpServer := forum.MakeServer()
// 		httpServer.Addr = ":80"
// 		// httpServer.Addr = ":8080"
// 		httpServer.Handler = certMan.HTTPHandler(nil)
// 		err := httpServer.ListenAndServe()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	mux := http.NewServeMux()
// 	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
// 	mux.HandleFunc("/", forum.HomeHandler)
// 	mux.HandleFunc("/login", forum.LoginHandler)
// 	mux.HandleFunc("/register", forum.RegisterHandler)
// 	mux.HandleFunc("/logout", forum.LogoutHandler)
// 	mux.HandleFunc("/postpage", forum.PostPageHandler)

// 	httpsServer := forum.MakeServer()
// 	httpsServer.Addr = ":443"
// 	// httpsServer.Addr = ":8080"
// 	httpsServer.Handler = mux
// 	httpsServer.TLSConfig = &tls.Config{GetCertificate: certMan.GetCertificate}

// 	fmt.Println("Starting server at port 443")
// ln, err := net.Listen("tcp", ":443")
// if err != nil {
// 	log.Fatal(err)
// }
// defer ln.Close()
// //rate limit here?
// httpsServer.ServeTLS(ln, "", "")

// fmt.Println("Starting server at port 443")
// // forum.RateLimit()
// // err := httpsServer.ListenAndServe()
// err := httpsServer.ListenAndServeTLS("", "")
// if err != nil {
// 	log.Fatal(err)
// }
// }
