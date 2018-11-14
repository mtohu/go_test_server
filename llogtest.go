package main
import (
	"bytes"
	"github.com/labstack/gommon/log"
)
func main(){
	l := log.New("test")
	b := new(bytes.Buffer)
	l.SetOutput(b)
	l.DisableColor()
	l.SetLevel(log.WARN)

	l.Print("print")
	l.Printf("print%s", "f")
	l.Debug("debug")
	l.Debugf("debug%s", "f")
	l.Info("info")
	l.Infof("info%s", "f")
	l.Warn("warn")
	l.Warnf("warn%s", "f")
	l.Error("error")
	l.Errorf("error%s", "f")

}
