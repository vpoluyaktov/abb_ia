package utils

import (
	"reflect"
	"time"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

type JobDispatcher struct {
	workers []worker
	jobs    []job
}

type worker struct {
	id   int
	busy bool
	job  *job
}

type Fn interface{}
type Params interface{}
type job struct {
	id       int
	jobFn    Fn
	params   []interface{}
	assigned bool
	complete bool
}

func (w *worker) run(j *job) {
	w.job = j
	w.busy = true
	j.assigned = true
	f := reflect.ValueOf(j.jobFn)
	if len(j.params) == f.Type().NumIn() {
		in := make([]reflect.Value, len(j.params))
		for k, param := range j.params {
			in[k] = reflect.ValueOf(param)
		}
		f.Call(in)
	} else {
		logger.Error("JobDispatcher Error: Wrong number of parameters for " + f.String() + " function")
	}
	j.complete = true
	w.busy = false
}

func NewJobDispatcher(workers int) *JobDispatcher {
	jd := &JobDispatcher{make([]worker, workers), make([]job, 0)}
	for i := range jd.workers {
		w := &jd.workers[i]
		w.id = i
		w.busy = false
	}
	return jd
}

func (d *JobDispatcher) AddJob(id int, jobFn interface{}, params ...interface{}) {
	t := reflect.TypeOf(jobFn)
	if t.Kind() == reflect.Func {
		j := job{
			id:       id,
			jobFn:    jobFn,
			params:   params,
			complete: false,
		}
		d.jobs = append(d.jobs, j)
	}
}

func (d *JobDispatcher) Start() {
	// assign jobs to workers
	for i := range d.jobs {
		j := &d.jobs[i]
		if !j.assigned {
			// find/wait for free worker
			var freeWorker *worker
			for freeWorker == nil {
				for ii := range d.workers {
					w := &d.workers[ii]
					if !w.busy {
						freeWorker = w
						break
					}
				}
				time.Sleep(200 * time.Microsecond)
			}
			freeWorker.busy = true
			j.assigned = true
			go freeWorker.run(j)
		}
	}

	// wait for all jobs to complete
	for {
		isWorking := false
		for _, j := range d.jobs {
			if !j.complete {
				isWorking = true
			}
		}
		if !isWorking {
			break
		}
		time.Sleep(200 * time.Microsecond)
	}
}
