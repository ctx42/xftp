package ftpsrv

// NOOPSuccess represents the ""NOOP" command success response.
var NOOPSuccess = Resp(200, "NOOP OK.")

type CmdNOOP struct{}

func (cmd CmdNOOP) Name() string { return "NOOP" }
func (cmd CmdNOOP) Flags() Flags { return FlagBase | FlagMinimal }

func (cmd CmdNOOP) Handle(cc *CtrlCon, args ...string) error {
	if len(args) != 0 {
		return cc.writeLine(ErrorArgNum)
	}
	return cc.writeLine(NOOPSuccess)
}
