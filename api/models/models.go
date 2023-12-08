package models

type Config struct {
	Machine    string `json:"machine"`
	Connection string `json:"connection"`
	PortName   string `json:"port"`
	BaudRate   int    `json:"baudrate"`
}

type Command struct {
	Command  string `json:"command"`
	Distance string `json:"distance"`
	Axe      string `json:"axe"`
}
type File struct {
	Path string `json:"path"`
}

type StateMachine struct {
	State string `json:"state"`
	MPos  string `json:"mpos"`
	BfR   int    `json:"bfr"`
	BfW   int    `json:"bfw"`
	FS    string `json:"fS"`
	WC0   string `json:"wc0"`
}
