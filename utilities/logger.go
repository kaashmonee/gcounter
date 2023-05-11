package utilities

import (
	"fmt"
	"log"
	"os"
)

func NewLogger(out *os.File) *log.Logger {
	if out != nil {
		return log.New(out, "", log.Lshortfile)
	}
	return log.New(os.Stdout, "", log.Lshortfile)
}

func NewInfoNodeLogger(nodeID int, out *os.File) *log.Logger {
	prefix := fmt.Sprintf("Node: %d", nodeID)
	if out != nil {
		return log.New(out, prefix, log.Lshortfile)
	}
	return log.New(os.Stdout, prefix, log.Lshortfile)
}
