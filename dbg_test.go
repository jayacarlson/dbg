package dbg

import (
	"errors"
	"testing"
)

var panicErr = errors.New("MyPanicErr")
var myErr = errors.New("MyErr")

func chkCloser() {
	Info("Closer was called")
}

func chkNonCloser() {
	Danger("Closer was called when it shouldn't have been")
}

func panicCatcher() {
	rcv := recover()
	if rcv != nil {
		str, _ := rcv.(string)
		Echo("panic:  '%s'", str)
	}
}

func panicFail() {
	rcv := recover()
	if rcv != nil {
		Error("SHOULD NOT BE HERE")
	}
}

func panicTest1() {
	defer panicCatcher()
	Panic("Panic, supplied text", chkCloser)
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
	Echo("Echo text")
	Note("Note text")
	Info("Info text")
	Message("Message text")
	Warning("Warning text")
	Caution("Caution text")
	Failed("Failed text")
	Error("Error text")
	Danger("Danger text")

	bug := Dbg{true}
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
	bug.Echo("bug{false} Echo text")
	bug.Note("bug{false} Note text")
	bug.Info("bug{false} Info text")
	bug.Message("bug{false} Message text")
	bug.Warning("bug{false} Warning text")
	bug.Caution("bug{false} Caution text")
	bug.Failed("bug{false} Failed text")
	bug.Error("bug{false} Error text")
	bug.Danger("bug{false} Danger text")

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
	lvl.Echo(1, "lvl{5} 1 - Echo text") // should not be output...
	lvl.Note(2, "lvl{5} 2 - Note text")
	lvl.Info(3, "lvl{5} 3 - Info text")
	lvl.Message(4, "lvl{5} 4 - Message text")
	lvl.Warning(5, "lvl{5} 5 - Warning text") // should be output...
	lvl.Caution(6, "lvl{5} 6 - Caution text")
	lvl.Failed(7, "lvl{5} 7 - Failed text")
	lvl.Error(8, "lvl{5} 8 - Error text")
	lvl.Danger(9, "lvl{5} 9 - Danger text")

	msk := DbgMsk{0xa}
	msk.Echo(1, "msk{xa} 1 - Echo text")       // should not be output
	msk.Note(2, "msk{xa} 2 - Note text")       // should be output
	msk.Info(3, "msk{xa} 3 - Info text")       // should be output
	msk.Message(4, "msk{xa} 4 - Message text") // should not be output
	msk.Warning(5, "msk{xa} 5 - Warning text") // should not be output
	msk.Caution(6, "msk{xa} 6 - Caution text") // should be output
	msk.Failed(7, "msk{xa} 7 - Failed text")   // should be output
	msk.Error(8, "msk{xa} 8 - Error text")     // should be output
	msk.Danger(9, "msk{xa} 9 - Danger text")   // should be output

	panicTest1()
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

	Note("Should see 'CHK @ ### in dbg/dbg_test.go  Check failed'")
	ChkTru(false)
	Note("Should see 'CHK @ ### in dbg/dbg_test.go  My check text'")
	ChkTru(false, "My check text")
	//	ChkTruX(false)

	Note("Should see 'ERR @ ### in dbg/dbg_test.go  MyErr'")
	ChkErr(myErr)
	Note("Should see 'ERR @ ### in dbg/dbg_test.go  My error text'")
	ChkErr(myErr, "My error text")
	//	ChkErrX(myErr)

	//	Fatal("Fatal")
	//	FatalIf(true)
	//	FatalIfErr(myErr)
}
