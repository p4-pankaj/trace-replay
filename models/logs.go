package models

import "time"

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Raw       string                 `json:"rawLog"`
	File      string                 `json:"file"`
	Line      int                    `json:"line"`
	Level     string                 `json:"level"`
	Fields    map[string]interface{} `json:"fields"`
}
