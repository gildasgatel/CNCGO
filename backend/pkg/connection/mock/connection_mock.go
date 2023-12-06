package mock

import "errors"

type MockService struct {
	Name     string
	WriteErr bool
	ReadErr  bool
	FlushErr bool
	CloseErr bool
}

func (m *MockService) GetName() string {
	return m.Name
}

func (m *MockService) Write(data []byte) error {
	if m.WriteErr {
		return errors.New(string(data))
	} else {
		return nil
	}
}

func (m *MockService) Read() []byte {
	if m.WriteErr {
		return []byte("data read mockService")
	} else {
		return []byte{}
	}
}

func (m *MockService) Flush() error {
	if m.WriteErr {
		return errors.New("error flush mockService")
	} else {
		return nil
	}
}
func (m *MockService) Close() error {
	if m.WriteErr {
		return errors.New("error close mockService")
	} else {
		return nil
	}
}
