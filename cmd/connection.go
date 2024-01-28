package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var conn *grpc.ClientConn

func init() {
	if os.Getenv("GRPC_PING_HOST") != "" {
		var err error
		conn, err = NewConn(os.Getenv("GRPC_PING_HOST"), os.Getenv("GRPC_PING_INSECURE") == "")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Starting without support for SendUpstream: configure with 'GRPC_PING_HOST' environment variable. E.g., example.com:443")
	}
}

// NewConn creates a new gRPC connection.
// host should be of the form domain:port, e.g., example.com:443
func NewConn(host string, insecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// Note: On the Windows platform, use of x509.SystemCertPool() requires
		// go version 1.18 or higher.
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}
	fmt.Sprintln("host : ", host)

	return grpc.Dial(host, opts...)
}
