package resolvers

import (
	"context"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/api/orgs"
	"github.com/les-cours/learning-service/api/users"
	"github.com/les-cours/learning-service/toGrpc"
	"github.com/les-cours/learning-service/types"
	"github.com/les-cours/learning-service/utils"
	"time"
)

func (s *Server) AddStudentToChatRoom(ctx context.Context, in *learning.IDRequest) (*learning.OperationStatus, error) {

	if !s.CanAccessToClassRoom(in.UserID, in.Id) {
		s.Logger.Error(ErrPermission.Error())
		return nil, ErrPermission
	}

	//test if  exist already then continue else add student.
	if s.MongoDB.IsStudentBelongRoom(ctx, in.Id, in.UserID) {
		return &learning.OperationStatus{Success: true}, nil
	}

	user, err := s.Users.GetUserByID(ctx, &users.GetUserByIDRequest{
		AccountID: in.UserID,
		UserRole:  "student",
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	err = s.MongoDB.AddStudentToRoom(ctx, in.Id, &types.User{
		ID:        user.Id,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	})

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	return &learning.OperationStatus{Success: true}, nil
}

func (s *Server) AddMessageToChatRoom(ctx context.Context, in *learning.AddMessage) (*learning.OperationStatus, error) {

	//if s.CanAccessToClassRoom(in.UserID, in.Id) {
	//	return nil, ErrPermission
	//}

	role := "student"
	if in.IsTeacher {
		role = "teacher"
	}
	user, err := s.Users.GetUserByID(ctx, &users.GetUserByIDRequest{
		AccountID: in.UserID,
		UserRole:  role,
	})

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	err = s.MongoDB.AddMessageToRoom(ctx, in.RoomID, &types.Message{
		ID:        utils.GenerateUUIDString(),
		RoomID:    in.RoomID,
		Message:   in.Content,
		Timestamp: time.Now().String(),
		IsTeacher: in.IsTeacher,
		Owner: &types.User{
			ID:        user.Id,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
		},
	})

	if err != nil {
		return nil, err
	}

	return &learning.OperationStatus{Success: true}, nil
}

func (s *Server) GetChatRoom(ctx context.Context, in *learning.IDRequest) (*learning.Room, error) {

	//if s.CanAccessToClassRoom(in.UserID, in.Id) {
	//	return nil, ErrPermission
	//}

	room, err := s.MongoDB.GetRoom(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return toGrpc.Room(room), nil

}

func (s *Server) GetMyChatRoom(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	//if s.CanAccessToClassRoom(in.UserID, in.Id) {
	//	return nil, ErrPermission
	//}

	var gradID string
	s.DB.QueryRow(`SELECT grade_id FROM students where student_id = $1`, in.Id).Scan(&gradID)
	res, err := s.Orgs.GetSubjectsByGrad(ctx, &orgs.GetSubjectByGradRequest{
		GradID: gradID,
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	s.Logger.Info("gradID" + gradID)

	var rooms = make([]*learning.ClassRoom, 0)
	for _, subject := range res.Subjects {
		s.Logger.Info("subject : " + subject.SubjectID)
		room, err := s.GetMyClassRooms(ctx, &learning.IDRequest{
			Id:     subject.SubjectID,
			UserID: in.UserID,
		})
		if err != nil {
			return nil, err
		}

		if room != nil {
			rooms = append(rooms, room.Classrooms...)
		}

	}

	return &learning.ClassRooms{
		Classrooms: rooms,
	}, nil

}
