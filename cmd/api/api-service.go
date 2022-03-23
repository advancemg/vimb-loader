package main

import (
	"context"
	"fmt"
	_ "github.com/advancemg/vimb-loader/docs"
	mq "github.com/advancemg/vimb-loader/pkg/mq-broker"
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
	/*dictionaries*/
	route.HandleFunc("/api/v1/channels", routes.PostGetChannels).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/adv-messages", routes.PostGetAdvMessages).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/ranks", routes.PostGetRanks).Methods("POST", "OPTIONS")
	/*networks*/
	route.HandleFunc("/api/v1/program-breaks", routes.PostGetProgramBreaks).Methods("POST", "OPTIONS")
	/*deals*/
	route.HandleFunc("/api/v1/budgets", routes.PostGetBudgets).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/customers-with-advertisers", routes.PostGetCustomersWithAdvertisers).Methods("POST", "OPTIONS")
	/*mediaPlans*/
	route.HandleFunc("/api/v1/mediaplan", routes.PutAddMPlan).Methods("PUT", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan", routes.PostGetMPLans).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan/film ", routes.PostAddMPlanFilm).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan/film", routes.DeleteMPlanFilm).Methods("DELETE", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan/change-film-planned-inventory", routes.PostChangeMPlanFilmPlannedInventory).Methods("POST", "OPTIONS")
	/*spots*/
	route.HandleFunc("/api/v1/spot", routes.PostGetSpots).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot", routes.PutAddSpot).Methods("PUT", "OPTIONS")
	route.HandleFunc("/api/v1/spot", routes.DeleteSpot).Methods("DELETE", "OPTIONS")
	route.HandleFunc("/api/v1/spot/change", routes.PostChangeSpot).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/set-position", routes.PostSetSpotPosition).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/change-films", routes.PostChangeFilms).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/deleted-info", routes.PostGetDeletedSpotInfo).Methods("POST", "OPTIONS")

	s := &http.Server{
		Addr:         port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 1,
		Handler:      route,
	}
	/* api */
	go func() {
		fmt.Print("Base... [http://localhost", port, "/api/v1", "]\n")
		fmt.Print("Base... [http://localhost", port, "/api/v1/docs/index.html", "]\n")
		utils.CheckErr(s.ListenAndServe())
	}()
	/* s3 */
	/* amqp server */
	go func() {
		config := mq.InitConfig()
		utils.CheckErr(config.ServerStart())
	}()
	/* amqp load services */

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}