package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	td := new(transferd)

	sltdLog := os.Getenv("SLTD_LOG")
	sltdLogLevel := os.Getenv("SLTD_LOG_LEVEL")
	td.awsRegion = os.Getenv("AWS_REGION")
	td.s3Bucket = os.Getenv("S3_BUCKET")
	td.s3KeyPrefix = os.Getenv("S3_KEY_PREFIX")
	td.slapdAccesslogDir = os.Getenv("SLAPD_ACCESSLOG_DIR")

	if len(sltdLog) != 0 {
		logfile, err := os.OpenFile(sltdLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			myLoggerInfo("Failed to open sltdLog: " + err.Error())
			os.Exit(1)
		}
		defer logfile.Close()
		log.SetOutput(logfile)
	}

	if len(td.slapdAccesslogDir) == 0 {
		td.slapdAccesslogDir = "/var/log/slapd/cn=accesslog/"
	}

	myLoggerInfo("sltd initialzing ...")
	myLoggerInfo("SLTD_LOG: " + sltdLog)
	myLoggerInfo("SLTD_LOG_LEVEL: " + sltdLogLevel)
	myLoggerInfo("SLAPD_ACCESSLOG_DIR: " + td.slapdAccesslogDir)
	myLoggerInfo("AWS_REGION: " + td.awsRegion)
	myLoggerInfo("S3_BUCKET: " + td.s3Bucket)
	myLoggerInfo("S3_KEY_PREFIX: " + td.s3KeyPrefix)

	td.daemonize()
}
