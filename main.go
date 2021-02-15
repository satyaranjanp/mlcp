package main

import (
	"mlcp/pkg/config"
	"mlcp/pkg/mlcp"
	"mlcp/pkg/signals"
	"github.com/golang/glog"
)

var (
	addr string = "127.0.0.1"
	port string = "8443"
	capath = "/etc/kubernetes/ssl/ca.crt"
	certPath = "/etc/certs/tls.crt"
	keyPath = "/etc/certs/tls.key"
)
func main() {
	config.Init()
	stopCh := signals.SetupSignalHandler()

	ms, err := mlcp.NewMlcpServer(addr, port, capath, certPath, keyPath)
	if err != nil {
		glog.Fatalf("Error initializing server: %v", err)
	}

	if err := ms.Run(stopCh); err != nil {
		glog.Fatalf("Error starting server: %v", err)
	}

	<-stopCh
}

