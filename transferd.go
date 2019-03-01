package main

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type transferd struct {
	logPrefix     string
	queues        *queues
	messages      []message
	mutex         sync.Mutex
	awsRegion     string
	s3Bucket      string
	s3KeyPrefix   string
	maxLines      int
	flushInterval int
}

func NewTransferd(qs *queues) *transferd {
	td := new(transferd)
	td.logPrefix = "[transferd] "
	td.queues = qs

	// default parameters
	td.maxLines = 100
	td.flushInterval = 60

	return td
}

func (td *transferd) run() {
	go td.receive()
	go td.timer()
}

func (td *transferd) receive() {
	for mes := range *td.queues.transfer {
		td.mutex.Lock()
		td.messages = append(td.messages, mes)
		td.mutex.Unlock()

		if td.maxLines > len(td.messages) {
			continue
		}

		go td.flush()
	}
}

func (td *transferd) timer() {
	for {
		mySleep(td.flushInterval)
		if 1 > len(td.messages) {
			continue
		}
		go td.flush()
	}
}

func (td *transferd) flush() error {
	td.mutex.Lock()
	defer td.mutex.Unlock()

	chunk := []string{}
	for _, mes := range td.messages {
		chunk = append(chunk, mes.content)
	}

	reader, writer := io.Pipe()
	var md5sum string
	hash := md5.New()
	go func() {
		gw := gzip.NewWriter(writer)
		buffer := []byte(strings.Join(chunk, "\n"))
		gw.Write(buffer)
		md5sum = base64.StdEncoding.EncodeToString(hash.Sum(buffer))
		gw.Close()
		writer.Close()
	}()

	n := time.Now()
	tp := n.Format("/2006/01/02/")
	tf := n.Format("20060102_150405")
	r := randString(10)
	s3Key := path.Clean(td.s3KeyPrefix + tp + "/slapd_access_log_" + tf + "_" + r + ".jsonl.gz")

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(td.awsRegion)}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:        reader,
		Bucket:      aws.String(td.s3Bucket),
		Key:         aws.String(s3Key),
		ContentType: aws.String("application/gzip"),
		ContentMD5:  aws.String(md5sum),
	})

	if err != nil {
		myLoggerInfo(td.logPrefix + "Failed to upload file: " + err.Error())
		return err
	}
	myLoggerInfo(td.logPrefix + "Succeeded to upload file to: " + result.Location)

	for _, mes := range td.messages {
		os.Remove(mes.filepath)
	}
	td.messages = []message{}
	return nil
}
