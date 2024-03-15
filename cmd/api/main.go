package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	api "github.com/felipevillarrealdaza/go-service-template/internal/api/http"
	"github.com/felipevillarrealdaza/go-service-template/internal/config"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	// Parse env variables into config
	apiConfig := config.ApiConfig{}
	if configErr := envconfig.Process(context.Background(), &apiConfig); configErr != nil {
		panic(fmt.Sprintf("could not parse API config: %+v\n", configErr))
	}
	dbConfig := config.DbConfig{}
	if configErr := envconfig.Process(context.Background(), &dbConfig); configErr != nil {
		panic(fmt.Sprintf("could not parse DB config: %+v\n", configErr))
	}

	// Create DB connection
	dbCtx, dbErr := sql.Open("postgres", dbConfig.RetrieveDBConnectionString())
	if dbErr != nil {
		panic("could not create db connection!")
	}
	defer dbCtx.Close()

	// Create channel to listen for SIGTERM event
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM)

	// Spawn goroutine to run the HTTP API
	serverErr := make(chan error, 1)
	go createAndStartHttpServer(dbCtx, apiConfig, serverErr)

	// Hold execution and listen for errors on both channels
	listenForErrorsAndHandleGracefully(serverErr, shutdown)
}

func createAndStartHttpServer(dbCtx *sql.DB, apiConfig config.ApiConfig, serverErr chan error) {
	handler := createHttpApiHandler(dbCtx)
	server := createHttpServer(apiConfig, handler)
	serverErr <- server.ListenAndServe()
}

func createHttpApiHandler(dbCtx *sql.DB) http.Handler {
	return api.NewRouter(dbCtx)
}

func createHttpServer(apiConfig config.ApiConfig, handler http.Handler) http.Server {
	return http.Server{
		Addr:    apiConfig.RetrieveApiAddress(),
		Handler: handler,
	}
}

func listenForErrorsAndHandleGracefully(serverErr chan error, shutdown chan os.Signal) {
	err := <-serverErr
	fmt.Print(err.Error())
	os.Exit(1)
}
