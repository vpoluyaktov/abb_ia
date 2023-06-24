package utils

import (
	"reflect"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

type JobDispatcher struct {
	workers []worker
	jobs    []job
}

type worker struct {
	id   int
	busy bool
}

type Fn interface{}
type Params interface{}
type job struct {
	id       int
	jobFn     Fn
	params 	 []interface{}
	assigned bool
	complete bool
}

func (w *worker) run(j job) {
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
	for i, w := range jd.workers {
		w.id = i
		w.busy = false
	}
	return jd
}

func (d *JobDispatcher) AddJob(id int, jobFn interface{}, params ...interface{}) {
	t := reflect.TypeOf(jobFn)
	if t.Kind() == reflect.Func {
		j := job {
			id: id,
			jobFn: jobFn, 
			params: params,
			complete: false,
		}
		d.jobs = append(d.jobs, j)
	}
}

func (d *JobDispatcher) Start() {
	for _, w := range d.workers {
		if !w.busy {
			for _, j := range d.jobs {
				if !j.assigned && !j.complete {
					go w.run(j)
				}
			}
		}
	}
}
