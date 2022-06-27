package utils

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	convert "github.com/advancemg/go-convert"
	cfg "github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	"github.com/buger/jsonparser"
	"golang.org/x/crypto/pkcs12"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//type Config struct {
//	Url      string `json:"url"`
//	Cert     string `json:"cert"`
//	CertFile string `json:"certFile"`
//	Password string `json:"password"`
//	Client   string `json:"client"`
//	Timeout  string `json:"timeout"`
//}

type Action struct {
	SOAPAction string
	Client     string
}

var Actions *Action

func init() {
	Actions = &Action{
		SOAPAction: "VIMBWebApplication2/GetVimbInfoStream",
		Client:     "vimb",
	}
}

func newClient() *http.Client {
	cfg := cfg.Config
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		panic(err)
	}
	var dataCert []byte
	if cfg.Cert != "" {
		dataCert, err = base64.StdEncoding.DecodeString(cfg.Cert)
		if err != nil {
			log.PrintLog("vimb-loader", "soap-client", "error", "base64.StdEncoding error", err.Error())
		}
	}
	if cfg.CertFile != "" {
		dataCert, err = FileToBase64(cfg.CertFile)
		if err != nil {
			log.PrintLog("vimb-loader", "soap-client", "error", "base64.StdEncoding error", err.Error())
		}
	}
	blocks, err := pkcs12.ToPEM(dataCert, cfg.Password)
	if err != nil {
		log.PrintLog("vimb-loader", "soap-client", "error", "pkcs12.ToPEM error", err.Error())
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		log.PrintLog("vimb-loader", "soap-client", "error", "err while converting to key pair: ", err.Error())
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

func (act *Action) Request(input []byte) ([]byte, error) {
	cfg := cfg.Config
	reqBody := vimbRequest(string(input))
	req, err := http.NewRequest("POST", cfg.Url, strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", act.SOAPAction)
	resp, err := newClient().Do(req)
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
		return nil, vimbErr
	}
	response, err := catchBody(res)
	if err != nil {
		return nil, err
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(string(response))
	if err != nil {
		return nil, err
	}
	return decodeBytes, nil
}

func (act *Action) RequestJson(input []byte) ([]byte, error) {
	res, err := act.Request(input)
	if err != nil {
		return nil, err
	}
	toJson, err := convert.ZipXmlToJson(res)
	if err != nil {
		return nil, err
	}
	return toJson, nil
}

func vimbRequest(inputXml string) string {
	input := strings.ReplaceAll(inputXml, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	//fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:vim="VIMBWebApplication2"><soapenv:Header/><soapenv:Body><vim:GetVimbInfoStream><vim:InputXML><![CDATA[%s]]></vim:InputXML></vim:GetVimbInfoStream></soapenv:Body></soapenv:Envelope>`, inputXml)
	return fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:vim="VIMBWebApplication2"><soapenv:Header/><soapenv:Body><vim:GetVimbInfoStream><vim:InputXML>%s</vim:InputXML></vim:GetVimbInfoStream></soapenv:Body></soapenv:Envelope>`, input)
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

func catchBody(resp []byte) ([]byte, error) {
	toJson, err := convert.XmlToJson(resp)
	if err != nil {
		return nil, err
	}
	msg, err := jsonparser.GetString(toJson, "Envelope", "Body", "GetVimbInfoStreamResponse", "GetVimbInfoStreamResult")
	if err != nil {
		return nil, err
	}
	return []byte(msg), nil
}

type VimbError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *VimbError) CheckTimeout(method string) {
	code := e.Code
	switch code {
	case 1001:
		wait(method, code, e.Message, time.Minute*5)
		return
	case 1003:
		wait(method, code, e.Message, time.Minute*5)
		return
	default:
		wait(method, code, e.Message, time.Minute*5)
		return
	}
}

type Timeout struct {
	IsTimeout bool `json:"id" bson:"_id"`
}

func wait(method string, code int, msg string, waitTime time.Duration) {
	log.PrintLog("vimb-loader", "soap-client", "error", method, " ", "timeout code:", code, " ", msg)
	db := store.OpenDb("db", "timeout")
	err := db.AddWithTTL("_id", Timeout{IsTimeout: true}, waitTime)
	if err != nil {
		panic(err)
	}
	time.Sleep(waitTime)
}

func (e VimbError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
