package resolvers

import (
	"context"
	"errors"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/database"
)

func (s *Server) CreateComment(ctx context.Context, in *learning.CreateVideoRequest) (*learning.Document, error) {

	err := s.MongoDB.AddComment(ctx, database.Comment{
		ID:         "1",
		DocumentID: "1",
		Message:    "test comment",
		Timestamp:  0,
		Owner:      "",
		IsEdited:   false,
		IsDeleted:  false,
	})

	if err != nil {
		return nil, err
	}

	return nil, errors.New("err")
}
