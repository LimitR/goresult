package goresult

import (
	"regexp"
	"runtime"
)

type trace struct {
	message  string
	fnName   string
	fileName string
	line     int
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
	res := "Trace: " + e.printTrace() + "\n"
	res += "Message: " + e.printMessage()
	return res
}

func (e *err) printTrace() string {
	str := ""
	for i, v := range reverseInts(e.trace) {
		if v.fnName != "" {
			if i == len(e.trace)-1 {
				str += v.fnName
			} else {
				str += v.fnName + " -> "
			}
		}
	}
	return str
}

func (e *err) printMessage() string {
	str := ""
	for _, v := range reverseInts(e.trace) {
		if v.message != "" {
			str += v.fnName + " '" + v.message + "', "
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
	pc, file, line, _ := runtime.Caller(lvl)
	functionObject := runtime.FuncForPC(pc)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	fnName := extractFnName.ReplaceAllString(functionObject.Name(), "$1")
	return trace{
		fnName:   fnName,
		fileName: file,
		line:     line,
		message:  "",
	}
}
