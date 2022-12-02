package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type server struct {
	httpServer *http.ServeMux
	grpcServer *grpc.Server
}

func main() {
	logrus.Println("Starting the http server: ", httpPort)
}
