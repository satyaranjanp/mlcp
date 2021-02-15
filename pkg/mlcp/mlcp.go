package mlcp

import (
	"crypto/tls"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"mlcp/pkg/cache"
	"mlcp/pkg/common"
	"mlcp/pkg/database"
	"mlcp/pkg/message_queue"
	"mlcp/pkg/worker"
	"net/http"
	"strings"
)

const (
	ErrorNoAvailableSlot = "NoAvailableSlot"
)

type MlcpServer struct {
	mlcp *Mlcp
	address string
	port string
	caPath string
	certPath string
	keyPath string
	handlers *http.ServeMux
	tlsConfig *tls.Config
}

func NewMlcpServer(addr, port, ca, cert, key string) (*MlcpServer, error) {
	s := &MlcpServer{
		address:   strings.TrimSuffix(addr, ":"),
		port:      port,
		caPath:    ca,
		certPath:  cert,
		keyPath:   key,
		handlers: &http.ServeMux{},
	}
	
	s.handlers.HandleFunc("/assignSlot", s.AssignSlot)
	return s, nil
}

func (ms *MlcpServer) Run(stopCh <-chan struct{}) error {
	mlcp, err := initMlcp(stopCh)
	if err != nil {
		return fmt.Errorf("Error starting server: %v", err)
	}
	ms.mlcp = mlcp

	srv := &http.Server{
		Addr: ms.address+":"+ms.port,
		TLSConfig: ms.tlsConfig,
		Handler: ms.handlers,
	}
	go func(ch <-chan struct{}) {
		err = srv.ListenAndServeTLS(ms.certPath, ms.keyPath)
		if err != nil {
			glog.Errorf("Error starting mlcp server: %v", err)
		}
	}(stopCh)
	glog.Errorf("Shuting down mlcp server")
	return nil
}

type Mlcp struct {
	c cache.Cache
	db *database.Database
	mq message_queue.MessageQueue
	wq *worker.WorkeQueue
}

func initMlcp(stopCh <-chan struct{}) (*Mlcp, error) {
	m := &Mlcp{}
	var err error
	if m.c, err = cache.SetupCache(); err != nil {
		return nil, fmt.Errorf("Error setting up cache: %v", err)
	}
	if m.db, err = database.InitializeDatabase(); err != nil {
		return nil, fmt.Errorf("Error initializing database: %v", err)
	}
	if m.mq, err = message_queue.InitMQ(); err != nil {
		return nil, fmt.Errorf("Error initializing message queue: %v", err)
	}
	wq := worker.NewWorkQueue(1, m.c, m.db)
	m.wq = wq
	go wq.Run(stopCh)
	go m.mq.ListenToQueue(m.wq)
	return m, nil
}

func (ms *MlcpServer) AssignSlot(w http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		glog.Errorf("Error parsing request body: %v", err)
		return
	}
	r, err := common.ParseRequest(reqBody)
	if err != nil {
		glog.Errorf("Error parsing request body: %v", err)
	}
	ms.mlcp.wq.Add(interface{}(common.NewSlot(common.NewCar(r.RegnNo), r.SlotId)))
}
