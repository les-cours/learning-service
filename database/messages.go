package database

import (
	"context"
	"github.com/les-cours/learning-service/types"
)

const MongoDBName = "karini"

func (db *MongoClient) AddComment(ctx context.Context, comment types.Comment) error {

	_, err := db.MongoDB.Database(MongoDBName).Collection(CommentsCollections).InsertOne(ctx, types.Comment{
		ID:         comment.ID,
		DocumentID: comment.DocumentID,
		Message:    comment.Message,
		Timestamp:  0,
		Owner:      "",
		IsEdited:   false,
		IsDeleted:  false,
	})

	return err

}

//func (db *MongoClient) GetComments(ctx context.Context, documentID string) ([]*types.Comment, error) {
//
//	res := types.RoomMessages{
//		Messages: []*types.Message{},
//	}
//	context := context.Background()
//	after := len(pagination.After) > 0
//	before := len(pagination.Before) > 0
//	if after && before {
//		return nil, errors.New("how is that?")
//	}
//	firstTime := !after && !before
//	filter := []bson.M{
//		{"roomID": roomID},
//		{"isDeleted": bson.M{"$exists": false}},
//	}
//	sort := bson.M{"timestamp": -1}
//	if !firstTime {
//		if after {
//			id, err := primitive.ObjectIDFromHex(pagination.After)
//			if err != nil {
//				return nil, errors.New("bad cursor")
//			}
//			filter = append(filter, bson.M{"_id": primitive.M{"$gt": id}})
//			res.Filter.After = pagination.After
//			res.Filter.Limit = pagination.Limit
//			sort = bson.M{"timestamp": 1}
//		} else if before {
//			id, err := primitive.ObjectIDFromHex(pagination.Before)
//			if err != nil {
//				return nil, errors.New("bad cursor")
//			}
//			filter = append(filter, bson.M{"_id": primitive.M{"$lt": id}})
//			res.Filter.Before = pagination.Before
//			res.Filter.Limit = pagination.Limit
//		}
//
//	} else {
//		count, err := db.MongoDB.Database(MongoDBName).Collection(accountID).CountDocuments(context, bson.M{"$and": filter})
//		if err != nil {
//			return nil, err
//		}
//		res.Total = count
//	}
//	cursor, err := db.MongoDB.Database(MongoDBName).Collection(accountID).Aggregate(context, []bson.M{
//		{"$match": bson.M{
//			"$and": filter,
//		}},
//		{"$sort": sort},
//		{"$limit": pagination.Limit},
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	for cursor.Next(context) {
//		var mongoMessage types.MongoMessageWithID
//		if err := cursor.Decode(&mongoMessage); err != nil {
//			continue
//		}
//		for _, c := range mongoMessage.Comments {
//			if c.IsDeleted {
//				continue
//			}
//			var message = types.Message{
//				ID:                 c.ID,
//				RoomID:             mongoMessage.RoomID,
//				Message:            c.Message,
//				Timestamp:          c.Timestamp,
//				FormattedTimestamp: c.FormattedTimestamp,
//				EditedTimestamp:    c.EditedTimestamp,
//				IsEdited:           c.IsEdited,
//				IsDeleted:          c.IsDeleted,
//				Owner:              mongoMessage.Owner,
//				Bubbles:            []string{},
//				Attachments:        c.Attachments,
//				EditHistory:        c.EditHistory,
//				IsNotification:     false,
//				IsComment:          true,
//				TargetMsgID:        mongoMessage.ID,
//			}
//			res.Messages = append(res.Messages, &message)
//		}
//		var message = types.Message{
//			ID:              mongoMessage.ID,
//			RoomID:          mongoMessage.RoomID,
//			Message:         mongoMessage.Message,
//			Timestamp:       mongoMessage.Timestamp,
//			EditedTimestamp: mongoMessage.EditedTimestamp,
//			IsEdited:        mongoMessage.IsEdited,
//			IsDeleted:       mongoMessage.IsDeleted,
//			Owner:           mongoMessage.Owner,
//			Bubbles:         mongoMessage.Bubbles,
//			Attachments:     mongoMessage.Attachments,
//			EditHistory:     mongoMessage.EditHistory,
//			IsNotification:  mongoMessage.IsNotification,
//			IsComment:       mongoMessage.IsComment,
//			IsFile:          mongoMessage.IsFile,
//			Cursor:          mongoMessage.DBID,
//		}
//
//		if mongoMessage.ReplyMessageID != nil {
//			var ID, _ = primitive.ObjectIDFromHex(*mongoMessage.ReplyMessageID)
//			var mongoMessage types.MongoMessage
//
//			err := db.MongoDB.Database(MongoDBName).Collection(accountID).FindOne(context, bson.M{
//				"$and": []bson.M{
//					{"_id": ID},
//					{"roomID": roomID},
//				}}).Decode(&mongoMessage)
//
//			if err == nil {
//				if mongoMessage.IsDeleted {
//					mongoMessage.Message = "Deleted message"
//				}
//
//				if mongoMessage.Message == "" && len(mongoMessage.Attachments) != 0 {
//					mongoMessage.Message = "Attachment"
//				}
//
//				message.Replied = &types.RepliedMessage{
//					ID:      mongoMessage.ID,
//					Message: mongoMessage.Message,
//					Owner:   mongoMessage.Owner,
//				}
//			}
//		}
//
//		res.Messages = append(res.Messages, &message)
//	}
//	return &res, nil
//
//
//	var mongoMessage types.MongoMessage
//	var err = db.MongoDB.Database(MongoDBName).
//		Collection(CommentsCollections).
//		FindOne(context.Background(), bson.M{
//			"$and": []bson.M{
//				{"documentID": documentID},
//			}}).
//		Decode(&mongoMessage)
//	if err != nil && err != mongo.ErrNoDocuments {
//		return nil, err
//	}
//
//	if err != nil {
//		return nil, ErrMessageDoesNotExist
//	}
//
//	var message = types.Comment{
//		ID:                 mongoMessage.ID,
//		RoomID:             mongoMessage.RoomID,
//		Message:            mongoMessage.Message,
//		Timestamp:          mongoMessage.Timestamp,
//		FormattedTimestamp: mongoMessage.FormattedTimestamp,
//		EditedTimestamp:    mongoMessage.EditedTimestamp,
//		IsComment:          mongoMessage.IsComment,
//		IsEdited:           mongoMessage.IsEdited,
//		IsDeleted:          mongoMessage.IsDeleted,
//		Owner:              mongoMessage.Owner,
//		Bubbles:            mongoMessage.Bubbles,
//		Attachments:        mongoMessage.Attachments,
//		EditHistory:        mongoMessage.EditHistory,
//	}
//	if mongoMessage.ReplyMessageID != nil {
//		var repliedMessage types.MongoMessage
//		var ID, _ = primitive.ObjectIDFromHex(*mongoMessage.ReplyMessageID)
//		var err = db.MongoDB.Database(MongoDBName).
//			Collection(accountID).
//			FindOne(context.Background(), bson.M{
//				"$and": []bson.M{
//					{"_id": ID},
//					{"roomID": roomID},
//				}}).
//			Decode(&repliedMessage)
//		if err != nil {
//			return &message, nil
//		}
//
//		if repliedMessage.IsDeleted {
//			repliedMessage.Message = "Deleted message"
//		}
//
//		if repliedMessage.Message == "" && len(mongoMessage.Attachments) != 0 {
//			repliedMessage.Message = "Attachment"
//		}
//
//		message.Replied = &types.RepliedMessage{
//			ID:      repliedMessage.ID,
//			Message: repliedMessage.Message,
//			Owner:   repliedMessage.Owner,
//		}
//	}
//
//	return &message, nil
//}

/*
// AddMessage Add message of a room, it returns a mongodb error if any error occured or it return
// ErrCanNotAddMessageToUnexistedRoom if the room does not exist or ErrCanNotAddMessageFromNoRoomMember
// if the message owner is not member of the room
func (db *MongoClient) AddMessage(ctx context.Context, accountID string, message Message) error {

	// ! TODO: should be deleted at some point
	switch {
	case message.FormattedTimestamp == "":
		message.FormattedTimestamp = utils.TimeFormattedWithTimezone()

	case message.Timestamp == 0:
		message.Timestamp = utils.TimestampMilli()
	}

	if message.Replied != nil {
		_, err := db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(ctx, types.MongoMessage{
			ID:                 message.ID,
			RoomID:             message.RoomID,
			Message:            message.Message,
			Timestamp:          message.Timestamp,
			FormattedTimestamp: message.FormattedTimestamp,
			Owner:              message.Owner,
			IsForm:             message.IsForm,
			IsNotification:     message.IsNotification,
			IsFile:             message.IsFile,
			IsPrivate:          message.IsPrivate,
			FormID:             message.FormID,
			Attachments:        message.Attachments,
			Replied:            message.Replied,
			ReplyMessageID:     &message.Replied.ID,
		})

		return err
	}

	// if message.IsNotification {
	// 	_, err := db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(ctx, types.MongoMessage{
	// 		ID:                 message.ID,
	// 		RoomID:             message.RoomID,
	// 		IsNotification:     true,
	// 		Message:            message.Message,
	// 		Timestamp:          message.Timestamp,
	// 		FormattedTimestamp: message.FormattedTimestamp,
	// 		Owner:              message.Owner,
	// 	})
	//
	// 	return err
	// }
	// if message.IsFile {
	// 	_, err := db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(ctx, types.MongoMessage{
	// 		ID:                 message.ID,
	// 		RoomID:             message.RoomID,
	// 		IsFile:             true,
	// 		Message:            message.Message,
	// 		Timestamp:          message.Timestamp,
	// 		FormattedTimestamp: message.FormattedTimestamp,
	// 		Attachments:        message.Attachments,
	// 		Owner:              message.Owner,
	// 	})
	//
	// 	return err
	// }
	//
	// if message.IsForm {
	// 	_, err := db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(ctx, types.MongoMessage{
	// 		ID:                 message.ID,
	// 		RoomID:             message.RoomID,
	// 		IsForm:             true,
	// 		FormID:             message.FormID,
	// 		Message:            message.Message,
	// 		Timestamp:          message.Timestamp,
	// 		FormattedTimestamp: message.FormattedTimestamp,
	// 		Attachments:        message.Attachments,
	// 		Owner:              message.Owner,
	// 	})
	//
	// 	return err
	// }

	var _, err = db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(ctx, types.MongoMessage{
		ID:                 message.ID,
		RoomID:             message.RoomID,
		Message:            message.Message,
		Timestamp:          message.Timestamp,
		FormattedTimestamp: message.FormattedTimestamp,
		Owner:              message.Owner,
		IsForm:             message.IsForm,
		IsNotification:     message.IsNotification,
		IsFile:             message.IsFile,
		IsPrivate:          message.IsPrivate,
		FormID:             message.FormID,
		Attachments:        message.Attachments,
	})

	return err
}

// GetMessageOfComment returns message of a comment
func (db *MongoClient) GetMessageOfComment(accountID, roomID, commentID string) (*types.Message, error) {
	var mongoMessage types.MongoMessage

	err := db.MongoDB.Database(MongoDBName).Collection(accountID).FindOne(context.Background(), bson.M{
		"$and": []bson.M{
			{"comments.id": commentID},
			{"roomID": roomID},
		}}).Decode(&mongoMessage)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err != nil {
		return nil, ErrMessageDoesNotExist
	}

	var comments []*types.Comment
	for _, c := range mongoMessage.Comments {
		comments = append(comments, &types.Comment{
			ID:              c.ID,
			Message:         c.Message,
			Timestamp:       c.Timestamp,
			Owner:           c.Owner,
			EditedTimestamp: c.EditedTimestamp,
			IsEdited:        c.IsEdited,
			IsDeleted:       c.IsDeleted,
		})
	}
	var message = types.Message{
		ID:              mongoMessage.ID,
		RoomID:          mongoMessage.RoomID,
		Message:         mongoMessage.Message,
		Timestamp:       mongoMessage.Timestamp,
		EditedTimestamp: mongoMessage.EditedTimestamp,
		IsComment:       mongoMessage.IsComment,
		IsEdited:        mongoMessage.IsEdited,
		IsDeleted:       mongoMessage.IsDeleted,
		Owner:           mongoMessage.Owner,
		EditHistory:     mongoMessage.EditHistory,
		Comments:        comments,
	}
	return &message, nil
}

// GetMessage Get message from specific room of an account, it returns a mongo db error if any error
// occurs or ErrMessageDoesNotExist if the message does not exist
func (db *MongoClient) GetMessage(accountID, roomID, messageID string) (*types.Message, error) {
	var mongoMessage types.MongoMessage
	var err = db.MongoDB.Database(MongoDBName).
		Collection(accountID).
		FindOne(context.Background(), bson.M{
			"$and": []bson.M{
				{"id": messageID},
				{"roomID": roomID},
			}}).
		Decode(&mongoMessage)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err != nil {
		return nil, ErrMessageDoesNotExist
	}

	var message = types.Message{
		ID:                 mongoMessage.ID,
		RoomID:             mongoMessage.RoomID,
		Message:            mongoMessage.Message,
		Timestamp:          mongoMessage.Timestamp,
		FormattedTimestamp: mongoMessage.FormattedTimestamp,
		EditedTimestamp:    mongoMessage.EditedTimestamp,
		IsComment:          mongoMessage.IsComment,
		IsEdited:           mongoMessage.IsEdited,
		IsDeleted:          mongoMessage.IsDeleted,
		Owner:              mongoMessage.Owner,
		Bubbles:            mongoMessage.Bubbles,
		Attachments:        mongoMessage.Attachments,
		EditHistory:        mongoMessage.EditHistory,
	}
	if mongoMessage.ReplyMessageID != nil {
		var repliedMessage types.MongoMessage
		var ID, _ = primitive.ObjectIDFromHex(*mongoMessage.ReplyMessageID)
		var err = db.MongoDB.Database(MongoDBName).
			Collection(accountID).
			FindOne(context.Background(), bson.M{
				"$and": []bson.M{
					{"_id": ID},
					{"roomID": roomID},
				}}).
			Decode(&repliedMessage)
		if err != nil {
			return &message, nil
		}

		if repliedMessage.IsDeleted {
			repliedMessage.Message = "Deleted message"
		}

		if repliedMessage.Message == "" && len(mongoMessage.Attachments) != 0 {
			repliedMessage.Message = "Attachment"
		}

		message.Replied = &types.RepliedMessage{
			ID:      repliedMessage.ID,
			Message: repliedMessage.Message,
			Owner:   repliedMessage.Owner,
		}
	}

	return &message, nil
}

type PaginationFilter struct {
	Limit  int    `json:"limit" validate:"required"`
	After  string `json:"after" validate:"required_without_all=Before"`
	Before string `json:"before" validate:"required_without_all=After"`
}

// GetMessages Get message from specific room
func (db *MongoClient) GetMessages(accountID, roomID string, pagination PaginationFilter) (*types.RoomMessages, error) {
	res := types.RoomMessages{
		Messages: []*types.Message{},
	}
	context := context.Background()
	after := len(pagination.After) > 0
	before := len(pagination.Before) > 0
	if after && before {
		return nil, errors.New("how is that?")
	}
	firstTime := !after && !before
	filter := []bson.M{
		{"roomID": roomID},
		{"isDeleted": bson.M{"$exists": false}},
	}
	sort := bson.M{"timestamp": -1}
	if !firstTime {
		if after {
			id, err := primitive.ObjectIDFromHex(pagination.After)
			if err != nil {
				return nil, errors.New("bad cursor")
			}
			filter = append(filter, bson.M{"_id": primitive.M{"$gt": id}})
			res.Filter.After = pagination.After
			res.Filter.Limit = pagination.Limit
			sort = bson.M{"timestamp": 1}
		} else if before {
			id, err := primitive.ObjectIDFromHex(pagination.Before)
			if err != nil {
				return nil, errors.New("bad cursor")
			}
			filter = append(filter, bson.M{"_id": primitive.M{"$lt": id}})
			res.Filter.Before = pagination.Before
			res.Filter.Limit = pagination.Limit
		}

	} else {
		count, err := db.MongoDB.Database(MongoDBName).Collection(accountID).CountDocuments(context, bson.M{"$and": filter})
		if err != nil {
			return nil, err
		}
		res.Total = count
	}
	cursor, err := db.MongoDB.Database(MongoDBName).Collection(accountID).Aggregate(context, []bson.M{
		{"$match": bson.M{
			"$and": filter,
		}},
		{"$sort": sort},
		{"$limit": pagination.Limit},
	})

	if err != nil {
		return nil, err
	}

	for cursor.Next(context) {
		var mongoMessage types.MongoMessageWithID
		if err := cursor.Decode(&mongoMessage); err != nil {
			continue
		}
		for _, c := range mongoMessage.Comments {
			if c.IsDeleted {
				continue
			}
			var message = types.Message{
				ID:                 c.ID,
				RoomID:             mongoMessage.RoomID,
				Message:            c.Message,
				Timestamp:          c.Timestamp,
				FormattedTimestamp: c.FormattedTimestamp,
				EditedTimestamp:    c.EditedTimestamp,
				IsEdited:           c.IsEdited,
				IsDeleted:          c.IsDeleted,
				Owner:              mongoMessage.Owner,
				Bubbles:            []string{},
				Attachments:        c.Attachments,
				EditHistory:        c.EditHistory,
				IsNotification:     false,
				IsComment:          true,
				TargetMsgID:        mongoMessage.ID,
			}
			res.Messages = append(res.Messages, &message)
		}
		var message = types.Message{
			ID:              mongoMessage.ID,
			RoomID:          mongoMessage.RoomID,
			Message:         mongoMessage.Message,
			Timestamp:       mongoMessage.Timestamp,
			EditedTimestamp: mongoMessage.EditedTimestamp,
			IsEdited:        mongoMessage.IsEdited,
			IsDeleted:       mongoMessage.IsDeleted,
			Owner:           mongoMessage.Owner,
			Bubbles:         mongoMessage.Bubbles,
			Attachments:     mongoMessage.Attachments,
			EditHistory:     mongoMessage.EditHistory,
			IsNotification:  mongoMessage.IsNotification,
			IsComment:       mongoMessage.IsComment,
			IsFile:          mongoMessage.IsFile,
			Cursor:          mongoMessage.DBID,
		}

		if mongoMessage.ReplyMessageID != nil {
			var ID, _ = primitive.ObjectIDFromHex(*mongoMessage.ReplyMessageID)
			var mongoMessage types.MongoMessage

			err := db.MongoDB.Database(MongoDBName).Collection(accountID).FindOne(context, bson.M{
				"$and": []bson.M{
					{"_id": ID},
					{"roomID": roomID},
				}}).Decode(&mongoMessage)

			if err == nil {
				if mongoMessage.IsDeleted {
					mongoMessage.Message = "Deleted message"
				}

				if mongoMessage.Message == "" && len(mongoMessage.Attachments) != 0 {
					mongoMessage.Message = "Attachment"
				}

				message.Replied = &types.RepliedMessage{
					ID:      mongoMessage.ID,
					Message: mongoMessage.Message,
					Owner:   mongoMessage.Owner,
				}
			}
		}

		res.Messages = append(res.Messages, &message)
	}
	return &res, nil
}

// DeleteMessage Delete specific message from room of an account, it returns a mongo db error
// if any error occurs or ErrNoOwnerMatch the owner is not the same or ErrMessageDeleted
// if the message is deleted
func (db *MongoClient) DeleteMessage(accountID, roomID, ownerID, messageID string) error {
	result, err := db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
		"$and": []bson.M{
			{"id": messageID},
			{"roomID": roomID},
			{"owner.id": ownerID},
		}}, bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	})
	if result.MatchedCount == 0 {
		return ErrMessageDoesNotExist
	}

	return err
}

func (db *MongoClient) DeleteRoomMessages(accountID, roomID string) error {
	_, err := db.MongoDB.Database(MongoDBName).Collection(accountID).DeleteMany(context.Background(), bson.M{
		"$and": []bson.M{
			{"roomID": roomID},
		}})

	if err != nil {
		return err
	}

	return nil
}

// DeleteComment delete comment
func (db *MongoClient) DeleteComment(accountID, roomID, ownerID, commentID string) error {
	result, err := db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
		"$and": []bson.M{
			{"roomID": roomID},
			{"comments.id": commentID},
			{"comments.owner.id": ownerID},
		}}, bson.M{
		"$set": bson.M{
			"comments.$.isDeleted": true,
		},
	})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		result, err := db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
			"$and": []bson.M{
				{"roomID": roomID},
				{"id": commentID},
				{"owner.id": ownerID},
			}}, bson.M{
			"$set": bson.M{
				"isDeleted": true,
			},
		})
		if err != nil {
			return err
		}

		if result.MatchedCount == 0 {
			return ErrCommentNotFound
		}
	}

	return nil
}

// EditComment edit comment
func (db *MongoClient) EditComment(accountID, ownerID string, input types.EditMessageInput) (*types.Message, error) {
	message, err := db.GetMessageOfComment(accountID, input.RoomID, input.MessageID)
	if err != nil {
		return nil, err
	}

	for _, c := range message.Comments {
		if c.ID == input.MessageID {
			if c.Owner.ID != ownerID {
				return nil, ErrNoOwnerMatch
			}

			c.IsEdited = true
			editedTime := time.Now().UnixMilli()
			c.EditHistory = append(c.EditHistory, types.MessageHistory{
				ID:        primitive.NewObjectID(),
				Message:   c.Message,
				Timestamp: editedTime,
			})
			_, err = db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
				"$and": []bson.M{
					{"id": message.ID},
					{"roomID": input.RoomID},
					{"comments.id": input.MessageID},
				}}, bson.M{
				"$set": bson.M{
					"comments.$.message":       input.Message,
					"comments.$.isEdited":      true,
					"comments.$.editTimestamp": editedTime,
				},
				"$push": bson.M{
					"comments.$.editHistory": types.MessageHistory{
						ID:        primitive.NewObjectID(),
						Message:   c.Message,
						Timestamp: c.Timestamp,
					},
				},
			})
			if err != nil {
				return nil, err
			}

			break
		}
	}

	return message, nil
}

// EditMessage Edit specific message from a room of an account, it returns a mongo db error
// if any error occurs, ErrNoOwnerMatch if the owner ID is not the same message owner or
// ErrMessageDoesNotExist if message does not exist
func (db *MongoClient) EditMessage(accountID, ownerID string, input types.EditMessageInput) (*types.Message, error) {
	msg, err := db.GetMessage(accountID, input.RoomID, input.MessageID)
	if err != nil {
		return nil, err
	}

	if msg.Owner.ID != ownerID {
		return nil, ErrNoOwnerMatch
	}

	_, err = db.MongoDB.Database(MongoDBName).Collection(accountID).UpdateOne(context.Background(), bson.M{
		"$and": []bson.M{
			{"id": input.MessageID},
			{"roomID": input.RoomID},
		}}, bson.M{
		"$set": bson.M{
			"message":       input.Message,
			"isEdited":      true,
			"editTimestamp": time.Now().UnixMilli(),
		},
		"$push": bson.M{
			"editHistory": types.MessageHistory{
				ID:        primitive.NewObjectID(),
				Message:   msg.Message,
				Timestamp: msg.EditedTimestamp,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return db.GetMessage(accountID, input.RoomID, input.MessageID)
}

// ReplyMessage Reply for a specific message from a room of an account.
// it returns a mongo db error if any error occurs.
func (db *MongoClient) ReplyMessage(accountID, roomID, messageID, message string, timestamp int64, owner types.Owner) (*types.Message, error) {
	err := db.checkMember(accountID, owner.ID, roomID)
	if err != nil {
		return nil, err
	}

	msg, err := db.GetMessage(accountID, roomID, messageID)
	if err != nil {
		return nil, err
	}

	if msg.IsDeleted {
		return nil, ErrMessageDeleted
	}

	insert, err := db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(context.Background(), types.MongoMessage{
		RoomID:          roomID,
		Message:         message,
		ReplyMessageID:  &messageID,
		Timestamp:       timestamp,
		EditedTimestamp: timestamp,
		Owner:           owner,
	})
	if err != nil {
		return nil, err
	}

	ID, _ := insert.InsertedID.(primitive.ObjectID)
	return db.GetMessage(accountID, roomID, ID.Hex())
}

func (db *MongoClient) AddBotMessage(accountID, roomID string, botID int, message string, bubbles []string) (*types.Message, error) {
	exist, err := db.RoomExists(accountID, roomID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrCanNotAddMessageToUnexistedRoom
	}

	var now = time.Now().UnixMilli()

	insert, err := db.MongoDB.Database(MongoDBName).Collection(accountID).InsertOne(context.Background(), types.MongoMessage{
		RoomID:          roomID,
		Message:         message,
		Timestamp:       now,
		EditedTimestamp: now,
		Bubbles:         bubbles,
		Owner:           types.Owner{ID: strconv.Itoa(botID), Kind: "bot"},
	})

	if err != nil {
		return nil, err
	}

	messageID, _ := insert.InsertedID.(primitive.ObjectID)
	return db.GetMessage(accountID, roomID, messageID.Hex())
}

func (db *MongoClient) checkMember(accountID, ownerID, roomID string) error {
	exist, err := db.IsRoomMember(accountID, ownerID, roomID)
	if err != nil && err != ErrRoomDoesNotExist {
		return err
	}

	if err == ErrRoomDoesNotExist {
		return ErrCanNotAddMessageToUnexistedRoom
	}

	if !exist {
		return ErrCanNotAddMessageFromNoRoomMember
	}

	return nil
}

func (db *MongoClient) GetMessagesCount(accountID, roomID string) (int64, error) {
	var filter = bson.D{{Key: "roomID", Value: roomID}}
	var count, err = db.MongoDB.Database(MongoDBName).Collection(accountID).CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *MongoClient) UpdateUserLastReadTimestamp(ctx context.Context, accountID, userID, roomID string, timestamp int64) error {
	var newTimestamp int64
	var err error
	if timestamp == 0 {
		var message types.MongoMessageWithID
		var options = options.FindOne()
		options.SetSort(bson.M{"_id": -1}).SetProjection(bson.D{
			primitive.E{Key: "timestamp", Value: 1},
		})

		err = db.MongoDB.Database(MongoDBName).
			Collection(accountID).FindOne(ctx, bson.M{"roomID": roomID}, options).
			Decode(&message)

		if err == mongo.ErrNoDocuments {
			newTimestamp = 0
		} else {
			if err != nil {
				return err
			}
			newTimestamp = message.Timestamp
		}
	}

	if timestamp != 0 {
		newTimestamp = timestamp
	}

	_, err = db.MongoDB.Database(MongoDBName).
		Collection(RoomsCollection).
		UpdateOne(ctx, bson.M{
			"roomID": roomID,
			"users":  bson.M{"$elemMatch": bson.M{"id": userID}},
		}, bson.M{"$set": bson.M{"users.$.lastReadTimestamp": newTimestamp}})

	return err
}

func (db *MongoClient) GetForm(ctx context.Context, id string) (*types.MongoForm, error) {
	var mongoForm types.MongoForm
	var coll = db.MongoDB.Database(MongoDBName).Collection(FormsCollections)
	var formID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = coll.FindOne(ctx, bson.M{"_id": formID}).Decode(&mongoForm)
	if err != nil {
		return nil, err
	}

	return &mongoForm, nil
}

func (db *MongoClient) AddForm(ctx context.Context, form *types.MongoFormWithoutID) (string, error) {
	var coll = db.MongoDB.Database(MongoDBName).Collection(FormsCollections)
	var res, err = coll.InsertOne(ctx, form)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (db *MongoClient) UpdateForm(ctx context.Context, form *types.MongoForm) error {
	var filter = bson.D{{Key: "_id", Value: form.ID}}
	var coll = db.MongoDB.Database(MongoDBName).Collection(FormsCollections)
	var _, err = coll.ReplaceOne(context.TODO(), filter, form.MongoFormWithoutID)

	return err
}


*/
