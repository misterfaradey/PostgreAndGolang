package default_logger

import (
	"log"
	"os"
)

func NewDefaultLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}
