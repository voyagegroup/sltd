package main

import (
	"io/ioutil"
	"log"
	"testing"
)

type testCase struct {
	input  []string
	expect []string
}

func testSetup() {
	log.SetOutput(ioutil.Discard)
}

func TestToJsonl(t *testing.T) {
	testSetup()

	testCases := []testCase{
		{[]string{"k1: v1\nk2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"k2\":[\"v2\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\nk1: v2\n"}, []string{"{\"k1\":[\"v1\",\"v2\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1:: djEgd2l0aCBzcGVjaWFsIGNoYXJz\n"}, []string{"{\"k1\":[\"v1 with special chars\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n#k2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n #k2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n-\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
	}

	for _, tc := range testCases {
		callToJsonl(t, tc)
	}
}

func callToJsonl(t *testing.T, tc testCase) {
	c := NewChunk()
	c.serverHostName = "localhost"
	c.ldif.text = tc.input[0]
	actual := c.toJsonl()

	if actual != tc.expect[0] {
		t.Errorf("\nactual: %v\nexpect: %v", actual, tc.expect[0])
	}
}

func TestParseLine(t *testing.T) {
	testSetup()

	testCases := []testCase{
		{[]string{"k: v"}, []string{"k", "v", ""}},
		{[]string{"k:"}, []string{"k", "", ""}},
		{[]string{""}, []string{"", "", "unexpected line format. line="}},
		{[]string{"-"}, []string{"", "", "unexpected line format. line=-"}},
		{[]string{"k:: djEgd2l0aCBzcGVjaWFsIGNoYXJz"}, []string{"k", "v1 with special chars", ""}},
	}

	for _, tc := range testCases {
		callParseLine(t, tc)
	}
}

func callParseLine(t *testing.T, tc testCase) {
	c := NewChunk()
	actual_key, actual_val, actual_err := c.parseLine(tc.input[0])

	if actual_key != tc.expect[0] {
		t.Errorf("\nactual: %v\nexpect: %v", actual_key, tc.expect[0])
	}

	if actual_val != tc.expect[1] {
		t.Errorf("\nactual: %v\nexpect: %v", actual_val, tc.expect[1])
	}

	if actual_err == nil && tc.expect[2] != "" {
		t.Errorf("\nactual: %v\nexpect: no error.", tc.expect[2])
	}

	if actual_err != nil && actual_err.Error() != tc.expect[2] {
		t.Errorf("\nactual: %v\nexpect: %v", actual_err.Error(), tc.expect[2])
	}
}
