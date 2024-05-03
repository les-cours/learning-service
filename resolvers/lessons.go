package resolvers

import (
	"context"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/utils"
)

func (s *Server) CreateLesson(ctx context.Context, in *learning.CreateLessonRequest) (*learning.Lesson, error) {

	//var canCreate bool
	//var err error
	//canCreate, err = checkCreateLessonPermission(s.DB, in.UserID, in.ChapterID)
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

	lessonID := utils.GenerateUUIDString()

	_, err := s.DB.Exec(`
INSERT INTO 
    lessons 
    (lesson_id,chapter_id, title, arabic_title,description)
VALUES ($1,$2,$3,$4,$5)`, lessonID, in.ChapterID, in.Title, in.ArabicTitle, in.Description)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return &learning.Lesson{
		LessonID:    lessonID,
		Title:       in.Title,
		ArabicTitle: in.ArabicTitle,
		Description: in.Description,
	}, nil
}

func (s *Server) GetLessonsByClassRoom(ctx context.Context, in *learning.IDRequest) (*learning.Lessons, error) {
	var classRoomID = in.Id
	lessons := &learning.Lessons{}
	lesson := &learning.Lesson{}
	rows, err := s.DB.Query(`SELECT lesson_id,title FROM lessons WHERE classroom_id = $1`, classRoomID)
	if err != nil {

		return nil, ErrInternal
	}

	for rows.Next() {
		err = rows.Scan(&lesson.LessonID, &lesson.Title)
		if err != nil {

			return nil, ErrInternal
		}
		lessons.Lessons = append(lessons.Lessons, lesson)
	}

	return lessons, nil
}

/*
	func checkCreateLessonPermission(db *sql.DB, userID, chapterID string) (bool, error) {
		var canCreate bool
		err := db.QueryRow(`SELECT EXISTS (
		    SELECT 1
		    FROM classrooms
		    WHERE teacher_id = $1
		    AND classroom_id = (SELECT classroom_id FROM chapters WHERE chapter_id = $2 LIMIT 1)
		);`, userID, chapterID).Scan(&canCreate)

		return canCreate, err
	}
*/
