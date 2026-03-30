package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Constants representing logging levels.
const (
	Debug   = "Debug"
	Info    = "Info"
	Warning = "Warning"
	Trace   = "Trace"
	Error   = "Error"
)

const (
	logInterval = 5
)

type logEntry struct {
	level   string
	message string
}

// Uses logger if ch is nil.
func log(level, message string, logger *logrus.Logger, ch chan logEntry) {
	if ch != nil {
		ch <- logEntry{level: level, message: message}
	} else {
		switch level {
		case Debug:
			logger.Debug(message)
		case Info:
			logger.Info(message)
		case Warning:
			logger.Warn(message)
		case Trace:
			logger.Trace(message)
		}
	}
}

func printChannelLogs(ip string, ch chan logEntry) {
	for len(ch) > 0 {
		entry := <-ch
		message := fmt.Sprintf("Node %s: %s", ip, entry.message)
		switch entry.level {
		case Debug:
			logrus.Debug(message)
		case Info:
			logrus.Info(message)
		case Warning:
			logrus.Warn(message)
		default:
			logrus.Info(message)
		}
	}
}

func printLogs(wg *sync.WaitGroup, ipChanMap map[string]chan logEntry) {
	defer wg.Done()
	for {
		if len(ipChanMap) == 0 { //nolint: staticcheck
			// no IPs to monitor or all channels are closed, exit loop
			break
		}
		for ip, ch := range ipChanMap {
			if len(ch) == 0 {
				// check if channel is closed
				_, ok := <-ch
				if !ok {
					// channel is closed, remove IP from map to stop checking for logs
					delete(ipChanMap, ip)
				}
			}
			printChannelLogs(ip, ch)
		}
		time.Sleep(logInterval * time.Second)
	}
}
