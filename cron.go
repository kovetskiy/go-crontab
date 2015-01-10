package cron

import (
	"errors"
	"regexp"
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
		Task func()
	}

	Jobs []Job
)

var (
	reSchedule = regexp.MustCompile(`(\*\/[0-9]{1,2})|(\*)|([0-9]{1,2})`)
)

func NewJob(schedule string, task func()) (*Job, error) {
	//six because should trigger error when params count not equals 5
	matches := reSchedule.FindAllString(schedule, 6)
	if len(matches) != 5 {
		return nil, errors.New(
			`Schedule should be specified as %min %hour %dom %mon %dow`)
	}

	job := &Job{
		Min: matches[0],
		Hour: matches[1],
		Dom: matches[2],
		Mon: matches[3],
		Dow: matches[4],
		Task: task,
	}

	return job, nil
}

func (jobs *Jobs) Process() {
	for {
		now := time.Now()
		for _, job := range *jobs {
			if job.IsMatchTime(now) {
				go job.Task()
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
