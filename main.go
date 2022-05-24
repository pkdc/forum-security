package main

import (
	"crypto/tls"
	"fmt"
	"forum/forum"
	"log"
	"net"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	forum.InitDB()
	// forum.ClearUsers()
	// forum.ClearPosts()
	// forum.ClearComments()

	dir := "./certs"
	certMan := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// HostPolicy: autocert.HostWhitelist("www.domain.com"),
		HostPolicy: nil,
		Cache:      autocert.DirCache(dir),
	}
	go func() {
		httpServer := forum.MakeServer()
		httpServer.Addr = ":80"
		httpServer.Handler = certMan.HTTPHandler(nil)

		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	mux.Handle("/", forum.RateLimiter(forum.HomeHandler))
	mux.Handle("/login", forum.RateLimiter(forum.LoginHandler))
	mux.Handle("/register", forum.RateLimiter(forum.RegisterHandler))
	mux.Handle("/logout", forum.RateLimiter(forum.LogoutHandler))
	mux.Handle("/postpage", forum.RateLimiter(forum.PostPageHandler))

	httpsServer := forum.MakeServer()
	httpsServer.Addr = ":443"

	var hello tls.ClientHelloInfo
	hello.ServerName = "instance-1@elephorum.com"
	fmt.Printf("https ClientHelloInfo: %s", hello.ServerName)
	// how to pass hello into httpsServer?
	httpsServer.Handler = mux
	httpsServer.TLSConfig = &tls.Config{GetCertificate: certMan.GetCertificate}

	fmt.Println("Starting server at port 443")
	ln, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	httpsServer.ServeTLS(ln, "", "")

	fmt.Println("Starting server at port 443")
	// err := httpsServer.ListenAndServe()
	err = httpsServer.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
