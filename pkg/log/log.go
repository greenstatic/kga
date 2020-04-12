package log

import (
	"fmt"
	"os"
)

func Infof(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func Info(a interface{}) {
	Infof("%s", a)
}

func Errorf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format+"\n", a...)
	if err != nil {
		panic(err)
	}
}

func Error(a interface{}) {
	Errorf("%s", a)
}

func Fatalf(format string, a ...interface{}) {
	Errorf(format, a)
	os.Exit(1)
}

func Fatal(a interface{}) {
	Errorf("%s", a)
	os.Exit(1)
}

func FatalOnError(err error) {
	if err != nil {
		Fatal(err)
	}
}
