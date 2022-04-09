package main

import (
	"context"
	"flag"
	"grpc_test/pkg/api"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func indexHandler(w http.ResponseWriter, r *http.Request, playlistId string) {
	query := r.URL.Query()
	var pageToken string = query.Get("page")
	log.Println("page", pageToken)
	var res *api.GetResponse

	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	tpl := template.Must(template.ParseFiles("cmd/webserver/templates/index.html"))
	if err == nil {
		c := api.NewPlaylistClient(conn)
		res, _ = c.Get(context.Background(), &api.GetRequest{
			PlaylistId: playlistId,
			PageToken:  pageToken,
		})
		log.Println(len(res.Items))
		for _, value := range res.Items {
			log.Println(value.GetId())
		}
	} else {
		log.Fatal(err)
	}
	tpl.Execute(w, res)
}

func main() {
	var playlistId string
	flag.Parse()
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("cmd/webserver/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	if flag.NArg() > 0 {
		playlistId = flag.Arg(0)
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file in package main")
		}
		playlistId = os.Getenv("PLAYLIST_ID")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, playlistId)
	})
	log.Println("Listening port :3000")
	http.ListenAndServe(":3000", mux)
}
