package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/learning-service/api/learning"
)

func (s *Server) GetDocuments(ctx context.Context, in *learning.IDRequest) (*learning.Documents, error) {

	err := userHasLesson(s.DB, in.UserID, in.Id)
	if err == nil {
		return GetDocumentsForTeacher(s.DB, in.Id)
	}

	//get documents for students
	rows, err := s.DB.Query(`SELECT document_id, document_type, title, arabic_title, description, description_ar, duration, lecture_number
FROM documents 
WHERE lesson_id = $1
AND deleted_at IS NULL
ORDER BY lecture_number;`, in.Id)
	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("documents")
		}
		return nil, ErrInternal
	}

	var documentID, documentType, title, arabicTitle, description, arabicDescription string
	var lectureNumber int32
	var duration sql.NullTime

	var documents = new(learning.Documents)

	for rows.Next() {
		err = rows.Scan(&documentID, &documentType, &title, &arabicTitle, &description, &arabicDescription, &duration, &lectureNumber)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}

		document := &learning.Document{
			DocumentID:        documentID,
			DocumentType:      documentType,
			Title:             title,
			ArabicTitle:       arabicTitle,
			Description:       description,
			ArabicDescription: arabicDescription,
			LectureNumber:     lectureNumber,
		}

		if documentType == "video" {
			d := duration.Time
			document.Duration = &learning.Duration{
				Hours:       int32(d.Hour()),
				Minutes:     int32(d.Minute()),
				Seconds:     int32(d.Second()),
				Nanoseconds: int32(d.Nanosecond()),
			}
		}
		documents.Documents = append(documents.Documents, document)
	}

	return documents, nil
}

func (s *Server) GetDocument(ctx context.Context, in *learning.IDRequest) (*learning.DocumentLink, error) {

	var lessonID string

	var documentLink *sql.NullString
	var err error
	err = s.DB.QueryRow(`SELECT lesson_id,document_link
FROM documents 
WHERE document_id = $1;`, in.Id).Scan(&lessonID, &documentLink)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrNotFound("document")
	}
	if !s.CanAccessToLesson(in.UserID, lessonID) {
		return nil, ErrClassroomNotPaid
	}

	if !documentLink.Valid {
		return nil, errors.New("coming soon")
	}
	return &learning.DocumentLink{
		DocumentLink: documentLink.String,
	}, nil

}

func (s *Server) DeleteDocument(ctx context.Context, in *learning.IDRequest) (*learning.OperationStatus, error) {

	var err error
	var lessonID string
	err = s.DB.QueryRow(`SELECT lesson_id FROM documents WHERE document_id = $1;`, in.Id).Scan(&lessonID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrNotFound("documents")
	}

	err = userHasLesson(s.DB, in.UserID, in.Id)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	_, err = s.DB.Exec(`UPDATE  documents SET deleted_at = now() WHERE document_id = $1;`, lessonID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	return &learning.OperationStatus{
		Success: true,
	}, nil

}

func GetDocumentsForTeacher(db *sql.DB, lessonsID string) (*learning.Documents, error) {

	var err error

	rows, err := db.Query(`SELECT document_id, document_type, title, arabic_title, description, description_ar, duration, lecture_number,document_link
FROM documents 
WHERE lesson_id = $1
ORDER BY lecture_number;`, lessonsID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &learning.Documents{}, nil
		}
		return nil, ErrInternal
	}

	var documentID, documentType, title, arabicTitle, description, arabicDescription, documentLink string
	var lectureNumber int32
	var duration sql.NullTime

	var documents = new(learning.Documents)

	for rows.Next() {
		err = rows.Scan(&documentID, &documentType, &title, &arabicTitle, &description, &arabicDescription, &duration, &lectureNumber, &documentLink)
		if err != nil {
			return nil, ErrInternal
		}

		document := &learning.Document{
			DocumentID:        documentID,
			DocumentType:      documentType,
			Title:             title,
			ArabicTitle:       arabicTitle,
			Description:       description,
			ArabicDescription: arabicDescription,
			LectureNumber:     lectureNumber,
			DocumentLink:      documentLink,
		}

		if documentType == "video" {
			d := duration.Time
			document.Duration = &learning.Duration{
				Hours:       int32(d.Hour()),
				Minutes:     int32(d.Minute()),
				Seconds:     int32(d.Second()),
				Nanoseconds: int32(d.Nanosecond()),
			}
		}
		documents.Documents = append(documents.Documents, document)
	}

	return documents, nil
}
