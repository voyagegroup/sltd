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
var reBase64 = regexp.MustCompile(`^:`)

func (c *chunk) parseLine(line string) (key string, val string, err error) {
	if !reColon.MatchString(line) {
		return "", "", errors.New("invalid format: '" + line + "'")
	}

	parts := strings.SplitN(line, ":", 2)
	key = strings.TrimSpace(parts[0])
	val = strings.TrimSpace(parts[1])

	if reBase64.MatchString(val) {
		val = c.decode(val)
	}
	return key, val, nil
}

func (c *chunk) decode(val string) string {
	val = strings.TrimPrefix(val, ":")
	val = strings.TrimSpace(val)

	p, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		myLoggerInfo("base64 decode failed. err=" + err.Error())
	}
	return string(p)
}
