package main

import (
	"io/ioutil"
	"os"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

type watcherd struct {
	logPrefix         string
	slapdAccesslogDir string
	queues            *queues
}

func NewWatcherd(qs *queues) *watcherd {
	wd := new(watcherd)
	wd.logPrefix = "[watcherd] "
	wd.queues = qs

	// default parameters
	wd.slapdAccesslogDir = "/var/log/slapd/cn=accesslog/"

	return wd
}

func (wd *watcherd) run() {
	fs := wd.getFileLists(wd.slapdAccesslogDir)
	for i := range fs {
		targetFilePath := wd.slapdAccesslogDir + fs[i].Name()
		mes, err := wd.makeMessage(targetFilePath)
		if err != nil {
			continue
		}
		*wd.queues.parser <- mes
	}

	go wd.daemonize()
}

var reFile = regexp.MustCompile(`\.ldif$`)

func (wd *watcherd) daemonize() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		myLoggerInfo(wd.logPrefix + err.Error())
	}
	defer watcher.Close()

	err = watcher.Add(wd.slapdAccesslogDir)
	if err != nil {
		myLoggerInfo(wd.logPrefix + err.Error())
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op.String() != "CREATE" || !reFile.MatchString(event.Name) {
				continue
			}

			mes, err := wd.makeMessage(event.Name)
			if err != nil {
				continue
			}
			*wd.queues.parser <- mes
		case err := <-watcher.Errors:
			myLoggerInfo(wd.logPrefix + "Error: " + err.Error())
		}
	}
}

func (wd *watcherd) makeMessage(filepath string) (message, error) {
	myLoggerInfo(wd.logPrefix + "New file found: " + filepath)

	file, err := os.Open(filepath)
	defer file.Close()

	mes := new(message)

	if err != nil {
		myLoggerInfo(wd.logPrefix + "Failed to open file: " + err.Error())
		return *mes, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		myLoggerInfo(wd.logPrefix + "Failed to read file: " + err.Error())
		return *mes, err
	}

	mes.filepath = filepath
	mes.content = string(content)

	return *mes, nil
}

func (wd *watcherd) getFileLists(path string) []os.FileInfo {
	fileLists, err := ioutil.ReadDir(path)
	if err != nil {
		myLoggerInfo(wd.logPrefix + "Directory cannot read: " + err.Error())
		return []os.FileInfo{}
	}

	nf := []os.FileInfo{}
	for _, f := range fileLists {
		if f.IsDir() || !reFile.MatchString(f.Name()) {
			continue
		}
		nf = append(nf, f)
	}
	return nf
}
