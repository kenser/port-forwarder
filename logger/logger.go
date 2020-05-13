package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

var Entry = logrus.NewEntry(logger)

var Error = Entry.Error
var Errorf = Entry.Errorf

var Info = Entry.Info
var Infof = Entry.Infof

var Print = Entry.Info
var Printf = Entry.Infof

var Debug = Entry.Debug
var Debugf = Entry.Debugf

var Panicf = Entry.Panicf
var Panic = Entry.Panic

var Trace = Entry.Trace
var Tracef = Entry.Tracef

var Warn = Entry.Warn
var Warnf = Entry.Warnf

var Fatal = Entry.Fatal
var Fatalf = Entry.Fatalf

func updateEntry() {
	Error = Entry.Error
	Errorf = Entry.Errorf

	Info = Entry.Info
	Infof = Entry.Infof

	Print = Entry.Info
	Printf = Entry.Infof

	Debug = Entry.Debug
	Debugf = Entry.Debugf

	Panicf = Entry.Panicf
	Panic = Entry.Panic

	Trace = Entry.Trace
	Tracef = Entry.Tracef

	Warn = Entry.Warn
	Warnf = Entry.Warnf

	Fatal = Entry.Fatal
	Fatalf = Entry.Fatalf
}

func initField() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	fields := logrus.Fields{
		"hostname": hostname,
	}
	Entry = logger.WithFields(fields)
	//logger.Info(fields)
}

func init() {
	//initField()
	//updateEntry()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true,
	})
	//logger.SetReportCaller(true)

	updateEntry()
	return
}
