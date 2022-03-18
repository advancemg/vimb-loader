package soap

import "C"
import (
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"fmt"
	convert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/utils"
	"github.com/buger/jsonparser"
	"golang.org/x/crypto/pkcs12"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var cfg *Config

func init() {
	cfg = loadConfig()
}

type Config struct {
	Url      string `json:"url"`
	Cert     string `json:"cert"`
	Password string `json:"password"`
}

func loadConfig() *Config {
	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config
}

func (cfg *Config) newClient() *http.Client {
	timeout := 30 * time.Second
	dataCert, err := utils.Base64DecodeString(cfg.Cert)
	if err != nil {
		log.Fatal("error:", err)
	}
	blocks, err := pkcs12.ToPEM(dataCert, cfg.Password)
	if err != nil {
		log.Fatal("error:", err)
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		log.Fatalf("err while converting to key pair: %v", err)
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
		Renegotiation:      tls.RenegotiateOnceAsClient,
		Certificates:       []tls.Certificate{cert},
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
		Timeout: timeout,
	}
}

func Request(input []byte) ([]byte, error) {
	reqBody := vimbRequest(string(input))
	req, err := http.NewRequest("POST", cfg.Url, strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "VIMBWebApplication2/GetVimbInfoStream")
	resp, err := cfg.newClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		vimbErr, err := catchError(res)
		if err != nil {
			return nil, err
		}
		fmt.Println("code:", vimbErr.Code)
		return nil, vimbErr
	}
	return res, nil
}

func vimbRequest(inputXml string) string {
	return fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:vim="VIMBWebApplication2"><soapenv:Header/><soapenv:Body><vim:GetVimbInfoStream><vim:InputXML><![CDATA[%s]]></vim:InputXML></vim:GetVimbInfoStream></soapenv:Body></soapenv:Envelope>`, inputXml)
}

func catchError(resp []byte) (*VimbError, error) {
	toJson, err := convert.XmlToJson(resp)
	if err != nil {
		return nil, err
	}
	statusCode, err := jsonparser.GetString(toJson, "Envelope", "Body", "Fault", "detail", "ErrorDescription", "Code")
	if err != nil {
		return nil, err
	}
	code, err := strconv.Atoi(statusCode)
	if err != nil {
		return nil, err
	}
	msg, err := jsonparser.GetString(toJson, "Envelope", "Body", "Fault", "detail", "ErrorDescription", "Message")
	if err != nil {
		return nil, err
	}
	return &VimbError{
		Code:    code,
		Message: msg,
	}, nil
}

type VimbError struct {
	Code    int
	Message string
}

func (err *VimbError) Error() string {
	return fmt.Sprintf("%s", err.Message)
}
