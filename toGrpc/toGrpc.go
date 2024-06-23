package toGrpc

import (
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/types"
)

func Room(room *types.Room) *learning.Room {
	return &learning.Room{
		Id:   room.ID,
		Name: room.Name,
		Owner: &learning.User{
			Id:        room.Teacher.ID,
			Username:  room.Teacher.Username,
			FirstName: room.Teacher.FirstName,
			LastName:  room.Teacher.LastName,
			Avatar:    room.Teacher.Avatar,
		},
		Users:    Users(room.Users),
		Messages: Messages(room.Messages),
	}
}

func Messages(messages []*types.Message) []*learning.Message {
	var grpcMessage = make([]*learning.Message, 0)
	for _, message := range messages {
		grpcMessage = append(grpcMessage, Message(message))
	}
	return grpcMessage
}

func Message(message *types.Message) *learning.Message {
	return &learning.Message{
		Id:        message.ID,
		RoomId:    message.RoomID,
		Content:   message.Message,
		Timestamp: message.Timestamp,
		IsTeacher: message.IsTeacher,
		User:      User(message.Owner),
	}
}

func Users(users []*types.User) []*learning.User {
	var grpcUsers = make([]*learning.User, 0)
	for _, user := range users {
		grpcUsers = append(grpcUsers, User(user))
	}
	return grpcUsers
}

func User(user *types.User) *learning.User {
	return &learning.User{
		Id:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Paid:      user.Paid,
	}
}
