package resolvers

import (
	"context"
	apivideosdk "github.com/apivideo/api.video-go-client"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/utils"
	"os"
	"time"
)

const (
	port       = ":50051"
	bucketName = "your-bucket-name"
)

func (s *Server) CreateVideo(ctx context.Context, in *learning.CreateVideoRequest) (*learning.Document, error) {

	//var canCreate bool
	//var err error
	//canCreate, err = checkCreateLessonPermission(s.DB, in.UserID, in.LessonID)
	//
	//switch {
	//case err != nil:
	//	if errors.Is(err, sql.ErrNoRows) {
	//		return nil, ErrPermission
	//	}
	//	return nil, ErrInternal
	//case !canCreate:
	//	return nil, ErrPermission
	//}

	//upload Video
	video, err := s.UploadVideo(&learning.UploadVideoRequest{
		Content:  in.Video.Content,
		Filename: in.Video.Filename,
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	//store the document
	documentID := utils.GenerateUUIDString()

	_, err = s.DB.Exec(`
INSERT INTO 
    Documents 
    (document_id, lesson_id, document_type, title, arabic_title, description, document_link, lecture_number)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		documentID, in.LessonID, "video", in.Title, in.ArabicTitle, in.Description, video.GetVideoId(), in.LectureNumber)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	//get video
	status, err := s.VideoApi.Videos.GetStatus(video.Source.GetUri())
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	t := time.Unix(int64(*status.Encoding.Metadata.Duration.Get()), 0).UTC()
	duration := &learning.Duration{
		Hours:       int32(t.Hour()),
		Minutes:     int32(t.Minute()),
		Seconds:     int32(t.Second()),
		Nanoseconds: int32(t.Nanosecond()),
	}

	return &learning.Document{
		DocumentID:    documentID,
		DocumentType:  "video",
		Title:         in.Title,
		ArabicTitle:   in.ArabicTitle,
		Description:   in.Description,
		Duration:      duration,
		LectureNumber: in.LectureNumber,
	}, nil
}

func (s *Server) UploadVideo(in *learning.UploadVideoRequest) (*apivideosdk.Video, error) {

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
