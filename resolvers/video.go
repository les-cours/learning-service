package resolvers

import (
	"context"
	apivideosdk "github.com/apivideo/api.video-go-client"
	"github.com/les-cours/learning-service/api/learning"
	"os"
)

const (
	port       = ":50051"
	bucketName = "your-bucket-name"
)

func (s *Server) UploadVideo(ctx context.Context, in *learning.UploadVideoRequest) (*apivideosdk.Video, error) {

	videoCreationPayload := *apivideosdk.NewVideoCreationPayload(in.Filename)
	video, err := s.VideoApi.Videos.Create(videoCreationPayload)

	file, err := os.Create(in.Filename)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	defer file.Close()

	// Write the uploaded content to the file
	_, err = file.Write(in.Content)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	video, err = s.VideoApi.Videos.UploadFile(video.GetVideoId(), file)

	return video, nil
}

//
//func (s *server) GetVideo(ctx context.Context, in *pb.GetVideoRequest) (*pb.GetVideoResponse, error) {
//	// Authenticate with Google Cloud Storage
//	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
//	if err != nil {
//		return nil, err
//	}
//	defer client.Close()
//
//	// Create a bucket handle
//	bucket := client.Bucket(bucketName)
//
//	// Read the video content from Google Cloud Storage
//	obj := bucket.Object(in.VideoId + ".mp4")
//	rc, err := obj.NewReader(ctx)
//	if err != nil {
//		return nil, err
//	}
//	defer rc.Close()
//
//	content, err := ioutil.ReadAll(rc)
//	if err != nil {
//		return nil, err
//	}
//
//	return &pb.GetVideoResponse{Content: content}, nil
//}
