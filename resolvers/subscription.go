package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/learning-service/api/learning"
	"log"
	"time"
)

func (s *Server) GetCurrentSubscription(ctx context.Context, in *learning.IDRequest) (*learning.CurrentSubscription, error) {
	//err := userHasClassRoom(s.DB, in.UserID, in.ClassRoomID)
	//if err != nil {
	//	s.Logger.Error(err.Error())
	//	return nil, err
	//}
	classRoomID := in.Id
	studentID := in.UserID

	currentTime := time.Now()
	var paidAt time.Time
	var subscriptionID string
	var monthID int32
	err := s.DB.QueryRow(`
SELECT subscription_id,month_id,paid_at FROM subscription where classroom_id = $1 AND student_id = $2 and paid_at between $3 and $4 `, classRoomID, studentID, currentTime.
		AddDate(0, -1, 0), currentTime).Scan(&subscriptionID, &monthID, &paidAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &learning.CurrentSubscription{
				Status: false,
			}, nil
		}
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	log.Println(paidAt)
	return &learning.CurrentSubscription{
		Status:   true,
		RestDays: int32(paidAt.AddDate(0, 1, 0).Sub(currentTime).Hours() / 24),
		Subscription: &learning.Subscription{
			Id:      subscriptionID,
			MonthId: monthID,
			PaidAt:  paidAt.Unix(),
		},
	}, nil
}

func (s *Server) GetSubscriptions(ctx context.Context, in *learning.IDRequest) (*learning.Subscriptions, error) {
	//err := userHasClassRoom(s.DB, in.UserID, in.ClassRoomID)
	//if err != nil {
	//	s.Logger.Error(err.Error())
	//	return nil, err
	//}
	classRoomID := in.Id
	studentID := in.UserID

	rows, err := s.DB.Query(`
SELECT subscription_id,month_id,paid_at FROM subscription where classroom_id = $1 AND student_id = $2`, classRoomID, studentID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	var subscriptions = make([]*learning.Subscription, 0)
	for rows.Next() {
		var subscription = new(learning.Subscription)
		var paidAt time.Time
		err = rows.Scan(&subscription.Id, &subscription.MonthId, &paidAt)
		if err != nil {
			return nil, ErrInternal
		}
		subscription.PaidAt = paidAt.Unix()
		subscriptions = append(subscriptions, subscription)
	}
	return &learning.Subscriptions{
		Subscriptions: subscriptions,
	}, nil
}
