package worker

import (
	"log"
	"stream/ent"
	"stream/internal/config"
	"stream/internal/tools"
)

// Job represents a job to be executed.
type Job struct {
	VideoID   string
	InputPath string
}

// JobQueue is a channel for sending jobs.
var JobQueue chan Job

// Worker represents a worker that executes jobs.
type Worker struct {
	ID         int
	JobQueue   chan Job
	Quit       chan bool
	Client     *ent.Client
	AppConfig  *config.Config
}

// NewWorker creates a new Worker.
func NewWorker(id int, jobQueue chan Job, client *ent.Client, appConfig *config.Config) Worker {
	return Worker{
		ID:         id,
		JobQueue:   jobQueue,
		Quit:       make(chan bool),
		Client:     client,
		AppConfig:  appConfig,
	}
}

// Start starts the worker's job processing loop.
func (w Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.JobQueue:
				log.Printf("Worker %d: received job for video %s", w.ID, job.VideoID)
				w.transcode(job)

			case <-w.Quit:
				log.Printf("Worker %d: stopping", w.ID)
				return
			}
		}
	}()
}

func (w Worker) transcode(job Job) {
	tools.TranscodeAndProcessVideo(job.VideoID, job.InputPath, w.Client, w.AppConfig)
}

// Stop signals the worker to stop.
func (w Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

// Dispatcher for the job queue
func StartDispatcher(nWorkers int, client *ent.Client, appConfig *config.Config) {
	JobQueue = make(chan Job, 100)

	for i := 1; i <= nWorkers; i++ {
		worker := NewWorker(i, JobQueue, client, appConfig)
		worker.Start()
	}
}
