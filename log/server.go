package log

import (
	"io"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type fileLog string

// fileLog 实现了io.Writer接口
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, nil
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("file closing failed, err:%v\n", err)
		}
	}()
	return f.Write(data)
}

func Run(destination string) {
	log = stlog.New(fileLog(destination), "[go easy distribute] ", stlog.LstdFlags)
}

func RegisterHandleFunc() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
