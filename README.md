# sltd (**s**lapd **l**og **t**ransfer **d**aemon) 

sltd transfers slapd [accesslog](http://www.openldap.org/doc/admin24/overlays.html#Access%20Logging) to Amazon S3.

```
$ sltd
2017/08/01 18:27:41 [Info] : sltd initialzing ...
2017/08/01 18:27:41 [Info] : SLTD_LOG:
2017/08/01 18:27:41 [Info] : SLTD_LOG_LEVEL: debug
2017/08/01 18:27:41 [Info] : SLAPD_ACCESSLOG_DIR: /var/log/slapd/cn=accesslog/
2017/08/01 18:27:41 [Info] : AWS_REGION: ap-northeast-1
2017/08/01 18:27:41 [Info] : S3_BUCKET: ACTUAL_BUCKET_NAME
2017/08/01 18:27:41 [Info] : S3_KEY_PREFIX: slapd_access_log/
```

## Index

* [Concepts](#concepts)
* [Requirements](#requirements)
* [Installation](#installation)
* [Configure](#configure)
* [Usage](#usage)       
* [License](#license)    

## Concepts

* write this later.

## Requirements

sltd requires the following to run:

* Golang

## Installation

```
$ write this later.
```

## Usage

### set slapd to enable accesslog

write this later

### configure sltd

Set your configuration as Environment Variables.
```
write this later
```
You can use .env file as well.

## License

[MIT](./LICENSE)

## Author

[Satoshi Tajima](https://github.com/s-tajima)
