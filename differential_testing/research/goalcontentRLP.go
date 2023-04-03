package research

// import (
// 	"io"
// 	"sort"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/rlp"
// )

// type GoalExchangeRLP struct {
// 	GoalBlock       uint64
// 	CodeCreateHash  common.Hash
// 	CodeRunTimeHash common.Hash
// 	Storage         [][2]common.Hash
// 	Pnonce          uint64
// }

// func NewGoalExchangeRLP(goal *GoalExchange) *GoalExchangeRLP {
// 	// fmt.Println("enter goalexchange newgoal")
// 	var goalExchangeRLP GoalExchangeRLP
// 	goalExchangeRLP.GoalBlock = goal.Goalblock
// 	goalExchangeRLP.CodeCreateHash = goal.CodeHashCreate()
// 	goalExchangeRLP.CodeRunTimeHash = goal.CodeHashCall()
// 	goalExchangeRLP.Pnonce = goal.Pnonce
// 	sortedKeys := []common.Hash{}
// 	for key := range goal.Storage {
// 		sortedKeys = append(sortedKeys, key)
// 	}
// 	sort.Slice(sortedKeys, func(i, j int) bool {
// 		return sortedKeys[i].Big().Cmp(sortedKeys[j].Big()) < 0
// 	})
// 	for _, key := range sortedKeys {
// 		value := goal.Storage[key]
// 		goalExchangeRLP.Storage = append(goalExchangeRLP.Storage, [2]common.Hash{key, value})
// 	}
// 	return &goalExchangeRLP
// }

// func (goal GoalExchange) EncodeRLP(w io.Writer) error {
// 	// fmt.Println("enter goalexchange encoderlp")
// 	return rlp.Encode(w, NewGoalExchangeRLP(&goal))
// }

// func (goal *GoalExchange) SetRLP(goalRLP *GoalExchangeRLP) {
// 	// fmt.Println("enter goalexchange setrlp")
// 	goal.Goalblock = goalRLP.GoalBlock
// 	goal.CodeCreate = GetModiCode(goalRLP.CodeCreateHash)
// 	goal.CodeRunTime = GetModiCode(goalRLP.CodeRunTimeHash)
// 	goal.Storage = make(map[common.Hash]common.Hash)
// 	goal.Pnonce = goalRLP.Pnonce
// 	for i := range goalRLP.Storage {
// 		goal.Storage[goalRLP.Storage[i][0]] = goalRLP.Storage[i][1]
// 	}
// }

// func (goal *GoalExchange) DecodeRLP(s *rlp.Stream) error {
// 	// fmt.Println("enter goalexchange decoderlp")
// 	var err error
// 	var goalRLP GoalExchangeRLP

// 	err = s.Decode(&goalRLP)
// 	if err != nil {
// 		return err
// 	}

// 	goal.SetRLP(&goalRLP)

// 	return nil
// }
