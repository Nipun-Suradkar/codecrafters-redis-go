package command

const (
	CmdEcho   = "echo"
	CmdPing   = "ping"
	CmdSet    = "set"
	CmdGet    = "get"
	CmdConfig = "config"
	CmdInfo   = "info"

	ConfigDir    = "dir"
	ConfigDbFile = "dbfilename"

	ErrWrongArgCount = "wrong number of arguments for '%s'"
)

var SupportedCommands = []string{
	CmdEcho,
	CmdPing,
	CmdSet,
	CmdGet,
	CmdConfig,
	CmdInfo,
}
