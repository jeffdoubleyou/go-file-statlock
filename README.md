# go-file-statlock
File locking based on age of lock file.  Basic file lock for NFS filesystem.

## Why

I am using this to handle services that run on multiple systems that need to know if another process is running on a different machine without any other access to the other system outside of a shared NFS filesystem.

Normal file locking isn't supported on NFS and I could not find any suitable library to handle this.

## Usage

```
go get github.com/jeffdoubleyou/go-file-statlock/statlock
```

```
package main

import (
	"fmt"

	"github.com/jeffdoubleyou/go-file-statlock/statlock"
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
```
