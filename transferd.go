package main

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type transferd struct {
	awsRegion         string
	s3Bucket          string
	s3KeyPrefix       string
	slapdAccesslogDir string
}

func (td *transferd) daemonize() {
	for {
		td.transfer()
	}
}

func (td *transferd) transfer() {
	awsSleepTime := 1
	fileSleepTime := 5

	fileLists := td.getFileLists(td.slapdAccesslogDir)
	if len(fileLists) == 0 {
		mySleep(fileSleepTime)
		return
	}

	for i := range fileLists {
		targetFilePath := td.slapdAccesslogDir + fileLists[i].Name()

		myLoggerDebug(targetFilePath)

		file, err := os.Open(targetFilePath)
		if err != nil {
			myLoggerInfo("Failed to open targetFilePath: " + err.Error())
			mySleep(fileSleepTime)
			return
		}

		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			myLoggerInfo("Failed to read file: " + err.Error())
			mySleep(fileSleepTime)
			return
		}
		file.Close()

		c := NewChunk()
		c.ldif.text = string(fileContent)

		myLoggerDebug(c.toJsonl())

		reader, writer := io.Pipe()
		go func() {
			gw := gzip.NewWriter(writer)
			gw.Write([]byte(c.toJsonl()))
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
			myLoggerInfo("Failed to upload file: " + err.Error())
			time.Sleep(time.Duration(awsSleepTime) * time.Second)
			awsSleepTime = awsSleepTime * 2
			return
		}
		myLoggerInfo("Succeeded to upload file to: " + result.Location)

		awsSleepTime = 1
		os.Remove(targetFilePath)
	}
}

func (td *transferd) getFileLists(path string) []os.FileInfo {
	fileLists, err := ioutil.ReadDir(path)
	if err != nil {
		myLoggerInfo("Directory cannot read: " + err.Error())
		return []os.FileInfo{}
	}

	r := regexp.MustCompile(`^\.ldif$`)
	nf := []os.FileInfo{}
	for _, f := range fileLists {
		if f.IsDir() || r.MatchString(f.Name()) {
			continue
		}
		nf = append(nf, f)
	}
	return nf
}
