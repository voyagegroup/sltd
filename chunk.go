package main

import (
	"encoding/base64"
	"encoding/json"
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
var reBase64 = regexp.MustCompile(`^:`)

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

			lines := strings.SplitN(dl, ":", 2)
			key := strings.TrimSpace(lines[0])
			val := strings.TrimSpace(lines[1])

			if reBase64.MatchString(val) {
				val = c.decode(val)
			}

			ej[key] = append(ej[key], val)
		}

		ej["serverHostName"] = append(ej["serverHostName"], c.serverHostName)

		s, _ := json.Marshal(ej)
		jsonSource[i] = string(s)
	}

	return strings.Join(jsonSource, "\n")
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
