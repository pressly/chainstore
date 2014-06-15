package logmgr

import (
	"fmt"
	"log"
)

type logManager struct {
	logger *log.Logger
	tag    string
}

func New(logger *log.Logger, tag string) *logManager {
	if tag != "" {
		tag = fmt.Sprintf(" [%s]", tag)
	}
	return &logManager{logger, tag}
}

func (m *logManager) Open() (err error)  { return }
func (m *logManager) Close() (err error) { return }

func (m *logManager) Put(key string, value []byte) (err error) {
	m.logger.Printf("chainstore%s: put %s of %d bytes", m.tag, key, len(value))
	return
}

func (m *logManager) Get(key string) (value []byte, err error) {
	m.logger.Printf("chainstore%s: get %s", m.tag, key)
	return
}

func (m *logManager) Del(key string) (err error) {
	m.logger.Printf("chainstore%s: del %s", m.tag, key)
	return
}
