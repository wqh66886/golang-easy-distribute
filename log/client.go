package log

import (
	"bytes"
	"fmt"
	"github.com/wqh/easy/distribute/registry"
	stlog "log"
	"net/http"
)

func SetClientLogger(serviceUrl string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] -", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{
		url: serviceUrl,
	})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	body := bytes.NewBuffer(data)
	resp, err := http.Post(cl.url+"/log", "text/plain", body)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("response body closed failed: %v\n", err)
		}
	}()
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log message. service returned %v\n", resp.Status)
	}
	return len(data), nil
}
