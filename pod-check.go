package main

import (

	"fmt"
	"context"
	"crypto/tls"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

var (
	tlscert, tlskey string
)

func main() {

	flag.StringVar(&tlscert, "tlsCertFile", "/etc/certs/cert.pem", "File contaains the X509 Certificate for HTTPS")
	flag.StringVar(&tlskey, "tlsKeyFile", "/etc/certs/key.pem", "File containing the X509 private key")

	flag.Parse()

	certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		glog.Errorf("Failed to load key pair: %v", err)
	}

	port := getEnv("PORT", "8080")

	server := &http.Server {
		Addr: fmt.Sprintf(":%v", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{certs}},
	}

	gs := myValidServerhandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", gs.serve)
	server.Handler = mux


	go func() {
		if err := server.ListenAndServeTLS("",""); err != nil {
			glog.Errorf("Failed to listen and serve WEB hook Server: %v", err)
		}
	}()

	glog.Infof("Server is running on port : %s", port)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	glog.Info("Get shutdown signal, sutting down webhook server gracefully...")
	server.Shutdown(context.Background())

}