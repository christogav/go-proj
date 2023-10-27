package hellov1

import (
	"context"

	pb "github.com/christogav/go-proj/pkg/api/hello/v1"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type HelloServer struct {
	pb.UnimplementedTestServiceServer
}

var (
	_ pb.TestServiceServer = (*HelloServer)(nil)

	Module = fx.Module(
		"hellov1",

		fx.Provide(newHelloServer),

		fx.Invoke(func(server *HelloServer) {
			// Required to pull the dependency and initialise the server.
		}),
	)
)

func newHelloServer(grpc *grpc.Server) *HelloServer {
	server := &HelloServer{}
	pb.RegisterTestServiceServer(grpc, server)

	return server
}

func (s *HelloServer) Hello(_ context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Info().Msg("Hello")
	return &pb.HelloResponse{Response: request.Message}, nil
}
