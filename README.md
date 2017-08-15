# sltd (**s**lapd **l**og **t**ransfer **d**aemon)

sltd transfers slapd [accesslog](http://www.openldap.org/doc/admin24/overlays.html#Access%20Logging) to Amazon S3.

```
$ sltd
2017/08/01 18:27:41 [Info] : sltd initialzing ...
2017/08/01 18:27:41 [Info] : SLTD_LOG:
2017/08/01 18:27:41 [Info] : SLTD_LOG_LEVEL: debug
2017/08/01 18:27:41 [Info] : SLAPD_ACCESSLOG_DIR: /var/log/slapd/cn=accesslog/
2017/08/01 18:27:41 [Info] : AWS_REGION: ap-northeast-1
2017/08/01 18:27:41 [Info] : S3_BUCKET: BUCKET_NAME
2017/08/01 18:27:41 [Info] : S3_KEY_PREFIX: slapd_access_log/
2017/08/01 18:27:42 [Info] : [watcherd] New file found: /var/log/slapd/cn=accesslog/reqStart=20170801092742\2E000001Z.ldif
.. snip ..
2017/08/01 13:28:42 [Info] : [watcherd] New file found: /var/log/slapd/cn=accesslog/reqStart=20170801092842\2E000005Z.ldif
2017/08/01 13:28:42 [Info] : [transferd] Succeeded to upload file to: https://BUCKET_NAME.s3-ap-northeast-1.amazonaws.com/slapd_access_log/2017/08/01/slapd_access_log_20170801_132842_y7uukhxdba.jsonl.gz
```

## Index

* [Concepts](#concepts)
* [Requirements](#requirements)
* [Installation](#installation)
* [Configure](#configure)
* [Usage](#usage)
* [License](#license)

## Concepts

* For auditing, store slapd access log to Amazon S3.
* Use Access Logging Overlay and LDIF Backend as log source.

## Requirements

sltd requires the following to run:

* Golang

## Installation

```
$ go get github.com/voyagegroup/sltd
```

or

Download from Releases Page (WIP)

## Usage

### set slapd to enable accesslog

```
# logging target database section.
database    mdb

.. snip ..
moduleload  accesslog
overlay     accesslog
logdb       cn=accesslog
logops      all
logsuccess  FALSE
logpurge    03:00:00 00:30:00
.. snip ..


# accesslog database section.
database    ldif
directory   /var/log/slapd/
suffix      cn=accesslog
rootdn      cn=XROOTDNX
rootpw      {SSHA}XROOTDNPWXXXXXXXXXXXXXXXXXXXXXXX
```

### configure sltd

Set your configuration as Environment Variables.
```
# require
AWS_REGION="XXX"
S3_BUCKET="XXX"

# optional
AWS_ACCESS_KEY_ID="XXX"
AWS_SECRET_ACCESS_KEY="XXX"
S3_KEY_PREFIX="XXX"
MAX_LINES="XXX"
SLTD_LOG_LEVEL="XXX"
```
You can use .env file as well.

### run

```
$ sltd
```

## License

[MIT](./LICENSE)

## Author

[Satoshi Tajima](https://github.com/s-tajima)
