package run

import (
	"crypto/tls"
	"github.com/DrmagicE/gmqtt/server"
	"net"
	"net/http"
	"time"
	"github.com/DrmagicE/gmqtt/logger"
	"os"
	"log"
)

func NewServer(config *Config) (*server.Server, error) {
	startProfile(config.ProfileConfig.CPUProfile,config.ProfileConfig.MemProfile)
	srv := server.NewServer()
	srv.SetDeliveryRetryInterval(time.Second * time.Duration(config.DeliveryRetryInterval))
	srv.SetMaxInflightMessages(config.MaxInflightMessages)
	srv.SetQueueQos0Messages(config.QueueQos0Messages)
	var l net.Listener
	var ws *server.WsServer
	var err error
	for _, v := range config.Listener {
		if v.Protocol == ProtocolMQTT {
			if v.KeyFile == "" {
				l, err = net.Listen("tcp", v.Addr)
				if err != nil {
					return nil, err
				}
			} else {
				crt, err := tls.LoadX509KeyPair(v.CertFile, v.KeyFile)
				if err != nil {
					return nil, err
				}
				tlsConfig := &tls.Config{}
				tlsConfig.Certificates = []tls.Certificate{crt}
				l, err = tls.Listen("tcp", v.Addr, tlsConfig)
			}
			srv.AddTCPListenner(l)
		} else {
			if v.KeyFile == "" {
				ws = &server.WsServer{
					Server: &http.Server{Addr: v.Addr},
				}
			} else {
				ws = &server.WsServer{
					Server:   &http.Server{Addr: v.Addr},
					CertFile: v.CertFile,
					KeyFile:  v.KeyFile,
				}
			}
			srv.AddWebSocketServer(ws)
		}
	}
	if config.Logging  {
		server.SetLogger(logger.NewLogger(os.Stderr, "", log.LstdFlags))
	}

	return srv, nil
}
