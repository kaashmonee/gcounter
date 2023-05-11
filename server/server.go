package server

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/kaashmonee/gcounter/model"
	"github.com/kaashmonee/gcounter/server/workers"
	"github.com/kaashmonee/gcounter/utilities"
)

const (
	display int = iota
	merge
)

const (
	defaultMergeIntervalMS = 100
)

type server struct {
	*log.Logger
	nodes           []workers.Worker
	requests        chan model.ServerRequest
	displayResponse chan int
	mergeIntervalMS int
}

type Server interface {
	Serve()
}

// displayPage - chooses a node to display the page
func (s *server) displayPage(w http.ResponseWriter, r *http.Request) {
	s.Println("got / request")
	s.requests <- model.ServerRequest{RequestType: display}
	numViews := <-s.displayResponse
	io.WriteString(w, fmt.Sprintf("This page has: %d views", numViews))
}

func (s *server) startProcessor() {
	for {
		select {
		case request := <-s.requests:
			switch request.RequestType {
			case display:
				chosenWorker := s.nodes[rand.Intn(len(s.nodes))]
				chosenWorker.Visit()
				numViews := chosenWorker.Value()
				s.displayResponse <- numViews
			case merge:

			}
		}
	}
}

func NewServer(numWorkers int) Server {
	workersSlice := make([]workers.Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		w := workers.NewWorker(numWorkers, i)
		workersSlice[i] = w
	}

	s := &server{
		Logger:          utilities.NewLogger(nil),
		nodes:           workersSlice,
		requests:        make(chan model.ServerRequest),
		displayResponse: make(chan int),
		mergeIntervalMS: defaultMergeIntervalMS,
	}

	go s.startProcessor()

	return s
}

func (s *server) periodicMerge() {

}

func (s *server) Serve() {
	http.HandleFunc("/", s.displayPage)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("err:", err)
	}
	return
}
