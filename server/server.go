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
	s.serverRequests <- model.ServerRequest{Type: constant.Request.Display()}
	serverResponseFromWorker := <-s.serverResponse
	numViews := serverResponseFromWorker.Payload.(int)
	io.WriteString(w, fmt.Sprintf("This page has: %d views\n", numViews))
	io.WriteString(w, fmt.Sprintf("Served by node: %d/%d", serverResponseFromWorker.WorkerNodeID+1, len(s.nodes)))
}

// this just checks if any of the workers have sent a signal
func (s *Server) workerListener(requestsFromWorker <-chan model.WorkerRequest) {
	for {
		request := <-requestsFromWorker
		serverRequest := model.ServerRequest{Type: request.Type, Payload: request.Payload}
		s.serverRequests <- serverRequest
	}
}

func (s *Server) startProcessor() {
	for {
		select {
		case request := <-s.serverRequests:
			switch request.Type {
			case constant.Request.Display():
				chosenWorker := s.nodes[rand.Intn(len(s.nodes))]
				chosenWorker.Visit()
				numViews := chosenWorker.Value()
				s.serverResponse <- model.ServerResponse{Type: constant.Request.Display(), Payload: numViews, WorkerNodeID: chosenWorker.ID}
			case constant.Request.Merge():
				newViewsAll := request.Payload.([]int)
				// Now send this to all the other nodes
				for _, node := range s.nodes {
					go func(n *workers.Worker) {
						n.MasterRequests <- model.ServerRequest{Type: constant.Request.Merge(), Payload: newViewsAll}
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
	for _, node := range s.nodes {
		go s.workerListener(node.WorkerRequestToMaster)
	}

	s.Println("started processor")
	return s
}

func (s *Server) Serve() {
	http.HandleFunc("/", s.displayPage)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("err:", err)
	}
	return
}
