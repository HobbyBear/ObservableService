package main

import (
	"ObservableService/servers/postserver/api"
	"ObservableService/servers/postserver/middlewire"
	"ObservableService/trace"
	"context"
	"log"
	"net"
	"net/http"
)

func main() {

	trace.Init()
	defer trace.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/get_post",middlewire.Trace(api.GetPostHandler))

	server := &http.Server{
		Handler: mux,
	}
	listener,err := net.Listen("tcp","127.0.0.1:8080")
	if err != nil{
		log.Println("init listener fail ", err)
		return
	}
	err = server.Serve(listener)
	if err != nil{
		log.Println("init server fail ", err)
		return
	}
	defer server.Shutdown(context.TODO())

}
