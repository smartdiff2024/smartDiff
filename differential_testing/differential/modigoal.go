package differential

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var DTraceResultDirProxy string

func CreateDTraceDir(ctx *cli.Context) {
	if DTraceResultDir == "" {
		panic("Empty DTrace Result Directory!")
	}
	DTraceResultDirProxy = DTraceResultDir + "/" + ctx.Args().Get(0) + ctx.Args().Get(1)
	fmt.Println(DTraceResultDirProxy)
	os.MkdirAll(DTraceResultDirProxy, 0755)
}

func ToDTraceResult(TraceResult *DTraceResult) {
	if DTraceResultDir == "" {
		panic("Empty DTrace Result Directory!")
	}
	var filename string
	if TraceResult.PostImplAdrress == OriVersion {
		filename = "origin"
	} else {
		filename = "modify"
	}
	resultFile, _ := os.OpenFile(DTraceResultDirProxy+"/"+filename+".json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer resultFile.Close()

	writer := bufio.NewWriter(resultFile)
	data, _ := json.Marshal(TraceResult)
	_, err := writer.Write(data)
	_, _ = writer.WriteString("\n")
	writer.Flush()
	if err != nil {
		panic(err)
	}
}
