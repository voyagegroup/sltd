package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
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

var reNewLine = regexp.MustCompile(`\r?\n`)
var reComment = regexp.MustCompile(`^\s*#`)
var reContinuation = regexp.MustCompile(`^ `)

func (c *chunk) toJsonl() string {
	myLoggerDebug(reNewLine.ReplaceAllString(c.ldif.text, " "))
	es := strings.Split(strings.TrimSpace(c.ldif.text), "\n\n")

	jsonSource := make([]string, len(es))
	for i, e := range es {
		dls := strings.Split(e, "\n")

		ej := make(map[string][]string)
		prevKey := ""
		for _, dl := range dls {
			if reComment.MatchString(dl) {
				continue
			}

			if reContinuation.MatchString(dl) {
				ej[prevKey][len(ej[prevKey])-1] = ej[prevKey][len(ej[prevKey])-1] + strings.TrimPrefix(dl, " ")
				continue
			}

			key, val, err := c.parseLine(dl)
			if err != nil {
				myLoggerInfo(err.Error())
				continue
			}
			ej[key] = append(ej[key], val)
			prevKey = key
		}

		ej["serverHostName"] = append(ej["serverHostName"], c.serverHostName)

		s, _ := json.Marshal(ej)
		jsonSource[i] = string(s)
	}

	jsonStr := strings.Join(jsonSource, "\n")
	myLoggerDebug(jsonStr)

	return jsonStr
}

var reColon = regexp.MustCompile(`:`)

func (c *chunk) parseLine(line string) (key string, val string, err error) {
	if !reColon.MatchString(line) {
		return "", "", errors.New("unexpected line format. line=" + line)
	}

	parts := strings.SplitN(line, ":", 2)
	key = strings.TrimSpace(parts[0])
	val = strings.TrimSpace(parts[1])

	val, err = c.decodeIfbase64(val)
	if err != nil {
		return "", "", errors.New("base64 decode failed. val=" + val + ", err=" + err.Error())
	}

	return key, val, nil
}

var reBase64 = regexp.MustCompile(`^:`)

func (c *chunk) decodeIfbase64(val string) (decoded string, err error) {
	if !reBase64.MatchString(val) {
		return val, nil
	}

	val = strings.TrimPrefix(val, ":")
	val = strings.TrimSpace(val)

	d, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
