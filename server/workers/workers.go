package workers

import (
	"log"
	"os"
	"time"

	"github.com/kaashmonee/gcounter/model"
	"github.com/kaashmonee/gcounter/utilities"
)

const (
	visit int = iota
	value
	update
)

type Worker interface {
	Visit()
	Value() int
}

type worker struct {
	*log.Logger
	id                             int
	views                          []int
	requestChan                    chan model.WorkerRequest
	valueChan                      chan int
	defaultMergeIntervalMS         float64
	intervalRandomizationOffsetPct float64
}

func NewWorker(totalWorkers int, id int) Worker {
	w := &worker{
		id:                             id,
		views:                          make([]int, totalWorkers),
		Logger:                         utilities.NewInfoNodeLogger(id, os.Stdout),
		requestChan:                    make(chan model.WorkerRequest),
		valueChan:                      make(chan int),
		defaultMergeIntervalMS:         100,
		intervalRandomizationOffsetPct: 20,
	}
	go w.startWorker()
	return w
}

func (w *worker) startWorker() {
	for {
		select {
		case request := <-w.requestChan:
			switch request.RequestType {
			case visit:
				w.views[w.id]++
			case value:
				total := 0
				for _, count := range w.views {
					total += count
				}
				w.valueChan <- total
			case update:
				data := request.Payload.([]int)
				for i := 0; i < len(w.views); i++ {
					w.views[i] = data[i]
				}
			}
		case <-time.After(w.getMergeDuration()):
			// send a merge request to the master node
		}
	}
}

func (w *worker) getMergeDuration() time.Duration {
	offset := utilities.RandIntInRange(
		int(w.defaultMergeIntervalMS)-int(w.defaultMergeIntervalMS*w.intervalRandomizationOffsetPct*0.01),
		int(w.defaultMergeIntervalMS)+int(w.defaultMergeIntervalMS*w.intervalRandomizationOffsetPct*0.01),
	)
	return time.Duration(int(w.defaultMergeIntervalMS)+offset) * time.Millisecond
}

func (w *worker) Visit() {
	w.requestChan <- model.WorkerRequest{RequestType: visit, Payload: w.id}
}

func (w *worker) Value() int {
	w.requestChan <- model.WorkerRequest{RequestType: value}
	result := <-w.valueChan
	return result
}
