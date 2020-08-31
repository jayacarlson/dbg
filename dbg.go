package dbg

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/jayacarlson/env"
)

/*
	A collection of debugging and tracking routines and some other convenience functions

	Color()									Enable colored text output (if system supports it)
	NoColor()								Disable colored text output

	Panic( "fmtStr", [chk_args] )			output colored text then PANIC! -- See dbg.Panic below
	Fatal( "fmtStr", [chk_args] )			output colored text then force exit
	PanicIf( bool, "fmtStr", [chk_args] )	PANIC only if true -- See dbg.Panic below
	FatalIf( bool, "fmtStr", [chk_args] )	force exit if true
	PanicIfErr( err, "fmtStr", [chk_args] )	PANIC only if err is not nil -- See dbg.Panic below
	FatalIfErr( err, "fmtStr", [chk_args] )	force exit if err is not nil

	SetLevel( int )							set the current debug output level
	SetMask( int )							set the current debug output bit mask

	LvlMsg( int, "fmtStr", [fmt_args] )
		output message if int val <= current debug level

	MaskMsg( int, "fmtStr", [fmt_args] )
		output message if val masks out as any non-zero value

	ChkTru( bool, [chk_args] ) bool
		if test value is false, output check failed message (see below)
		 returns TRUE on failure allowing this to be wrapped as part of 'if'

	ChkTru[PX]( bool, [chk_args] )
		if test value is false, output check failed message (see below)
		 then can either Panic or force Exit -- See dbg.Panic below

	ChkErr( error, [chk_args] ) bool
		if error non-nil, output check failed message (see below)
		 returns TRUE on non-nil allowing this to be wrapped as part of 'if'

	ChkErr[PX]( error, [chk_args] )
		if error non-nil, output check failed message (see below) then
		 either Panic or force Exit -- See dbg.Panic below

	ChkErrI( error, []error, [chk_args]) bool
		if error non-nil, output check failed message (see below) as long
		 as it's not in the ignore list of errors
		 returns TRUE on non-nil allowing this to be wrapped as part of 'if'

	TRC( [trc_args] )			output calling func file & line number followed by any arg data
	TRCIF( tst [, trc_args] )	conditional TRC
	TRCFROM( [trc_args] )		output func calling func file & line number followed by any arg data

	Trace()						output call stack (up to ten levels deep)

	Echo( "fmtStr", [fmt_args] )			output normal text (quick way to do output w/o 'fmt' if you want)
	Note( "fmtStr", [fmt_args] )			output colored text
	Info( "fmtStr", [fmt_args] )			output colored text
	Message( "fmtStr", [fmt_args] )			output colored text
	Warning( "fmtStr", [fmt_args] )			output colored text
	Error( "fmtStr", [fmt_args] )			output colored text
	Danger( "fmtStr", [fmt_args] )			output colored text

	Different argument options:
		args:			arguments for fmt.Printf format staring
		chk_args:		["fmtStr", [args]], [CLOSER()]
							any CLOSER() func is called before doing
							 any panic or exit for a failure case
		trc_args:		[error] | ["fmtStr", [args]]

	dbg.Panic() always passes the built-in panic a STRING, even if given an error
	If using defer, the returned value for 'recover' will therefore always be a string
*/

var (
	// Can redirect debug output to logging by changing this to log.Printf
	output = fmt.Printf

	dbgLevel = 0
	dbgMask  = 0

	normColor, msgColor, infoColor, noteColor string
	warnColor, ccnColor, errColor, fatalColor string
)

func init() {
	if env.IsLinux() {
		Color() // enable color output on linux systems
	}
}

// dummy func to allow external use / not-use
//	have dbg.Link() at start of file and you can enable / disable dbg code
//	without getting the pesky build errors for import use of non-use
//	-- should remove after any debug
func Link() {}

// ----------------------------------------------------------------------------

// enable color output for debug text
func Color() {
	normColor = "\033[0m"     // reset to normal text
	msgColor = "\033[33m"     // ORANGE - DIM YELLOW
	infoColor = "\033[32m"    // GREEN
	noteColor = "\033[34m"    // BLUE
	warnColor = "\033[35m"    // MAGENTA
	ccnColor = "\033[93m"     // YELLOW - BRIGHT ORANGE
	errColor = "\033[31m"     // RED
	fatalColor = "\033[1;41m" // WHITE on RED
}

// disable color output for debug text
func NoColor() {
	normColor = ""
	msgColor = ""
	infoColor = ""
	noteColor = ""
	warnColor = ""
	ccnColor = ""
	errColor = ""
	fatalColor = ""
}

// Simple output functions that give colored text -- can be redirected to logging if desired

// simply echo to output, no color hilites
func Echo(fstr string, a ...interface{}) {
	output(fstr+"\n", a...)
}

// orange (dim yellow) text to output
func Message(fstr string, a ...interface{}) {
	output(msgColor+fstr+normColor+"\n", a...)
}

// green text to output
func Info(fstr string, a ...interface{}) {
	output(infoColor+fstr+normColor+"\n", a...)
}

// blue text to output
func Note(fstr string, a ...interface{}) {
	output(noteColor+fstr+normColor+"\n", a...)
}

// yellow (bright orange) text to output
func Warning(fstr string, a ...interface{}) {
	output(warnColor+fstr+normColor+"\n", a...)
}

// magenta text to output
func Caution(fstr string, a ...interface{}) {
	output(ccnColor+fstr+normColor+"\n", a...)
}

// red text to output
func Error(fstr string, a ...interface{}) {
	output(errColor+fstr+normColor+"\n", a...)
}

// bold white on red background text to output
func Danger(fstr string, a ...interface{}) {
	output(fatalColor+fstr+normColor+"\n", a...)
}

// panic with any optional chk_args
func Panic(a ...interface{}) {
	panic(genText_Closer(a...))
}

// fatal error (exit) with any optional chk_args
func Fatal(a ...interface{}) {
	output(fatalColor + genText_Closer(a...) + normColor + "\n")
	os.Exit(-1)
}

// conditional panic
func PanicIf(b bool, a ...interface{}) {
	if b {
		panic(failed(true, a...))
	}
}

// conditional fatal
func FatalIf(b bool, a ...interface{}) {
	if b {
		output(fatalColor + failed(true, a...) + normColor + "\n")
		os.Exit(-1)
	}
}

// conditional panic
func PanicIfErr(e error, a ...interface{}) {
	if e != nil {
		panic(errored(true, e, a...))
	}
}

// conditional fatal
func FatalIfErr(e error, a ...interface{}) {
	if e != nil {
		output(fatalColor + errored(true, e, a...) + normColor + "\n")
		os.Exit(-1)
	}
}

// ------------------------------------------------------------------------- //

// set the current debugging level -- values <= this value will be output
func SetLevel(l int) {
	dbgLevel = l
}

// set the current debugging mask -- values masking as non-zero are output
func SetMask(m int) {
	dbgMask = m
}

// conditionally output debug text
func LvlMsg(l int, fstr string, a ...interface{}) {
	if l <= dbgLevel {
		Message(fstr, a...)
	}
}

// conditionally output debug text
func MaskMsg(m int, fstr string, a ...interface{}) {
	if 0 != (m & dbgMask) {
		Message(fstr, a...)
	}
}

// ------------------------------------------------------------------------- //

// output err message if test not true
func ChkTru(tst bool, a ...interface{}) bool {
	if !tst {
		output(errColor + "CHK " + at() + normColor + failed(false, a...) + "\n")
	}
	return !tst
}

// output err message if test not true, then PANIC
func ChkTruP(tst bool, a ...interface{}) {
	if !tst {
		panic(failed(true, a...))
	}
}

// output err message if test not true, then EXIT
func ChkTruX(tst bool, a ...interface{}) {
	if !tst {
		output(errColor + "CHK " + at() + normColor + failed(true, a...) + "\n")
		os.Exit(-1)
	}
}

// output err message if given error isn't nil - returns testable boolean
func ChkErr(e error, a ...interface{}) bool {
	if e != nil {
		output(errColor + "ERR " + at() + normColor + errored(false, e, a...) + "\n")
	}
	return (e != nil)
}

// output err message and PANIC if given error isn't nil
func ChkErrP(e error, a ...interface{}) {
	if e != nil {
		panic(errored(true, e, a...))
	}
}

// output err message and EXIT if given error isn't nil
func ChkErrX(e error, a ...interface{}) {
	if e != nil {
		output(errColor + "ERR " + at() + normColor + errored(true, e, a...) + "\n")
		os.Exit(-1)
	}
}

// output err message if there are any errors in the given list
func ChkErrList(errs []error, a ...interface{}) {
	for _, e := range errs {
		if e != nil {
			output(errColor + "ERR " + at() + normColor + errored(false, e, a...) + "\n")
		}
	}
}

// output err message if error, but ignore any in the 'i' slice
func ChkErrI(e error, i []error, a ...interface{}) bool {
	if e != nil {
		for _, t := range i {
			if t == e {
				return true // error still occured, just not reported
			}
		}
		output(errColor + "ERR " + at() + normColor + errored(false, e, a...) + "\n")
	}
	return (e != nil)
}

// output err message if expected error not matched
func ExpErr(e, x error) bool {
	if e != x {
		output(errColor + "ERR " + at() + normColor + errored(false, e, "Expected error (%v) not given", x) + "\n")
	}
	return (e != x)
}

// ------------------------------------------------------------------------- //

// a quick 'I am here' function for debugging & tracking, takes optional trc_args
func TRC(a ...interface{}) {
	trcAt(a...)
}

// a quick conditional 'I am here' function for debugging & tracking, takes optional trc_args
func TRCIF(c bool, a ...interface{}) {
	if c {
		trcAt(a...)
	}
}

// a quick 'I came from' function for debugging & tracking, takes optional trc_args
func TRCFROM(a ...interface{}) {
	trcBefore(a...)
}

// output a stack trace to aid in debugging
func StackTrace() {
	callers := make([]uintptr, 10)
	d := runtime.Callers(0, callers)
	Message("Depth: %d", d)

	frames := runtime.CallersFrames(callers)
	for {
		frame, more := frames.Next()
		if 0 == frame.Line {
			break
		}
		Warning("  Func: %s - %d   %s", frame.Function, frame.Line, path.Dir(frame.File))
		if !more {
			break
		}
	}
}

// ========================================================================= //

/*
	Possible trc_args sequences:
		"fmtstr", [args]				generates string from "fmtstr" & any args
		error							generates string from error
		nil								outputs 'nil' -- called with a nil error
*/

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
	output(s + "\n")
}

// returns location of CHK caller
func at() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		file = shortName(file)
		return fmt.Sprintf("@ %d in %s  ", line, file)
	}
	return ""
}

// return arg text after calling any possible CLOSER()
func failed(c bool, a ...interface{}) string {
	if len(a) > 0 {
		if c {
			if cl, ok := a[len(a)-1].(func()); ok {
				cl()             // call closer function
				a = a[:len(a)-1] // remove it from arg list
			}
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
