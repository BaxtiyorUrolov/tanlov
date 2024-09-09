package service

import (
	"it-tanlov/pkg/logger"
	"it-tanlov/storage"
)

type IServiceManager interface {
	Partner() partnerService
	User() 	  userService
}

type Service struct {
	partnerService partnerService
	userService    userService
}

func New(storage storage.IStorage, log logger.ILogger) Service {
	services := Service{}

	services.partnerService = NewPartnerService(storage, log)
	services.userService = NewUserService(storage, log)

	return services
}

func (s Service) Partner() partnerService {
	return s.partnerService
}

func (s Service) User() userService {
	return s.userService
}
