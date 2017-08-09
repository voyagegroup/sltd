package main

import (
	"io/ioutil"
	"log"
)

type testCase struct {
	input  []string
	expect []string
}

func testSetup() {
	log.SetOutput(ioutil.Discard)
}
