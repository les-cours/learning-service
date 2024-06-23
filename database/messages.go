package database

import (
	"context"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/types"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

const MongoDBName = "karini"

func (db *MongoClient) AddComment(ctx context.Context, comment *types.Comment) error {

	_, err := db.MongoDB.Database(MongoDBName).Collection(CommentsCollections).InsertOne(ctx, types.Comment{
		ID:         comment.ID,
		UserID:     comment.UserID,
		RepliedTo:  comment.RepliedTo,
		Content:    comment.Content,
		DocumentID: comment.DocumentID,
		Timestamp:  comment.Timestamp,
		IsTeacher:  comment.IsTeacher,
	})

	return err

}

func (db *MongoClient) GetComments(ctx context.Context, id string, replied bool) ([]*types.Comment, error) {

	filter := bson.D{{"documentid", id}}
	if replied {
		filter = bson.D{{"repliedto", id}}
	}

	var comments []*types.Comment
	cur, err := db.MongoDB.Database(MongoDBName).Collection(CommentsCollections).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var comment *types.Comment
		cur.Decode(&comment)
		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *MongoClient) CreateRoom(ctx context.Context, room *learning.ClassRoom, user *learning.User) error {

	_, err := db.MongoDB.Database(MongoDBName).Collection(RoomsCollections).InsertOne(ctx, types.Room{
		ID:   room.ClassRoomID,
		Name: room.ArabicTitle,
		Teacher: &types.User{
			ID:        user.Id,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
		},
		Users:    []*types.User{},
		Messages: []*types.Message{},
	})

	return err
}

func (db *MongoClient) GetRoom(ctx context.Context, classroomID string) (*types.Room, error) {

	filter := bson.D{{"id", classroomID}}

	var room = new(types.Room)
	err := db.MongoDB.Database(MongoDBName).Collection(RoomsCollections).FindOne(ctx, filter).Decode(room)
	if err != nil {
		log.Println("mongo err messages.go:75 ", err)
		return nil, err
	}

	return room, nil
}

func (db *MongoClient) AddMessageToRoom(ctx context.Context, roomID string, message *types.Message) error {

	filter := bson.D{{"id", roomID}}
	update := bson.M{
		"$push": bson.M{
			"messages": &message,
		},
	}

	_, err := db.MongoDB.Database(MongoDBName).Collection(RoomsCollections).UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("err messages:86 ", err)
		return err
	}
	return nil
}

func (db *MongoClient) AddStudentToRoom(ctx context.Context, roomID string, user *types.User) error {

	filter := bson.D{{"id", roomID}}
	update := bson.M{
		"$push": bson.M{
			"users": &user,
		},
	}

	_, err := db.MongoDB.Database(MongoDBName).Collection(RoomsCollections).UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("err messages:110 ", err)
		return err
	}
	return nil
}

func (db *MongoClient) IsStudentBelongRoom(ctx context.Context, roomID string, studentID string) bool {

	// Filter document with $in operator
	filter := bson.D{
		{"id", roomID},
		{"users", bson.M{"$elemMatch": bson.M{"id": studentID}}},
	}

	// Count documents matching the filter
	count, err := db.MongoDB.Database(MongoDBName).Collection(RoomsCollections).CountDocuments(ctx, filter)
	if err != nil {
		return false
	}
	if count == 0 {
		return false
	}

	return true
}

func (db *MongoClient) ChangeStudentStatus(ctx context.Context, roomID string, studentID string, paid bool) error {

	// Update filter
	filter := bson.M{"id": roomID, "users.id": studentID}

	// Update operation
	update := bson.M{
		"$set": bson.M{"users.$.paid": paid},
	}

	// Perform the update
	_, err := db.MongoDB.Database(MongoDBName).Collection(RoomsCollections).UpdateOne(ctx, filter, update)

	return err

}
