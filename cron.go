package cron

import (
	"strconv"
	"time"
)

type (
	Job struct {
		Min  string
		Hour string
		Dom  string
		Mon  string
		Dow  string
		Run  func()
	}

	Jobs []Job
)

func (jobs *Jobs) Process() {
	for {
		now := time.Now()
		for _, job := range *jobs {
			if job.IsMatchTime(now) {
				go job.Run()
			}
		}

		time.Sleep(time.Second)
	}
}

func (job *Job) IsMatchTime(t time.Time) (result bool) {
	if t.Second() != 0 {
		return false
	}

	result = compare(int64(t.Minute()), job.Min) &&
		compare(int64(t.Hour()), job.Hour) &&
		compare(int64(t.Day()), job.Dom) &&
		compare(int64(t.Month()), job.Mon) &&
		compare(int64(t.Weekday()), job.Dow)

	return
}

func compare(origin int64, scheduled string) bool {
	if scheduled == "*" || scheduled == "" {
		return true
	}

	if numeric, err := strconv.ParseInt(scheduled, 10, 64); err == nil {
		return numeric == origin
	}

	if scheduled[:2] != "*/" {
		return false
	}

	divider, err := strconv.ParseInt(scheduled[2:], 10, 64)
	if err != nil || divider == 0 {
		return false
	}

	return origin%divider == 0
}
