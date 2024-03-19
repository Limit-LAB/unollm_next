package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.limit.dev/unollm/grpcServer"
	"go.limit.dev/unollm/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	// start grpc server
	svr := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
	))

	reflection.Register(svr)
	logger.Info("start grpc server", "addr", "0.0.0.0:19198")
	tcpList, err := net.Listen("tcp", "0.0.0.0:19198")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service := grpcServer.UnoForwardServer{}
	model.RegisterUnoLLMv1Server(svr, &service)
	svr.Serve(tcpList)
}
