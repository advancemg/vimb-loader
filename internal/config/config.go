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
	Host       string `json:"Host"`
	Port       string `json:"Port"`
	AuthDB     string `json:"AuthDb"`
	DB         string `json:"Db"`
	Username   string `json:"Username"`
	Password   string `json:"Password"`
	CronBackup string `json:"CronBackup"`
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
	Token                    string                                `json:"token"`
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

func EditConfig() error {
	flagEdit := flag.Bool("config", false, "a bool")
	flag.Parse()
	if *flagEdit {
		err := enterConfig()
		if err != nil {
			return fmt.Errorf("EditConfig(): %w", err)
		}
	}
	return nil
}

func enterConfig() error {
	open, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			marshal, err := json.MarshalIndent(Config, "", "  ")
			if err != nil {
				return err
			}
			err = os.WriteFile("config.json", marshal, 0666)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	fmt.Println("Default Config:")
	fmt.Println(string(open))
	fmt.Printf("%s", "Edit config? (Y/n):")
	cfg := &Configuration{}
	n, err := readLine()
	if err != nil {
		return err
	}
	if n == "y" || n == "Y" || n == "yes" || n == "Yes" {
		/*Budget*/
		fmt.Printf("%s", "Budget loading? (docker false):")
		line, err := readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err := strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.Budget.Loading = loading
		if cfg.Budget.Loading {
			fmt.Printf("%s", "Enter Budget cron(docker 0/46 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/46 * * * *"
			}
			cfg.Budget.Cron = line
		}
		fmt.Printf("%s", "Enter Budget sellingDirection(docker 23):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "23"
		}
		cfg.Budget.SellingDirection = line
		/*ProgramBreaks*/
		fmt.Printf("%s", "ProgramBreaks loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.ProgramBreaks.Loading = loading
		if cfg.ProgramBreaks.Loading {
			fmt.Printf("%s", "Enter ProgramBreaks cron(docker 0/8 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/8 * * * *"
			}
			cfg.ProgramBreaks.Cron = line
		}
		fmt.Printf("%s", "Enter ProgramBreaks sellingDirection(docker 23):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "23"
		}
		cfg.ProgramBreaks.SellingDirection = line
		/*ProgramBreaksLight*/
		fmt.Printf("%s", "ProgramBreaksLight loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.ProgramBreaksLight.Loading = loading
		if cfg.ProgramBreaksLight.Loading {
			fmt.Printf("%s", "Enter ProgramBreaksLight cron(docker 0/2 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/2 * * * *"
			}
			cfg.ProgramBreaksLight.Cron = line
		}
		fmt.Printf("%s", "Enter ProgramBreaksLight sellingDirection(docker 23):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "23"
		}
		cfg.ProgramBreaksLight.SellingDirection = line
		/*Mediaplan*/
		cfg.Mediaplan.SellingDirection = line
		fmt.Printf("%s", "Mediaplan loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.Mediaplan.Loading = loading
		if cfg.Mediaplan.Loading {
			fmt.Printf("%s", "Enter Mediaplan cron(docker 0/20 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/20 * * * *"
			}
			cfg.Mediaplan.Cron = line
		}
		fmt.Printf("%s", "Enter Mediaplan sellingDirection(docker 23):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "23"
		}
		cfg.Mediaplan.SellingDirection = line
		/*AdvMessages*/
		fmt.Printf("%s", "AdvMessages loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.AdvMessages.Loading = loading
		if cfg.AdvMessages.Loading {
			fmt.Printf("%s", "Enter AdvMessages cron(docker 0/2 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/2 * * * *"
			}
			cfg.AdvMessages.Cron = line
		}
		/*Spots*/
		fmt.Printf("%s", "Spots loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.Spots.Loading = loading
		if cfg.Spots.Loading {
			fmt.Printf("%s", "Enter Spots cron(docker 0/59 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/59 * * * *"
			}
			cfg.Spots.Cron = line
		}
		fmt.Printf("%s", "Enter Spots sellingDirection(docker 23):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "23"
		}
		cfg.Spots.SellingDirection = line
		/*DeletedSpotInfo*/
		fmt.Printf("%s", "DeletedSpotInfo loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.DeletedSpotInfo.Loading = loading
		if cfg.DeletedSpotInfo.Loading {
			fmt.Printf("%s", "Enter DeletedSpotInfo cron(docker 0/12 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/12 * * * *"
			}
			cfg.DeletedSpotInfo.Cron = line
		}
		/*Channel*/
		fmt.Printf("%s", "Channels loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.Channel.Loading = loading
		if cfg.Channel.Loading {
			fmt.Printf("%s", "Enter Channels cron(docker 0/18 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/18 * * * *"
			}
			cfg.Channel.Cron = line
		}
		fmt.Printf("%s", "Enter Channels sellingDirection(docker 23):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "23"
		}
		cfg.Channel.SellingDirection = line
		/*CustomersWithAdvertisers*/
		fmt.Printf("%s", "CustomersWithAdvertisers loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.CustomersWithAdvertisers.Loading = loading
		if cfg.CustomersWithAdvertisers.Loading {
			fmt.Printf("%s", "Enter CustomersWithAdvertisers cron(docker 0/16 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/16 * * * *"
			}
			cfg.CustomersWithAdvertisers.Cron = line
		}
		fmt.Printf("%s", "Enter CustomersWithAdvertisers sellingDirection(docker 23):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "23"
		}
		cfg.CustomersWithAdvertisers.SellingDirection = line
		/*Rank*/
		fmt.Printf("%s", "Rank loading? (docker false):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "true" {
			line = "false"
		}
		loading, err = strconv.ParseBool(line)
		if err != nil {
			return err
		}
		cfg.Rank.Loading = loading
		if cfg.Rank.Loading {
			fmt.Printf("%s", "Enter Rank cron(docker 0/23 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/23 * * * *"
			}
			cfg.Rank.Cron = line
		}
		/*AMQP*/
		fmt.Printf("%s", "Enter amqp host(docker localhost):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "localhost"
		}
		cfg.Amqp.MqHost = line
		fmt.Printf("%s", "Enter amqp port(docker 5555):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "5555"
		}
		cfg.Amqp.MqPort = line
		fmt.Printf("%s", "Enter amqp username(docker guest):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "guest"
		}
		cfg.Amqp.MqUsername = line
		fmt.Printf("%s", "Enter amqp password(docker guest):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "guest"
		}
		cfg.Amqp.MqPassword = line
		/*S3*/
		fmt.Printf("%s", "Enter S3 AccessKeyId(docker minioadmin):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "minioadmin"
		}
		cfg.S3.S3AccessKeyId = line
		fmt.Printf("%s", "Enter S3 SecretAccessKey(docker minioadmin):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "minioadmin"
		}
		cfg.S3.S3SecretAccessKey = line
		fmt.Printf("%s", "Enter S3 Region(docker us-west-0):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "us-west-0"
		}
		cfg.S3.S3Region = line
		fmt.Printf("%s", "Enter S3 Endpoint(docker 127.0.0.1:9999):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "127.0.0.1:9999"
		}
		cfg.S3.S3Endpoint = line
		fmt.Printf("%s", "Enter S3 Debug(docker true):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line != "false" {
			line = "true"
		}
		cfg.S3.S3Debug = line
		fmt.Printf("%s", "Enter S3 Bucket(docker storage):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "storage"
		}
		cfg.S3.S3Bucket = line
		fmt.Printf("%s", "Enter S3 LocalDir(docker s3-buckets):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "s3-buckets"
		}
		cfg.S3.S3LocalDir = line
		/*Choose database*/
		fmt.Printf("%s", "Enter database, mongodb or badger(default badger):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "badger"
		}
		cfg.Database = line
		if cfg.Database == "mongodb" {
			/*Mongodb*/
			fmt.Printf("%s", "Enter MongoDB Host(docker localhost):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "localhost"
			}
			cfg.Mongo.Host = line
			fmt.Printf("%s", "Enter MongoDB Port(docker 27017):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "27017"
			}
			cfg.Mongo.Port = line
			fmt.Printf("%s", "Enter MongoDB DB(docker db):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "db"
			}
			cfg.Mongo.DB = line
			fmt.Printf("%s", "Enter MongoDB Username(docker root):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "root"
			}
			cfg.Mongo.Username = line
			fmt.Printf("%s", "Enter MongoDB Password(docker qwerty):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "qwerty"
			}
			cfg.Mongo.Password = line
			fmt.Printf("%s", "Enter MongoDB Cron Backup(docker 0/24 * * * *):")
			line, err = readLine()
			if err != nil {
				return err
			}
			if line == "" {
				line = "0/24 * * * *"
			}
			cfg.Mongo.CronBackup = line
		}
		/*VIMB*/
		fmt.Printf("%s", "Enter url(docker https://vimb-svc.vitpc.com:436/VIMBService.asmx):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "https://vimb-svc.vitpc.com:436/VIMBService.asmx"
		}
		cfg.Url = line
		fmt.Printf("%s", "Enter certificate format. 1 - cert file, 2 - cert base64?:")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "1" {
			fmt.Printf("%s", "Enter cert file:")
			line, err = readLine()
			if err != nil {
				return err
			}
			cfg.CertFile = line
		} else {
			fmt.Printf("%s", "Enter base64 cert:")
			line, err = readLine()
			if err != nil {
				return err
			}
			cfg.Cert = line
		}
		fmt.Printf("%s", "Enter password:")
		line, err = readLine()
		if err != nil {
			return err
		}
		fmt.Printf("%s", "Enter client(docker test_client):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "test_client"
		}
		cfg.Client = line
		fmt.Printf("%s", "Enter timeout(docker 120s):")
		line, err = readLine()
		if err != nil {
			return err
		}
		if line == "" {
			line = "120s"
		}
		cfg.Timeout = line
		fmt.Printf("%s", "Enter token for API:")
		line, err = readLine()
		if err != nil {
			return err
		}
		cfg.Token = line
		Config = cfg
		marshal, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile("config.json", marshal, 0666)
		if err != nil {
			return err
		}
		fmt.Println(string(marshal))
	}
	return nil
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	result := ""
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return result, fmt.Errorf("readLine(): %w", err)
		}
		result += string(line)
		if !isPrefix {
			return result, nil
		}
	}
}
