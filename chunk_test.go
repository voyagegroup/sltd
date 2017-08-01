package main

import (
	"testing"
)

type testCase struct {
	input  []string
	expect []string
}

func TestToJsonl(t *testing.T) {
	testCases := []testCase{
		{[]string{"k1: v1\nk2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"k2\":[\"v2\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\nk1: v2\n"}, []string{"{\"k1\":[\"v1\",\"v2\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1:: djEgd2l0aCBzcGVjaWFsIGNoYXJz\n"}, []string{"{\"k1\":[\"v1 with special chars\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n#k2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n #k2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n -\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
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
		t.Errorf("\nexpect: %v\nactual: %v", tc.expect[0], actual)
	}
}

func TestParseLine(t *testing.T) {
	testCases := []testCase{
		{[]string{"k: v"}, []string{"k", "v", ""}},
		{[]string{"k:"}, []string{"k", "", ""}},
		{[]string{""}, []string{"", "", "invalid format: ''"}},
		{[]string{"-"}, []string{"", "", "invalid format: '-'"}},
		{[]string{"k:: djEgd2l0aCBzcGVjaWFsIGNoYXJz"}, []string{"k", "v1 with special chars", ""}},
	}

	for _, tc := range testCases {
		callParseLine(t, tc)
	}
}

func callParseLine(t *testing.T, tc testCase) {
	c := NewChunk()
	actual_key, actual_val, _ := c.parseLine(tc.input[0])

	if actual_key != tc.expect[0] {
		t.Errorf("\nexpect: %v\nactual: %v", tc.expect[0], actual_key)
	}

	if actual_val != tc.expect[1] {
		t.Errorf("\nexpect: %v\nactual: %v", tc.expect[1], actual_val)
	}
}
