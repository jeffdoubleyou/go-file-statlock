package main

import (
	"fmt"
	"github.com/jeffdoubleyou/go-file-statlock"
)

func main() {
	lock := statlock.NewLock("lockfile.lck", 5)
	if locked, err := lock.Lock(); err != nil {
		fmt.Printf("Could not lock file: %s", err.Error())
	} else {
		fmt.Printf("Locked: %t\n", locked)
	}

	// Do stuff
	if _, err := lock.Unlock(); err != nil {
		panic(err)
	} else {
		fmt.Println("Unlocked!")
	}

	// Wait for 10 iterations of 5 seconds for a lock.  Set to 0 to wait forever...
	wait := statlock.NewLock("lockfile.lck", 5, 10)
	if _, err := wait.WaitForLock(); err != nil {
		panic(err)
	} else {
		fmt.Println("Got a lock!")
		wait.Unlock()
	}

}
