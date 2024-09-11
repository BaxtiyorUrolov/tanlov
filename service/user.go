package service

import (
	"context"
	"fmt"
	"it-tanlov/api/models"
	"it-tanlov/pkg/logger"
	"it-tanlov/storage"
)

type userService struct {
	storage storage.IStorage
	log     logger.ILogger
}

func NewUserService(storage storage.IStorage, log logger.ILogger) userService {
	return userService{
		storage: storage,
		log:     log,
	}
}

func (u userService) Create(ctx context.Context, telegramID int) error {

	err := u.storage.User().Create(context.Background(), telegramID)
	if err != nil {
		u.log.Error("Error while creating user", logger.Error(err))
		return err
	}

	return nil
}

func (u userService) AddScore(ctx context.Context, id string) (models.Partner, error) {
	
	fmt.Println("ID: ", id)

	err := u.storage.User().AddScore(ctx, id)
	if err != nil {
		u.log.Error("error in service layer while updating score", logger.Error(err))
		return models.Partner{}, err
	}

	partner, err := u.storage.Partner().GetByID(ctx, models.PrimaryKey{ID: id})
	if err != nil {
		u.log.Error("error in service layer while getting by id", logger.Error(err))
		return models.Partner{}, err
	}

	return partner, nil

}

func (u userService) IUserTelegramIDExist(ctx context.Context, telegramID int) (bool, error) {

    exists, err := u.storage.User().IUserTelegramIDExist(ctx, telegramID)
    if err != nil {
        fmt.Println("Error occurred in TelegramIDExist:", err)
        return false, fmt.Errorf("error while checking telegramID existence: %w", err)
    }

    fmt.Println("TelegramIDExist result:", exists)
    return exists, nil
}


