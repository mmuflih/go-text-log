package datalog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mmuflih/datelib"
)

type DataLog interface {
	Write(err error, data ...interface{})
}

type dataLog struct {
	fileName string
}

func New(fileName string, daily bool, logPath ...string) DataLog {
	path := ""
	for _, p := range logPath {
		if strings.TrimSpace(p) == "" {
			continue
		}
		path += p + "/"
	}
	if daily {
		os.Mkdir(path+"logs", 0600)
		fileName = path + "logs/" + time.Now().Format(datelib.YMD) + "-" + fileName
	}
	return &dataLog{fileName}
}

func (dl dataLog) Write(e error, data ...interface{}) {
	if data == nil {
		return
	}
	f, err := os.OpenFile(dl.fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		f, err = os.Create(dl.fileName)
	}
	if err != nil {
		fmt.Println("Error while createng file", err)
	}

	defer f.Close()
	d, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err, "marshal")
	}
	text := ""
	if e == nil {
		text = time.Now().Format(datelib.YMD_HMS_WS) + ` : data => ` + string(d)
	} else {
		text = time.Now().Format(datelib.YMD_HMS_WS) + ` : error with exception "` + e.Error() + `" data => ` + string(d)
	}

	if _, err = f.WriteString(text + "\n"); err != nil {
		fmt.Println("Error writing log", err)
	}
}
