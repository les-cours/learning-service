package resolvers

import (
	"context"
	"github.com/les-cours/learning-service/api/learning"
)

func (s *Server) CreateClassRoom(ctx context.Context, in *learning.CreateClassRoomRequest) (*learning.ClassRoom, error) {

	_, err := s.DB.Exec(`INSERT INTO classrooms (teacher_id,subject_id) VALUES ($1,$2)`, in.TeacherID, in.SubjectID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
