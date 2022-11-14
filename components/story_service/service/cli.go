package service

// type ServerGRPC struct {
// 	logger  *logrus.Logger
// 	server  *grpc.Server
// 	Address string

// 	certificate string
// 	key         string
// 	destDir     string
// }

// type ServerGRPCConfig struct {
// 	Certificate string
// 	Key         string
// 	Address     string
// 	DestDir     string
// }

// func NewServerGRPC(cfg ServerGRPCConfig) (server ServerGRPC, err error) {
// 	server.logger = logrus.New()

// 	if cfg.Address == "" {
// 		err = errors.Errorf("Address must be specified")
// 		server.logger.Error("Address must be specified")
// 		return
// 	}

// 	server.Address = cfg.Address
// 	server.certificate = cfg.Certificate
// 	server.key = cfg.Key

// 	if _, err = os.Stat(cfg.DestDir); err != nil {
// 		server.logger.Error("Directory doesn't exist")
// 		return
// 	}

// 	server.destDir = cfg.DestDir
// 	return
// }

// func StartServerCommand() cli.Command {

// 	return cli.Command{
// 		Name:  "serve",
// 		Usage: "initiates a gRPC server",

// 		Flags: []cli.Flag{
// 			&cli.StringFlag{
// 				Name:  "a",
// 				Usage: "Address to listen",
// 				Value: "localhost:80",
// 			},

// 			&cli.StringFlag{
// 				Name:  "key",
// 				Usage: "path to TLS certificate",
// 			},
// 			&cli.StringFlag{
// 				Name:  "certificate",
// 				Usage: "path to TLS certificate",
// 			},
// 			&cli.StringFlag{
// 				Name:  "d",
// 				Usage: "Destination directory Default is /cloud",
// 				Value: "/cloud",
// 			},
// 		},
// 		Action: func(c *cli.Context) error {
// 			grpcServer, err := NewServerGRPC(ServerGRPCConfig{
// 				Address:     c.String("a"),
// 				Certificate: c.String("certificate"),
// 				Key:         c.String("key"),
// 				DestDir:     c.String("d"),
// 			})
// 			if err != nil {
// 				fmt.Println("error is creating server")

// 				return err
// 			}
// 			serverGrpc := &grpcServer
// 			err = serverGrpc.Listen()

// 			defer serverGrpc.Close()
// 			return nil
// 		},
// 	}

// }

// func (s *ServerGRPC) Listen() (err error) {
// 	var (
// 		listener  net.Listener
// 		grpcOpts  = []grpc.ServerOption{}
// 		grpcCreds credentials.TransportCredentials
// 	)

// 	listener, err = net.Listen("tcp", s.Address)
// 	if err != nil {
// 		err = errors.Wrapf(err,
// 			"failed to listen on  %d",
// 			s.Address)
// 		return
// 	}

// 	if s.certificate != "" && s.key != "" {
// 		grpcCreds, err = credentials.NewServerTLSFromFile(
// 			s.certificate, s.key)
// 		if err != nil {
// 			err = errors.Wrapf(err,
// 				"failed to create tls grpc server using cert %s and key %s",
// 				s.certificate, s.key)
// 			return
// 		}

// 		grpcOpts = append(grpcOpts, grpc.Creds(grpcCreds))
// 	}

// 	s.server = grpc.NewServer(grpcOpts...)
// 	// pb.RegisterStoryServiceHandler(context.Background(), s.server, s)
// 	pb.RegisterStoryServiceServer(s.server, s)

// 	err = s.server.Serve(listener)
// 	if err != nil {
// 		err = errors.Wrapf(err, "errored listening for grpc connections")
// 		return
// 	}

// 	return
// }

// func (s *ServerGRPC) Close() {
// 	if s.server != nil {
// 		s.server.Stop()
// 	}

// 	return
// }
