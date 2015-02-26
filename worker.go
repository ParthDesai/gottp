package gottp

import (
	"gopkg.in/simversity/gotracer.v1"
	"sync"
)

var workerRunning bool

var worker func(chan bool)

var errChan = make(chan bool)

var exitChan = make(chan bool, 1)

var wg = new(sync.WaitGroup)

func spawner() {
	go workerWrapper()

	s := <-errChan
	if s {
		go spawner()
	}
}

func workerWrapper() {

	wg.Add(1)

	defer wg.Done()
	defer gotracer.Tracer{
		Dummy:         settings.Gottp.EmailDummy,
		EmailHost:     settings.Gottp.EmailHost,
		EmailPort:     settings.Gottp.EmailPort,
		EmailPassword: settings.Gottp.EmailPassword,
		EmailUsername: settings.Gottp.EmailUsername,
		EmailSender:   settings.Gottp.EmailSender,
		EmailFrom:     settings.Gottp.EmailFrom,
		ErrorTo:       settings.Gottp.ErrorTo,
	}.Notify(func() string {
		errChan <- true
		return "Exception in worker"
	})

	worker(exitChan)
	errChan <- false
}

func RunWorker(wk func(chan bool)) {
	if workerRunning {
		panic("Worker already running.")
	}
	worker = wk
	workerRunning = true
	go spawner()
}

func StopWorker() {
	if workerRunning {
		exitChan <- true
		wg.Wait()
		workerRunning = false
	}
}
