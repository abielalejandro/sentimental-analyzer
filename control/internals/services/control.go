package services

import (
	"context"
	"time"

	"github.com/abielalejandro/control/config"
	"github.com/abielalejandro/control/internals/storage"
	"github.com/google/uuid"
)

type ControlService struct {
	config  *config.Config
	storage storage.Storage
}

func (svc *ControlService) ProcessMsg(
	ctx context.Context,
	msg string) (string, error) {

	t := time.Now()
	message := &storage.Message{
		Id:        uuid.New().String(),
		Msg:       msg,
		CreatedAt: t,
		UpdatedAt: t,
		ExpiresAt: t.Add(time.Duration(svc.config.Ttl) * time.Minute),
	}
	_, err := svc.storage.Create(ctx, message)
	if err != nil {
		return "", err
	}
	return message.Id, nil
}

func (svc *ControlService) UpdateSentimentalMsg(
	ctx context.Context,
	id string,
	result *storage.SentimentalResult) error {

	_, err := svc.storage.Update(ctx, id, result)
	if err != nil {
		return err
	}

	return nil
}
