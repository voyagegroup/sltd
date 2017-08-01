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

var reComment = regexp.MustCompile(`^\s*#`)

func NewChunk() *chunk {
	c := chunk{}
	c.serverHostName, _ = os.Hostname()
	return &c
}

func (c *chunk) toJsonl() string {
	es := strings.Split(strings.TrimSpace(c.ldif.text), "\n\n")

	jsonSource := make([]string, len(es))
	for i, e := range es {
		dls := strings.Split(e, "\n")

		ej := make(map[string][]string)
		for _, dl := range dls {

			if reComment.MatchString(dl) {
				continue
			}

			key, val, err := c.parseLine(dl)
			if err != nil {
				myLoggerInfo(err.Error())
				continue
			}
			ej[key] = append(ej[key], val)
		}

		ej["serverHostName"] = append(ej["serverHostName"], c.serverHostName)

		s, _ := json.Marshal(ej)
		jsonSource[i] = string(s)
	}

	return strings.Join(jsonSource, "\n")
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
