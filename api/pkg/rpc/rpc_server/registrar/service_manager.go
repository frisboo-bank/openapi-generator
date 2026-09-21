package registrar

import (
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type ServiceManager struct {
	logger   loggercontracts.Logger
	services []Service
}

func NewServiceManager(logger loggercontracts.Logger) *ServiceManager {
	validation.AssertNotNil("logger", logger)

	return &ServiceManager{
		logger: logger,
	}
}

func (m *ServiceManager) Register(services ...Service) {
	m.services = append(m.services, services...)
}

func (m *ServiceManager) Services() []Service {
	return m.services
}
