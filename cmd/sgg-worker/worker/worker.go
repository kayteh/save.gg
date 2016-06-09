// Library for managing work queues
package worker

import (
	//"encoding/json"
	"errors"
)

var (
	ErrNoWorker = errors.New("sgg/worker: job has no handler")
)

var wrs workerRoutes

type workerRoutes map[string]workerRoute

func (w *workerRoutes) Exec(job *Job) error {
	fn, ok := wrs[job.Name]

	if !ok {
		return ErrNoWorker
	}

	fn(job)
	return nil
}

type workerRoute func(job *Job)

type Job struct {
	Name string
	Data map[string]interface{}
}

func (j *Job) Failed(err error) {
	// handle a failure (log, measure, etc)
}

func Enqueue(job *Job) {

}
