package logs

import (
	"bytes"
	"github.com/labstack/gommon/log"
	"io"
)

type Logerr struct {
	Level int
	Out  io.Writer
	OutLog *log.Logger
	Debugs bool
}
type ILogerr interface {
	Debug(err []byte)
	Info(err []byte)
	Warn(err []byte)
}

func NewLogerr(d bool)  ILogerr{
	return &Logerr{Debugs:d,OutLog:log.New("huhu")}
}
func (r *Logerr) Debug(err []byte)  {
	b := new(bytes.Buffer)
	r.OutLog.SetOutput(b)
	r.OutLog.SetLevel(log.WARN)
	r.Level=1
	r.OutLog.Debug(string(err))
}

func (r *Logerr) Info(err []byte)  {
	r.Level=3
	r.OutLog.Info(err)
}

func (r *Logerr) Warn(err []byte)  {
	r.Level=2
	r.OutLog.Warn(err)

}
