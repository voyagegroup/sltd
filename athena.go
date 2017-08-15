package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

type athenaDDL struct {
	location string
}

func (ad *athenaDDL) initialize() {
	s3Bucket := os.Getenv("S3_BUCKET")
	s3KeyPrefix := os.Getenv("S3_KEY_PREFIX")
	ad.location = "s3://" + path.Clean(s3Bucket+"/"+s3KeyPrefix) + "/"
}

func (ad *athenaDDL) printCreateTable() {
	ddl := `CREATE EXTERNAL TABLE IF NOT EXISTS slapd_accesslog_day (
		serverHostName array<string>,
		createTimestamp array<string>,
		creatorsName array<string>,
		dn array<string>,
		entryCSN array<string>,
		entryUUID array<string>,
		modifiersName array<string>,
		modifyTimestamp array<string>,
		objectClass array<string>,
		reqAssertion array<string>,
		reqAttr array<string>,
		reqAttrsOnly array<string>,
		reqAuthzID array<string>,
		reqControls array<string>,
		reqDN array<string>,
		reqData array<string>,
		reqDeleteOldRDN array<string>,
		reqDerefAliases array<string>,
		reqEnd array<string>,
		reqEntries array<string>,
		reqFilter array<string>,
		reqId array<string>,
		reqMessage array<string>,
		reqMethod array<string>,
		reqMod array<string>,
		reqNewRDN array<string>,
		reqNewSuperior array<string>,
		reqOld array<string>,
		reqReferral array<string>,
		reqRespControls array<string>,
		reqResult array<string>,
		reqScope array<string>,
		reqSession array<string>,
		reqSizeLimit array<string>,
		reqStart array<string>,
		reqTimeLimit array<string>,
		reqType array<string>,
		reqVersion array<string>,
		structuralObjectClass array<string>
	)
	PARTITIONED BY ( "year" int, "month" int, "day" int )
	ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
	LOCATION '%s'
`

	fmt.Println("## CREATE TABLE DDL....")
	fmt.Printf(ddl, ad.location)
	fmt.Println()
}

func (ad *athenaDDL) printAddPartition() {
	ddl := "ALTER TABLE slapd_accesslog_day ADD"

	template := `
	PARTITION (year=%d, month=%d, day=%d)
	LOCATION '%s%04d/%02d/%02d/'`

	fmt.Println("## ADD PARTITION DDL....")
	fmt.Printf(ddl)
	start, _ := time.Parse("2006-1-2", strconv.Itoa(time.Now().Year())+"-1-1")
	for d := start; d.Year() == start.Year(); d = d.AddDate(0, 0, 1) {
		fmt.Printf(template, d.Year(), d.Month(), d.Day(), ad.location, d.Year(), d.Month(), d.Day())
	}
	fmt.Println()
}
