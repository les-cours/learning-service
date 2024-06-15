package database

import (
	"context"
	"github.com/les-cours/learning-service/types"
	"go.mongodb.org/mongo-driver/bson"
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

func (db *MongoClient) GetComments(ctx context.Context, documentID string) ([]*types.Comment, error) {

	filter := bson.D{{"documentid", documentID}}
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

//
//// DeleteMessage Delete specific message from room of an account, it returns a mongo db error
//// if any error occurs or ErrNoOwnerMatch the owner is not the same or ErrMessageDeleted
//// if the message is deleted
//func (db *MongoClient) DeleteMessage(accountID, roomID, ownerID, messageID string) error {
//	result, err := db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
//		"$and": []bson.M{
//			{"id": messageID},
//			{"roomID": roomID},
//			{"owner.id": ownerID},
//		}}, bson.M{
//		"$set": bson.M{
//			"isDeleted": true,
//		},
//	})
//	if result.MatchedCount == 0 {
//		return ErrMessageDoesNotExist
//	}
//
//	return err
//}
//
//func (db *MongoClient) DeleteRoomMessages(accountID, roomID string) error {
//	_, err := db.MongoDB.Database(MongoDBName).Collection(accountID).DeleteMany(context.Background(), bson.M{
//		"$and": []bson.M{
//			{"roomID": roomID},
//		}})
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// DeleteComment delete comment
//func (db *MongoClient) DeleteComment(accountID, roomID, ownerID, commentID string) error {
//	result, err := db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
//		"$and": []bson.M{
//			{"roomID": roomID},
//			{"comments.id": commentID},
//			{"comments.owner.id": ownerID},
//		}}, bson.M{
//		"$set": bson.M{
//			"comments.$.isDeleted": true,
//		},
//	})
//	if err != nil {
//		return err
//	}
//
//	if result.MatchedCount == 0 {
//		result, err := db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
//			"$and": []bson.M{
//				{"roomID": roomID},
//				{"id": commentID},
//				{"owner.id": ownerID},
//			}}, bson.M{
//			"$set": bson.M{
//				"isDeleted": true,
//			},
//		})
//		if err != nil {
//			return err
//		}
//
//		if result.MatchedCount == 0 {
//			return ErrCommentNotFound
//		}
//	}
//
//	return nil
//}
//
//// EditComment edit comment
//func (db *MongoClient) EditComment(accountID, ownerID string, input types.EditMessageInput) (*types.Message, error) {
//	message, err := db.GetMessageOfComment(accountID, input.RoomID, input.MessageID)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, c := range message.Comments {
//		if c.ID == input.MessageID {
//			if c.Owner.ID != ownerID {
//				return nil, ErrNoOwnerMatch
//			}
//
//			c.IsEdited = true
//			editedTime := time.Now().UnixMilli()
//			c.EditHistory = append(c.EditHistory, types.MessageHistory{
//				ID:        primitive.NewObjectID(),
//				Message:   c.Message,
//				Timestamp: editedTime,
//			})
//			_, err = db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
//				"$and": []bson.M{
//					{"id": message.ID},
//					{"roomID": input.RoomID},
//					{"comments.id": input.MessageID},
//				}}, bson.M{
//				"$set": bson.M{
//					"comments.$.message":       input.Message,
//					"comments.$.isEdited":      true,
//					"comments.$.editTimestamp": editedTime,
//				},
//				"$push": bson.M{
//					"comments.$.editHistory": types.MessageHistory{
//						ID:        primitive.NewObjectID(),
//						Message:   c.Message,
//						Timestamp: c.Timestamp,
//					},
//				},
//			})
//			if err != nil {
//				return nil, err
//			}
//
//			break
//		}
//	}
//
//	return message, nil
//}
//
//// EditMessage Edit specific message from a room of an account, it returns a mongo db error
//// if any error occurs, ErrNoOwnerMatch if the owner ID is not the same message owner or
//// ErrMessageDoesNotExist if message does not exist
//func (db *MongoClient) EditMessage(accountID, ownerID string, input types.EditMessageInput) (*types.Message, error) {
//	msg, err := db.GetMessage(accountID, input.RoomID, input.MessageID)
//	if err != nil {
//		return nil, err
//	}
//
//	if msg.Owner.ID != ownerID {
//		return nil, ErrNoOwnerMatch
//	}
//
//	_, err = db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
//		"$and": []bson.M{
//			{"id": input.MessageID},
//			{"roomID": input.RoomID},
//		}}, bson.M{
//		"$set": bson.M{
//			"message":       input.Message,
//			"isEdited":      true,
//			"editTimestamp": time.Now().UnixMilli(),
//		},
//		"$push": bson.M{
//			"editHistory": types.MessageHistory{
//				ID:        primitive.NewObjectID(),
//				Message:   msg.Message,
//				Timestamp: msg.EditedTimestamp,
//			},
//		},
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return db.GetMessage(accountID, input.RoomID, input.MessageID)
//}
//
//// ReplyMessage Reply for a specific message from a room of an account.
//// it returns a mongo db error if any error occurs.
//func (db *MongoClient) ReplyMessage(accountID, roomID, messageID, message string, timestamp int64, owner types.Owner) (*types.Message, error) {
//	err := db.checkMember(accountID, owner.ID, roomID)
//	if err != nil {
//		return nil, err
//	}
//
//	msg, err := db.GetMessage(accountID, roomID, messageID)
//	if err != nil {
//		return nil, err
//	}
//
//	if msg.IsDeleted {
//		return nil, ErrMessageDeleted
//	}
//
//	insert, err := db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(context.Background(), types.MongoMessage{
//		RoomID:          roomID,
//		Message:         message,
//		ReplyMessageID:  &messageID,
//		Timestamp:       timestamp,
//		EditedTimestamp: timestamp,
//		Owner:           owner,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	ID, _ := insert.InsertedID.(primitive.ObjectID)
//	return db.GetMessage(accountID, roomID, ID.Hex())
//}
//
//func (db *MongoClient) AddBotMessage(accountID, roomID string, botID int, message string, bubbles []string) (*types.Message, error) {
//	exist, err := db.RoomExists(accountID, roomID)
//	if err != nil {
//		return nil, err
//	}
//
//	if !exist {
//		return nil, ErrCanNotAddMessageToUnexistedRoom
//	}
//
//	var now = time.Now().UnixMilli()
//
//	insert, err := db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(context.Background(), types.MongoMessage{
//		RoomID:          roomID,
//		Message:         message,
//		Timestamp:       now,
//		EditedTimestamp: now,
//		Bubbles:         bubbles,
//		Owner:           types.Owner{ID: strconv.Itoa(botID), Kind: "bot"},
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	messageID, _ := insert.InsertedID.(primitive.ObjectID)
//	return db.GetMessage(accountID, roomID, messageID.Hex())
//}
//
//func (db *MongoClient) checkMember(accountID, ownerID, roomID string) error {
//	exist, err := db.IsRoomMember(accountID, ownerID, roomID)
//	if err != nil && err != ErrRoomDoesNotExist {
//		return err
//	}
//
//	if err == ErrRoomDoesNotExist {
//		return ErrCanNotAddMessageToUnexistedRoom
//	}
//
//	if !exist {
//		return ErrCanNotAddMessageFromNoRoomMember
//	}
//
//	return nil
//}
//
//func (db *MongoClient) GetMessagesCount(accountID, roomID string) (int64, error) {
//	var filter = bson.D{{Key: "roomID", Value: roomID}}
//	var count, err = db.MongoDB.Database(MongoDBName).Collection(accountID).CountDocuments(context.Background(), filter)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
//
//func (db *MongoClient) UpdateUserLastReadTimestamp(ctx context.Context, accountID, userID, roomID string, timestamp int64) error {
//	var newTimestamp int64
//	var err error
//	if timestamp == 0 {
//		var message types.MongoMessageWithID
//		var options = options.FindOne()
//		options.SetSort(bson.M{"_id": -1}).SetProjection(bson.D{
//			primitive.E{Key: "timestamp", Value: 1},
//		})
//
//		err = db.MongoDB.Database(MongoDBName).
//			Collection(accountID).FindOne(ctx, bson.M{"roomID": roomID}, options).
//			Decode(&message)
//
//		if err == mongo.ErrNoDocuments {
//			newTimestamp = 0
//		} else {
//			if err != nil {
//				return err
//			}
//			newTimestamp = message.Timestamp
//		}
//	}
//
//	if timestamp != 0 {
//		newTimestamp = timestamp
//	}
//
//	_, err = db.MongoDB.Database(MongoDBName).
//		Collection(RoomsCollection).
//		UpdateOne(ctx, bson.M{
//			"roomID": roomID,
//			"users":  bson.M{"$elemMatch": bson.M{"id": userID}},
//		}, bson.M{"$set": bson.M{"users.$.lastReadTimestamp": newTimestamp}})
//
//	return err
//}
//
//func (db *MongoClient) GetForm(ctx context.Context, id string) (*types.MongoForm, error) {
//	var mongoForm types.MongoForm
//	var coll = db.MongoDB.Database(MongoDBName).Collection(FormsCollections)
//	var formID, err = primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return nil, err
//	}
//
//	err = coll.FindOne(ctx, bson.M{"_id": formID}).Decode(&mongoForm)
//	if err != nil {
//		return nil, err
//	}
//
//	return &mongoForm, nil
//}
//
//func (db *MongoClient) AddForm(ctx context.Context, form *types.MongoFormWithoutID) (string, error) {
//	var coll = db.MongoDB.Database(MongoDBName).Collection(FormsCollections)
//	var res, err = coll.InsertOne(ctx, form)
//	if err != nil {
//		return "", err
//	}
//
//	return res.InsertedID.(primitive.ObjectID).Hex(), nil
//}
//
//func (db *MongoClient) UpdateForm(ctx context.Context, form *types.MongoForm) error {
//	var filter = bson.D{{Key: "_id", Value: form.ID}}
//	var coll = db.MongoDB.Database(MongoDBName).Collection(FormsCollections)
//	var _, err = coll.ReplaceOne(context.TODO(), filter, form.MongoFormWithoutID)
//
//	return err
//}
//
//
//*/
