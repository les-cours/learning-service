package resolvers

import (
	"context"
	"database/sql"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/utils"
)

func (s *Server) CreateChapter(ctx context.Context, in *learning.CreateChapterRequest) (*learning.Chapter, error) {

	//var canCreate bool
	//var err error
	//canCreate, err = checkCreateChapterPermission(s.DB, in.UserID, in.ClassRoomID)
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

	ChapterID := utils.GenerateUUIDString()

	_, err := s.DB.Exec(`
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

func checkCreateChapterPermission(db *sql.DB, userID, classRoomID string) (bool, error) {
	var canCreate bool
	err := db.QueryRow(`SELECT EXISTS (
	    SELECT 1 
	    FROM classrooms 
	    WHERE teacher_id = $1 
	    AND classroom_id = $2
	);`, userID, classRoomID).Scan(&canCreate)

	return canCreate, err
}
