package glog

import (
	"fmt"
	"time"
)

// flushInterval
const flushInterval = 1 * time.Second

// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 1024

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
func logName(tag string, t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%s.%04d%02d%02d-%02d%02d%02d.%d",
		program,
		tag,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)
	return name, program + "." + tag
}
