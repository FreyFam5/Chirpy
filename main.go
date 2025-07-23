package main

import "net/http"


func main() {
	sm := http.ServeMux{}
	
	sm.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Handler: &sm,
		Addr: ":8080",
	}

	server.ListenAndServe()
}