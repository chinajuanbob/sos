/*
 * api_server.go
 *
 * Trivial rewrite of the api-server in #golang.
 *
 * Presents two end-points, on two different ports:
 *
 *    http://127.0.0.1:9991/upload
 *
 *    http://127.0.0.1:9992/fetch/:id
 */

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

/**
 * This is the list of servers we know about.
 *
 * It is populated by reading ~/.sos.conf and /etc/sos.conf
 */
var SERVERS []string

/**
 * Upload a file to to the public-root.
 */
func UploadHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Upload FAILED")
	return

}

/**
 * Download a file.
 */
func DownloadHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(res, "Object not found.")
}

/**
 * Entry point to our code.
 */
func main() {

	//
	// Parse our command-line arguments.
	//
	host := flag.String("host", "0.0.0.0", "The IP to listen upon.")
	dport := flag.Int("download-port", 9992, "The port to bind upon for downloading objects.")
	uport := flag.Int("upload-port", 9991, "The port to bind upon for uploading objects.")
	flag.Parse()

        //
        // Pretend we have only a pair of servers
        //
        SERVERS = append( SERVERS, "http://127.0.0.1:4000" )
        SERVERS = append( SERVERS, "http://127.0.0.1:4001" )

	//
	// Show a banner.
	//
	fmt.Printf("Launching API-server\n")
	fmt.Printf("\nUpload service\nhttp://127.0.0.1:%d/upload\n", *uport)
	fmt.Printf("\nDownload service\nhttp://127.0.0.1:%d/fetch/:id\n", *dport)

	//
	// Create a route for uploading.
	//
	up_router := mux.NewRouter()
	up_router.HandleFunc("/upload", UploadHandler).Methods("POST")

	//
	// Create a route for downloading.
	//
	down_router := mux.NewRouter()
	down_router.HandleFunc("/fetch/{id}", DownloadHandler).Methods("GET")

	//
	// The following code is a hack to allow us to run two distinct
	// HTTP-servers on different ports.
	//
	// It's almost sexy the way it worked the first time :)
	//
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *uport),
			up_router)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *dport),
			down_router)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	wg.Wait()
}