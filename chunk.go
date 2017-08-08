package main

import (
	//	"encoding/base64"
	//	"encoding/json"
	//	"errors"
	"os"
	//	"regexp"
	//	"strings"
)

type chunk struct {
	serverHostName string
	ldif           ldif
}

type ldif struct {
	text string
}

func NewChunk() *chunk {
	c := chunk{}
	c.serverHostName, _ = os.Hostname()
	return &c
}
