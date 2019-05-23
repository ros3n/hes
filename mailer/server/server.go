package server

import "net/http"

func StartHealthCheckServer(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(addr, nil)
}
