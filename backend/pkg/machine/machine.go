package machine

import "cncgo/backend/internal/api/models"

type Service interface {
	SendCommand(data models.Command) ([]byte, error)
	SendFile(path string) error
	GetName() string
	GetState() *models.StateMachine
}
