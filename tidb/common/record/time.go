package record

import (
	"log"
	"time"
)

func RecordTime(taskName string, task func()) {
	start := time.Now().UnixNano()
	log.Printf("start to %s", taskName)
	task()
	log.Printf("finish to %s", taskName)
	t := time.Now().UnixNano() - start
	log.Printf("spend %d nanoseconds or %d microseconds to execute %s", t, t/1000000, taskName)
}
