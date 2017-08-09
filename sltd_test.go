package main

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestInitialize(t *testing.T) {
	testSetup()

	testCases := []testCase{
		{[]string{".env.tests/.env.test1"}, []string{"/var/log/slapd/cn=accesslog/", "100", "60"}},
		{[]string{".env.tests/.env.test2"}, []string{"/tmp/slapd/cn=accesslog/", "2017", "89"}},
	}

	for _, tc := range testCases {
		callInitialize(t, tc)
	}
}

func callInitialize(t *testing.T, tc testCase) {
	godotenv.Load(tc.input[0])

	s := new(sltd)
	s.initialize()

	actual := fmt.Sprintf("slapdAccesslogDir:%s,maxLines:%d,flushInterval:%d", s.wd.slapdAccesslogDir, s.td.maxLines, s.td.flushInterval)
	expect := fmt.Sprintf("slapdAccesslogDir:%s,maxLines:%s,flushInterval:%s", tc.expect[0], tc.expect[1], tc.expect[2])
	if actual != expect {
		t.Errorf("\nactual: %v\nexpect: %v", actual, expect)
	}
}
