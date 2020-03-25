package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"zorm/constant"
)

var (
	errLog = log.New(os.Stdout,constant.ERR_COLOR,log.LstdFlags|log.Lshortfile)
	infoLog = log.New(os.Stdout,constant.INOF_COLOR,log.LstdFlags|log.Lshortfile)
	logs = []*log.Logger{errLog,infoLog}
	mu sync.Mutex
)
// log func
var (
	Error = errLog.Println
	Errorf = errLog.Printf
	Info = infoLog.Println
	Infof = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrLevel
	Disabled
)

func SetLevel(level int)  {
	mu.Lock()
	defer mu.Unlock()

	for _, log := range logs {
		log.SetOutput(os.Stdout)
	}

	if level > ErrLevel {
		errLog.SetOutput(ioutil.Discard)
	}
	if level > InfoLevel {
		infoLog.SetOutput(ioutil.Discard)
	}
}