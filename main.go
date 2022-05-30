package main

import (
	"crypto/tls"
	"fmt"
	"forum/forum"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

// func MyGetCertificate(man *autocert.Manager) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
// 	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
// 		hello.ServerName = "www.elephorum.com"
// 		hello.CipherSuites = [TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256]
// 		fmt.Printf("https ClientHelloInfo's ServerName: %s\n", hello.ServerName)
// 		// cipher suite
// 		return man.GetCertificate(hello)
// 	}
// }

func main() {
	forum.InitDB()
	// forum.ClearUsers()
	// forum.ClearPosts()
	// forum.ClearComments()

	dir := "./forum/certs"
	certMan := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("www.elephorum.com"),
		// HostPolicy: nil,
		Cache: autocert.DirCache(dir),
	}

	// httpServer is to communicate to the CA to get the certs
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

	// create https server to serve content of our website
	httpsServer := forum.MakeServer()
	httpsServer.Addr = ":443"
	httpsServer.Handler = mux

	// write a custom GetCertificate func and put a ServerName and
	// a list of cipher suites into the server hello msg
	manTlsConfig := certMan.TLSConfig()
	manTlsConfig.ServerName = "www.elephorum.com"
	manTlsConfig.GetCertificate = certMan.GetCertificate
	manTlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	}
	httpsServer.TLSConfig = manTlsConfig

	fmt.Println("Starting server at port 443")
	err := httpsServer.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
