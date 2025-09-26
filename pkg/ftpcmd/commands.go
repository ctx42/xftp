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

// List of FTP commands documented in [RFC 2228].
//
// [RFC 2228]: https://datatracker.ietf.org/doc/html/rfc2228.
const (
	AUTH = "AUTH"
	PBSZ = "PBSZ"
	PROT = "PROT"
)

// List of FTP commands documented in [RFC 2389].
//
// [RFC 2389]: https://datatracker.ietf.org/doc/html/rfc2389.
const (
	FEAT = "FEAT"
	OPTS = "OPTS"
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

// TODO(rz):

// // commands is a list of all FTP commands.
// //
// // [RFC 959]: https://datatracker.ietf.org/doc/html/rfc959
// // [RFC 2389]: https://datatracker.ietf.org/doc/html/rfc2389.
// // [RFC 2228]: https://datatracker.ietf.org/doc/html/rfc2228.
// // [RFC 2428]: https://datatracker.ietf.org/doc/html/rfc2428.
// // [RFC 3659]: https://datatracker.ietf.org/doc/html/rfc3659.
// var commands = []Command{
// 	{ABOR, "RFC 959", Base},
// 	{ACCT, "RFC 959", Base},
// 	{ALLO, "RFC 959", Base},
// 	{APPE, "RFC 959", Base},
// 	{CDUP, "RFC 959", Base},
// 	{CWD, "RFC 959", Base},
// 	{DELE, "RFC 959", Base},
// 	{HELP, "RFC 959", Base},
// 	{LIST, "RFC 959", Base},
// 	{MKD, "RFC 959", Base},
// 	{MODE, "RFC 959", Base | Minimal},
// 	{NLST, "RFC 959", Base},
// 	{PASS, "RFC 959", Base},
// 	{PASV, "RFC 959", Base},
// 	{PORT, "RFC 959", Base | Minimal},
// 	{PWD, "RFC 959", Base},
// 	{QUIT, "RFC 959", Base | Minimal},
// 	{REIN, "RFC 959", Base},
// 	{REST, "RFC 959", Base},
// 	{RETR, "RFC 959", Base | Minimal},
// 	{RMD, "RFC 959", Base},
// 	{RNFR, "RFC 959", Base},
// 	{RNTO, "RFC 959", Base},
// 	{SITE, "RFC 959", Base},
// 	{SMNT, "RFC 959", Base},
// 	{STAT, "RFC 959", Base},
// 	{STOR, "RFC 959", Base | Minimal},
// 	{STOU, "RFC 959", Base},
// 	{STRU, "RFC 959", Base | Minimal},
// 	{SYST, "RFC 959", Base},
// 	{TYPE, "RFC 959", Base | Minimal},
// 	{USER, "RFC 959", Base | Minimal},
// 	{AUTH, "RFC 2228", Extended},
// 	{PBSZ, "RFC 2228", Extended},
// 	{PROT, "RFC 2228", Extended},
// 	{FEAT, "RFC 2389", Extended},
// 	{OPTS, "RFC 2389", Extended},
// 	{EPSV, "RFC 2428", Extended},
// 	{MDTM, "RFC 3659", Extended},
// 	{MLSD, "RFC 3659", Extended},
// 	{MLST, "RFC 3659", Extended},
// 	{SIZE, "RFC 3659", Extended},
// }
