package goresult

import (
	"fmt"
	"regexp"
	"runtime"
)

type trace struct {
	Message  string
	FnName   string
	FileName string
	Line     int
}

type err struct {
	trace     []trace
	TimeStamp int64
	Err       error
}

func (e *err) AddTrace() {
	e.trace = append(e.trace, getTrace(3))
}

func (e *err) print() string {
	res := "\nTrace: " + e.printTrace() + "\n"
	res += "Message: " + e.printMessage()
	return res
}

func (e *err) printTrace() string {
	str := ""
	l := len(e.trace)
	for i, _ := range e.trace {
		if e.trace[l-i-1].FnName != "" {
			if i == l-1 {
				str += e.trace[l-i-1].FnName + ":" + fmt.Sprint(e.trace[l-i-1].Line)
			} else {
				str += e.trace[l-i-1].FnName + ":" + fmt.Sprint(e.trace[l-i-1].Line) + " -> "
			}
		}
	}
	return str
}

func (e *err) printMessage() string {
	str := ""
	l := len(e.trace)
	for i, _ := range e.trace {
		if e.trace[l-i-1].Message != "" {
			if i == l-1 {
				str += e.trace[l-i-1].FnName + " '" + e.trace[l-i-1].Message + "'"
			} else {
				str += e.trace[l-i-1].FnName + " '" + e.trace[l-i-1].Message + "', "
			}
		}
	}
	return str
}

func getTrace(lvl int) trace {
	pc, file, Line, _ := runtime.Caller(lvl)
	functionObject := runtime.FuncForPC(pc)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	FnName := extractFnName.ReplaceAllString(functionObject.Name(), "$1")
	return trace{
		FnName:   FnName,
		FileName: file,
		Line:     Line,
		Message:  "",
	}
}
