package main

import (
	"context"
	"fmt"
	"net/http"

	grpcHealth "github.com/harish908/Golang_Micro_B/proto/gen/health_check"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var healthServer grpcHealth.HealthCheckServiceClient

func main() {
	registerGRPCServer()
	registerHTTPServer()
}

func registerHTTPServer() {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service A up")
	})
	e.GET("/grpc/health", grpcHealthCheck)

	e.Logger.Fatal(e.Start("0.0.0.0:80"))
}

func registerGRPCServer() {
	conn, err := grpc.Dial("microservices.service-b:5001", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Server error", err)
	}
	healthServer = grpcHealth.NewHealthCheckServiceClient(conn)
}

func grpcHealthCheck(c echo.Context) error {
	res, err := healthServer.CheckGRPCConnection(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Errorf("Error for grpc health %s", err)
	}
	return c.String(http.StatusOK, res.Connection)
}
