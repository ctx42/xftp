package ftpsrv

// Control connection responses.
var (
	ServerReady = Resp(220, "FTP Server ready.")
	Goodbye     = Resp(221, "Goodbye.")
)
