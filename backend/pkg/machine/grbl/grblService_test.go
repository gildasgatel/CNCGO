package grbl

import (
	"bytes"
	"cncgo/backend/internal/api/models"
	connection "cncgo/backend/pkg/connection/mock"
	"testing"

	"github.com/go-playground/assert"
)

func TestNewFail(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Error("Conn nil, want error got nil")
	}
}
func TestNewGood(t *testing.T) {
	grbl, err := New(&connection.MockService{})
	if err != nil {
		t.Errorf("Conn ok, want nil got %s\n", err.Error())
	}
	_, ok := grbl.(*Grbl)
	if !ok {
		t.Errorf("Expected an instance of *Grbl, got %T", grbl)
	}

}
func TestHandelState(t *testing.T) {
	var grbl = &Grbl{}
	grbl.state.State = "test"

	tests := []struct {
		name string
		line []byte
		want *models.StateMachine
	}{
		{name: "do not match", line: []byte(""),
			want: &models.StateMachine{
				State: "test",
				MPos:  "",
				BfR:   0,
				BfW:   0,
				FS:    "",
				WC0:   ""},
		},
		{name: "match ok", line: []byte("<State:Idle|MPos: mypos|Bf:12,53|Fs:testok|WC0:boom>"),
			want: &models.StateMachine{
				State: "Idle",
				MPos:  " mypos",
				BfR:   12,
				BfW:   53,
				FS:    "testok",
				WC0:   "boom"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grbl.HandleState([]byte(tt.line))
			got := grbl.state
			assert.Equal(t, tt.want.State, got.State)
			assert.Equal(t, tt.want.MPos, got.MPos)
			assert.Equal(t, tt.want.BfR, got.BfR)
			assert.Equal(t, tt.want.BfW, got.BfW)
			assert.Equal(t, tt.want.FS, got.FS)
			assert.Equal(t, tt.want.WC0, got.WC0)

		})
	}
}

func TestSendCommand(t *testing.T) {
	mockConn := &connection.MockService{}
	grbl, _ := New(mockConn)
	type want struct {
		err bool
		out []byte
	}
	tests := []struct {
		name string
		f    func()
		arg  models.Command
		want want
	}{
		{name: "Send move",
			f: func() {
				mockConn.WriteErr = true
			},
			arg: models.Command{
				Command:  "move",
				Distance: "10",
				Axe:      "Y"},
			want: want{err: true, out: []byte("Y10\n")},
		},
		{name: "Send G0",
			f: func() {
				mockConn.WriteErr = true
			},
			arg: models.Command{
				Command: "G0"},
			want: want{err: true, out: []byte("G0\n")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f()
			out, err := grbl.SendCommand(tt.arg)
			if err != nil {
				out = []byte(err.Error())
			}
			gotErr := err != nil
			if tt.want.err != gotErr {
				t.Errorf("%s want %v got %v\n", tt.name, tt.want.err, gotErr)
			}
			if ok := bytes.Compare(tt.want.out, out); ok != 0 {
				t.Errorf("%s want %v got %v\n", tt.name, tt.want.out, out)
			}
		})
	}
}
