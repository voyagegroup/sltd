package main

import (
	"testing"
)

func TestToJsonl(t *testing.T) {
	testSetup()

	testCases := []testCase{
		{[]string{"k1: v1\nk2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"k2\":[\"v2\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\nk1: v2\n"}, []string{"{\"k1\":[\"v1\",\"v2\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1:: djEgd2l0aCBzcGVjaWFsIGNoYXJz\n"}, []string{"{\"k1\":[\"v1 with special chars\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n continue line"}, []string{"{\"k1\":[\"v1continue line\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n  continue line"}, []string{"{\"k1\":[\"v1 continue line\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\ninvalid line"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n#k2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n #k2: v2\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
		{[]string{"k1: v1\n-\n"}, []string{"{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"}},
	}

	for _, tc := range testCases {
		callToJsonl(t, tc)
	}
}

func callToJsonl(t *testing.T, tc testCase) {
	pd := new(parserd)
	pd.serverHostName = "localhost"
	actual := pd.toJsonl(tc.input[0])

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
	pd := new(parserd)
	actual_key, actual_val, actual_err := pd.parseLine(tc.input[0])

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

func TestDecodeIfbase64(t *testing.T) {
	testSetup()

	testCases := []testCase{
		{[]string{": YmFzZTY0"}, []string{"base64", ""}},
		{[]string{":  YmFzZTY0"}, []string{"base64", ""}},
		{[]string{":  i!!3ga! bytes"}, []string{"", "illegal base64 data at input byte 1"}},
		{[]string{"rawstr"}, []string{"rawstr", ""}},
	}

	for _, tc := range testCases {
		callDecodeIfbase64(t, tc)
	}
}

func callDecodeIfbase64(t *testing.T, tc testCase) {
	pd := new(parserd)
	actual_val, actual_err := pd.decodeIfbase64(tc.input[0])

	if actual_val != tc.expect[0] {
		t.Errorf("\nactual: %v\nexpect: %v", actual_val, tc.expect[0])
	}

	if actual_err == nil && tc.expect[1] != "" {
		t.Errorf("\nactual: %v\nexpect: no error.", tc.expect[1])
	}

	if actual_err != nil && actual_err.Error() != tc.expect[1] {
		t.Errorf("\nactual: %v\nexpect: %v", actual_err.Error(), tc.expect[1])
	}
}
