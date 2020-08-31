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

func TestMain(t *testing.T) {
	Echo("Echo text")
	Note("Note text")
	Info("Info text")
	Message("Message text")
	Warning("Warning text")
	Caution("Caution text")
	Error("Error text")
	Danger("Danger text")

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

	ChkTru(false)
	ChkTru(false, "My check text")
	//	ChkTruX(false)

	ChkErr(myErr)
	ChkErr(myErr, "My error text")
	//	ChkErrX(myErr)

	//	Fatal("Fatal")
	//	FatalIf(true)
	//	FatalIfErr(myErr)
}
