package replay

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/differential"
	"github.com/ethereum/go-ethereum/differential/util"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/research"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

// record-replay: substate-cli replay command
var ReplayCommand = cli.Command{
	Action:    replayAction,
	Name:      "replay",
	Usage:     "executes full state transitions and check output consistency",
	ArgsUsage: "<blockNumFirst> <blockNumLast>",
	Flags: []cli.Flag{
		&research.WorkersFlag,
		&research.SkipTransferTxsFlag,
		&research.SkipCallTxsFlag,
		&research.SkipCreateTxsFlag,
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
		research.SubstateDirFlag,
		&differential.FirstPositionFlag,
		&differential.LastPositionFlag,
		&differential.ImplBlockFlag,
		&differential.ImplPostionFlag,
	},
	Description: `
The substate-cli replay command requires two arguments:
<blockNumFirst> <blockNumLast>

<blockNumFirst> and <blockNumLast> are the first and
last block of the inclusive range of blocks to replay transactions.`,
}

// replayTask replays a transaction substate
func replayTask(block uint64, tx int, substate *research.Substate, taskPool *research.SubstateTaskPool, hasNewContract bool) error {
	inputAlloc := substate.InputAlloc
	inputEnv := substate.Env
	inputMessage := substate.Message

	outputAlloc := substate.OutputAlloc
	outputResult := substate.Result

	var (
		vmConfig    vm.Config
		chainConfig *params.ChainConfig
		getTracerFn func(txIndex int, txHash common.Hash) (tracer vm.EVMLogger, err error)
	)

	vmConfig = vm.Config{}

	chainConfig = &params.ChainConfig{}
	*chainConfig = *params.MainnetChainConfig
	// disable DAOForkSupport, otherwise account states will be overwritten
	chainConfig.DAOForkSupport = false

	getTracerFn = func(txIndex int, txHash common.Hash) (tracer vm.EVMLogger, err error) {
		return nil, nil
	}

	var hashError error
	getHash := func(num uint64) common.Hash {
		if inputEnv.BlockHashes == nil {
			hashError = fmt.Errorf("getHash(%d) invoked, no blockhashes provided", num)
			return common.Hash{}
		}
		h, ok := inputEnv.BlockHashes[num]
		if !ok {
			hashError = fmt.Errorf("getHash(%d) invoked, blockhash for that block not provided", num)
		}
		return h
	}

	var (
		statedb   = MakeOffTheChainStateDB(inputAlloc)
		gaspool   = new(core.GasPool)
		blockHash = common.Hash{0x01}
		txHash    = common.Hash{0x02}
		txIndex   = tx
	)

	gaspool.AddGas(inputEnv.GasLimit)
	blockCtx := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Coinbase:    inputEnv.Coinbase,
		BlockNumber: new(big.Int).SetUint64(inputEnv.Number),
		Time:        new(big.Int).SetUint64(inputEnv.Timestamp),
		Difficulty:  inputEnv.Difficulty,
		GasLimit:    inputEnv.GasLimit,
		GetHash:     getHash,
	}
	// If currentBaseFee is defined, add it to the vmContext.
	if inputEnv.BaseFee != nil {
		blockCtx.BaseFee = new(big.Int).Set(inputEnv.BaseFee)
	}

	msg := inputMessage.AsMessage()

	tracer, err := getTracerFn(txIndex, txHash)
	if err != nil {
		return err
	}
	vmConfig.Tracer = tracer
	vmConfig.Debug = (tracer != nil)
	statedb.Prepare(txHash, txIndex)

	txCtx := vm.TxContext{
		GasPrice: msg.GasPrice(),
		Origin:   msg.From(),
	}

	evm := vm.NewEVM(blockCtx, txCtx, statedb, chainConfig, vmConfig)
	snapshot := statedb.Snapshot()
	var msgResult *core.ExecutionResult
	// gas overflow is not a err in this condition, but it result in an execution failure
	if hasNewContract && outputResult.Status == types.ReceiptStatusSuccessful {
		msgResult, err = core.ApplyMessage2(evm, msg, gaspool)
	} else {
		msgResult, err = core.ApplyMessage(evm, msg, gaspool)
	}
	if err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}

	if hashError != nil {
		return hashError
	}

	if chainConfig.IsByzantium(blockCtx.BlockNumber) {
		statedb.Finalise(true)
	} else {
		statedb.IntermediateRoot(chainConfig.IsEIP158(blockCtx.BlockNumber))
	}

	evmResult := &research.SubstateResult{}
	var result bool
	if msgResult.Failed() {
		evmResult.Status = types.ReceiptStatusFailed
		result = false
	} else {
		evmResult.Status = types.ReceiptStatusSuccessful
		result = true
	}
	evmResult.Logs = statedb.GetLogs(txHash, blockHash)
	evmResult.Bloom = types.BytesToBloom(types.LogsBloom(evmResult.Logs))
	if to := msg.To(); to == nil {
		evmResult.ContractAddress = crypto.CreateAddress(evm.TxContext.Origin, msg.Nonce())
	}
	evmResult.GasUsed = msgResult.UsedGas

	evmAlloc := statedb.ResearchPostAlloc

	// update the storage and record dtrace result
	if differential.DTrace {
		if hasNewContract {
			dtraceresult := differential.NewDTraceResult(block, tx, result)
			if block > differential.ImplBlock || (block == differential.ImplBlock && tx >= differential.ImplPosi) {
				research.UpdateStorage(&evmAlloc)
				if differential.ProxyInclude {
					if _, ok := evmAlloc[differential.ProxyAddress]; ok {
						dtraceresult.ProxyStorage = (*evmAlloc[differential.ProxyAddress]).Storage
					} else {
						logrus.Info("Something wrong about proxyaddress ", differential.ProxyAddress, " in impl ", differential.ModiVersion, " in block_tx: ", block, tx)
						unsubstate := taskPool.DB.GetSubstate(block, tx)
						research.DtraceResultLogrus(unsubstate, substate, &evmAlloc)
					}
				}
			} else {
				logrus.Info("Skip update Proxyaddress Storage: ", block, "_", tx)
			}
			if differential.ImplInclude {
				if _, ok := evmAlloc[differential.OriVersion]; ok {
					dtraceresult.Storage = (*evmAlloc[differential.OriVersion]).Storage
				} else {
					logrus.Info("Something wrong about modiaddress ", differential.ModiVersion, " in block_tx: ", block, tx)
					unsubstate := taskPool.DB.GetSubstate(block, tx)
					research.DtraceResultLogrus(unsubstate, substate, &evmAlloc)
				}
				differential.ToDTraceResult(dtraceresult)
			}
			return nil
		} else {
			dtraceresult := differential.NewDTraceResult(block, tx, result)
			if block > differential.ImplBlock || (block == differential.ImplBlock && tx >= differential.ImplPosi) {
				research.UpdateStorage(&evmAlloc)
				if differential.ProxyInclude {
					if _, ok := evmAlloc[differential.ProxyAddress]; ok {
						dtraceresult.ProxyStorage = (*evmAlloc[differential.ProxyAddress]).Storage
					} else {
						logrus.Info("Something wrong about proxyaddress ", differential.ProxyAddress, " in ori ", differential.OriVersion, " in block_tx: ", block, tx)
						unsubstate := taskPool.DB.GetSubstate(block, tx)
						research.DtraceResultLogrus(unsubstate, substate, &evmAlloc)
					}
				}
			} else {
				logrus.Info("Skip record dtrace Proxyaddress Storage: ", block, "_", tx)
			}
			if differential.ImplInclude {
				if _, ok := evmAlloc[differential.OriVersion]; ok {
					dtraceresult.Storage = (*evmAlloc[differential.OriVersion]).Storage
				} else {
					logrus.Info("Something wrong about ori address ", differential.OriVersion, " in block_tx: ", block, tx)
					unsubstate := taskPool.DB.GetSubstate(block, tx)
					research.DtraceResultLogrus(unsubstate, substate, &evmAlloc)
				}
				differential.ToDTraceResult(dtraceresult)
			}
		}
	}

	r := outputResult.Equal(evmResult)
	a := outputAlloc.Equal(evmAlloc)
	if !(r && a) {

		if !r {
			fmt.Printf("inconsistent output: result\n")
		}
		if !a {
			fmt.Printf("inconsistent output: alloc\n")
		}
		var jbytes []byte
		jbytes, _ = json.MarshalIndent(inputAlloc, "", " ")
		fmt.Printf("inputAlloc:\n%s\n", jbytes)
		jbytes, _ = json.MarshalIndent(inputEnv, "", " ")
		fmt.Printf("inputEnv:\n%s\n", jbytes)
		jbytes, _ = json.MarshalIndent(inputMessage, "", " ")
		fmt.Printf("inputMessage:\n%s\n", jbytes)
		jbytes, _ = json.MarshalIndent(outputAlloc, "", " ")
		fmt.Printf("outputAlloc:\n%s\n", jbytes)
		jbytes, _ = json.MarshalIndent(evmAlloc, "", " ")
		fmt.Printf("evmAlloc:\n%s\n", jbytes)
		jbytes, _ = json.MarshalIndent(outputResult, "", " ")
		fmt.Printf("outputResult:\n%s\n", jbytes)
		jbytes, _ = json.MarshalIndent(evmResult, "", " ")
		fmt.Printf("evmResult:\n%s\n", jbytes)
		return fmt.Errorf("inconsistent output")
	}
	return nil
}

// record-replay: func replayAction for replay command
func replayAction(ctx *cli.Context) error {

	var err error

	if ctx.Args().Len() != 2 {
		return fmt.Errorf("substate-cli replay command requires exactly 2 arguments")
	}

	first, ferr := strconv.ParseInt(ctx.Args().Get(0), 10, 64)
	last, lerr := strconv.ParseInt(ctx.Args().Get(1), 10, 64)
	if ferr != nil || lerr != nil {
		return fmt.Errorf("substate-cli replay: error in parsing parameters: block number not an integer")
	}
	if first < 0 || last < 0 {
		return fmt.Errorf("substate-cli replay: error: block number must be greater than 0")
	}
	if first > last {
		return fmt.Errorf("substate-cli replay: error: first block has larger number than last block")
	}

	research.SetSubstateFlags(ctx)
	research.OpenSubstateDBReadOnly()
	defer research.CloseSubstateDB()

	// Initialize the flags to determine the purpost and required parameters for a particular replay.
	//  There are several types of step control symbols, including:
	// Impls: get the logic contract addresses of the proxy address
	// UpdatePosition: get the position in block of the upgrade transaction
	//  GetByteCode: get the bytecode of the logic contracts
	// DTrace: differential test
	differential.Impls = ctx.Bool(differential.ImplsFlag.Name)
	differential.UpdatePosition = ctx.Bool(differential.UpadtePositionFlag.Name)
	differential.GetByteCode = ctx.Bool(differential.GetByteCodeFlag.Name)
	differential.DTrace = ctx.Bool(differential.DTraceFlag.Name)

	differential.DirStoreJSONS = ctx.String(differential.DirStoreJSONsFlag.Name)
	differential.FileToReadModi = ctx.String(differential.FilePathFlag.Name)
	differential.CSVfileToStoreResult = ctx.String(differential.ResultCSVfileFlag.Name)
	differential.DTraceFile = ctx.String(differential.DTraceFileFlag.Name)
	differential.DTraceResultDir = ctx.String(differential.DTraceResultDirFlag.Name)
	differential.OriImpl = ctx.String(differential.DTraceOriImplFlag.Name)
	differential.ReplaceImpl = ctx.String(differential.DTraceReplaceFlag.Name)

	if differential.Impls {
		differential.ReadProxyList()
	}
	if differential.UpdatePosition {
		differential.UpdatePositionRead()
	}
	if differential.GetByteCode {
		differential.ReadBytecodeCSV()
	}
	if differential.DTrace {
		logDir := "./dtrace_log/" + ctx.String(differential.DTraceFileFlag.Name)
		os.MkdirAll(logDir, 0755)
		logFile, errlog := os.OpenFile(logDir+"/"+ctx.String(differential.DTraceOriImplFlag.Name)+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if errlog != nil {
			fmt.Println("something wrong happen in create logFile: ", errlog)
		}
		defer logFile.Close()
		output := io.Writer(logFile)
		logrus.SetOutput(output)
		logrus.SetFormatter(&util.MyFormatter)
		logrus.SetLevel(logrus.InfoLevel)
		differential.OriVersion = common.HexToAddress(differential.OriImpl)
		differential.UpdateVersions(differential.ReplaceImpl)
		differential.ProxyAddress = common.HexToAddress(differential.DTraceFile)
		differential.ModiVersion = differential.ModiVersions[0]
		differential.FirstBlock = uint64(first)
		differential.FirstPosi = ctx.Int(differential.FirstPositionFlag.Name)
		differential.LastBlock = uint64(last)
		differential.LastPosi = ctx.Int(differential.LastPositionFlag.Name)
		differential.ImplBlock = ctx.Uint64(differential.ImplBlockFlag.Name)
		differential.ImplPosi = ctx.Int(differential.ImplPostionFlag.Name)
		differential.CreateDTraceDir(ctx)
		if differential.DTraceFile == "" {
			differential.DirReadBytecodeJSONS()
		} else {
			differential.FileReadBytecodeJSONS()

		}
	}

	taskPool := research.NewSubstateTaskPool("substate-cli replay", replayTask, uint64(first), uint64(last), ctx)
	err = taskPool.Execute()
	if differential.Impls {
		differential.WriteProxyList()
	}
	if differential.UpdatePosition {
		differential.UpdatePositionWrite()
	}
	if differential.GetByteCode {
		differential.WriteBytecodeJSON()
	}
	if differential.DTrace {
		differential.ModiVersion = differential.ModiVersions[0]
		differential.UpdateBytecodeJSON()
	}
	return err
}
