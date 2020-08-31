# dbg
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
