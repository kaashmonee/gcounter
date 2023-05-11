package workers

import (
	"log"
	"os"

	"github.com/kaashmonee/gcounter/utilities"
)

const (
	visit int = iota
	value
)

type requestData struct {
	requestType int
	payload     int
}

type Worker interface {
	Visit()
	Value() int
}

type worker struct {
	*log.Logger
	id          int
	views       []int
	requestChan chan requestData
	valueChan   chan int
}

func NewWorker(totalWorkers int, id int) Worker {
	w := &worker{
		id:          id,
		views:       make([]int, totalWorkers),
		Logger:      utilities.NewInfoNodeLogger(id, os.Stdout),
		requestChan: make(chan requestData),
		valueChan:   make(chan int),
	}
	go w.startWorker()
	return w
}

func (w *worker) startWorker() {
	for {
		select {
		case request := <-w.requestChan:
			switch request.requestType {
			case visit:
				w.views[w.id]++
			case value:
				total := 0
				for _, count := range w.views {
					total += count
				}
				w.valueChan <- total
			}
		}
	}
}

func (w *worker) Visit() {
	w.requestChan <- requestData{requestType: visit, payload: w.id}
}

func (w *worker) Value() int {
	w.requestChan <- requestData{requestType: value}
	result := <-w.valueChan
	return result
}
