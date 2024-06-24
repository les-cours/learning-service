package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/utils"
	"log"
)

func (s *Server) CreateChapter(ctx context.Context, in *learning.CreateChapterRequest) (*learning.Chapter, error) {
	err := userHasClassRoom(s.DB, in.UserID, in.ClassRoomID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	ChapterID := utils.GenerateUUIDString()

	_, err = s.DB.Exec(`
INSERT INTO 
    Chapters 
    (Chapter_id,classroom_id, title, arabic_title,description)
VALUES ($1,$2,$3,$4,$5)`, ChapterID, in.ClassRoomID, in.Title, in.ArabicTitle, in.Description)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	return &learning.Chapter{
		ChapterID:   ChapterID,
		Title:       in.Title,
		ArabicTitle: in.ArabicTitle,
		Description: in.Description,
	}, nil
}

func (s *Server) GetChaptersByClassRoom(ctx context.Context, in *learning.IDRequest) (*learning.Chapters, error) {
	rows, err := s.DB.Query(`SELECT 
    chapter_id, title, arabic_title, description,description_ar
FROM chapters 
WHERE classroom_id = $1 AND deleted_at IS NULL;`, in.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("chapter")
		}
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	var chapters = new(learning.Chapters)
	for rows.Next() {
		var chapter = new(learning.Chapter)
		err = rows.Scan(&chapter.ChapterID, &chapter.Title, &chapter.ArabicTitle, &chapter.Description, &chapter.ArabicDescription)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}
		/*
			Get lessons
		*/
		var lessons = new(learning.Lessons)
		lessons, err = s.GetLessonsByChapter(ctx, &learning.IDRequest{
			Id:     chapter.ChapterID,
			UserID: in.UserID,
		})
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}
		chapter.Lessons = lessons
		chapters.Chapters = append(chapters.Chapters, chapter)
	}
	return chapters, nil
}

func (s *Server) UpdateChapter(ctx context.Context, in *learning.UpdateChapterRequest) (*learning.Chapter, error) {

	err := userHasChapter(s.DB, in.UserID, in.ChapterID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	_, err = s.DB.Exec(`
UPDATE chapters SET title =$2, arabic_title =$3,description =$4 WHERE chapter_id =$1;
`, in.ChapterID, in.Title, in.ArabicTitle, in.Description)

	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("chapter")
		}
		return nil, ErrInternal
	}

	return &learning.Chapter{
		ChapterID:   in.ChapterID,
		Title:       in.Title,
		ArabicTitle: in.ArabicTitle,
		Description: in.Description,
	}, nil
}

func (s *Server) DeleteChapter(ctx context.Context, in *learning.IDRequest) (*learning.OperationStatus, error) {

	var err error
	err = userHasChapter(s.DB, in.UserID, in.Id)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err //not ErrInternal
	}

	res, err := s.DB.Exec(`
UPDATE chapters SET deleted_at = CURRENT_TIMESTAMP WHERE chapter_id = $1;
`, in.Id)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	if rowsAffected == 0 {
		return nil, ErrNotFound("lesson")
	}

	return &learning.OperationStatus{
		Success: true,
	}, nil
}

func (s *Server) GetChapter(ctx context.Context, in *learning.IDRequest) (*learning.Chapters, error) {

	rows, err := s.DB.Query(`SELECT 
    chapter_id, title, arabic_title, description
FROM chapters 
WHERE classroom_id = $1 AND deleted_at IS NULL;`, in.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("chapter")
		}
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	chapters := &learning.Chapters{}
	for rows.Next() {
		chapter := &learning.Chapter{}
		err = rows.Scan(&chapter.ChapterID, &chapter.Title, &chapter.ArabicTitle, &chapter.Description)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}
		chapters.Chapters = append(chapters.Chapters, chapter)
	}

	return chapters, nil
}

func userHasChapter(db *sql.DB, userID, chapterID string) error {
	var has bool
	err := db.QueryRow(`SELECT EXISTS (
	    SELECT 1
	    FROM classrooms
	    INNER JOIN chapters on classrooms.classroom_id = chapters.classroom_id
	    WHERE teacher_id = $1
	    AND chapters.chapter_id = $2
	);`, userID, chapterID).Scan(&has)

	if err != nil {
		log.Println(err.Error())
		return ErrInternal
	}

	if has {
		return nil
	}

	return ErrPermission
}
