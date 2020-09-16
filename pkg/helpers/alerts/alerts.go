package alerts

import (
	"fmt"
	"log"
)

type Alerter interface {
	AlertError(err error, format string, args ...interface{})
}

var alerter Alerter = NewLogAlerter()

func SetGlobalAlerter(a Alerter) {
	alerter = a
}

func AlertError(err error, format string, args ...interface{}) {
	alerter.AlertError(err, format, args...)
}

// -----------------------------------------------------------------------------
type LogAlerter struct{}

func NewLogAlerter() *LogAlerter {
	return &LogAlerter{}
}

func (la *LogAlerter) AlertError(err error, s string, args ...interface{}) {
	log.Printf("ALERT: [%v] %s", err, fmt.Sprintf(s, args...))
}
