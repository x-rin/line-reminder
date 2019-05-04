package factory

import (
	"github.com/kutsuzawa/line-reminder/reminder"
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
	channelToken, err := reminder.GetChannelToken(sf.id, sf.secret)
	if err != nil {
		return nil, err
	}
	client, err := linebot.New(sf.secret, *channelToken)
	if err != nil {
		return nil, err
	}
	return service.NewLineService(client), nil
}
