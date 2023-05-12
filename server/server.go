package server

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/kaashmonee/gcounter/constant"
	"github.com/kaashmonee/gcounter/model"
	"github.com/kaashmonee/gcounter/server/workers"
	"github.com/kaashmonee/gcounter/utilities"
)

const (
	defaultMergeIntervalMS = 100
)

type Server struct {
	*log.Logger
	nodes           []*workers.Worker
	serverRequests  chan model.ServerRequest
	serverResponse  chan model.ServerResponse
	workerResponse  chan model.WorkerResponse
	mergeIntervalMS int
}

// displayPage - chooses a node to display the page
func (s *Server) displayPage(w http.ResponseWriter, r *http.Request) {
	s.serverRequests <- model.ServerRequest{Type: constant.ServerRequest.Display()}
	workerResponse := <-s.serverResponse
	numViews := workerResponse.Payload.(int)
	io.WriteString(w, fmt.Sprintf("This page has: %d views", numViews))
}

func (s *Server) startProcessor() {
	for {
		select {
		case request := <-s.serverRequests:
			switch request.Type {
			case constant.ServerRequest.Display():
				chosenWorker := s.nodes[rand.Intn(len(s.nodes))]
				chosenWorker.Visit()
				numViews := chosenWorker.Value()
				s.serverResponse <- model.ServerResponse{Type: constant.ServerRequest.Display(), Payload: numViews}
			case constant.WorkerRequest.Merge():
				newViewsAll := request.Payload.([]int)
				// Now send this to all the other nodes
				for _, node := range s.nodes {
					go func(n *workers.Worker) {
						n.MasterRequests <- model.ServerRequest{Type: constant.ServerRequest.Merge(), Payload: newViewsAll}
					}(node)
				}
			}
		}
	}
}

func NewServer(numWorkers int) *Server {
	workersSlice := make([]*workers.Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		masterRequestChan := make(chan model.ServerRequest)
		w := workers.NewWorker(numWorkers, i, masterRequestChan)
		workersSlice[i] = w
	}

	s := &Server{
		Logger:          utilities.NewLogger(nil),
		nodes:           workersSlice,
		serverRequests:  make(chan model.ServerRequest),
		serverResponse:  make(chan model.ServerResponse),
		workerResponse:  make(chan model.WorkerResponse),
		mergeIntervalMS: defaultMergeIntervalMS,
	}

	s.Printf("initialized %d workers\n", numWorkers)

	go s.startProcessor()

	s.Println("started processor")
	return s
}

func (s *Server) Serve() {
	http.HandleFunc("/", s.displayPage)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("err:", err)
	}
	return
}
