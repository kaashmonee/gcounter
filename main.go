package main

import "github.com/kaashmonee/gcounter/server"

// just need a couple server nodes that serve a site
// the server nodes need a server master that chooses which node actually should serve the site
// the server master should also trigger the nodes periodically so that they share data

func main() {
	srv := server.NewServer(10)
	srv.Serve()
}
