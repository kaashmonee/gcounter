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
	WorkerRequestToMaster          chan model.WorkerRequest
	id                             int
	views                          []int
	workerRequests                 chan model.WorkerRequest
	workerResponses                chan model.WorkerResponse
	defaultMergeIntervalMS         float64
	intervalRandomizationOffsetPct float64
}

func NewWorker(totalWorkers int, id int, masterRequests chan model.ServerRequest) *Worker {
	w := &Worker{
		id:                             id,
		views:                          make([]int, totalWorkers),
		Logger:                         utilities.NewInfoNodeLogger(id, os.Stdout),
		workerRequests:                 make(chan model.WorkerRequest),
		workerResponses:                make(chan model.WorkerResponse),
		MasterRequests:                 masterRequests,
		WorkerRequestToMaster:          make(chan model.WorkerRequest),
		defaultMergeIntervalMS:         100,
		intervalRandomizationOffsetPct: 20,
	}
	go w.startWorker()
	return w
}

func (w *Worker) startWorker() {
	for {
		select {
		case request := <-w.workerRequests:
			switch request.Type {
			case constant.Request.Visit():
				w.views[w.id]++
			case constant.Request.Value():
				total := 0
				for _, count := range w.views {
					total += count
				}
				w.workerResponses <- model.WorkerResponse{Type: constant.Request.Display(), Payload: total}
			}
		// request to merge after a somewhat random amount of time
		case <-time.After(w.getMergeDuration()):
			w.WorkerRequestToMaster <- model.WorkerRequest{Type: constant.Request.Merge(), Payload: w.views}
		case masterRequest := <-w.MasterRequests:
			switch masterRequest.Type {
			case constant.Request.Merge():
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
	w.workerRequests <- model.WorkerRequest{Type: constant.Request.Visit(), Payload: w.id}
}

func (w *Worker) Value() int {
	w.workerRequests <- model.WorkerRequest{Type: constant.Request.Value()}
	result := <-w.workerResponses
	return result.Payload.(int)
}
