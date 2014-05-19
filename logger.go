package chainstore

import (
	"log"
	"os"
)

var lg = log.New(os.Stderr, "", log.LstdFlags)

func SetLogger(l *log.Logger) {
	lg = l
}
