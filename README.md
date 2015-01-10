# What is it?

It's yet cron package written in Go.

# Why not other?

I was looking for pretty good `*go*cron*` package, but all of these packages
does not work with cron format "*/10" and does not allows to use groups of jobs.

### All `*go*cron*` packages sucks, this package sucks less.

# Usage
go-crontab allows to use group of jobs, just create instance `Jobs`.

```
jobs := cron.Jobs{}
```

Actually, `Jobs` it's just `[]Job`, so you can append other jobs with `append`.

```
job := cron.NewJob("* */2 * * *", func() {
        //stuff
    },
)

jobs = append(jobs, *job)
```

To run the schedule you should write:
```
go jobs.Process()
```
