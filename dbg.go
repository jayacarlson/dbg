package dbg

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jayacarlson/env"
)

/*
	A collection of debugging and tracking routines and some other convenience functions

		See Debug.go for functions & descriptions
*/

var (
	// Can redirect debug output to logging by changing this to log.Printf
	output = fmt.Printf
	outerr = errout

	normColor, msgColor, infoColor    string
	noteColor, warnColor, ccnColor    string
	failColor, errColor, fatalColor   string
	WARNColor, CAUTNColor, ERRORColor string
)

// ========================================================================= //

func init() {
	if env.IsLinux() {
		Color() // enable color output on linux systems
	}
}

func errout(f string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, f, a...) // why not going to Stderr?
}

// returns a shortened file name with minimal leading path
func shortName(s string) string {
	p, l := 0, 0
	for n, c := range s { // cheap way to find 2nd to last '/'
		if c == '/' { // TODO_TODO_TODO : change to os.PathSeparator if needed
			p = l
			l = n + 1
		}
	}
	return s[p:]
}

// outputs location information, 2 steps back (who called the dbg.func)
func trcAt(a ...interface{}) {
	if _, file, line, ok := runtime.Caller(2); ok {
		file = shortName(file)
		output("TRC @ %d in %s ", line, file)
	}
	trc(a...)
}

// outputs location information, 3 steps back (who called the function calling dbg.func)
func trcBefore(a ...interface{}) {
	if _, file, line, ok := runtime.Caller(3); ok {
		file = shortName(file)
		output("WAS @ %d in %s ", line, file)
	}
	trc(a...)
}

// output trc info -- see trc_args
func trc(a ...interface{}) {
	s := ""
	if len(a) > 0 {
		if f, ok := a[0].(string); ok { // string with possible args
			s = fmt.Sprintf(msgColor+f+normColor, a[1:]...)
		} else if e, ok := a[0].(error); ok { // error, output error text
			s = fmt.Sprintf(errColor+"%v"+normColor, e)
		} else if nil == a[0] { // condition where given error is NIL
			s = fmt.Sprintf(infoColor + "nil" + normColor)
		}
	}
	output("%s\n", s)
}

// returns location of CHK caller
func at() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		file = shortName(file)
		return fmt.Sprintf("@ %d in %s  ", line, file)
	}
	return ""
}

// return location line, file & func as string
func funcAt(d int) string {
	if uptr, file, line, ok := runtime.Caller(d + 1); ok {
		name := runtime.FuncForPC(uptr).Name()
		return fmt.Sprintf("@ %d in %s - %s()", line, shortName(file), name[strings.LastIndex(name, "/")+1:])
	}
	return "@ <UNKNOWN>"
}

// return arg text after calling any possible CLOSER()
func failed(c bool, a ...interface{}) string {
	if len(a) > 0 && c {
		if cl, ok := a[len(a)-1].(func()); ok {
			cl()             // call closer function
			a = a[:len(a)-1] // remove it from arg list
		}
	}
	return genText(a...)
}

// return error text or arg text after calling any possible CLOSER()
func errored(c bool, e error, a ...interface{}) string {
	var txt string
	if len(a) > 0 {
		if c {
			if cl, ok := a[len(a)-1].(func()); ok {
				cl()             // call closer function
				a = a[:len(a)-1] // remove it from arg list
			}
		}
		txt = genText(a...)
	} else {
		txt = fmt.Sprintf("%v", e)
	}
	return txt
}

// generates text for output and calls any CLOSER function before error processing continues
func genText_Closer(a ...interface{}) string {
	if len(a) > 0 { // check for CLOSER -- pull last interface off and see if a 'func'
		if cl, ok := a[len(a)-1].(func()); ok {
			cl()             // call closer function
			a = a[:len(a)-1] // remove it from arg list
		}
	}
	return genText(a...)
}

// generates text for output, supplying a 'Check failed' if none given
func genText(a ...interface{}) string {
	s := "Check failed"
	if len(a) > 0 {
		if f, ok := a[0].(string); ok {
			s = fmt.Sprintf(f, a[1:]...)
		} else if e, ok := a[0].(error); ok {
			s = fmt.Sprintf("%v", e)
			if len(a) > 1 {
				if f, ok := a[1].(string); ok {
					s += fmt.Sprintf(f, a[2:]...)
				}
			}
		}
	}
	return s
}
