package workers

import (
	"log"
	"os"
	"time"

	"github.com/kaashmonee/gcounter/constant"
	"github.com/kaashmonee/gcounter/model"
	"github.com/kaashmonee/gcounter/utilities"
)

type Worker struct {
	*log.Logger
	MasterRequests                 chan model.ServerRequest // Only pull from this channel
	id                             int
	views                          []int
	requestChan                    chan model.WorkerRequest
	valueChan                      chan int
	defaultMergeIntervalMS         float64
	intervalRandomizationOffsetPct float64
}

func NewWorker(totalWorkers int, id int, masterRequests chan model.ServerRequest) *Worker {
	w := &Worker{
		id:                             id,
		views:                          make([]int, totalWorkers),
		Logger:                         utilities.NewInfoNodeLogger(id, os.Stdout),
		requestChan:                    make(chan model.WorkerRequest),
		valueChan:                      make(chan int),
		MasterRequests:                 masterRequests,
		defaultMergeIntervalMS:         100,
		intervalRandomizationOffsetPct: 20,
	}
	go w.startWorker()
	return w
}

func (w *Worker) startWorker() {
	for {
		select {
		case request := <-w.requestChan:
			switch request.Type {
			case constant.WorkerRequest.Visit():
				w.views[w.id]++
			case constant.WorkerRequest.Value():
				total := 0
				for _, count := range w.views {
					total += count
				}
				w.valueChan <- total
			case constant.WorkerRequest.Update():
				data := request.Payload.([]int)
				for i := 0; i < len(w.views); i++ {
					w.views[i] = data[i]
				}
			}
		case <-time.After(w.getMergeDuration()):
			// instead of doing the merge as an operation, just keep doing it
		case masterRequest := <-w.MasterRequests:
			switch masterRequest.Type {
			case constant.ServerRequest.Merge():
				masterViews := masterRequest.Payload.([]int)
				for i := 0; i < len(masterViews); i++ {
					w.views[i] = utilities.Max(w.views[i], masterViews[i])
				}
			}
		}
	}
}

func (w *Worker) getMergeDuration() time.Duration {
	offset := utilities.RandIntInRange(
		int(w.defaultMergeIntervalMS)-int(w.defaultMergeIntervalMS*w.intervalRandomizationOffsetPct*0.01),
		int(w.defaultMergeIntervalMS)+int(w.defaultMergeIntervalMS*w.intervalRandomizationOffsetPct*0.01),
	)
	return time.Duration(int(w.defaultMergeIntervalMS)+offset) * time.Millisecond
}

func (w *Worker) Visit() {
	w.requestChan <- model.WorkerRequest{Type: constant.WorkerRequest.Visit(), Payload: w.id}
}

func (w *Worker) Value() int {
	w.requestChan <- model.WorkerRequest{Type: constant.WorkerRequest.Value()}
	result := <-w.valueChan
	return result
}
