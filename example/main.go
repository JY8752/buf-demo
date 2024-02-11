package main

import (
	weatherv1 "buf-demo/example/gen/go/junichiy/weather/v1"
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":8080"

type weatherService struct {
	weatherv1.UnimplementedWeatherServiceServer
}

func (w *weatherService) GetWeather(ctx context.Context, req *weatherv1.GetWeatherRequest) (*weatherv1.GetWeatherResponse, error) {
	return &weatherv1.GetWeatherResponse{
		Temperature: 1.0,
		Conditions:  weatherv1.Condition_CONDITION_SUNNY,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	weatherv1.RegisterWeatherServiceServer(s, &weatherService{})

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
