package workers

import (
	"log"
	"os"

	"github.com/kaashmonee/gcounter/utilities"
)

type Worker interface {
	Visit()
	Value() int
	Logger() *log.Logger
}

type worker struct {
	lg    *log.Logger
	id    int
	views []int
}

func NewWorker(totalWorkers int, id int) Worker {
	return &worker{
		id:    id,
		views: make([]int, totalWorkers),
		lg:    utilities.NewInfoNodeLogger(id, os.Stdout),
	}
}

func (w *worker) Visit() {
	w.views[w.id]++
}

func (w *worker) Value() int {
	sum := 0
	for _, cnt := range w.views {
		sum += cnt
	}
	return sum
}

func (w *worker) Logger() *log.Logger {
	return w.lg
}
