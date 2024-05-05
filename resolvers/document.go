package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/utils"
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
	video, err := s.UploadVideo(ctx, &learning.UploadVideoRequest{
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
	return &learning.Document{
		DocumentID:    documentID,
		DocumentType:  "video",
		Title:         in.Title,
		ArabicTitle:   in.ArabicTitle,
		Description:   in.Description,
		Duration:      0,
		LectureNumber: in.LectureNumber,
	}, nil
}

func (s *Server) GetDocumentsByLesson(ctx context.Context, in *learning.IDRequest) (*learning.Documents, error) {

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

	rows, err := s.DB.Query(`SELECT document_id, title, arabic_title, description, lecture_number, document_link
FROM documents 
WHERE lesson_id = $1 AND document_type = 'video'
ORDER BY lecture_number;`, in.Id)
	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("documents")
		}
		return nil, ErrInternal
	}

	documents := &learning.Documents{}
	for rows.Next() {
		videoDocument := &learning.Document{}
		err = rows.Scan(&videoDocument.DocumentID, &videoDocument.Title, &videoDocument.ArabicTitle, &videoDocument.Description, &videoDocument.LectureNumber, &videoDocument.DocumentLink)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}

		//get video
		video, err := s.VideoApi.Videos.Get(videoDocument.DocumentLink)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}
		status, err := s.VideoApi.Videos.GetStatus(videoDocument.DocumentLink)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}

		videoDocument.DocumentLink = video.Assets.GetPlayer()
		videoDocument.Duration = *status.Encoding.Metadata.Duration.Get()
		documents.Documents = append(documents.Documents, videoDocument)
	}

	/*
	   PDFS ...
	*/
	return documents, nil
}
