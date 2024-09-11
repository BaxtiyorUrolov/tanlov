package storage

import (
	"context"
	"it-tanlov/api/models"
)

type IStorage interface {
	Close()
	Partner() IPartnerStorage
	User() IUserStorage
}

type IPartnerStorage interface {
	Create(context.Context, models.CreatePartner) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.Partner, error)
	GetList(context.Context, models.GetListRequest)(models.PartnerResponse, error)
	Update(context.Context, string, string) error
	Delete(context.Context, string) error
	PhoneExist(context.Context, string) (bool, error)
	IEmailExist(context.Context, string) (bool, error)
	IVideoLinkExist(context.Context, string) (bool, error)
}

type IUserStorage interface {
	Create(context.Context, int) error
	AddScore(context.Context, string) error
	IUserTelegramIDExist(context.Context, int) (bool, error)
}

