package factory

import (
	authorizer "github.com/kutsuzawa/line-authorizer"
	"github.com/kutsuzawa/line-reminder/service"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ServiceFactory defines method to initialize service
type ServiceFactory interface {
	LineService() (service.LineService, error)
}

type serviceFactory struct {
	id     string
	secret string
}

// NewServiceFactory initializes ServiceFactory
func NewServiceFactory(id, secret string) ServiceFactory {
	return &serviceFactory{
		id:     id,
		secret: secret,
	}
}

// LineService initializes service.LineService
func (sf *serviceFactory) LineService() (service.LineService, error) {
	config := authorizer.Config{
		ID:     sf.id,
		Secret: sf.secret,
	}
	apiClient := authorizer.NewClient(config)
	token, err := apiClient.PublishChannelToken()
	if err != nil {
		return nil, err
	}
	client, err := linebot.New(sf.secret, *token)
	if err != nil {
		return nil, err
	}
	return service.NewLineService(client), nil
}
