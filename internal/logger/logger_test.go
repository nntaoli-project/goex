package logger

import (
	"errors"
	"testing"
)

func Test_Logger(t *testing.T) {
	Log.SetLevel(DEBUG)
	Log.Debug("debug log")
	Log.Debugf("%.8f", 0.2912101221212)
	Log.Info("info log")
	Log.Warn(errors.New("test error"))
	//Log.Fatal("fatal log")
	//Log.Panicf("%s","panicf log")
}
