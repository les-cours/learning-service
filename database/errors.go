package database

import (
	"errors"
	"fmt"
)

var ErrRoomExist = errors.New("room already exists")
var ErrRoomDoesNotExist = errors.New("room does not exist")
var ErrRoomDoesNotHaveResumedHistory = errors.New("room does not have resumed history")
var ErrRoomAgentDoesNotResolve = errors.New("agent room does not resolve")
var ErrRoomAgent = errors.New("agent room")
var ErrCantRemoveVisitorRoom = errors.New("Cannot remove visitor room")
var ErrAgentRoomExceedAllowedMembers = errors.New("agent room can not have more than two agent on it")
var ErrNoEnoughMembers = errors.New("no enough members to create room")
var ErrUnknownRoomKind = errors.New("unknown room kind")
var ErrAllAgentCanNotBePrivate = errors.New("all agents can not be in private mode we must have at least one agent in public")
var ErrCanNotCreateRoomWithoutMembers = errors.New("can not create a room without having members")
var ErrAgentRoomMustNotHaveVisitor = errors.New("agent room must not have a visitor")
var ErrVisitorRoomMustHaveVisitor = errors.New("visitor room must have a visitor")
var ErrOnlyOneVisitorAllowed = errors.New("must have only one visitor in room")
var ErrOnlyAgentMorgueRoom = errors.New("only agent who can morgue a room")
var ErrRoomMemberDoesNotExist = errors.New("room member does not exist")
var ErrRoomOwnerMustBeInUserList = errors.New("room owner must be in user list")

var ErrCanNotAddMessageToUnexistedRoom = errors.New("can not add message to room which does not exist")
var ErrCanNotAddMessageFromNoRoomMember = errors.New("can not add message from no room member")
var ErrMessageDoesNotExist = errors.New("message does not exist")
var ErrNoOwnerMatch = errors.New("no owner match")
var ErrMessageDeleted = errors.New("message is deleted")
var ErrUnacceptedOwnerKind = errors.New("unaccepted message owner kind")
var ErrCommentNotFound = errors.New("comment not found")
var ErrVisitorNotFound = fmt.Errorf("visitor not found")
