package model

type WorkerRequest struct {
	Type    int
	Payload any
}

type WorkerResponse struct {
	NodeID  int
	Type    int
	Payload any
}
