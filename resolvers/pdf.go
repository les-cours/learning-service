package resolvers

import (
	"context"
	"github.com/les-cours/learning-service/api/learning"
)

func (s *Server) CreatePdf(ctx context.Context, in *learning.CreatePdfRequest) (*learning.Document, error) {

	err := userHasLesson(s.DB, in.UserID, in.LessonID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	_, err = s.DB.Exec(`INSERT INTO documents (lesson_id, document_type, title, arabic_title, description, description_ar, lecture_number, document_link) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, in.LessonID, "pdf", in.Title, in.ArabicTitle, in.Description, in.ArabicDescription, in.LectureNumber, in.Url)
	if err != nil {
		return nil, ErrInternal
	}
	return &learning.Document{
		DocumentType:      "pdf",
		Title:             in.Title,
		ArabicTitle:       in.ArabicTitle,
		Description:       in.Description,
		ArabicDescription: in.ArabicDescription,
		LectureNumber:     in.LectureNumber,
		DocumentLink:      in.Url,
	}, nil

}
