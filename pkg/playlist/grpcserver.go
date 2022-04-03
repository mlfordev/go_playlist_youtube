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
	api.UnimplementedPlaylistServer
}

func (s *GRPCServer) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	grpsResponse := GetPlaylistItems(req.PlaylistId, req.PageToken)
	pageInfo := &api.GetResponse_Pageinfo{
		TotalResults: uint32(grpsResponse.PageInfo.TotalResults),
		ResultsPerPage: uint32(grpsResponse.PageInfo.ResultsPerPage),
	}
	var item *api.GetResponse_Items
	items := []*api.GetResponse_Items{}
	for _, value := range grpsResponse.Items {
		// log.Println(value.Snippet.Title, value.Snippet.Thumbnails.Maxres)
		thumbnails := &api.GetResponse_Thumbnails{}
		if value.Snippet.Thumbnails.Default != nil {
			thumbnails.Default = &api.GetResponse_Default{
				Url: value.Snippet.Thumbnails.Default.Url,
				Width: uint32(value.Snippet.Thumbnails.Default.Width),
				Height: uint32(value.Snippet.Thumbnails.Default.Height),
			}
		}
		if value.Snippet.Thumbnails.Medium != nil {
			thumbnails.Medium = &api.GetResponse_Medium{
				Url: value.Snippet.Thumbnails.Medium.Url,
				Width: uint32(value.Snippet.Thumbnails.Medium.Width),
				Height: uint32(value.Snippet.Thumbnails.Medium.Height),
			}
		}
		if value.Snippet.Thumbnails.High != nil {
			thumbnails.High = &api.GetResponse_High{
				Url: value.Snippet.Thumbnails.High.Url,
				Width: uint32(value.Snippet.Thumbnails.High.Width),
				Height: uint32(value.Snippet.Thumbnails.High.Height),
			}
		}
		if value.Snippet.Thumbnails.Standard != nil {
			thumbnails.Standard = &api.GetResponse_Standard{
				Url: value.Snippet.Thumbnails.Standard.Url,
				Width: uint32(value.Snippet.Thumbnails.Standard.Width),
				Height: uint32(value.Snippet.Thumbnails.Standard.Height),
			}
		}
		if value.Snippet.Thumbnails.Maxres != nil {
			thumbnails.Maxres = &api.GetResponse_Maxres{
				Url:    value.Snippet.Thumbnails.Maxres.Url,
				Width:  uint32(value.Snippet.Thumbnails.Maxres.Width),
				Height: uint32(value.Snippet.Thumbnails.Maxres.Height),
			}
		}
		resourceId := &api.GetResponse_Resourceid{
			Kind:    value.Snippet.ResourceId.Kind,
			VideoId: value.Snippet.ResourceId.VideoId,
		}
		snippet := &api.GetResponse_Snippet {
			PublishedAt: value.Snippet.PublishedAt,
			ChannelId: value.Snippet.ChannelId,
			Title: value.Snippet.Title,
			Description: value.Snippet.Description,
			Thumbnails: thumbnails,
			ChannelTitle: value.Snippet.ChannelTitle,
			PlaylistId: value.Snippet.PlaylistId,
			Position: uint32(value.Snippet.Position),
			ResourceId: resourceId,
			VideoOwnerChannelTitle: value.Snippet.VideoOwnerChannelTitle,
			VideoOwnerChannelId: value.Snippet.VideoOwnerChannelId,
		}
		item = &api.GetResponse_Items {
			Kind: value.Kind,
			Etag: value.Etag,
			Id: value.Id,
			Snippet: snippet,
		}
		items = append(items, item)
	}
	response := &api.GetResponse{
		Kind: grpsResponse.Kind,
		Etag: grpsResponse.Etag,
		NextPageToken: grpsResponse.NextPageToken,
		Items: items,
		PrevPageToken: grpsResponse.PrevPageToken,
		PageInfo: pageInfo,
	}
	return response, nil
}

func playlistItemsList(service *youtube.Service, part string, playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	p := []string{part}
	call := service.PlaylistItems.List(p)
	call = call.PlaylistId(playlistId).MaxResults(10)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}
	return response
}

func GetPlaylistItems(playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in package playlst")
	}
	apiKey := os.Getenv("API_KEY")
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	return playlistItemsList(service, "snippet", playlistId, pageToken)
}
