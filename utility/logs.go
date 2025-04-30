package utility

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"time"

	"github.com/p4-pankaj/trace-replay/models"
)

type InMemoryLogWriter struct {
	Logs   []models.LogEntry
	Writer chan func(r *models.TraceRecord)
}

func (w *InMemoryLogWriter) Write(p []byte) (n int, err error) {
	var file string
	var line int
	var found bool
	i := 1
	for {
		// find the most recent file in runtime that called the log which is not from zerolog pkg
		_, file, line, found = runtime.Caller(i)
		if !found {
			file = "unknown"
			line = 0
			break
		}
		if strings.Contains(file, "zerolog") || strings.Contains(file, "log/zerolog") {
			i += 1
			continue
		}
		break
	}

	var fields map[string]any
	level := "unknown"
	if err := json.Unmarshal(p, &fields); err != nil {
		fields = map[string]any{"message": "invalid"}
	}

	if v, exist := fields["level"]; exist {
		level = v.(string)
	}

	logMsg := string(p)
	logEntry := models.LogEntry{
		Timestamp: time.Now(),
		Raw:       logMsg,
		File:      file,
		Line:      line,
		Level:     level,
		Fields:    fields,
	}
	w.Logs = append(w.Logs, logEntry)
	//todo feed to chan , and sync chan to close

	w.Writer <- func(r *models.TraceRecord) {
		if r.Logs != nil {
			r.Logs = append(r.Logs, &logEntry)
			return
		}
		r.Logs = []*models.LogEntry{&logEntry}
	}
	fmt.Println(string(p))
	return len(p), nil
}
