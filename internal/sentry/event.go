package sentry

import "time"


type LogEvent struct {
	Line string `json:"line"`
	Inode uint64 `json:"inode"`
	Timestamp time.Time `json:"time_stamp"`
}

type FileState struct {
	Path string `json:"path"`
	Inode uint64 `json:"inode"`
	Offset int64 `json:"offset"`
}