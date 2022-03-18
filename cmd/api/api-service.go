package main

import (
	"context"
	"fmt"
	_ "github.com/advancemg/vimb-loader/docs"
	"github.com/advancemg/vimb-loader/pkg/routes"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title ВИМБ API
// @version 1.0
// @description Документация
// @BasePath /
func main() {
	port := utils.GetEnv("PORT", ":8000")
	route := mux.NewRouter()
	route.PathPrefix("/api/v1/docs").Handler(httpSwagger.WrapHandler)
	route.HandleFunc("/api/v1", routes.Health).Methods("GET", "OPTIONS")
	route.HandleFunc("/api/v1/get-program-brakes", routes.PostGetProgramBreaks).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/get-adv-messages", routes.PostGetAdvMessages).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/get-budgets", routes.PostGetBudgets).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/get-channels", routes.PostGetChannels).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/customers-with-advertisers", routes.PostGetCustomersWithAdvertisers).Methods("POST", "OPTIONS")
	s := &http.Server{
		Addr:         port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 1,
		Handler:      route,
	}
	go func() {
		fmt.Print("Base... [http://localhost", port, "/api/v1", "]\n")
		fmt.Print("Base... [http://localhost", port, "/api/v1/docs/index.html", "]\n")
		utils.CheckErr(s.ListenAndServe())
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
