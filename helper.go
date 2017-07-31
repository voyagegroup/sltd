package main

import (
	"log"
	"math/rand"
	"os"
	"time"
)

func mySleep(sleepTime int) {
	time.Sleep(time.Duration(sleepTime) * time.Second)
}

// Logging policy refs:http://dave.cheney.net/2015/11/05/lets-talk-about-logging
func myLoggerInfo(message string) {
	log.Println("[Info] : " + message)
}
func myLoggerDebug(message string) {
	if "debug" == os.Getenv("SLTD_LOG_LEVEL") {
		log.Println("[Debug]: " + message)
	}
}

func randString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
