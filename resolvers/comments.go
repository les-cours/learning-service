package resolvers

import (
	"context"
	"github.com/google/uuid"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/api/users"
	"github.com/les-cours/learning-service/types"
	"time"
)

func (s *Server) CreateComment(ctx context.Context, in *learning.CreateCommentRequest) (*learning.OperationStatus, error) {

	commentID := uuid.NewString()

	comment := types.Comment{
		ID:         commentID,
		UserID:     in.UserID,
		Content:    in.Content,
		DocumentID: in.DocumentID,
		Timestamp:  time.Now(),
		IsTeacher:  in.IsTeacher,
	}

	if in.RepliedTo != "" {
		comment.RepliedTo = in.RepliedTo
	}

	err := s.MongoDB.AddComment(ctx, &comment)

	return &learning.OperationStatus{
		Success: true,
	}, err
}

func (s *Server) GetComments(ctx context.Context, in *learning.IDRequest) (*learning.Comments, error) {

	comments, err := s.MongoDB.GetComments(ctx, in.Id, false)

	if err != nil {
		return nil, err
	}

	var apiComments = new(learning.Comments)

	for _, comment := range comments {
		s.Logger.Info(comment.UserID)
		var role = "student"
		if comment.IsTeacher {
			role = "teacher"
		}
		user, err := s.Users.GetUserByID(ctx, &users.GetUserByIDRequest{
			AccountID: comment.UserID,
			UserRole:  role,
		})
		if err != nil {
			s.Logger.Error(err.Error())
			user = &users.User{
				Id:        "",
				Username:  "deleted account",
				FirstName: "deleted account",
				LastName:  "deleted account",
				Avatar:    "",
			}
		}
		apiComments.Comments = append(apiComments.Comments, &learning.Comment{
			Id: comment.ID,
			User: &learning.User{
				Id:        user.Id,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Avatar:    user.Avatar,
			},
			RepliedTo:  comment.RepliedTo,
			Content:    comment.Content,
			DocumentID: comment.DocumentID,
			Timestamp:  comment.Timestamp.Unix(),
			Edited:     comment.Edited,
			IsTeacher:  comment.IsTeacher,
		})
	}

	return apiComments, nil
}

func (s *Server) GetRepliedComments(ctx context.Context, in *learning.IDRequest) (*learning.Comments, error) {

	comments, err := s.MongoDB.GetComments(ctx, in.Id, true)

	if err != nil {
		return nil, err
	}

	var apiComments = new(learning.Comments)

	for _, comment := range comments {
		user, err := s.Users.GetUserByID(ctx, &users.GetUserByIDRequest{
			AccountID: comment.UserID,
		})
		if err != nil {
			user = &users.User{
				Id:        "",
				Username:  "deleted account",
				FirstName: "deleted account",
				LastName:  "deleted account",
				Avatar:    "",
			}
		}
		apiComments.Comments = append(apiComments.Comments, &learning.Comment{
			Id: comment.ID,
			User: &learning.User{
				Id:        user.Id,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Avatar:    user.Avatar,
			},
			RepliedTo:  comment.RepliedTo,
			Content:    comment.Content,
			DocumentID: comment.DocumentID,
			Timestamp:  comment.Timestamp.Unix(),
			Edited:     comment.Edited,
			IsTeacher:  comment.IsTeacher,
		})
	}

	return apiComments, nil
}
