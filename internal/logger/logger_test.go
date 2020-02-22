package logger

import (
	"errors"
	"os"
	"testing"
)

func Test_Logger(t *testing.T) {
	f , _ := os.Create("logger.log")
	Log.SetOut(f)
	Log.SetLevel(DEBUG)
	Log.Debug("debug log")
	Log.Debugf("%.8f", 0.2912101221212)
	Log.Info("info log")
	Log.Warn(errors.New("test error"))
	//	Log.Fatal("fatal log")
	//Log.Panicf("%s","panicf log")
	Debug("debug log2")
	Info("info log2")
}
