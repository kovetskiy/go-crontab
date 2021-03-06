package cron

import (
	"testing"
	"time"
)

func TestCompare(t *testing.T) {
	comparing := []struct {
		Origin    int64
		Scheduled string
		Expected  bool
	}{
		{
			10,
			"10",
			true,
		},
		{
			20,
			"20",
			true,
		},
		{
			0,
			"0",
			true,
		},
		{
			10,
			"*",
			true,
		},
		{
			10,
			"",
			true,
		},
		{
			10,
			"12",
			false,
		},
		{
			10,
			"1a",
			false,
		},
		{
			10,
			"*/10",
			true,
		},
		{
			10,
			"*/20",
			false,
		},
		{
			10,
			"*/0", //heheh
			false,
		},
		{
			10,
			"*/2",
			true,
		},
		{
			10,
			"*/5",
			true,
		},
		{
			10,
			"*/3",
			false,
		},
	}

	for _, c := range comparing {
		got := compare(c.Origin, c.Scheduled)
		if got != c.Expected {
			t.Errorf("For %v and '%v' expected %v got %v",
				c.Origin, c.Scheduled, c.Expected, got)
		}
	}
}

func TestJobIsMatchTime(t *testing.T) {
	form := "2006-01-02 15:04:05"

	var job *Job
	var ts time.Time

	ts, _ = time.Parse(form, "2014-01-10 22:00:00")
	job = &Job{}
	if !job.IsMatchTime(ts) {
		t.Error("Something went wrong")
	}

	ts, _ = time.Parse(form, "2014-01-10 22:00:01")
	job = &Job{}
	if job.IsMatchTime(ts) {
		t.Error("Jobs should runned only in 0 seconds")
	}
}

var TestJobsProcessTick = 0

func TestJobsProcess(t *testing.T) {
	jobs := &Jobs{
		Job{
			Task: func() {
				TestJobsProcessTick++
			},
		},
	}

	go jobs.Process()

	time.Sleep(time.Minute)

	if TestJobsProcessTick != 1 {
		t.Fatalf("Something went wrong, got %d ticks", TestJobsProcessTick)
	}
}

func TestNewJob(t *testing.T) {
	tests := []struct {
		Schedule    string
		ErrExpected bool
	}{
		{
			"1 2 3 4 5",
			false,
		},
		{
			"1 2 3 4 5 6 7",
			true,
		},
		{
			"1 2 3 4 5 6",
			true,
		},
		{
			"1 2 3 4",
			true,
		},
		{
			"10/20 2 3 4 5",
			true,
		},
		{
			"* * 3 * *",
			false,
		},
		{
			"*/10 * 3 */2 *",
			false,
		},
	}

	for _, test := range tests {
		_, err := NewJob(test.Schedule, func() {})
		if err == nil && test.ErrExpected {
			t.Errorf("For %v error must not be nil", test.Schedule)
		}
		if err != nil && !test.ErrExpected {
			t.Errorf("For %v error must be nil", test.Schedule)
		}
	}
}
