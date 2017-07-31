package main

import (
	"testing"
)

type testCase struct {
	input  string
	expect string
}

func TestToJsonl(t *testing.T) {
	testCases := []testCase{
		{"k1: v1\nk2: v2\n", "{\"k1\":[\"v1\"],\"k2\":[\"v2\"],\"serverHostName\":[\"localhost\"]}"},
		{"k1: v1\nk1: v2\n", "{\"k1\":[\"v1\",\"v2\"],\"serverHostName\":[\"localhost\"]}"},
		{"k1:: djEgd2l0aCBzcGVjaWFsIGNoYXJz\n", "{\"k1\":[\"v1 with special chars\"],\"serverHostName\":[\"localhost\"]}"},
		{"k1: v1\n#k2: v2\n", "{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"},
		{"k1: v1\n #k2: v2\n", "{\"k1\":[\"v1\"],\"serverHostName\":[\"localhost\"]}"},
	}

	for _, tc := range testCases {
		callToJsonl(t, tc)
	}
}

func callToJsonl(t *testing.T, tc testCase) {
	c := NewChunk()
	c.serverHostName = "localhost"
	c.ldif.text = tc.input
	actual := c.toJsonl()

	if actual != tc.expect {
		t.Errorf("\nexpect: %v\nactual: %v", tc.expect, actual)
	}
}
