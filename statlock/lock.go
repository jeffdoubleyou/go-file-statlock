package statlock

import (
	"fmt"
	"os"
	"time"
)

const (
	StatusLocked   = 1
	StatusWaiting  = 2
	StatusUnlocked = 3
	StatusOrphaned = 4
	StatusTimeout  = 5
	StatusFailed   = 6
	StatusRelease  = 99
)

type Lock struct {
	Path            string
	Duration        int
	Status          int
	MaxWaitInterval int
	file            *os.File
}

type StatLock interface {
	Lock() (bool, error)
	Unlock() (bool, error)
	WaitForLock() (bool, error)
	openLock() (bool, error)
	keepLock()
	refreshLock() (bool, error)
}

func NewLock(path string, seconds int, maxWaitInterval ...int) *Lock {
	waitInterval := 10
	if len(maxWaitInterval) == 1 {
		waitInterval = maxWaitInterval[0]
	}

	return &Lock{
		Path:            path,
		Duration:        seconds,
		MaxWaitInterval: waitInterval,
	}
}

func (l *Lock) Lock() (bool, error) {
	if _, err := l.openLock(); err != nil {
		return false, err
	}
	l.Status = StatusLocked
	go l.keepLock()
	return true, nil
}

func (l *Lock) Unlock() (bool, error) {
	l.Status = StatusRelease
	unlockAttempts := 0
	maxAttempts := 10
	for l.Status == StatusRelease && unlockAttempts < maxAttempts {
		time.Sleep(time.Duration(500) * time.Millisecond)
		unlockAttempts++
	}
	if l.Status == StatusUnlocked {
		return true, nil
	}
	return false, fmt.Errorf("unable to release lock err %d", l.Status)
}

func (l *Lock) WaitForLock() (bool, error) {
	l.Status = StatusWaiting
	unlockAttempts := 0

	for l.Status != StatusLocked && (unlockAttempts < l.MaxWaitInterval || l.MaxWaitInterval == 0) {
		if locked, err := l.Lock(); locked {
			return true, nil
		} else {
			fmt.Println(err.Error())
		}
		unlockAttempts++
		time.Sleep(time.Duration(l.Duration) * time.Second)
	}
	l.Status = StatusTimeout
	return false, fmt.Errorf("timeout waiting for lock")
}

func (l *Lock) openLock() (bool, error) {
	if fileInfo, err := os.Stat(l.Path); err != nil {
		if file, err := os.Create(l.Path); err != nil {
			return false, err
		} else {
			l.file = file
			return true, nil
		}
	} else {
		now := time.Now().Local()
		if fileInfo.ModTime().Local().Add(time.Duration(l.Duration*2) * time.Second).Before(now) {
			return l.refreshLock()
		} else {
			l.Status = StatusFailed
			return false, fmt.Errorf("could not take ownership of lock")
		}
	}
}

func (l *Lock) keepLock() {
	for l.Status != StatusRelease {
		time.Sleep(time.Duration(l.Duration) * time.Second)
		if _, err := l.refreshLock(); err != nil {
			l.Status = StatusRelease
			fmt.Printf("Unable to update lock file: %s", err.Error())
		}
	}

	if _, err := os.Stat(l.Path); err == nil {
		if err := os.Remove(l.Path); err != nil {
			fmt.Printf("could not remove lock file '%s': %s\n", l.Path, err.Error())
			l.Status = StatusOrphaned
		} else {
			l.Status = StatusUnlocked
		}
	} else if os.IsNotExist(err) { // This could occur if another process is waiting for a lock and acquires it while we are unlocking
		fmt.Printf("lock file '%s' does not exist\n", l.Path)
		l.Status = StatusUnlocked
	} else {
		fmt.Printf("failed to access lock file '%s': %s\n", l.Path, err.Error())
		l.Status = StatusOrphaned
	}
}

func (l *Lock) refreshLock() (bool, error) {
	t := time.Now().Local()
	if err := os.Chtimes(l.Path, t, t); err != nil {
		return false, err
	}
	return true, nil
}
