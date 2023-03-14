package main

import (
	"context"
	"fmt"
	_ "github.com/advancemg/vimb-loader/docs"
	cfg "github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/pkg/logging/zap"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/routes"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/services"
	"github.com/advancemg/vimb-loader/pkg/storage/badger-client"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"
)

var (
	localHostRegex, _ = regexp.Compile("localhost|127.0.0.1|0.0.0.0")
)

// @title ВИМБ API
// @version 1.0
// @description Документация
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token
// @BasePath /
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := cfg.EditConfig()
	if err != nil {
		return err
	}
	err = cfg.Load()
	if err != nil {
		return err
	}
	err = zap.Init()
	if err != nil {
		return err
	}
	port := utils.GetEnv("PORT", ":8000")
	route := mux.NewRouter()
	route.PathPrefix("/api/v1/docs").Handler(httpSwagger.WrapHandler)
	route.HandleFunc("/api/v1", routes.Health).Methods("GET", "OPTIONS")
	/*dictionaries*/
	route.HandleFunc("/api/v1/channels", routes.AuthRequired(routes.PostGetChannels)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/channels/load", routes.AuthRequired(routes.PostLoadChannels)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/channels/query", routes.AuthRequired(routes.PostChannelsQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/adv-messages", routes.AuthRequired(routes.PostGetAdvMessages)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/adv-messages/load", routes.AuthRequired(routes.PostLoadAdvMessages)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/adv-messages/query", routes.AuthRequired(routes.PostAdvMessagesQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/ranks", routes.AuthRequired(routes.PostGetRanks)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/ranks/load", routes.AuthRequired(routes.PostLoadRanks)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/ranks/query", routes.AuthRequired(routes.PostRanksQuery)).Methods("POST", "OPTIONS")
	/*networks*/
	route.HandleFunc("/api/v1/program-breaks", routes.AuthRequired(routes.PostGetProgramBreaks)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/program-breaks/load", routes.AuthRequired(routes.PostLoadProgramBreaks)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/program-breaks/query", routes.AuthRequired(routes.PostProgramBreaksQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/program-breaks/pro-master/query", routes.AuthRequired(routes.PostProgramBreaksProMasterQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/program-breaks-light/load", routes.AuthRequired(routes.PostLoadProgramLightModeBreaks)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/program-breaks-light/query", routes.AuthRequired(routes.PostProgramLightModeBreaksQuery)).Methods("POST", "OPTIONS")
	/*deals*/
	route.HandleFunc("/api/v1/budgets", routes.AuthRequired(routes.PostGetBudgets)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/budgets/load", routes.AuthRequired(routes.PostLoadBudgets)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/budgets/query", routes.AuthRequired(routes.PostBudgetsQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/customers-with-advertisers", routes.AuthRequired(routes.PostGetCustomersWithAdvertisers)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/customers-with-advertisers/load", routes.AuthRequired(routes.PostLoadCustomersWithAdvertisers)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/customers-with-advertisers/query", routes.AuthRequired(routes.PostCustomersWithAdvertisersQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/customers-with-advertisers-data/query", routes.AuthRequired(routes.PostCustomersWithAdvertisersDataQuery)).Methods("POST", "OPTIONS")
	/*mediaPlans*/
	route.HandleFunc("/api/v1/mediaplan", routes.AuthRequired(routes.PutAddMPlan)).Methods("PUT", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan", routes.AuthRequired(routes.PostGetMPLans)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan/load", routes.AuthRequired(routes.PostLoadMPLans)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan/query", routes.AuthRequired(routes.PostMPLansQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/agg-mediaplan/query", routes.AuthRequired(routes.PostAggMPLansQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan/film ", routes.AuthRequired(routes.PostAddMPlanFilm)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan/film", routes.AuthRequired(routes.DeleteMPlanFilm)).Methods("DELETE", "OPTIONS")
	route.HandleFunc("/api/v1/mediaplan/change-film-planned-inventory", routes.AuthRequired(routes.PostChangeMPlanFilmPlannedInventory)).Methods("POST", "OPTIONS")
	/*spots*/
	route.HandleFunc("/api/v1/spot", routes.AuthRequired(routes.PostGetSpots)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/load", routes.AuthRequired(routes.PostLoadSpots)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/query", routes.AuthRequired(routes.PostSpotsQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/order-block/query", routes.AuthRequired(routes.PostSpotsOrderBlockQuery)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot", routes.AuthRequired(routes.PutAddSpot)).Methods("PUT", "OPTIONS")
	route.HandleFunc("/api/v1/spot", routes.AuthRequired(routes.DeleteSpot)).Methods("DELETE", "OPTIONS")
	route.HandleFunc("/api/v1/spot/change", routes.AuthRequired(routes.PostChangeSpot)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/set-position", routes.AuthRequired(routes.PostSetSpotPosition)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/change-films", routes.AuthRequired(routes.PostChangeFilms)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/spot/deleted-info", routes.AuthRequired(routes.PostGetDeletedSpotInfo)).Methods("POST", "OPTIONS")
	/*mq-metrics*/
	route.HandleFunc("/api/v1/mq/queues", routes.AuthRequired(routes.GetQueuesMetrics)).Methods("GET", "OPTIONS")
	/*backup*/
	route.HandleFunc("/api/v1/backup", routes.AuthRequired(routes.PostMongoBackup)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/backup-list", routes.AuthRequired(routes.PostListBackups)).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/v1/backup-restore", routes.AuthRequired(routes.PostMongoRestore)).Methods("POST", "OPTIONS")
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
	/* s3 server start*/
	s3Config := s3.InitConfig()
	if localHostRegex.MatchString(s3Config.S3Endpoint) {
		err = os.Mkdir(s3Config.S3LocalDir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
		//go func() {
		//	utils.CheckErr(s3Config.ServerStart())
		//}()
		/* s3 CreateDefaultBucket */
		for !s3Config.Ping() {
		}
		utils.CheckErr(s3.CreateDefaultBucket())
	}
	for !s3Config.Ping() {
	}
	/* Clean BadgerGC every 15 min*/
	if cfg.Config.Database != "mongodb" {
		go func() {
			badger_client.CleanGC()
		}()
	}
	/* amqp server */
	mqConfig := mq_broker.InitConfig()
	if localHostRegex.MatchString(mqConfig.MqHost) {
		go func() {
			utils.CheckErr(mqConfig.ServerStart())
		}()
		for !mqConfig.Ping() {
		}
	}
	for !mqConfig.Ping() {
	}
	/* amqp load services */
	go func() {
		service := services.LoadService{}
		utils.CheckErr(service.Start())
	}()
	/* amqp update services */
	go func() {
		service := services.UpdateService{}
		utils.CheckErr(service.Start())
	}()
	/*mongo backup service*/
	if cfg.Config.Database == "mongodb" && cfg.Config.Mongo.CronBackup != "" {
		go func() {
			service := services.BackupService{}
			utils.CheckErr(service.Start())
		}()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(ctx)
	return nil
}
