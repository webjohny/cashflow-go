package objects

import (
	"sync"
	"time"
)

type MutexTimeout struct {
	Mutex    sync.Mutex
	LastTime time.Time
}

func NewMutexTimeout() *MutexTimeout {
	return &MutexTimeout{}
}

func (tm *MutexTimeout) TryLock(cooldown time.Duration) bool {
	tm.Mutex.Lock()
	defer tm.Mutex.Unlock()

	if time.Since(tm.LastTime) < cooldown {
		return false
	}

	tm.LastTime = time.Now()
	return true
}
