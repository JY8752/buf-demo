package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"buf.build/gen/go/jyapp/weather/connectrpc/go/jyapp/weather/v1/weatherv1connect"
	weatherv1 "buf.build/gen/go/jyapp/weather/protocolbuffers/go/jyapp/weather/v1"
	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"connectrpc.com/grpcreflect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const port = ":8080"

type weatherService struct {
}

func (w *weatherService) GetWeather(ctx context.Context, req *connect.Request[weatherv1.GetWeatherRequest]) (*connect.Response[weatherv1.GetWeatherResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&weatherv1.GetWeatherResponse{
		Temperature: 1.0,
		Conditions:  weatherv1.Condition_CONDITION_SUNNY,
	})
	res.Header().Set("Weather-Version", "v1")
	return res, nil
}

func withCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}

func main() {
	ws := &weatherService{}
	mux := http.NewServeMux()
	path, handler := weatherv1connect.NewWeatherServiceHandler(ws)
	mux.Handle(path, handler)

	// サーバーリフレクション対応
	reflector := grpcreflect.NewStaticReflector("jyapp.weather.v1.WeatherSerice")
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	http.ListenAndServe(
		fmt.Sprintf("localhost%s", port),
		withCORS(h2c.NewHandler(mux, &http2.Server{})),
	)
}
