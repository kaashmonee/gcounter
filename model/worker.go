package model

type WorkerRequest struct {
	Type    int
	Payload any
}

type WorkerResponse struct {
	Type    int
	Payload any
}
