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
// 	// forum.ClearUsers()
// 	// forum.ClearPosts()
// 	// forum.ClearComments()

// 	dir := "./forum/certs"
// 	certMan := &autocert.Manager{
// 		Prompt:     autocert.AcceptTOS,
// 		HostPolicy: autocert.HostWhitelist("www.elephorum.com"),
// 		// HostPolicy: nil,
// 		Cache: autocert.DirCache(dir),
// 	}

// 	// httpServer is for communicating with the CA to get the certs
// 	go func() {
// 		httpServer := forum.MakeServer()
// 		httpServer.Addr = ":80"
// 		httpServer.Handler = certMan.HTTPHandler(nil)

// 		err := httpServer.ListenAndServe()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	mux := http.NewServeMux()
// 	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
// 	mux.Handle("/", forum.RateLimiter(forum.HomeHandler))
// 	mux.Handle("/login", forum.RateLimiter(forum.LoginHandler))
// 	mux.Handle("/register", forum.RateLimiter(forum.RegisterHandler))
// 	mux.Handle("/logout", forum.RateLimiter(forum.LogoutHandler))
// 	mux.Handle("/postpage", forum.RateLimiter(forum.PostPageHandler))

// 	// create a https server to serve the content of the website
// 	httpsServer := forum.MakeServer()
// 	httpsServer.Addr = ":443"
// 	httpsServer.Handler = mux

// 	// put a ServerName and a list of cipher suites into the tls config
// 	manTlsConfig := certMan.TLSConfig()
// 	manTlsConfig.ServerName = "www.elephorum.com"
// 	manTlsConfig.GetCertificate = certMan.GetCertificate
// 	manTlsConfig.CipherSuites = []uint16{
// 		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
// 		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
// 		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
// 	}
// 	httpsServer.TLSConfig = manTlsConfig

// 	fmt.Println("Starting server at port 443")
// 	err := httpsServer.ListenAndServeTLS("", "")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
