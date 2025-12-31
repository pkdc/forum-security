package main

import (
	"crypto/tls"
	"fmt"
	"forum/forum"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	forum.InitDB()

	// Setup routes
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	mux.Handle("/", forum.RateLimiter(forum.HomeHandler))
	mux.Handle("/login", forum.RateLimiter(forum.LoginHandler))
	mux.Handle("/register", forum.RateLimiter(forum.RegisterHandler))
	mux.Handle("/logout", forum.RateLimiter(forum.LogoutHandler))
	mux.Handle("/postpage", forum.RateLimiter(forum.PostPageHandler))

	// Check if running on PaaS (PORT env var is set by Heroku, Render, Railway, Fly.io, etc.)
	// or if USE_AUTOCERT is explicitly disabled
	port := os.Getenv("PORT")
	useAutocert := os.Getenv("USE_AUTOCERT")

	if port != "" && useAutocert != "true" {
		// PaaS mode - Let PaaS handle TLS
		runOnPaaS(mux, port)
	} else {
		// VPS/Bare Metal mode - Use autocert for Let's Encrypt
		runWithAutocert(mux)
	}
}

// runOnPaaS runs the server for PaaS environments
// The PaaS provider handles TLS termination
func runOnPaaS(mux *http.ServeMux, port string) {
	server := forum.MakeServer()
	server.Addr = ":" + port
	server.Handler = mux

	fmt.Printf("Starting server on port %s (PaaS mode - TLS handled by platform)\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// runWithAutocert runs the server with automatic Let's Encrypt certificates
// Requires: control of ports 80 and 443, and a domain pointing to this server
func runWithAutocert(mux *http.ServeMux) {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "www.elephorum.com" // Default domain
	}

	dir := "./forum/certs"
	certMan := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(dir),
	}

	// HTTP server for ACME challenges (port 80)
	go func() {
		httpServer := forum.MakeServer()
		httpServer.Addr = ":80"
		httpServer.Handler = certMan.HTTPHandler(nil)

		fmt.Println("Starting HTTP server on port 80 for ACME challenges")
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// HTTPS server (port 443)
	httpsServer := forum.MakeServer()
	httpsServer.Addr = ":443"
	httpsServer.Handler = mux

	// Configure TLS
	manTlsConfig := certMan.TLSConfig()
	manTlsConfig.ServerName = domain
	manTlsConfig.GetCertificate = certMan.GetCertificate
	manTlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	}
	httpsServer.TLSConfig = manTlsConfig

	fmt.Printf("Starting HTTPS server on port 443 for domain: %s\n", domain)
	err := httpsServer.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
