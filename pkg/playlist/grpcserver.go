package playlist

import (
	"context"
	"grpc_test/pkg/api"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type GRPCServer struct {
	api.UnimplementedAdderServer
}

func (s *GRPCServer) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	return &api.GetResponse{}, nil
}

func playlistItemsList(service *youtube.Service, part string, playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	_part := []string{part}
	call := service.PlaylistItems.List(_part)
	call = call.PlaylistId(playlistId)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}
	return response
}

func GetPlaylistItems(pageToken string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")
	playlistId := os.Getenv("PLAYLIST_ID")
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	playlistResponse := playlistItemsList(service, "snippet", playlistId, pageToken)
	log.Println(playlistResponse)
}
