package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/joho/godotenv"
	"go.limit.dev/unollm/grpcServer"
	"go.limit.dev/unollm/httpHandler"
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
	// start openai style server
	go func() {
		godotenv.Load("./.env")
		g := gin.New()
		g.Use(gin.Logger())
		corsCfg := cors.DefaultConfig()
		corsCfg.AllowAllOrigins = true
		corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "Authorization")
		g.Use(cors.New(corsCfg))
		g.Use(gin.Recovery())

		zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
		httpHandler.RegisterRoute(g, httpHandler.RegisterOpt{
			ChatGLMKey: zhipuaiApiKey,
		})
		g.Run("0.0.0.0:11451")
	}()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	// start grpc server
	svr := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
	))

	reflection.Register(svr)
	tcpList, err := net.Listen("tcp", "0.0.0.0:19198")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service := grpcServer.UnoForwardServer{}
	model.RegisterUnoLLMv1Server(svr, &service)
	svr.Serve(tcpList)
}
