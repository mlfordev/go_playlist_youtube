package main

import (
	"context"
	"fmt"
	"grpc_test/pkg/api"
	"html/template"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Result struct {
	R int32
	E error
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var x int32 = 2
	var y int32 = 2
	var res *api.AddResponse
	var rs = Result{}

	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	tpl := template.Must(template.ParseFiles("cmd/webserver/templates/index.html"))
	if err == nil {
		c := api.NewAdderClient(conn)
		res, _ = c.Add(context.Background(), &api.AddRequest{X: x, Y: y})
		rs.R = res.GetResult()
	} else {
		rs.E = err
	}
	tpl.Execute(w, rs)
}

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("cmd/webserver/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", indexHandler)
	fmt.Println("Listening port :3000")
	http.ListenAndServe(":3000", mux)
}
