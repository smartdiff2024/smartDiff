package research

import (
	"github.com/sirupsen/logrus"
)

func DtraceResultLogrus(oriSubstate *Substate, modiSubstate *Substate, evmalloc *SubstateAlloc) {
	logrus.Info("\nOri InputAlloc")
	oriinalloc := oriSubstate.InputAlloc
	for addr := range oriinalloc {
		logrus.Info("addr: ", addr.String())
		logrus.Info("nonce: ", oriinalloc[addr].Nonce, " Balance:", oriinalloc[addr].Balance)
		for storagehash := range oriinalloc[addr].Storage {
			logrus.Info("key: ", storagehash.String(), " value: ", oriinalloc[addr].Storage[storagehash].String())
		}
		if len(oriinalloc[addr].Code) != 0 {
			logrus.Info("ori input code: ", TurnBytesIntoCode(oriinalloc[addr].Code))
		}
	}
	logrus.Info("\n Modi InputAlloc")
	modiinalloc := modiSubstate.InputAlloc
	for addr := range modiinalloc {
		logrus.Info("addr: ", addr.String())
		logrus.Info("nonce: ", modiinalloc[addr].Nonce, " Balance:", modiinalloc[addr].Balance)
		for storagehash := range modiinalloc[addr].Storage {
			logrus.Info("key: ", storagehash.String(), " value: ", modiinalloc[addr].Storage[storagehash].String())
		}
		if len(modiinalloc[addr].Code) != 0 {
			logrus.Info("modi input code: ", TurnBytesIntoCode(modiinalloc[addr].Code))
		}
	}
	logrus.Info("\nOri OutputAlloc")
	orioutalloc := oriSubstate.OutputAlloc
	for addr := range orioutalloc {
		logrus.Info("addr: ", addr.String())
		logrus.Info("nonce: ", orioutalloc[addr].Nonce, " Balance:", orioutalloc[addr].Balance)
		for storagehash := range orioutalloc[addr].Storage {
			logrus.Info("key: ", storagehash.String(), " value: ", orioutalloc[addr].Storage[storagehash].String())
		}
		if len(orioutalloc[addr].Code) != 0 {
			logrus.Info("ori output code: ", TurnBytesIntoCode(orioutalloc[addr].Code))
		}
	}
	logrus.Info("\n Modi OutputAlloc")
	modioutalloc := *evmalloc
	for addr := range modioutalloc {
		logrus.Info("addr: ", addr.String())
		logrus.Info("nonce: ", modioutalloc[addr].Nonce, " Balance:", modioutalloc[addr].Balance)
		for storagehash := range modioutalloc[addr].Storage {
			logrus.Info("key: ", storagehash.String(), " value: ", modioutalloc[addr].Storage[storagehash].String())
		}
		if len(modioutalloc[addr].Code) != 0 {
			logrus.Info("modi output code: ", TurnBytesIntoCode(modioutalloc[addr].Code))
		}
	}
}
