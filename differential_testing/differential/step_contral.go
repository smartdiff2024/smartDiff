package differential

import cli "github.com/urfave/cli/v2"

var (
	ImplsFlag = cli.BoolFlag{
		Name:  "impls",
		Usage: "Determines whether to get the first implementation address",
		Value: false,
	}
	UpadtePositionFlag = cli.BoolFlag{
		Name:  "updatep",
		Usage: "Get the update position of proxy update",
		Value: false,
	}
	GetByteCodeFlag = cli.BoolFlag{
		Name:  "getcode",
		Usage: "Determines whether to get the bytecode of impl addresses",
	}
	DTraceFlag = cli.BoolFlag{
		Name:  "dtrace",
		Usage: "Determines whether to replay impl versions to get traces",
	}
	FilePathFlag = &cli.StringFlag{
		Name:  "withoutimpl",
		Usage: "The file to do read",
		Value: "",
	}
	ResultCSVfileFlag = &cli.StringFlag{
		Name:  "withimpls",
		Usage: "the file is used to save result",
		Value: "",
	}
	DirStoreJSONsFlag = &cli.StringFlag{
		Name:  "jsondir",
		Usage: "the directory is used to save json result",
		Value: "",
	}
	DTraceFileFlag = &cli.StringFlag{
		Name:  "dtracefiles",
		Usage: "the files that save impl jsons, split by ',' ",
		Value: "",
	}
	DTraceResultDirFlag = &cli.StringFlag{
		Name:  "dtraceresultdir",
		Usage: "the directory that save trace result",
		Value: "",
	}
	DTraceOriImplFlag = &cli.StringFlag{
		Name:  "oriimpl",
		Usage: "the goal address that to be replaced",
		Value: "",
	}
	DTraceReplaceFlag = &cli.StringFlag{
		Name:  "replaceimpl",
		Usage: "the addresses that replace original address, split by ',' ",
		Value: "",
	}
	FirstPositionFlag = cli.IntFlag{
		Name:  "firstp",
		Usage: "the num imply first position of the first block to restore proxyaddress storage",
		Value: 0,
	}
	LastPositionFlag = cli.IntFlag{
		Name:  "lastp",
		Usage: "the num imply the position of the last block to restore proxyaddress storage",
		Value: 0,
	}
	ImplBlockFlag = cli.Uint64Flag{
		Name:  "implblock",
		Usage: "the num imply the block changing the impl version",
		Value: 0,
	}
	ImplPostionFlag = cli.IntFlag{
		Name:  "implp",
		Usage: "the num imply the position changing the impl version",
		Value: 0,
	}
	Impls          bool
	UpdatePosition bool
	GetByteCode    bool
	DTrace         bool

	FileToReadModi       string
	CSVfileToStoreResult string
	DirStoreJSONS        string
	DTraceFile           string
	DTraceResultDir      string

	OriImpl     string
	ReplaceImpl string

	FirstBlock uint64
	FirstPosi  int
	LastBlock  uint64
	LastPosi   int
	ImplBlock  uint64
	ImplPosi   int
)
