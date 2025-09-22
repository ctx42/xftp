// Package ftpcmd provides a list of all FTP commands and their options.
package ftpcmd

// List of FTP commands documented in [RFC 959].
//
// [RFC 959]: https://datatracker.ietf.org/doc/html/rfc959.
const (
	ABOR = "ABOR"
	ACCT = "ACCT"
	ALLO = "ALLO"
	APPE = "APPE"
	CDUP = "CDUP"
	CWD  = "CWD"
	DELE = "DELE"
	HELP = "HELP"
	LIST = "LIST"
	MKD  = "MKD"
	MODE = "MODE"
	NLST = "NLST"
	NOOP = "NOOP"
	PASS = "PASS"
	PASV = "PASV"
	PORT = "PORT"
	PWD  = "PWD"
	QUIT = "QUIT"
	REIN = "REIN"
	REST = "REST"
	RETR = "RETR"
	RMD  = "RMD"
	RNFR = "RNFR"
	RNTO = "RNTO"
	SITE = "SITE"
	SMNT = "SMNT"
	STAT = "STAT"
	STOR = "STOR"
	STOU = "STOU"
	STRU = "STRU"
	SYST = "SYST"
	TYPE = "TYPE"
	USER = "USER"
)

// List of FTP commands documented in [RFC 2389].
//
// [RFC 2389]: https://datatracker.ietf.org/doc/html/rfc2389.
const (
	FEAT = "FEAT"
	OPTS = "OPTS"
)

// List of FTP commands documented in [RFC 2228].
//
// [RFC 2228]: https://datatracker.ietf.org/doc/html/rfc2228.
const (
	AUTH = "AUTH"
	PBSZ = "PBSZ"
	PROT = "PROT"
)

// EPSV FTP passive connection command documented in [RFC 2428].
//
// [RFC 2428]: https://datatracker.ietf.org/doc/html/rfc2428.
const EPSV = "EPSV"

// List of FTP commands documented in [RFC 3659].
//
// [RFC 3659]: https://datatracker.ietf.org/doc/html/rfc3659.
const (
	MDTM = "MDTM"
	MLSD = "MLSD"
	MLST = "MLST"
	SIZE = "SIZE"
)

// RFC959 is a list of FTP commands defined in [RFC 959].
//
// [RFC 959]: https://datatracker.ietf.org/doc/html/rfc959
var RFC959 = []string{
	ABOR, ACCT, ALLO,
	APPE, CDUP, CWD,
	DELE, HELP, LIST,
	MKD, MODE, NLST,
	NOOP, PASS, PASV,
	PORT, PWD, QUIT,
	REIN, REST, RETR,
	RMD, RNFR, RNTO,
	SITE, SMNT, STAT,
	STOR, STOU, STRU,
	SYST, TYPE, USER,
}

// Minimal is a list of FTP commands that must be in a minimal FTP server
// implementation, according to [RFC 959].
//
// [RFC 959]: https://datatracker.ietf.org/doc/html/rfc959
var Minimal = []string{
	MODE, NOOP, // PORT, // TODO(rz): implement.
	QUIT, RETR, STOR,
	STRU, TYPE, USER,
}
