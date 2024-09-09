package service

import (
	"context"
	"it-tanlov/api/models"
	"it-tanlov/pkg/logger"
	"it-tanlov/storage"
)

type partnerService struct {
	storage storage.IStorage
	log     logger.ILogger
}

func NewPartnerService(storage storage.IStorage, log logger.ILogger) partnerService {
	return partnerService{
		storage: storage,
		log:     log,
	}
}

func (p partnerService) Create(ctx context.Context, createUser models.CreatePartner) (models.Partner, error) {
	p.log.Info("user create service layer", logger.Any("user", createUser))

	pKey, err := p.storage.Partner().Create(context.Background(), createUser)
	if err != nil {
		p.log.Error("Error while creating user", logger.Error(err))
		return  models.Partner{}, err
	}


	partner, err := p.storage.Partner().GetByID(context.Background(), models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		p.log.Error("Error in service layer when getting partner by id", logger.Error(err))
		return partner, err
	}

	return partner, nil
}

func (p partnerService) Get(ctx context.Context, id models.PrimaryKey) (models.Partner, error) {
	partner, err := p.storage.Partner().GetByID(ctx, id)
	if err != nil {
		p.log.Error("error in service layer while getting by id", logger.Error(err))
		return models.Partner{}, err
	}

	return partner, nil
}

func (p partnerService) GetList(ctx context.Context, request models.GetListRequest) (models.PartnerResponse, error) {

	partners, err := p.storage.Partner().GetList(ctx, request)
	if err != nil {
		p.log.Error("error in service layer  while getting list", logger.Error(err))
		return models.PartnerResponse{}, err
	}

	return partners, nil
}

func (p partnerService) Update(ctx context.Context, key string, imageURL string) error {
	err := p.storage.Partner().Update(ctx, key, imageURL)
	if err != nil {
		p.log.Error("error in service layer while updating", logger.Error(err))
		return  err
	}

	return  nil
}

func (p partnerService) Delete(ctx context.Context, key string) error {
	err := p.storage.Partner().Delete(ctx, key)

	return err
}
