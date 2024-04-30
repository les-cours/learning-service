package resolvers

import (
	"context"
	"github.com/les-cours/learning-service/api/learning"
)

func (s *Server) GetLessonsByClassRoom(ctx context.Context, in *learning.GetByUUIDRequest) (*learning.Lessons, error) {
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
