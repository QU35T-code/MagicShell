package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Iface  string `short:"i" long:"host" description:"Interface address on which to bind" default:"127.0.0.1" required:"true"`
	Port   string `short:"p" long:"port" description:"Port on which to bind" default:"9000" required:"true"`
	Keys   string `short:"k" long:"keys" description:"Path to folder with server.{pem,key}" default:"./certs" required:"true"`
	Socket string `short:"s" long:"socket" description:"Domain socket from which the program reads"`
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}
}

func Server() {
	var listener net.Listener
	var err error

	fmt.Println("[Start] Listener...")
	listener, err = newTLSListener()
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func newTLSListener() (net.Listener, error) {
	pem := path.Join(opts.Keys, "server.pem")
	key := path.Join(opts.Keys, "server.key")
	cer, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	connStr := fmt.Sprintf("%s:%s", opts.Iface, opts.Port)
	return tls.Listen("tcp", connStr, config)
}

func handleConnection(conn net.Conn) {
	fmt.Println("[EVENT] Connection received...")
}
