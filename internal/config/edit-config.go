package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/models"
	mq "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"io/ioutil"
	"os"
	"strconv"
)

type configuration struct {
	Mediaplan                models.MediaplanConfiguration                `json:"mediaplan"`
	Budget                   models.BudgetConfiguration                   `json:"budget"`
	Channel                  models.ChannelConfiguration                  `json:"channel"`
	AdvMessages              models.AdvMessagesConfiguration              `json:"advMessages"`
	CustomersWithAdvertisers models.CustomersWithAdvertisersConfiguration `json:"customersWithAdvertisers"`
	DeletedSpotInfo          models.DeletedSpotInfoConfiguration          `json:"deletedSpotInfo"`
	Rank                     models.RanksConfiguration                    `json:"rank"`
	ProgramBreaks            models.ProgramBreaksConfiguration            `json:"programBreaks"`
	ProgramBreaksLight       models.ProgramBreaksLightConfiguration       `json:"programBreaksLight"`
	Spots                    models.SpotsConfiguration                    `json:"spots"`
	S3Cfg                    s3.Config                                    `json:"s3"`
	AmqpConfig               mq.Config                                    `json:"amqp"`
	Url                      string                                       `json:"url"`
	Cert                     string                                       `json:"cert"`
	Password                 string                                       `json:"password"`
	Client                   string                                       `json:"client"`
	Timeout                  string                                       `json:"timeout"`
}

func EditConfig() {
	flagEdit := flag.Bool("config", false, "a bool")
	flag.Parse()
	if *flagEdit {
		enterConfig()
	}
}

func enterConfig() {
	cfg := &configuration{}
	open, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			marshal, err := json.MarshalIndent(cfg, "", "  ")
			checkErr(err)
			err = os.WriteFile("config.json", marshal, 0666)
			checkErr(err)
		} else {
			checkErr(err)
		}
	}
	fmt.Println("Default Config:")
	fmt.Println(string(open))
	fmt.Printf("%s", "Edit config? (Y/n):")
	n, err := readLine()
	checkErr(err)
	if n == "y" || n == "Y" || n == "yes" || n == "Yes" {
		/*Budget*/
		fmt.Printf("%s", "Enter Budget cron(docker 0 0/46 * * *):")
		line, err := readLine()
		checkErr(err)
		if line == "" {
			line = "0 0/46 * * *"
		}
		cfg.Budget.Cron = line
		fmt.Printf("%s", "Enter Budget sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		cfg.Budget.SellingDirection = line
		fmt.Printf("%s", "Budget loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err := strconv.ParseBool(line)
		checkErr(err)
		cfg.Budget.Loading = loading
		/*ProgramBreaks*/
		fmt.Printf("%s", "Enter ProgramBreaks cron(docker 0 0 0/8 * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0 0 0/8 * *"
		}
		cfg.ProgramBreaks.Cron = line
		fmt.Printf("%s", "Enter ProgramBreaks sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		cfg.ProgramBreaks.SellingDirection = line
		fmt.Printf("%s", "ProgramBreaks loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.ProgramBreaks.Loading = loading
		/*ProgramBreaksLight*/
		fmt.Printf("%s", "Enter ProgramBreaksLight cron(docker 0/2 * * * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0/2 * * * *"
		}
		cfg.ProgramBreaksLight.Cron = line
		fmt.Printf("%s", "Enter ProgramBreaksLight sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		cfg.ProgramBreaksLight.SellingDirection = line
		fmt.Printf("%s", "ProgramBreaksLight loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.ProgramBreaksLight.Loading = loading
		/*Mediaplan*/
		fmt.Printf("%s", "Enter Mediaplan cron(docker 0 0/20 * * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0 0/20 * * *"
		}
		cfg.Mediaplan.Cron = line
		fmt.Printf("%s", "Enter Mediaplan sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		cfg.Mediaplan.SellingDirection = line
		fmt.Printf("%s", "Mediaplan loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.Mediaplan.Loading = loading
		/*AdvMessages*/
		fmt.Printf("%s", "Enter AdvMessages cron(docker 0 0/2 * * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0 0/2 * * *"
		}
		cfg.AdvMessages.Cron = line
		fmt.Printf("%s", "AdvMessages loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.AdvMessages.Loading = loading
		/*Spots*/
		fmt.Printf("%s", "Enter Spots cron(docker 0 0/59 * * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0 0/59 * * *"
		}
		cfg.Spots.Cron = line
		fmt.Printf("%s", "Enter Spots sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		cfg.Spots.SellingDirection = line
		fmt.Printf("%s", "Spots loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.Spots.Loading = loading
		/*DeletedSpotInfo*/
		fmt.Printf("%s", "Enter DeletedSpotInfo cron(docker 0 0 0/12 * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0 0 0/12 * *"
		}
		cfg.DeletedSpotInfo.Cron = line
		fmt.Printf("%s", "DeletedSpotInfo loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.DeletedSpotInfo.Loading = loading
		/*Channel*/
		fmt.Printf("%s", "Enter Channels cron(docker 0 0/18 * * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0 0/18 * * *"
		}
		cfg.Channel.Cron = line
		fmt.Printf("%s", "Enter Channels sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		cfg.Channel.SellingDirection = line
		fmt.Printf("%s", "Channels loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.Channel.Loading = loading
		/*CustomersWithAdvertisers*/
		fmt.Printf("%s", "Enter CustomersWithAdvertisers cron(docker 0 0/16 * * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0 0/16 * * *"
		}
		cfg.CustomersWithAdvertisers.Cron = line
		fmt.Printf("%s", "Enter CustomersWithAdvertisers sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		cfg.CustomersWithAdvertisers.SellingDirection = line
		fmt.Printf("%s", "CustomersWithAdvertisers loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.CustomersWithAdvertisers.Loading = loading
		/*Rank*/
		fmt.Printf("%s", "Enter Rank cron(docker 0 0 0/23 * *):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "0 0 0/23 * *"
		}
		cfg.Rank.Cron = line
		fmt.Printf("%s", "Rank loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		cfg.Rank.Loading = loading
		/*AMQP*/
		fmt.Printf("%s", "Enter amqp host(docker localhost):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "localhost"
		}
		cfg.AmqpConfig.MqHost = line
		fmt.Printf("%s", "Enter amqp port(docker 5555):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "5555"
		}
		cfg.AmqpConfig.MqPort = line
		fmt.Printf("%s", "Enter amqp username(docker guest):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "guest"
		}
		cfg.AmqpConfig.MqUsername = line
		fmt.Printf("%s", "Enter amqp password(docker guest):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "guest"
		}
		cfg.AmqpConfig.MqPassword = line
		/*S3*/
		fmt.Printf("%s", "Enter S3 AccessKeyId(docker minioadmin):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "minioadmin"
		}
		cfg.S3Cfg.S3AccessKeyId = line
		fmt.Printf("%s", "Enter S3 SecretAccessKey(docker minioadmin):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "minioadmin"
		}
		cfg.S3Cfg.S3SecretAccessKey = line
		fmt.Printf("%s", "Enter S3 Region(docker us-west-0):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "us-west-0"
		}
		cfg.S3Cfg.S3Region = line
		fmt.Printf("%s", "Enter S3 Endpoint(docker 127.0.0.1:9999):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "127.0.0.1:9999"
		}
		cfg.S3Cfg.S3Endpoint = line
		fmt.Printf("%s", "Enter S3 Debug(docker true):")
		line, err = readLine()
		checkErr(err)
		if line != "false" {
			line = "true"
		}
		cfg.S3Cfg.S3Debug = line
		fmt.Printf("%s", "Enter S3 Bucket(docker storage):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "storage"
		}
		cfg.S3Cfg.S3Bucket = line
		fmt.Printf("%s", "Enter S3 LocalDir(docker s3-buckets):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "s3-buckets"
		}
		cfg.S3Cfg.S3LocalDir = line
		/*VIMB*/
		fmt.Printf("%s", "Enter url(docker https://vimb-svc.vitpc.com:436/VIMBService.asmx):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "https://vimb-svc.vitpc.com:436/VIMBService.asmx"
		}
		cfg.Url = line
		fmt.Printf("%s", "Enter cert:")
		line, err = readLine()
		checkErr(err)
		cfg.Cert = line
		fmt.Printf("%s", "Enter password:")
		line, err = readLine()
		checkErr(err)
		fmt.Printf("%s", "Enter client(docker test_client):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "test_client"
		}
		cfg.Client = line
		fmt.Printf("%s", "Enter timeout(docker 120s):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "120s"
		}
		cfg.Timeout = line
		marshal, err := json.MarshalIndent(cfg, "", "  ")
		checkErr(err)
		err = os.WriteFile("config.json", marshal, 0666)
		checkErr(err)
		fmt.Println(string(marshal))
	}
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	result := ""
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return result, err
		}
		result += string(line)
		if !isPrefix {
			return result, nil
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
