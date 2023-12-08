package connection

type Service interface {
	GetName() string
	Read() []byte
	Write(data []byte) error
	Flush() error
	Close() error
}
