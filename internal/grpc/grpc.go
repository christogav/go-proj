package grpc

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newGrpc(lifecycle fx.Lifecycle, logger *zerolog.Logger, config *Config) (*grpc.Server, error) {
	uriString := fmt.Sprintf("tcp://%v:%v", config.Host, config.Port)

	uri, err := url.Parse(uriString)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen(uri.Scheme, uri.Host)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				reflection.Register(server)

				logger.Info().Msgf("Listening on: %s", uri.String())
				if err := server.Serve(listener); err != nil {
					logger.Fatal().Err(err).Send()
				}
			}()

			logger.Info().
				Str("address", uri.String()).
				Msg("gRPC server listening")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()
			return nil
		},
	})

	return server, nil
}
