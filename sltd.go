package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type sltd struct {
	wd           *watcherd
	pd           *parserd
	td           *transferd
	sltdLog      string
	sltdLogLevel string
}

func main() {
	godotenv.Load()

	s := new(sltd)
	s.initialize()
	s.run()

}

func (s *sltd) initialize() {

	s.sltdLog = os.Getenv("SLTD_LOG")
	s.sltdLogLevel = os.Getenv("SLTD_LOG_LEVEL")

	if len(s.sltdLog) != 0 {
		logfile, err := os.OpenFile(s.sltdLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			myLoggerInfo("Failed to open sltdLog: " + err.Error())
			os.Exit(1)
		}
		defer logfile.Close()
		log.SetOutput(logfile)
	}

	qs := NewQueues()
	s.wd = NewWatcherd(qs)
	s.pd = NewParserd(qs)
	s.td = NewTransferd(qs)

	s.td.awsRegion = os.Getenv("AWS_REGION")
	s.td.s3Bucket = os.Getenv("S3_BUCKET")
	s.td.s3KeyPrefix = os.Getenv("S3_KEY_PREFIX")

	if len(os.Getenv("MAX_LINES")) != 0 {
		s.td.maxLines, _ = strconv.Atoi(os.Getenv("MAX_LINES"))
	}

	if len(os.Getenv("FLUSH_INTERVAL")) != 0 {
		s.td.flushInterval, _ = strconv.Atoi(os.Getenv("FLUSH_INTERVAL"))
	}

	if len(os.Getenv("SLAPD_ACCESSLOG_DIR")) != 0 {
		s.wd.slapdAccesslogDir = os.Getenv("SLAPD_ACCESSLOG_DIR")
	}
}

func (s *sltd) run() {
	myLoggerInfo("sltd initialzing ...")
	myLoggerInfo("SLTD_LOG: " + s.sltdLog)
	myLoggerInfo("SLTD_LOG_LEVEL: " + s.sltdLogLevel)
	myLoggerInfo("SLAPD_ACCESSLOG_DIR: " + s.wd.slapdAccesslogDir)
	myLoggerInfo("AWS_REGION: " + s.td.awsRegion)
	myLoggerInfo("S3_BUCKET: " + s.td.s3Bucket)
	myLoggerInfo("S3_KEY_PREFIX: " + s.td.s3KeyPrefix)

	done := make(chan bool)
	s.td.run()
	s.pd.run()
	s.wd.run()
	<-done
}
