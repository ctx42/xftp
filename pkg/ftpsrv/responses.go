package ftpsrv

// Session connection responses.
var (
	ServerReady = Resp(220, "FTP Server ready.")
	Goodbye     = Resp(221, "Goodbye.")

	ErrorUnkCmd = Resp(500, "Unknown command %s.")
)
