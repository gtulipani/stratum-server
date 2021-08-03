package service

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"stratum-server/repository"
	"time"
)

type subscription struct {
	extraNonce1   int64
	extraNonce2   int64
	setDifficulty string
	notify        string
	subscriber    string
	createdAt     time.Time
	activeSession bool
}

func (s *service) getExistingSubscription(subscriber string, extraNonce1 int64) (*subscription, error) {
	sqlStatement := fmt.Sprintf(`
	SELECT extra_nonce_1, extra_nonce_2, set_difficulty, notify, subscriber, created_at, active_session
	FROM %s.%s
	WHERE extra_nonce_1 = $1
	AND subscriber = $2`, s.subscriptionsTable.Schema, s.subscriptionsTable.Name)

	sub := &subscription{}
	if err := s.repository.Query(repository.QueryRequest{
		Query: sqlStatement,
		Args: []interface{}{
			extraNonce1,
			subscriber,
		},
	}, &sub.extraNonce1, &sub.extraNonce2, &sub.setDifficulty, &sub.notify, &sub.subscriber, &sub.createdAt, &sub.activeSession); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no subscription found for extraNonce1: %d and subscriber: %s", extraNonce1, subscriber)
			return nil, nil
		}
		log.Printf("error getting subscription: %v", err)
		return nil, err
	}

	return sub, nil
}

func (s *service) createSubscription(subscriber string) (*subscription, error) {
	sqlStatement := fmt.Sprintf(`
	INSERT INTO %s.%s (extra_nonce_2, set_difficulty, notify, subscriber)
	VALUES ($1, $2, $3, $4)
	RETURNING extra_nonce_1, extra_nonce_2, set_difficulty, notify, subscriber, created_at, active_session`, s.subscriptionsTable.Schema, s.subscriptionsTable.Name)

	sub := &subscription{}
	if err := s.repository.Insert(repository.InsertRequest{
		Query: sqlStatement,
		Args: []interface{}{
			s.GetExtraNonce2(),
			uuid.NewString(),
			uuid.NewString(),
			subscriber,
		},
	}, &sub.extraNonce1, &sub.extraNonce2, &sub.setDifficulty, &sub.notify, &sub.subscriber, &sub.createdAt, &sub.activeSession); err != nil {
		log.Printf("error creating subscription: %v", err)
		return nil, err
	}

	return sub, nil
}

func (s *service) inactiveSubscription(subscription *subscription) {

	sqlStatement := fmt.Sprintf(`
	UPDATE %s.%s
	SET active_session = $1
	WHERE extra_nonce_1 = $2
	RETURNING active_session
	`, s.subscriptionsTable.Schema, s.subscriptionsTable.Name)
	var activeSession bool

	if err := s.repository.Update(repository.UpdateRequest{
		Query: sqlStatement,
		Args: []interface{}{
			false,
			subscription.extraNonce1,
		},
	}, &activeSession); err != nil {
		log.Printf("error inactivating subscription: %v", err)
	}
}
