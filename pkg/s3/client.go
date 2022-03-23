package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/zenthangplus/goccm"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var cfg *Config

func init() {
	cfg = loadConfig()
}

type Config struct {
	S3AccessKeyId     string `json:"s3AccessKeyId"`
	S3SecretAccessKey string `json:"s3SecretAccessKey"`
	S3Region          string `json:"s3Region"`
	S3Endpoint        string `json:"s3Endpoint"`
	S3Debug           string `json:"s3Debug"`
	S3Bucket          string `json:"s3Bucket"`
	S3Session         *session.Session
}

func loadConfig() *Config {
	var config Config
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config
}

type ProgressWriter struct {
	Written int64
	Writer  io.WriterAt
	Size    int64
}

func (pw *ProgressWriter) WriteAt(p []byte, off int64) (int, error) {
	return pw.Writer.WriteAt(p, off)
}

func GetFileSize(svc *s3.S3, bucket string, prefix string) (filesize int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}
	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}
	return *resp.ContentLength, nil
}

func connection() error {
	if cfg.S3Debug == `true` {
		cfg.S3Session = session.Must(session.NewSession(&aws.Config{
			MaxRetries:       aws.Int(3),
			Region:           &cfg.S3Region,
			Endpoint:         &cfg.S3Endpoint,
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
			Credentials:      credentials.NewStaticCredentials(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
		}))
		return nil
	}
	cfg.S3Session = session.Must(session.NewSession(&aws.Config{
		MaxRetries:  aws.Int(3),
		Region:      &cfg.S3Region,
		Endpoint:    &cfg.S3Endpoint,
		Credentials: credentials.NewStaticCredentials(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
	}))
	return nil
}

func CreateBucket(name string) error {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"CreateBucket(s3Client)": "error", "error": err.Error()})
		return err
	}
	client := s3.New(cfg.S3Session)
	request := s3.CreateBucketInput{
		Bucket: aws.String(name),
	}
	_, err = client.CreateBucket(&request)
	if err != nil {
		log.Print(map[string]string{"CreateBucket(s3Client)": "error", "error": err.Error()})
		return err
	}
	return nil
}
func CheckFile(bucket, key string) (map[string]string, error) {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"CreateBucket(s3Client)": "error", "error": err.Error()})
		return nil, err
	}
	client := s3.New(cfg.S3Session)
	request := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := client.HeadObject(&request)
	if err != nil {
		log.Print(map[string]string{"CheckFile(s3Client)": "error", "error": err.Error()})
		return nil, err
	}
	state := map[string]string{
		"ContentLength": strconv.FormatInt(aws.Int64Value(result.ContentLength), 10),
		"ETag":          aws.StringValue(result.ETag),
	}
	return state, nil
}
func Exist(bucket, key string) bool {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"Exist(s3Client)": "error", "error": err.Error()})
		return false
	}
	client := s3.New(cfg.S3Session)
	request := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := client.HeadObject(&request)
	if err != nil {
		log.Print(map[string]string{"Exist(s3Client)": "error", "error": err.Error()})
		return false
	}
	b := nil != result
	return b
}
func ListDirectories(bucket, prefix string) map[string]string {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"ListDirectories(s3Client)": "error", "error": err.Error()})
		return nil
	}
	client := s3.New(cfg.S3Session)
	request := s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(100),
	}
	var data = map[string]string{}
	var result *s3.ListObjectsV2Output
	var errData error
	for {
		result, errData = client.ListObjectsV2(&request)
		if errData != nil {
			log.Print(map[string]string{"ListDirectories(s3Client)": "error", "error": err.Error()})
			return nil
		}
		for _, file := range result.Contents {
			fileResult := *file.Key
			data[path.Dir(fileResult)] = path.Dir(fileResult)
		}
		request = s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(prefix),
			MaxKeys:           aws.Int64(100),
			ContinuationToken: result.NextContinuationToken,
		}
		if !*result.IsTruncated {
			break
		}
	}
	return data
}
func ListKeys(bucket, prefix string) map[string]string {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"ListKeys(s3Client)": "error", "error": err.Error()})
		return nil
	}
	client := s3.New(cfg.S3Session)
	request := s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}
	result, err := client.ListObjects(&request)
	if err != nil {
		log.Print(map[string]string{"ListKeys(s3Client)": "error", "error": err.Error()})
		return nil
	}
	var data = map[string]string{}
	for _, file := range result.Contents {
		data[*file.Key] = *file.ETag
	}
	return data
}
func ListKeysWithCred(bucket, prefix string) map[string]string {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"ListKeysWithCred(s3Client)": "error", "error": err.Error()})
		return nil
	}
	client := s3.New(cfg.S3Session)
	request := s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}
	result, err := client.ListObjects(&request)
	if err != nil {
		log.Print(map[string]string{"ListKeysWithCred(s3Client)": "error", "error": err.Error()})
		return nil
	}
	var data = map[string]string{}
	for _, file := range result.Contents {
		data[*file.Key] = *file.ETag
	}
	return data
}
func CopyBatch(bucket, inputPrefix, outputPrefix string) error {
	keys := ListKeys(bucket, inputPrefix)
	err := connection()
	if err != nil {
		log.Print(map[string]string{"CopyBatch(s3Client)": "error", "error": err.Error()})
		return err
	}
	s3Client := s3.New(cfg.S3Session)
	for s3Key, _ := range keys {
		outputS3Key := strings.ReplaceAll(s3Key, inputPrefix, outputPrefix)
		_, err = s3Client.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			CopySource: aws.String(fmt.Sprintf("%s/%s", bucket, s3Key)),
			Key:        aws.String(outputS3Key),
		})
		if err != nil {
			log.Print(map[string]string{"CopyBatch(s3Client)": "error", "error": err.Error()})
			return err
		}
	}
	return nil
}
func DownloadBatch(bucket, prefix string) (string, error) {
	keys := ListKeys(bucket, prefix)
	err := connection()
	if err != nil {
		log.Print(map[string]string{"DownloadBatch(s3Client)": "error", "error": err.Error()})
		return "", err
	}
	s3Client := s3.New(cfg.S3Session)
	sessionDataDir, err := ioutil.TempDir(``, `s3_client-dir-Session-`)
	if err != nil {
		log.Print(map[string]string{"DownloadBatch(s3Client)": "error", "error": err.Error()})
		return "", err
	}
	errorCh := make(chan error)
	successCh := make(chan interface{})
	lenKeys := len(keys)
	concurrency := goccm.New(lenKeys)
	for s3Key, _ := range keys {
		concurrency.Wait()
		go func(inputS3Key string) {
			defer concurrency.Done()
			size, err := GetFileSize(s3Client, bucket, inputS3Key)
			if err != nil {
				concurrency.Done()
				errorCh <- err
				return
			}
			temp, err := ioutil.TempFile(sessionDataDir, "s3_client-load-file-tmp-")
			if err != nil {
				concurrency.Done()
				errorCh <- err
				return
			}
			writer := &ProgressWriter{Writer: temp, Size: size, Written: 0}
			svc := s3manager.NewDownloader(cfg.S3Session, func(d *s3manager.Downloader) {
				d.PartSize = 5 * 1024 * 1024
				d.Concurrency = lenKeys
			})
			_, err = svc.Download(writer, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(inputS3Key),
			})
			if err != nil {
				os.Remove(temp.Name())
				concurrency.Done()
				errorCh <- err
				return
			}
			concurrency.Done()
			successCh <- inputS3Key
		}(s3Key)
	}
	concurrency.WaitAllDone()
	for {
		select {
		case <-time.After(time.Second * 1):
			return sessionDataDir, nil
		case errC := <-errorCh:
			return "", errC
		}
	}
}
func Download(bucket, key string) (string, error) {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"Download(s3Client)": "error", "error": err.Error()})
		return "", err
	}
	s3Client := s3.New(cfg.S3Session)
	size, err := GetFileSize(s3Client, bucket, key)
	if err != nil {
		log.Print(map[string]string{"Download(s3Client)": "error", "error": err.Error()})
		return ``, err
	}
	temp, err := ioutil.TempFile(``, "s3_client-load-file-tmp-")
	if err != nil {
		panic(err)
	}
	tempfileName := temp.Name()
	writer := &ProgressWriter{Writer: temp, Size: size, Written: 0}
	svc := s3manager.NewDownloader(cfg.S3Session, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024
		d.Concurrency = 10
	})
	_, err = svc.Download(writer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Print(map[string]string{"Download(s3Client)": "error", "error": err.Error()})
		os.Remove(tempfileName)
		return ``, err
	}
	return tempfileName, nil
}
func UploadFileWithBucket(filePathInput, s3Key string) (*s3manager.UploadOutput, error) {
	file, err := os.Open(filePathInput)
	if err != nil {
		log.Print(map[string]string{"UploadFileWithBucket(s3Client)": "error", "error": err.Error()})
		file.Close()
	}
	err = connection()
	if err != nil {
		log.Print(map[string]string{"UploadFileWithBucket(s3Client)": "error", "error": err.Error()})
		return nil, err
	}
	defer file.Close()
	svc := s3manager.NewUploader(cfg.S3Session)
	return svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(s3Key),
		Body:   file,
	})
}

func UploadBytesWithBucket(s3Key string, data []byte) (*s3manager.UploadOutput, error) {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"UploadFileWithBucket(s3Client)": "error", "error": err.Error()})
		return nil, err
	}
	svc := s3manager.NewUploader(cfg.S3Session)
	return svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(data),
	})
}

func DeleteWithBucketPrefix(bucket string, prefix string) error {
	keys := ListKeys(bucket, prefix)
	err := connection()
	if err != nil {
		log.Print(map[string]string{"DeleteWithBucketPrefix(s3Client)": "error", "error": err.Error()})
		return err
	}
	svc := s3manager.NewBatchDelete(cfg.S3Session, func(d *s3manager.BatchDelete) {
		d.BatchSize = len(keys)
	})
	var objects = []s3manager.BatchDeleteObject{}
	for key, _ := range keys {
		objects = append(objects, s3manager.BatchDeleteObject{
			Object: &s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			},
		})
	}
	err = svc.Delete(aws.BackgroundContext(), &s3manager.DeleteObjectsIterator{
		Objects: objects,
	})
	if err != nil {
		log.Print(map[string]string{"DeleteWithBucketPrefix(s3Client)": "error", "error": err.Error()})
		return err
	}
	return nil
}
func DeleteWithBucket(bucket string, s3Keys []string) error {
	err := connection()
	if err != nil {
		log.Print(map[string]string{"DeleteWithBucket(s3Client)": "error", "error": err.Error()})
		return err
	}
	svc := s3manager.NewBatchDelete(cfg.S3Session, func(d *s3manager.BatchDelete) {
		d.BatchSize = len(s3Keys)
	})
	var objects = []s3manager.BatchDeleteObject{}
	for _, key := range s3Keys {
		objects = append(objects, s3manager.BatchDeleteObject{
			Object: &s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			},
		})
	}
	err = svc.Delete(aws.BackgroundContext(), &s3manager.DeleteObjectsIterator{
		Objects: objects,
	})
	if err != nil {
		log.Print(map[string]string{"DeleteWithBucket(s3Client)": "error", "error": err.Error()})
		return err
	}
	return nil
}