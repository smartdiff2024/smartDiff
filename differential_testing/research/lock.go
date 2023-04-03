package research

import (
	"fmt"
	"sync"
)

var lockMap sync.Map
var lockcondMap map[sync.Locker]*sync.Cond

func init() {
	// fmt.Println("invoke sync.cond init")
	lockcondMap = make(map[sync.Locker]*sync.Cond)
}

func lockByName(name string) {
	lock, _ := lockMap.LoadOrStore(name, new(sync.Mutex))
	lock.(*sync.Mutex).Lock()
}

func unlockByName(name string) {
	lock, exist := lockMap.Load(name)
	if !exist {
		panic("Trying to unlock a nonexistent mutex by name: " + name)
	}
	lock.(*sync.Mutex).Unlock()
}

func lockbyAddr(name string, nonce uint64, noncereal uint64) bool {
	fmt.Println("try to lock: ", name)
	lock, _ := lockMap.LoadOrStore(name, new(sync.Mutex))
	fmt.Println("test1")
	if _, exist := lockcondMap[lock.(*sync.Mutex)]; !exist {
		lockcondMap[lock.(*sync.Mutex)] = sync.NewCond(lock.(*sync.Mutex))
	}
	fmt.Println("test2")
	lock.(*sync.Mutex).Lock()
	fmt.Println("test3")
	if nonce != noncereal {
		fmt.Println("enter lock by nonce")
		return false
		lockcondMap[lock.(*sync.Mutex)].Wait()
		return true
	}
	return false
}

// func unlockByAddr(name string) {
// 	fmt.Println("try to unlock ", name)
// 	lock, exist := lockMap.Load(name)
// 	if !exist {
// 		panic("Trying to unlock a nonexistent mutex by name: " + name)
// 	}
// 	lock.(*sync.Mutex).Unlock()
// 	lockcondMap[lock.(*sync.Mutex)].Broadcast()
// }

// func UnlockSubstate(inalloc *SubstateAlloc) {
// 	var addrlist Addrlist
// 	for addr := range *inalloc {
// 		if ok := HasModiSubstate(addr); ok {
// 			addrlist = append(addrlist, addr)
// 		}
// 	}
// 	sort.Sort(addrlist)
// 	for _, addr := range addrlist {
// 		unlockByAddr(addr.String())
// 	}
// }
