package dbg

import (
	"os"
	"path"
	"runtime"
	"strings"
)

/*
	A collection of debugging and tracking routines and some other convenience functions

	Different argument options for the different functions:
		fmtStr:			fmt.Sprintf format string
		fmt_args:		[fmtStr, [fmt.Sprintf arguments...]]
		chk_args:		[fmt_args,] [CLOSER()]
							any CLOSER() func is called before doing
							the panic or exit for a failure case
		trc_args:		[error] | [nil(error)] | [fmt_args]

	These are simple text output functions that will output colored text
	-- there are multiple versions of these:
		one through the package 'dbg.Info' dbg.Info( [fmt_args] )
		one through the 'Dbg struct{bool}' bug.Info( [fmt_args] )
		one through the 'DbgLvl struct'    dlvl.Info( 5, [fmt_args] )
		one through the 'dbgMsk struct'    dmsk.Info( 0x8, [fmt_args] )

	Echo( [fmt_args] )						output normal text (quick way to do output w/o 'fmt' if you want)
	Note( [fmt_args] )						output colored text (Blue)
	Info( [fmt_args] )						output colored text (Green)
	Message( [fmt_args] )					output colored text (Cyan)
	Warning( [fmt_args] )					output colored text (Orange)
	Caution( [fmt_args] )					output colored text (Yellow)
	Failed( [fmt_args] )					output colored text (Magenta)
	Error( [fmt_args] )						output colored text (Red)
	Danger( [fmt_args] )					output colored text (White on Red)

	Color()									Enable colored text output (if system supports it)
	NoColor()								Disable colored text output

	ExpErr( err, err ) bool					output error if expected error is not given

	ChkTru( bool, [fmt_args] ) bool
		if test value is false, output check failed message (see below)
		 returns TRUE on failure allowing this to be wrapped as part of 'if'

	ChkErr( error, [fmt_args] ) bool
		if error non-nil, output check failed message (see below)
		 returns TRUE on non-nil allowing this to be wrapped as part of 'if'

	ChkErrI( error, []error, [fmt_args]) bool
		if error non-nil, output check failed message (see below) as long
		 as it's not in the ignore list of errors
		 returns TRUE on non-nil allowing this to be wrapped as part of 'if'

	ChkErrList( []error, [fmt_args]) bool
		output check failed message (see below) if there are any non-nil
		 values in the error list
		 returns TRUE on non-nil allowing this to be wrapped as part of 'if'

	ChkTru[PX]( bool, [chk_args] )
		if test value is false, output check failed message (see below)
		 then either Panic or force Exit -- See dbg.Panic below

	ChkErr[PX]( error, [chk_args] )
		if error is non-nil, output check failed message (see below) then
		 either Panic or force Exit -- See dbg.Panic below

	Panic( [chk_args] )						output colored text then PANIC! -- See dbg.Panic below
	Fatal( [chk_args] )						output colored text then force exit
	PanicIf( bool [, chk_args] )			PANIC only if true -- See dbg.Panic below
	FatalIf( bool [, chk_args] )			force exit if true
	PanicIfErr( err [, chk_args] )			PANIC only if err is not nil -- See dbg.Panic below
	FatalIfErr( err [, chk_args] )			force exit if err is not nil

	TRC( [trc_args] )						output calling func file & line number
											 followed by any arg data
	Dbg.TRC()								conditional TRC based off of Dbg flag
	TRCIF( bool [, trc_args] )				conditional TRC based off of given bool
	TRCFROM( [trc_args] )					output func calling func file & line number
											 followed by any arg data
	Dbg.TRCFROM()							conditional TRCFROM based off of Dbg flag

	IAm() string							returns callers func name
	ImAt() string							returns callers file & line number
	WasAt() string							returns callers caller file & line number
	ErrAt() (string, int)					returns callers file & line number
	ErrWasAt() (string, int)				returns callers caller file & line number

	StackTrace()							output call stack (up to ten levels deep)

	NOTE:
	dbg.Panic() always passes the built-in panic a STRING, even if given an error
	If using defer, the returned value for 'recover' will therefore always be a string
*/

type (
	// Debug output that can work off of a simple bool flag
	Dbg struct {
		Enabled bool
	}

	// Debug output that can work off of an output level:
	//	if Level == 0, all output is disabled
	//  if Level > level, output is disabled
	DbgLvl struct {
		Level int
	}

	// Debug output that can work off of an output mask:
	//	if Mask & mask == 0, output is disabled
	DbgMsk struct {
		Mask uint32
	}
)

// dummy func to allow external use / non-use
//	have dbg.Link() at start of file and you can enable / disable dbg code
//	without getting the pesky build errors for import use of non-use
//	-- should remove after any debug along with the import
func Link() {}

// ----------------------------------------------------------------------------

// enable color output for debug text
func Color() {
	normColor = "\033[0m"     // reset to normal text
	msgColor = "\033[36m"     // CYAN
	infoColor = "\033[32m"    // GREEN
	noteColor = "\033[34m"    // BLUE
	warnColor = "\033[33m"    // ORANGE - YELLOW
	ccnColor = "\033[93m"     // YELLOW - BRIGHT ORANGE
	failColor = "\033[35m"    // MAGENTA
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
	failColor = ""
	errColor = ""
	fatalColor = ""
}

// ------------------------------------------------------------------------- //
// Simple output functions that give colored text -- can be redirected to logging if desired

// simply echo to output, no color hilites
func Echo(fstr string, a ...interface{}) {
	output(fstr+"\n", a...)
}

// cyan text to output
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

// orange text to output
func Warning(fstr string, a ...interface{}) {
	output(warnColor+fstr+normColor+"\n", a...)
}

// yellow (bright orange) text to output
func Caution(fstr string, a ...interface{}) {
	output(ccnColor+fstr+normColor+"\n", a...)
}

// magenta text to output
func Failed(fstr string, a ...interface{}) {
	output(failColor+fstr+normColor+"\n", a...)
}

// red text to output
func Error(fstr string, a ...interface{}) {
	output(errColor+fstr+normColor+"\n", a...)
}

// bold white on red background text to output
func Danger(fstr string, a ...interface{}) {
	output(fatalColor+fstr+normColor+"\n", a...)
}

// ------------------------------------------------------------------------- //

// simply echo to output, no color hilites
func (d Dbg) Echo(fstr string, a ...interface{}) {
	if d.Enabled {
		output(fstr+"\n", a...)
	}
}

// cyan text to output
func (d Dbg) Message(fstr string, a ...interface{}) {
	if d.Enabled {
		output(msgColor+fstr+normColor+"\n", a...)
	}
}

// green text to output
func (d Dbg) Info(fstr string, a ...interface{}) {
	if d.Enabled {
		output(infoColor+fstr+normColor+"\n", a...)
	}
}

// blue text to output
func (d Dbg) Note(fstr string, a ...interface{}) {
	if d.Enabled {
		output(noteColor+fstr+normColor+"\n", a...)
	}
}

// orange text to output
func (d Dbg) Warning(fstr string, a ...interface{}) {
	if d.Enabled {
		output(warnColor+fstr+normColor+"\n", a...)
	}
}

// yellow (bright orange) text to output
func (d Dbg) Caution(fstr string, a ...interface{}) {
	if d.Enabled {
		output(ccnColor+fstr+normColor+"\n", a...)
	}
}

// magenta text to output
func (d Dbg) Failed(fstr string, a ...interface{}) {
	if d.Enabled {
		output(failColor+fstr+normColor+"\n", a...)
	}
}

// red text to output
func (d Dbg) Error(fstr string, a ...interface{}) {
	if d.Enabled {
		output(errColor+fstr+normColor+"\n", a...)
	}
}

// bold white on red background text to output
func (d Dbg) Danger(fstr string, a ...interface{}) {
	if d.Enabled {
		output(fatalColor+fstr+normColor+"\n", a...)
	}
}

// ------------------------------------------------------------------------- //

// simply echo to output, no color hilites
func (d DbgLvl) Echo(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(fstr+"\n", a...)
	}
}

// cyan text to output
func (d DbgLvl) Message(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(msgColor+fstr+normColor+"\n", a...)
	}
}

// green text to output
func (d DbgLvl) Info(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(infoColor+fstr+normColor+"\n", a...)
	}
}

// blue text to output
func (d DbgLvl) Note(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(noteColor+fstr+normColor+"\n", a...)
	}
}

// orange text to output
func (d DbgLvl) Warning(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(warnColor+fstr+normColor+"\n", a...)
	}
}

// yellow (bright orange) text to output
func (d DbgLvl) Caution(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(ccnColor+fstr+normColor+"\n", a...)
	}
}

// magenta text to output
func (d DbgLvl) Failed(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(failColor+fstr+normColor+"\n", a...)
	}
}

// red text to output
func (d DbgLvl) Error(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(errColor+fstr+normColor+"\n", a...)
	}
}

// bold white on red background text to output
func (d DbgLvl) Danger(l int, fstr string, a ...interface{}) {
	if d.Level > 0 && d.Level <= l {
		output(fatalColor+fstr+normColor+"\n", a...)
	}
}

// ------------------------------------------------------------------------- //

// simply echo to output, no color hilites
func (d DbgMsk) Echo(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(fstr+"\n", a...)
	}
}

// cyan text to output
func (d DbgMsk) Message(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(msgColor+fstr+normColor+"\n", a...)
	}
}

// green text to output
func (d DbgMsk) Info(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(infoColor+fstr+normColor+"\n", a...)
	}
}

// blue text to output
func (d DbgMsk) Note(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(noteColor+fstr+normColor+"\n", a...)
	}
}

// orange text to output
func (d DbgMsk) Warning(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(warnColor+fstr+normColor+"\n", a...)
	}
}

// yellow (bright orange) text to output
func (d DbgMsk) Caution(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(ccnColor+fstr+normColor+"\n", a...)
	}
}

// magenta text to output
func (d DbgMsk) Failed(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(failColor+fstr+normColor+"\n", a...)
	}
}

// red text to output
func (d DbgMsk) Error(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(errColor+fstr+normColor+"\n", a...)
	}
}

// bold white on red background text to output
func (d DbgMsk) Danger(m uint32, fstr string, a ...interface{}) {
	if 0 != d.Mask&m {
		output(fatalColor+fstr+normColor+"\n", a...)
	}
}

// ------------------------------------------------------------------------- //
// These functions do not allow the 'closer' func as they always return

// output err message if expected error not matched
func ExpErr(e, x error) bool {
	if e != x {
		output("%s\n", errColor + "ERR " + at() + normColor + errored(false, e, "Expected error (%v) not given", x))
	}
	return (e != x)
}

// output err message if test not true
func ChkTru(tst bool, a ...interface{}) bool {
	if !tst {
		output("%s\n", failColor + "CHK " + at() + normColor + failed(false, a...))
	}
	return !tst
}

// output err message if given error isn't nil - returns testable boolean
func ChkErr(e error, a ...interface{}) bool {
	if e != nil {
		output("%s\n", errColor + "ERR " + at() + normColor + errored(false, e, a...))
	}
	return (e != nil)
}

// output err message if error, but ignore (don't output) any in the 'i' slice
func ChkErrI(e error, i []error, a ...interface{}) bool {
	if e != nil {
		for _, t := range i {
			if t == e {
				return true // error still occured, just not reported
			}
		}
		output("%s\n", errColor + "ERR " + at() + normColor + errored(false, e, a...))
	}
	return (e != nil)
}

// output err message if there are any errors in the given list
func ChkErrList(errs []error, a ...interface{}) bool {
	failed := false
	for _, e := range errs {
		if e != nil {
			output("%s\n", errColor + "ERR " + at() + normColor + errored(false, e, a...))
			failed = true
		}
	}
	return failed
}

// ------------------------------------------------------------------------- //
// These functions can work with a 'closer'

// output err message if test not true, then PANIC
func ChkTruP(tst bool, a ...interface{}) {
	if !tst {
		panic(failed(true, a...))
	}
}

// output err message if test not true, then EXIT
func ChkTruX(tst bool, a ...interface{}) {
	if !tst {
		output("%s\n", failColor + "CHK " + at() + normColor + failed(true, a...))
		os.Exit(-1)
	}
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
		output("%s\n", errColor + "ERR " + at() + normColor + errored(true, e, a...))
		os.Exit(-1)
	}
}

// panic with any optional chk_args
func Panic(a ...interface{}) {
	panic(genText_Closer(a...))
}

// fatal error (exit) with any optional chk_args
func Fatal(a ...interface{}) {
	output("%s\n", fatalColor + genText_Closer(a...) + normColor)
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
		output("%s\n", fatalColor + failed(true, a...) + normColor)
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
		output("%s\n", fatalColor + errored(true, e, a...) + normColor)
		os.Exit(-1)
	}
}

// ------------------------------------------------------------------------- //

// a quick 'I am here' function for debugging & tracking, takes optional trc_args
func TRC(a ...interface{}) {
	trcAt(a...)
}

// use Dbg interface for TRC
func (d Dbg) TRC(a ...interface{}) {
	if d.Enabled {
		trcAt(a...)
	}
}

// a quick conditional 'I am here' function for debugging & tracking, takes optional trc_args
//  Remove because we now have (b Dbg) TRC?
func TRCIF(b bool, a ...interface{}) {
	if b {
		trcAt(a...)
	}
}

// a quick 'I came from' function for debugging & tracking, takes optional trc_args
func TRCFROM(a ...interface{}) {
	trcBefore(a...)
}

// use Dbg interface for TRCFROM
func (d Dbg) TRCFROM(a ...interface{}) {
	if d.Enabled {
		trcBefore(a...)
	}
}

// ------------------------------------------------------------------------- //
// Some simple utility routines

// return the callers func name
func IAm() string {
	pc := make([]uintptr, 4)
	runtime.Callers(2, pc)
	nm := runtime.FuncForPC(pc[0]).Name()
	return nm[strings.LastIndex(nm, ".")+1:]
}

func IWas() string {
	pc := make([]uintptr, 4)
	runtime.Callers(3, pc)
	nm := runtime.FuncForPC(pc[0]).Name()
	return nm[strings.LastIndex(nm, ".")+1:]
}

// a quick func to output location information (file & line#)
func ImAt() string {
	return funcAt(1)
}

// a quick func to output location information (file & line#)
func WasAt() string {
	return funcAt(2)
}

// return the callers location information (file & line#)
func ErrAt() (string, int) {
	if _, file, line, ok := runtime.Caller(1); ok {
		return shortName(file), line
	}
	return "", 0
}

// return the callers caller location information (file & line#)
func ErrWasAt() (string, int) {
	if _, file, line, ok := runtime.Caller(2); ok {
		return shortName(file), line
	}
	return "", 0
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
