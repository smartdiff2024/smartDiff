package research

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/differential"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var (
	WorkersFlag = cli.IntFlag{
		Name:  "workers",
		Usage: "Number of worker threads that execute in parallel",
		Value: 1,
	}
	SkipTransferTxsFlag = cli.BoolFlag{
		Name:  "skip-transfer-txs",
		Usage: "Skip executing transactions that only transfer ETH",
	}
	SkipCallTxsFlag = cli.BoolFlag{
		Name:  "skip-call-txs",
		Usage: "Skip executing CALL transactions to accounts with contract bytecode",
	}
	SkipCreateTxsFlag = cli.BoolFlag{
		Name:  "skip-create-txs",
		Usage: "Skip executing CREATE transactions",
	}
)

type SubstateTaskFunc func(block uint64, tx int, substate *Substate, taskPool *SubstateTaskPool, hasModifiedContract bool) error

type SubstateTaskPool struct {
	Name     string
	TaskFunc SubstateTaskFunc

	First uint64
	Last  uint64

	Workers         int
	SkipTransferTxs bool
	SkipCallTxs     bool
	SkipCreateTxs   bool

	Ctx *cli.Context // CLI context required to read additional flags

	DB *SubstateDB
}

func NewSubstateTaskPool(name string, taskFunc SubstateTaskFunc, first, last uint64, ctx *cli.Context) *SubstateTaskPool {
	return &SubstateTaskPool{
		Name:     name,
		TaskFunc: taskFunc,

		First: first,
		Last:  last,

		Workers:         ctx.Int(WorkersFlag.Name),
		SkipTransferTxs: ctx.Bool(SkipTransferTxsFlag.Name),
		SkipCallTxs:     ctx.Bool(SkipCallTxsFlag.Name),
		SkipCreateTxs:   ctx.Bool(SkipCreateTxsFlag.Name),

		Ctx: ctx,

		DB: staticSubstateDB,
	}
}

// ExecuteBlock function iterates on substates of a given block call TaskFunc
func (pool *SubstateTaskPool) ExecuteBlock(block uint64) (numTx int64, err error) {
	fmt.Println("replay: ", block)
	var value differential.BlockProxy
	if differential.UpdatePosition {
		if _, ok := differential.UpdatePositions[block]; ok && len(differential.UpdatePositions[block].UpdatePositionStructs) != 0 {
			differential.ImplBlockNum = block
		} else {
			return numTx, nil
		}
	}
	if differential.GetByteCode {
		var ok bool
		value, ok = differential.BlockImpl[block]
		if !ok {
			return 0, nil
		}
	}
	unsubstates := pool.DB.GetBlockSubstates(block)
	substateLength := len(unsubstates)
	for tx := 0; tx < substateLength; tx++ {

		unsubstate := unsubstates[tx]
		alloc := unsubstate.InputAlloc
		msg := unsubstate.Message

		to := msg.To
		if pool.SkipTransferTxs && to != nil {
			// skip regular transactions (ETH transfer)
			if account, exist := alloc[*to]; !exist || len(account.Code) == 0 {
				continue
			}
		}
		if pool.SkipCallTxs && to != nil {
			// skip CALL trasnactions with contract bytecode
			if account, exist := alloc[*to]; exist && len(account.Code) > 0 {
				continue
			}
		}
		if pool.SkipCreateTxs && to == nil {
			// skip CREATE transactions
			continue
		}

		if differential.Impls {
			if needtx := NeedImpl(unsubstate.OutputAlloc); !needtx {
				continue
			}
			differential.ImplBlockNum = block
			differential.ImplPositionNum = tx
		}
		if differential.UpdatePosition {
			differential.ImplPositionNum = tx
		}

		if differential.GetByteCode {
			addresses, ok := value[tx]
			if ok {
				implAddress := addresses.ImplAddress
				for _, proxyaddress := range addresses.ProxyAddress {
					(*differential.ProxyByte[proxyaddress])[implAddress].CreateBin = unsubstate.Message.Data
					if _, ok := unsubstate.OutputAlloc[implAddress]; ok {
						(*differential.ProxyByte[proxyaddress])[implAddress].RuntimeBin = unsubstate.OutputAlloc[implAddress].Code
					}
				}
			} else {
				continue
			}
		}

		if differential.DTrace {
			if needtx := NeedTx(unsubstate.OutputAlloc); !needtx {
				continue
			}

			// smartDiff will be executed at the transaction level,
			// with the scope of execution being from the transaction at FirstPosition in the FirstBlock to the transaction at LastPosition in the LastBlock.
			if block == differential.FirstBlock && tx < differential.FirstPosi {
				continue
			}
			if block == differential.LastBlock && tx > differential.LastPosi {
				return numTx, nil
			}
			logrus.Info("replay: ", block, tx)
			ProxyIncludeCheck(unsubstate.OutputAlloc)
			ImplIncludeCheck(unsubstate.OutputAlloc)

			for _, modiversion := range differential.ModiVersions {
				differential.ModiVersion = modiversion
				substate := pool.DB.GetSubstate(block, tx)
				substate, _ = CreateOrCallCheck(substate, unsubstate)
				errcheck := pool.TaskFunc(block, tx, substate, pool, true)
				if errcheck != nil {
					return numTx, fmt.Errorf("stage1-substate: transitionSubstateTransaction %v_%v: %v in modiversion: %v", block, tx, errcheck, modiversion)
				}
			}
			// smartDiff will execute the original transaction. Replacing the modiVersion into the oriVersion is convenient for the execution of the TODtraceResult()
			differential.ModiVersion = differential.OriVersion
		}
		err = pool.TaskFunc(block, tx, unsubstate, pool, false)
		if differential.UpdatePosition {
			if _, ok := differential.UpdatePositions[block]; ok && len(differential.UpdatePositions[block].UpdatePositionStructs) != 0 {
				if tx == substateLength-1 {
					for _, update := range differential.UpdatePositions[block].UpdatePositionStructs {
						if !update.HasPosition {
							differential.ProxyInits[update.ProxyAddress].ImplAddresses = append(differential.ProxyInits[update.ProxyAddress].ImplAddresses, differential.ImplMap{
								ImplAddress: update.LogicAddress,
								BlockNum:    block,
								PositionNum: substateLength,
							})
						}
					}
				}
			}
		}
		if err != nil {
			return numTx, fmt.Errorf("%s: %v_%v: %v", pool.Name, block, tx, err)
		}

		numTx++
	}

	return numTx, nil
}

// Execute function spawns worker goroutines and schedule tasks.
func (pool *SubstateTaskPool) Execute() error {
	start := time.Now()

	var totalNumBlock, totalNumTx int64
	defer func() {
		duration := time.Since(start) + 1*time.Nanosecond
		sec := duration.Seconds()

		nb, nt := atomic.LoadInt64(&totalNumBlock), atomic.LoadInt64(&totalNumTx)
		blkPerSec := float64(nb) / sec
		txPerSec := float64(nt) / sec
		fmt.Printf("%s: block range = %v %v\n", pool.Name, pool.First, pool.Last)
		fmt.Printf("%s: total #block = %v\n", pool.Name, nb)
		fmt.Printf("%s: total #tx    = %v\n", pool.Name, nt)
		fmt.Printf("%s: %.2f blk/s, %.2f tx/s\n", pool.Name, blkPerSec, txPerSec)
		fmt.Printf("%s done in %v\n", pool.Name, duration.Round(1*time.Millisecond))
	}()

	numWorker := 1
	runtime.GOMAXPROCS(1)
	fmt.Printf("stage1-substate: TransitionSubstate: #CPU = %v, #worker = %v, #thread = %v\n", runtime.NumCPU(), numWorker, runtime.GOMAXPROCS(0))

	fmt.Printf("%s: block range = %v %v\n", pool.Name, pool.First, pool.Last)
	fmt.Printf("%s: #CPU = %v, #worker = %v\n", pool.Name, runtime.NumCPU(), pool.Workers)

	for block := pool.First; block <= pool.Last; block++ {
		nt, err := pool.ExecuteBlock(block)
		if err != nil {
			fmt.Println("Execute error: ", err)
		}
		atomic.AddInt64(&totalNumTx, nt)
		atomic.AddInt64(&totalNumBlock, 1)
		var lastSec float64
		var lastNumBlock, lastNumTx int64
		duration := time.Since(start) + 1*time.Nanosecond
		sec := duration.Seconds()
		if block == pool.Last {
			nb, nt := atomic.LoadInt64(&totalNumBlock), atomic.LoadInt64(&totalNumTx)
			blkPerSec := float64(nb-lastNumBlock) / (sec - lastSec)
			txPerSec := float64(nt-lastNumTx) / (sec - lastSec)
			fmt.Printf("%s: elapsed time: %v, number = %v\n", pool.Name, duration.Round(1*time.Millisecond), block)
			fmt.Printf("%s: %.2f blk/s, %.2f tx/s\n", pool.Name, blkPerSec, txPerSec)

			lastSec, lastNumBlock, lastNumTx = sec, nb, nt
		}
	}

	return nil
}
