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
	for i, v := range reverseInts(e.trace) {
		if v.FnName != "" {
			if i == len(e.trace)-1 {
				str += v.FnName + ":" + fmt.Sprint(v.Line)
			} else {
				str += v.FnName + ":" + fmt.Sprint(v.Line) + " -> "
			}
		}
	}
	return str
}

func (e *err) printMessage() string {
	str := ""
	for i, v := range reverseInts(e.trace) {
		if v.Message != "" {
			if i == len(e.trace)-1 {
				str += v.FnName + " '" + v.Message + "'"
			} else {
				str += v.FnName + " '" + v.Message + "', "
			}
		}
	}
	return str
}

func reverseInts(input []trace) []trace {
	if len(input) == 0 {
		return input
	}
	return append(reverseInts(input[1:]), input[0])
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
