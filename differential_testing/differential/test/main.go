package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/differential"
	"github.com/ethereum/go-ethereum/research"
	cli "github.com/urfave/cli/v2"
)

var (
	app = &cli.App{
		Name:  "flags",
		Usage: "flag example",
		Flags: []cli.Flag{
			&differential.ImplsFlag,
			&differential.UpadtePositionFlag,
			&differential.GetByteCodeFlag,
			differential.FilePathFlag,
			differential.ResultCSVfileFlag,
			differential.DirStoreJSONsFlag,
			&differential.DTraceFlag,
			differential.DTraceFileFlag,
			differential.DTraceResultDirFlag,
			differential.DTraceOriImplFlag,
			differential.DTraceReplaceFlag,
		},
		Action: UpdateTest,
	}
)

func ProxyTest(ctx *cli.Context) error {
	differential.FileToReadModi = ctx.String(differential.FilePathFlag.Name)
	differential.Impls = ctx.Bool(differential.ImplsFlag.Name)
	differential.CSVfileToStoreResult = ctx.String(differential.ResultCSVfileFlag.Name)
	if differential.Impls {
		differential.ReadProxyList()
		for index, proxyinit := range differential.ProxyInits {
			fmt.Println(proxyinit)
			differential.ProxyInits[index].ImplAddresses = []differential.ImplMap{{ImplAddress: common.HexToAddress("0x4402"), BlockNum: 10}}
		}
		differential.WriteProxyList()
	}
	return nil
}

func UpdateTest(ctx *cli.Context) error {
	differential.FileToReadModi = ctx.String(differential.FilePathFlag.Name)
	differential.UpdatePosition = ctx.Bool(differential.UpadtePositionFlag.Name)
	differential.CSVfileToStoreResult = ctx.String(differential.ResultCSVfileFlag.Name)
	if differential.UpdatePosition {
		differential.UpdatePositionRead()
		for index := range differential.ProxyInits {
			differential.ProxyInits[index].ImplAddresses = []differential.ImplMap{{ImplAddress: common.HexToAddress("0x4402"), BlockNum: 10, PositionNum: 11}}
		}
		differential.UpdatePositionWrite()
	}
	return nil
}

func GetByteCodeTest(ctx *cli.Context) error {
	differential.CSVfileToStoreResult = ctx.String(differential.ResultCSVfileFlag.Name)
	differential.DirStoreJSONS = ctx.String(differential.DirStoreJSONsFlag.Name)
	differential.DTraceFile = ctx.String(differential.DTraceFileFlag.Name)
	differential.GetByteCode = ctx.Bool(differential.GetByteCodeFlag.Name)
	if differential.GetByteCode {
		differential.ReadBytecodeCSV()
		for proxy, impls := range differential.ProxyByte {
			bytecodes := *impls
			for impl, stage := range bytecodes {
				(*differential.ProxyByte[proxy])[impl].CreateBin = research.TurnCodeIntoBytes("63714402")
				(*differential.ProxyByte[proxy])[impl].RuntimeBin = research.TurnCodeIntoBytes("63714402")
				fmt.Println("proxyAddress: ", proxy, "impl: ", impl, "bytecode: ", stage)
			}
		}
		differential.WriteBytecodeJSON()
	}
	return nil
}

func DTraceTest(ctx *cli.Context) error {
	fmt.Println("ctx: ", ctx.Args())
	differential.DTrace = ctx.Bool(differential.DTraceFlag.Name)
	differential.DirStoreJSONS = ctx.String(differential.DirStoreJSONsFlag.Name)
	differential.DTraceFile = ctx.String(differential.DTraceFileFlag.Name)
	differential.DTraceResultDir = ctx.String(differential.DTraceResultDirFlag.Name)
	differential.OriImpl = ctx.String(differential.DTraceOriImplFlag.Name)
	differential.ReplaceImpl = ctx.String(differential.DTraceReplaceFlag.Name)
	differential.OriVersion = common.HexToAddress(differential.OriImpl)
	differential.UpdateVersions(differential.ReplaceImpl)
	differential.ProxyAddress = common.HexToAddress(differential.DTraceFile)

	if differential.DTrace {
		fmt.Println("differential Dtrace")
		differential.CreateDTraceDir(ctx)
		if differential.DTraceFile == "" {
			fmt.Println("DTrace Files by Dir")
			differential.DirReadBytecodeJSONS()
		} else {
			fmt.Println("DTrace files by Files")
			differential.FileReadBytecodeJSONS()
		}
		// differential.BytecodeToGoal()
		// fmt.Println("After BytecodeToGoal")

		for proxy, impls := range differential.ProxyByte {
			bytecodes := *impls
			for impl, stage := range bytecodes {
				fmt.Println("proxyAddress: ", proxy, "impl: ", impl, "bytecode: ", stage)
				(*differential.ProxyByte[proxy])[impl].Storage = map[common.Hash]common.Hash{
					common.HexToHash("0x123"): common.HexToHash("0x456"),
				}
				fmt.Println("second time::proxyAddress: ", proxy, "impl: ", impl, "bytecode: ", stage)
			}
		}

		differential.UpdateBytecodeJSON()

		for _, modiversion := range differential.ModiVersions {
			differential.ModiVersion = modiversion
			dtraceResult := differential.NewDTraceResult(1, 2, false)
			fmt.Println("differential.ProxyAddress: ", differential.ProxyAddress, "modiversion: ", differential.ModiVersion)
			dtraceResult.Storage = (*differential.ProxyByte[differential.ProxyAddress])[modiversion].Storage
			differential.ToDTraceResult(dtraceResult)
		}

	}
	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		code := 1
		fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}
}
