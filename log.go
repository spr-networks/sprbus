package sprbus

import (
	"github.com/sirupsen/logrus"
	//"github.com/spr-networks/sprbus"
	"io/ioutil"
	"os"
	"strings"
)

// test logging

type Hook struct {
	client    *Client
	prefix    string
	LogLevels []logrus.Level
}

func NewSprBusHook(prefix string) (*Hook, error) {
	client, err := NewClient(ServerEventSock)

	return &Hook{client, prefix, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}}, err
}

// Fire will be called when some logging function is called with current hook
func (hook *Hook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}

	_, err = hook.client.Publish(hook.prefix, strings.Trim(string(line), "\n"))
	//write to stderr if we cant dial
	if err != nil {
		//transport: Error while dialing: dial unix
		//os.Stderr.WriteString("[error] "+err.Error()+"\n")
		os.Stderr.Write(line)
		return nil
	}

	return err
}

// logrus.Levels define on which log logrus.Levels this hook would trigger
func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}

// shortcut to get a logrus instance with hook
func NewLog(prefix string) *logrus.Logger {
	if prefix == "" {
		prefix = "log"
	}

	log := logrus.New()
	log.SetOutput(ioutil.Discard)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	//log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetReportCaller(true)

	h, _ := NewSprBusHook(prefix)
	log.Hooks.Add(h)
	//h, err

	return log
}
