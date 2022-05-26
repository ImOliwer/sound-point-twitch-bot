package scheduler

import "time"

type internal_task struct {
	cancel   chan bool
	finished bool
}

type RepeatingTask struct {
	internal_task
	ticker *time.Ticker
}

type LaterTask struct {
	internal_task
	timer *time.Timer
}

func Every(duration time.Duration, handle func(*RepeatingTask)) *RepeatingTask {
	task := &RepeatingTask{
		internal_task: internal_task{
			cancel: make(chan bool),
		},
		ticker: time.NewTicker(duration),
	}
	go func() {
		for {
			select {
			case <-task.ticker.C:
				handle(task)
			case <-task.cancel:
				return
			}
		}
	}()
	return task
}

func After(duration time.Duration, handle func(*LaterTask)) *LaterTask {
	task := &LaterTask{
		internal_task: internal_task{
			cancel: make(chan bool),
		},
		timer: time.NewTimer(duration),
	}
	go func() {
		select {
		case <-task.timer.C:
			handle(task)
			task.cancel_internally(nil)
		case <-task.cancel:
			return
		}
	}()
	return task
}

func (r *internal_task) Wait() {
	<-r.cancel
}

func (r *internal_task) cancel_internally(before func()) {
	if r.finished {
		return
	}

	if before != nil {
		before()
	}

	close(r.cancel)
	r.finished = true
}

func (r *RepeatingTask) Cancel() {
	r.cancel_internally(func() { r.ticker.Stop() })
}

func (r *LaterTask) Cancel() {
	r.cancel_internally(func() { r.timer.Stop() })
}
