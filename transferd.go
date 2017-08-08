package main

import (
	"compress/gzip"
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
	logPrefix         string
	queues            *queues
	awsRegion         string
	s3Bucket          string
	s3KeyPrefix       string
	slapdAccesslogDir string
	messages          []message
	mutex             sync.Mutex
}

func NewTransferd(qs *queues) *transferd {
	td := new(transferd)
	td.logPrefix = "[transferd] "
	td.queues = qs
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

		if 10 > len(td.messages) {
			continue
		}

		go td.flush()
	}
}

func (td *transferd) timer() {
	for {
		mySleep(10)
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
	go func() {
		gw := gzip.NewWriter(writer)
		gw.Write([]byte(strings.Join(chunk, "\n")))
		gw.Close()
		writer.Close()
	}()

	n := time.Now()
	tp := n.Format("/2006/01/02/")
	tf := n.Format("20060102_030405")
	r := randString(10)
	s3Key := path.Clean(td.s3KeyPrefix + tp + "/slapd_access_log_" + tf + "_" + r + ".jsonl.gz")

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(td.awsRegion)}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:        reader,
		Bucket:      aws.String(td.s3Bucket),
		Key:         aws.String(s3Key),
		ContentType: aws.String("application/gzip"),
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
