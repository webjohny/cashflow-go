package objects

import (
	"strconv"
	"sync"
	"time"
)

type MutexMap struct {
	mutexes sync.Map
}

func (lm *MutexMap) LockMethodRace(methodKey string, raceId uint64, tm time.Duration) bool {
	mtx, _ := lm.mutexes.LoadOrStore(methodKey+strconv.Itoa(int(raceId)), NewMutexTimeout())
	timedMutex := mtx.(*MutexTimeout)

	return timedMutex.TryLock(tm)
}

func (lm *MutexMap) UnlockMethodRace(methodKey string, raceId uint64) {
	key := methodKey + strconv.Itoa(int(raceId))
	_, ok := lm.mutexes.Load(key)
	if !ok {
		return
	}

	lm.mutexes.Delete(key) // Remove the mutex from the map
}
