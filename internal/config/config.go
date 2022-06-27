package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Amqp struct {
	MqHost     string `json:"Host"`
	MqPort     string `json:"Port"`
	MqUsername string `json:"Username"`
	MqPassword string `json:"Password"`
}

type Mongo struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	DB       string `json:"Db"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type S3 struct {
	S3AccessKeyId     string `json:"s3AccessKeyId"`
	S3SecretAccessKey string `json:"s3SecretAccessKey"`
	S3Region          string `json:"s3Region"`
	S3Endpoint        string `json:"s3Endpoint"`
	S3Debug           string `json:"s3Debug"`
	S3Bucket          string `json:"s3Bucket"`
	S3LocalDir        string `json:"s3LocalDir"`
}

type MediaplanConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

type BudgetConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

type ChannelConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

type AdvMessagesConfiguration struct {
	Cron    string `json:"cron"`
	Loading bool   `json:"loading"`
}

type CustomersWithAdvertisersConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

type DeletedSpotInfoConfiguration struct {
	Cron    string `json:"cron"`
	Loading bool   `json:"loading"`
}

type RanksConfiguration struct {
	Cron    string `json:"cron"`
	Loading bool   `json:"loading"`
}

type ProgramBreaksConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

type ProgramBreaksLightConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

type SpotsConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

type Configuration struct {
	Mediaplan                MediaplanConfiguration                `json:"mediaplan"`
	Budget                   BudgetConfiguration                   `json:"budget"`
	Channel                  ChannelConfiguration                  `json:"channel"`
	AdvMessages              AdvMessagesConfiguration              `json:"advMessages"`
	CustomersWithAdvertisers CustomersWithAdvertisersConfiguration `json:"customersWithAdvertisers"`
	DeletedSpotInfo          DeletedSpotInfoConfiguration          `json:"deletedSpotInfo"`
	Rank                     RanksConfiguration                    `json:"rank"`
	ProgramBreaks            ProgramBreaksConfiguration            `json:"programBreaks"`
	ProgramBreaksLight       ProgramBreaksLightConfiguration       `json:"programBreaksLight"`
	Spots                    SpotsConfiguration                    `json:"spots"`
	S3                       S3                                    `json:"s3"`
	Mongo                    Mongo                                 `json:"mongodb"`
	Amqp                     Amqp                                  `json:"amqp"`
	Database                 string                                `json:"database"`
	Url                      string                                `json:"url"`
	Cert                     string                                `json:"cert"`
	CertFile                 string                                `json:"certFile"`
	Password                 string                                `json:"password"`
	Client                   string                                `json:"client"`
	Timeout                  string                                `json:"timeout"`
}

var Config *Configuration

func Load() error {
	configFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)
	if err != nil {
		return err
	}
	return nil
}

func EditConfig() {
	flagEdit := flag.Bool("config", false, "a bool")
	flag.Parse()
	if *flagEdit {
		enterConfig()
	}
}

func enterConfig() {
	open, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			marshal, err := json.MarshalIndent(Config, "", "  ")
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
		fmt.Printf("%s", "Budget loading? (docker false):")
		line, err := readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err := strconv.ParseBool(line)
		checkErr(err)
		Config.Budget.Loading = loading
		if Config.Budget.Loading {
			fmt.Printf("%s", "Enter Budget cron(docker 0 0/46 * * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0/46 * * *"
			}
			Config.Budget.Cron = line
		}
		fmt.Printf("%s", "Enter Budget sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		Config.Budget.SellingDirection = line
		/*ProgramBreaks*/
		fmt.Printf("%s", "ProgramBreaks loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.ProgramBreaks.Loading = loading
		if Config.ProgramBreaks.Loading {
			fmt.Printf("%s", "Enter ProgramBreaks cron(docker 0 0 0/8 * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0 0/8 * *"
			}
			Config.ProgramBreaks.Cron = line
		}
		fmt.Printf("%s", "Enter ProgramBreaks sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		Config.ProgramBreaks.SellingDirection = line
		/*ProgramBreaksLight*/
		fmt.Printf("%s", "ProgramBreaksLight loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.ProgramBreaksLight.Loading = loading
		if Config.ProgramBreaksLight.Loading {
			fmt.Printf("%s", "Enter ProgramBreaksLight cron(docker 0/2 * * * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0/2 * * * *"
			}
			Config.ProgramBreaksLight.Cron = line
		}
		fmt.Printf("%s", "Enter ProgramBreaksLight sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		Config.ProgramBreaksLight.SellingDirection = line
		/*Mediaplan*/
		Config.Mediaplan.SellingDirection = line
		fmt.Printf("%s", "Mediaplan loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.Mediaplan.Loading = loading
		if Config.Mediaplan.Loading {
			fmt.Printf("%s", "Enter Mediaplan cron(docker 0 0/20 * * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0/20 * * *"
			}
			Config.Mediaplan.Cron = line
		}
		fmt.Printf("%s", "Enter Mediaplan sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		Config.Mediaplan.SellingDirection = line
		/*AdvMessages*/
		fmt.Printf("%s", "AdvMessages loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.AdvMessages.Loading = loading
		if Config.AdvMessages.Loading {
			fmt.Printf("%s", "Enter AdvMessages cron(docker 0 0/2 * * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0/2 * * *"
			}
			Config.AdvMessages.Cron = line
		}
		/*Spots*/
		fmt.Printf("%s", "Spots loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.Spots.Loading = loading
		if Config.Spots.Loading {
			fmt.Printf("%s", "Enter Spots cron(docker 0 0/59 * * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0/59 * * *"
			}
			Config.Spots.Cron = line
		}
		fmt.Printf("%s", "Enter Spots sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		Config.Spots.SellingDirection = line
		/*DeletedSpotInfo*/
		fmt.Printf("%s", "DeletedSpotInfo loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.DeletedSpotInfo.Loading = loading
		if Config.DeletedSpotInfo.Loading {
			fmt.Printf("%s", "Enter DeletedSpotInfo cron(docker 0 0 0/12 * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0 0/12 * *"
			}
			Config.DeletedSpotInfo.Cron = line
		}
		/*Channel*/
		fmt.Printf("%s", "Channels loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.Channel.Loading = loading
		if Config.Channel.Loading {
			fmt.Printf("%s", "Enter Channels cron(docker 0 0/18 * * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0/18 * * *"
			}
			Config.Channel.Cron = line
		}
		fmt.Printf("%s", "Enter Channels sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		Config.Channel.SellingDirection = line
		/*CustomersWithAdvertisers*/
		fmt.Printf("%s", "CustomersWithAdvertisers loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.CustomersWithAdvertisers.Loading = loading
		if Config.CustomersWithAdvertisers.Loading {
			fmt.Printf("%s", "Enter CustomersWithAdvertisers cron(docker 0 0/16 * * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0/16 * * *"
			}
			Config.CustomersWithAdvertisers.Cron = line
		}
		fmt.Printf("%s", "Enter CustomersWithAdvertisers sellingDirection(docker 23):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "23"
		}
		Config.CustomersWithAdvertisers.SellingDirection = line
		/*Rank*/
		fmt.Printf("%s", "Rank loading? (docker false):")
		line, err = readLine()
		checkErr(err)
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		checkErr(err)
		Config.Rank.Loading = loading
		if Config.Rank.Loading {
			fmt.Printf("%s", "Enter Rank cron(docker 0 0 0/23 * *):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "0 0 0/23 * *"
			}
			Config.Rank.Cron = line
		}
		/*AMQP*/
		fmt.Printf("%s", "Enter amqp host(docker localhost):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "localhost"
		}
		Config.Amqp.MqHost = line
		fmt.Printf("%s", "Enter amqp port(docker 5555):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "5555"
		}
		Config.Amqp.MqPort = line
		fmt.Printf("%s", "Enter amqp username(docker guest):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "guest"
		}
		Config.Amqp.MqUsername = line
		fmt.Printf("%s", "Enter amqp password(docker guest):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "guest"
		}
		Config.Amqp.MqPassword = line
		/*S3*/
		fmt.Printf("%s", "Enter S3 AccessKeyId(docker minioadmin):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "minioadmin"
		}
		Config.S3.S3AccessKeyId = line
		fmt.Printf("%s", "Enter S3 SecretAccessKey(docker minioadmin):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "minioadmin"
		}
		Config.S3.S3SecretAccessKey = line
		fmt.Printf("%s", "Enter S3 Region(docker us-west-0):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "us-west-0"
		}
		Config.S3.S3Region = line
		fmt.Printf("%s", "Enter S3 Endpoint(docker 127.0.0.1:9999):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "127.0.0.1:9999"
		}
		Config.S3.S3Endpoint = line
		fmt.Printf("%s", "Enter S3 Debug(docker true):")
		line, err = readLine()
		checkErr(err)
		if line != "false" {
			line = "true"
		}
		Config.S3.S3Debug = line
		fmt.Printf("%s", "Enter S3 Bucket(docker storage):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "storage"
		}
		Config.S3.S3Bucket = line
		fmt.Printf("%s", "Enter S3 LocalDir(docker s3-buckets):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "s3-buckets"
		}
		Config.S3.S3LocalDir = line
		/*Choose database*/
		fmt.Printf("%s", "Enter database, mongodb or badger(default badger):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "badger"
		}
		Config.Database = line
		if Config.Database == "mongodb" {
			/*Mongodb*/
			fmt.Printf("%s", "Enter MongoDB Host(docker localhost):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "localhost"
			}
			Config.Mongo.Host = line
			fmt.Printf("%s", "Enter MongoDB Port(docker 27017):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "27017"
			}
			Config.Mongo.Port = line
			fmt.Printf("%s", "Enter MongoDB DB(docker db):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "db"
			}
			Config.Mongo.DB = line
			fmt.Printf("%s", "Enter MongoDB Username(docker root):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "root"
			}
			Config.Mongo.Username = line
			fmt.Printf("%s", "Enter MongoDB Password(docker qwerty):")
			line, err = readLine()
			checkErr(err)
			if line == "" {
				line = "qwerty"
			}
			Config.Mongo.Password = line
		}
		/*VIMB*/
		fmt.Printf("%s", "Enter url(docker https://vimb-svc.vitpc.com:436/VIMBService.asmx):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "https://vimb-svc.vitpc.com:436/VIMBService.asmx"
		}
		Config.Url = line
		fmt.Printf("%s", "Enter certificate format. 1 - cert file, 2 - cert base64?:")
		line, err = readLine()
		checkErr(err)
		if line == "1" {
			fmt.Printf("%s", "Enter cert file:")
			line, err = readLine()
			checkErr(err)
			Config.CertFile = line
		} else {
			fmt.Printf("%s", "Enter base64 cert:")
			line, err = readLine()
			checkErr(err)
			Config.Cert = line
		}
		fmt.Printf("%s", "Enter password:")
		line, err = readLine()
		checkErr(err)
		fmt.Printf("%s", "Enter client(docker test_client):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "test_client"
		}
		Config.Client = line
		fmt.Printf("%s", "Enter timeout(docker 120s):")
		line, err = readLine()
		checkErr(err)
		if line == "" {
			line = "120s"
		}
		Config.Timeout = line
		marshal, err := json.MarshalIndent(Config, "", "  ")
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
