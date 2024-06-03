package types

//
//// RepliedMessage Describe the replied message
//type RepliedMessage struct {
//	ID      string `json:"id"`
//	Message string `json:"message"`
//	Owner   Owner  `json:"owner"`
//}
//
//// Message Describe a message to send to users
//type Message struct {
//	ID                 string           `json:"id"                        validate:"required"`
//	RoomID             string           `json:"roomID"                    validate:"required_unless=IsNotification true,omitempty"`
//	Message            string           `json:"message"                   validate:"required"`
//	Timestamp          int64            `json:"timestamp"                 validate:"required"`
//	FormattedTimestamp string           `json:"formattedTimestamp"`
//	EditedTimestamp    int64            `json:"editedTimestamp,omitempty"`
//	IsFile             bool             `json:"isFile,omitempty"`
//	IsEdited           bool             `json:"isEdited"`
//	IsDeleted          bool             `json:"isDeleted"`
//	IsComment          bool             `json:"isComment"`
//	IsNotification     bool             `json:"isNotification"`
//	IsForm             bool             `json:"isForm"`
//	IsPrivate          bool             `json:"isPrivate"`
//	TargetMsgID        string           `json:"targetMessageID,omitempty"`
//	Bubbles            []string         `json:"bubbles,omitempty"`
//	Owner              Owner            `json:"owner,omitempty"`
//	Replied            *RepliedMessage  `json:"replied,omitempty"`
//	Comments           []*Comment       `json:"comments,omitempty"`
//	Attachments        []Attachment     `json:"files,omitempty"`
//	EditHistory        []MessageHistory `json:"editHistory,omitempty"`
//	Cursor             string           `json:"cursor"                                                                             bson:"-"`
//	FormID             string           `json:"formID,omitempty"`
//}

type Comment struct {
	ID         string `bson:"id"                        json:"id"`
	DocumentID string `bson:"documentID"                json:"documentID"`
	Message    string `bson:"message"                   json:"message"`
	Timestamp  int64  `bson:"timestamp"                 json:"timestamp"`
	Owner      string `bson:"owner"                     json:"owner"`
	IsEdited   bool   `bson:"isEdited,omitempty"        json:"isEdited,omitempty"`
	IsDeleted  bool   `bson:"isDeleted,omitempty"       json:"isDeleted,omitempty"`
}

//type MongoMessage struct {
//	ID                 string           `bson:"id"`
//	RoomID             string           `bson:"roomID,omitempty"`
//	ReplyMessageID     *string          `bson:"replyMessageID,omitempty"`
//	Replied            *RepliedMessage  `json:"replied,omitempty"`
//	Message            string           `bson:"message,omitempty" `
//	Timestamp          int64            `bson:"timestamp,omitempty" `
//	FormattedTimestamp string           `bson:"formattedTimestamp"`
//	EditedTimestamp    int64            `bson:"editedTimestamp,omitempty"`
//	IsComment          bool             `bson:"isComment,omitempty"`
//	IsEdited           bool             `bson:"isEdited,omitempty"`
//	IsDeleted          bool             `bson:"isDeleted,omitempty"`
//	IsNotification     bool             `bson:"isNotification,omitempty"`
//	IsFile             bool             `bson:"isFile,omitempty"`
//	IsForm             bool             `bson:"isForm,omitempty"`
//	IsPrivate          bool             `bson:"isPrivate,omitempty"`
//	Bubbles            []string         `bson:"bubbles,omitempty"`
//	Attachments        []Attachment     `bson:"attachments,omitempty"`
//	Owner              Owner            `bson:"owner,omitempty"`
//	EditHistory        []MessageHistory `bson:"editHistory,omitempty"`
//	Comments           []Comment        `bson:"comments,omitempty"`
//	Cursor             string           `bson:"-"`
//	FormID             string           `bson:"formID,omitempty"`
//}
