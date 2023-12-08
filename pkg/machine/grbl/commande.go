package grbl

var commandGrbl = map[string]string{
	"info":  "$$",
	"state": "?",
	"play":  "M3",
	"pause": "M0",
	"stop":  "M4",
}
