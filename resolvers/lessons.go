package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/utils"
)

func (s *Server) CreateLesson(ctx context.Context, in *learning.CreateLessonRequest) (*learning.Lesson, error) {

	err := userHasChapter(s.DB, in.UserID, in.ChapterID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	lessonID := utils.GenerateUUIDString()

	_, err = s.DB.Exec(`
INSERT INTO 
    lessons 
    (lesson_id,chapter_id, title, arabic_title,description,description_ar,lesson_order)
VALUES ($1,$2,$3,$4,$5,$6,$7)`, lessonID, in.ChapterID, in.Title, in.ArabicTitle, in.Description, in.Description, in.LessonOrder)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return &learning.Lesson{
		LessonID:          lessonID,
		Title:             in.Title,
		ArabicTitle:       in.ArabicTitle,
		Description:       in.Description,
		ArabicDescription: in.Description,
		LessonOrder:       in.LessonOrder,
	}, nil
}

func (s *Server) GetLessonsByChapter(ctx context.Context, in *learning.IDRequest) (*learning.Lessons, error) {

	var chapterID = in.Id
	var lessons = new(learning.Lessons)
	rows, err := s.DB.Query(`SELECT lesson_id, title, arabic_title, description,description_ar,lesson_order
	FROM lessons WHERE chapter_id = $1 AND deleted_at IS  NULL;`, chapterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("lesson")
		}
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	for rows.Next() {
		var lesson = new(learning.Lesson)
		err = rows.Scan(&lesson.LessonID, &lesson.Title, &lesson.ArabicTitle, &lesson.Description, &lesson.ArabicDescription, &lesson.LessonOrder)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}

		var documents = new(learning.Documents)
		lesson.Documents = documents

		documents, err = s.GetDocuments(ctx, &learning.IDRequest{
			Id:     lesson.LessonID,
			UserID: in.UserID,
		})

		if err != nil {
			s.Logger.Error(err.Error())
			goto APPEND
		}

		lesson.Documents = documents

		/*
			Append
		*/

	APPEND:
		lessons.Lessons = append(lessons.Lessons, lesson)
	}

	return lessons, nil
}

func (s *Server) UpdateLesson(ctx context.Context, in *learning.UpdateLessonRequest) (*learning.Lesson, error) {

	err := userHasLesson(s.DB, in.UserID, in.LessonID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	_, err = s.DB.Exec(`
UPDATE lessons SET chapter_id = $2, title = $3, arabic_title = $4,description = $5 WHERE lesson_id = $1;
`, in.LessonID, in.ChapterID, in.Title, in.ArabicTitle, in.Description)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("lesson")
		}
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return &learning.Lesson{
		LessonID:    in.LessonID,
		Title:       in.Title,
		ArabicTitle: in.ArabicTitle,
		Description: in.Description,
	}, nil
}

func (s *Server) DeleteLesson(ctx context.Context, in *learning.IDRequest) (*learning.OperationStatus, error) {

	err := userHasLesson(s.DB, in.UserID, in.Id)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	res, err := s.DB.Exec(`
UPDATE lessons SET deleted_at = CURRENT_TIMESTAMP WHERE lesson_id = $1;
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

func userHasLesson(db *sql.DB, userID, lessonID string) error {

	var has bool
	err := db.QueryRow(`SELECT EXISTS (
    SELECT 1
    FROM classrooms
    JOIN chapters ON classrooms.classroom_id = chapters.classroom_id
    JOIN lessons ON chapters.chapter_id = lessons.chapter_id
    WHERE classrooms.teacher_id = $1
    AND lessons.lesson_id = $2
);
`, userID, lessonID).Scan(&has)

	if err != nil {
		return err
	}

	if has {
		return nil
	}

	return ErrPermission
}
func (s *Server) CanAccessToLesson(studentID, lessonID string) bool {

	var classroomID string
	err := s.DB.QueryRow(`SELECT classrooms.classroom_id FROM lessons 
                  INNER JOIN chapters  on lessons.chapter_id = chapters.chapter_id
                  INNER JOIN classrooms  on chapters.classroom_id = classrooms.classroom_id
                  where lesson_id = $1;`, lessonID).Scan(&classroomID)
	if err != nil {
		s.Logger.Error(err.Error())
		return false
	}

	return s.CanAccessToClassRoom(studentID, classroomID)

}
