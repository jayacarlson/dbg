package dbg

import (
	"errors"
	"flag"
	"testing"
)

var (
	panicErr = errors.New("MyPanicErr")
	myErr    = errors.New("MyErr")

	tx, ex        bool
	fx, fif, ferr bool
	cdx           bool
)

func init() {
	flag.BoolVar(&tx, "tx", false, "Test ChkTruX")
	flag.BoolVar(&ex, "ex", false, "Test ChkErrX")
	flag.BoolVar(&fx, "fx", false, "Test Fatal")
	flag.BoolVar(&fif, "fif", false, "Test FatalIf")
	flag.BoolVar(&ferr, "ferr", false, "Test FatalIfErr")
	flag.BoolVar(&cdx, "cdx", false, "Test countdown dbg")
}

func chkCloser() {
	Message("Closer was called")
}

func chkNonCloser() {
	Danger("Closer was called when it shouldn't have been")
}

func panicCatcher() {
	rcv := recover()
	if rcv != nil {
		err, _ := rcv.(error)
		Echo("panic:  '%s'", err.Error())
	}
}

func panicFail() {
	rcv := recover()
	if rcv != nil {
		Error("SHOULD NOT BE HERE")
	}
}

func panicTest1a() {
	defer panicCatcher()
	Panic("Panic, supplied text", chkCloser)
}

func panicTest1b() {
	defer panicCatcher()
	Panic("Panic, %s", "supplied text", chkCloser)
}

func panicTest2() {
	defer panicCatcher()
	Panic(panicErr, chkCloser)
}

func panicTest3t() {
	defer panicCatcher()
	PanicIf(true, "test true")
}

func panicTest3f() {
	defer panicFail()
	PanicIf(false, "test false")
	Info("panicTest3f good")
}

func panicTest4e() {
	defer panicCatcher()
	PanicIfErr(panicErr, "test error")
}

func panicTest4n() {
	defer panicFail()
	PanicIfErr(nil, "test nil")
	Info("panicTest4n good")
}

func panicTest5() {
	defer panicCatcher()
	ChkTruP(false, chkCloser)
}

func panicTest6() {
	defer panicCatcher()
	ChkTruP(false, "ChkTruP: Failed text", chkCloser)
}

func panicTest7() {
	defer panicCatcher()
	ChkErrP(myErr, chkCloser)
}

func panicTest8() {
	defer panicCatcher()
	ChkErrP(myErr, "ChkErrP: Failed text", chkCloser)
}

func panicTest9() {
	defer panicFail()
	ChkTruP(true, "Ooops", chkNonCloser)
	Info("panicTest9 good")
}

func panicTest10() {
	defer panicFail()
	ChkErrP(nil, "Ooops", chkNonCloser)
	Info("panicTest10 good")
}

func TestDbg(t *testing.T) {
	flag.Parse()

	bug := Dbg{}

	if tx {
		ChkTruX(false, "False, exiting")
	}
	if ex {
		ChkErrX(myErr, "Error, exiting")
	}
	if fx {
		Fatal("Fatal, exiting")
	}
	if fif {
		FatalIf(true, "True, exiting")
	}
	if ferr {
		FatalIfErr(myErr, "Error, exiting")
	}
	if cdx {
		bug.Enabled = true
		bug.MaxOut = 4

		bug.Echo("3")
		bug.Echo("2")
		bug.Echo("1")
		bug.Echo("0")
	}
	if tx || ex || fx || fif || ferr || cdx {
		ERROR("--- Well, we shouldn't be here...")
	}

	Echo("Echo text")
	Note("Note text")
	Info("Info text")
	Message("Message text")
	Warning("Warning text")
	Caution("Caution text")
	Failed("Failed text")
	Error("Error text")
	Danger("Danger text")
	WARNING("WARNING text")
	CAUTION("CAUTION text")
	ERROR("ERROR text")
	FAULT("FAULT text")

	bug.Enabled = true // debug messages should be visable
	bug.Echo("bug{true} Echo text")
	bug.Note("bug{true} Note text")
	bug.Info("bug{true} Info text")
	bug.Message("bug{true} Message text")
	bug.Warning("bug{true} Warning text")
	bug.Caution("bug{true} Caution text")
	bug.Failed("bug{true} Failed text")
	bug.Error("bug{true} Error text")
	bug.Danger("bug{true} Danger text")

	bug.Enabled = false // debug messages should not be visable
	bug.Echo("bug{false} Echo text FAILED")
	bug.Note("bug{false} Note text FAILED")
	bug.Info("bug{false} Info text FAILED")
	bug.Message("bug{false} Message text FAILED")
	bug.Warning("bug{false} Warning text FAILED")
	bug.Caution("bug{false} Caution text FAILED")
	bug.Failed("bug{false} Failed text FAILED")
	bug.Error("bug{false} Error text FAILED")
	bug.Danger("bug{false} Danger text FAILED")

	lvl := DbgLvl{0}
	lvl.Echo(1, "lvl{0} 1 - Echo text") // should not be output...
	lvl.Note(2, "lvl{0} 2 - Note text")
	lvl.Info(3, "lvl{0} 3 - Info text")
	lvl.Message(4, "lvl{0} 4 - Message text")
	lvl.Warning(5, "lvl{0} 5 - Warning text")
	lvl.Caution(6, "lvl{0} 6 - Caution text")
	lvl.Failed(7, "lvl{0} 7 - Failed text")
	lvl.Error(8, "lvl{0} 8 - Error text")
	lvl.Danger(9, "lvl{0} 9 - Danger text")

	lvl.Level = 5
	lvl.Echo(1, "lvl{5} 1 - Echo text") // should be output...
	lvl.Note(2, "lvl{5} 2 - Note text")
	lvl.Info(3, "lvl{5} 3 - Info text")
	lvl.Message(4, "lvl{5} 4 - Message text")
	lvl.Warning(5, "lvl{5} 5 - Warning text")
	lvl.Caution(6, "lvl{5} 6 - Caution text FAILED") // should not be output...
	lvl.Failed(7, "lvl{5} 7 - Failed text FAILED")
	lvl.Error(8, "lvl{5} 8 - Error text FAILED")
	lvl.Danger(9, "lvl{5} 9 - Danger text FAILED")

	msk := DbgMsk{0xA}
	msk.Echo(1, "msk{xA} 1 - Echo text FAILED")       // should not be output
	msk.Note(2, "msk{xA} 2 - Note text")              // should be output
	msk.Info(3, "msk{xA} 3 - Info text")              // should be output
	msk.Message(4, "msk{xA} 4 - Message text FAILED") // should not be output
	msk.Warning(5, "msk{xA} 5 - Warning text FAILED") // should not be output
	msk.Caution(6, "msk{xA} 6 - Caution text")        // should be output
	msk.Failed(7, "msk{xA} 7 - Failed text")          // should be output
	msk.Error(8, "msk{xA} 8 - Error text")            // should be output
	msk.Danger(9, "msk{xA} 9 - Danger text")          // should be output

	panicTest1a()
	panicTest1b()
	panicTest2()
	panicTest3t()
	panicTest3f()
	panicTest4e()
	panicTest4n()
	panicTest5()
	panicTest6()
	panicTest7()
	panicTest8()
	panicTest9()
	panicTest10()

	Message("Should see\nCHK @ ### in dbg/dbg_test.go  Check failed")
	ChkTru(false)
	Message("Should see\nCHK @ ### in dbg/dbg_test.go  My check text")
	ChkTru(false, "My check text")

	Message("Should see\nERR @ ### in dbg/dbg_test.go  MyErr")
	ChkErr(myErr)
	Message("Should see\nERR @ ### in dbg/dbg_test.go  My error text")
	ChkErr(myErr, "My error text")
}
