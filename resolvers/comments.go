package resolvers

import (
	"context"
	"github.com/les-cours/learning-service/types"
)

/*
add WS
*/
func (s *Server) CreateComment(ctx context.Context) error {

	err := s.MongoDB.AddComment(ctx, types.Comment{
		ID:         "1",
		DocumentID: "1",
		Message:    "test comment",
		Timestamp:  0,
		Owner:      "",
		IsEdited:   false,
		IsDeleted:  false,
	})

	return err
}

func (s *Server) GetComments(ctx context.Context, documentID string) ([]*types.Comment, error) {

	//s.MongoDB.GetComments(ctx, documentID)
	//err := s.MongoDB.AddComment(ctx, types.Comment{
	//	ID:         "1",
	//	DocumentID: "1",
	//	Message:    "test comment",
	//	Timestamp:  0,
	//	Owner:      "",
	//	IsEdited:   false,
	//	IsDeleted:  false,
	//})
	//
	//return err
	return nil, nil
}
