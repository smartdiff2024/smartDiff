package differential

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

type UpdatePositionStruct struct {
	ProxyAddress common.Address
	LogicAddress common.Address
	Position     int
	HasPosition  bool
}

type UpdatePositionStructs struct {
	UpdatePositionStructs []UpdatePositionStruct
}

var UpdatePositions map[uint64]*UpdatePositionStructs

func UpdatePositionRead() {
	ProxyInits = make(map[common.Address]*ProxyInitStruct)
	UpdatePositions = make(map[uint64]*UpdatePositionStructs)
	if FileToReadModi == "" {
		panic("Emtpy File")
	}
	csvFile, _ := os.Open(FileToReadModi)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	reader.FieldsPerRecord = -1
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		initblocknum, _ := strconv.Atoi(line[6])
		blockPosition, _ := strconv.Atoi(line[7])
		knownblock, _ := strconv.Atoi(line[9])
		proxyaddress := common.HexToAddress(line[5])
		ProxyInits[proxyaddress] = &ProxyInitStruct{
			Initblocknum:  uint64(initblocknum),
			BlockPosition: blockPosition,
			KnownBlock:    uint64(knownblock),
		}

		lineLength := len(line)
		for linenum := 8; linenum < lineLength; linenum += 5 {
			UpdateBlockNumTemp, _ := strconv.Atoi(line[linenum+1])

			//
			if _, ok := UpdatePositions[uint64(UpdateBlockNumTemp)]; !ok {
				UpdatePositions[uint64(UpdateBlockNumTemp)] = &UpdatePositionStructs{}
			}
			UpdatePositions[uint64(UpdateBlockNumTemp)].UpdatePositionStructs = append(UpdatePositions[uint64(UpdateBlockNumTemp)].UpdatePositionStructs, UpdatePositionStruct{
				ProxyAddress: proxyaddress,
				LogicAddress: common.HexToAddress(line[linenum]),
			})
		}
	}
}

func UpdatePositionWrite() {
	if FileToReadModi == "" {
		panic("Emtpy File")
	}
	csvFile, _ := os.Open(FileToReadModi)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	reader.FieldsPerRecord = -1

	File, _ := os.OpenFile(CSVfileToStoreResult, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	WriterCsv := csv.NewWriter(bufio.NewWriter(File))
	defer File.Close()

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		proxyaccount := ProxyInits[common.HexToAddress(line[5])]
		if len(proxyaccount.ImplAddresses) == 0 {
			WriterCsv.Write(line)
		} else {
			writerString := []string{}
			writerString = append(writerString, line[0:8]...)
			for implLength := 0; implLength < len(ProxyInits[common.HexToAddress(line[5])].ImplAddresses); implLength++ {
				writerString = append(writerString, proxyaccount.ImplAddresses[implLength].ImplAddress.String())
				writerString = append(writerString, strconv.FormatUint(proxyaccount.ImplAddresses[implLength].BlockNum, 10))
				writerString = append(writerString, strconv.Itoa(proxyaccount.ImplAddresses[implLength].PositionNum))
				writerString = append(writerString, []string{"-1", "-1"}...)
			}
			WriterCsv.Write(writerString)
		}
	}
	WriterCsv.Flush()
}
