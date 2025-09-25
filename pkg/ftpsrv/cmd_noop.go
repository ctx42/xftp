package ftpsrv

// NOOPSuccess represents the ""NOOP" command success response.
var NOOPSuccess = Resp(200, "NOOP OK.")

// handleNOOP handles "NOOP" command.
func handleNOOP(cc *CtrlCon, args ...string) error {
	if len(args) != 0 {
		return cc.writeLine(ErrorArgNum)
	}
	return cc.writeLine(NOOPSuccess)
}
