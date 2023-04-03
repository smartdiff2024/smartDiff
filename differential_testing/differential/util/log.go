package util

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Myformatter struct {
	Suffix string
}

var MyFormatter Myformatter

func (mf *Myformatter) Format(entry *logrus.Entry) ([]byte, error) {
	MyFormatter = Myformatter{
		Suffix: "\n",
	}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("[%s] %s %s", timestamp, entry.Message, mf.Suffix))
	return b.Bytes(), nil
}
