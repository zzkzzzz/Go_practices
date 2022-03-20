package main

import (
	"fmt"
	"math/rand"
	"time"
)

// However, using a daemon has some downsides. Since the daemon runs on a goroutine and
// has references to the water factory object, even if all other references to the water
// factory are dropped, the water factory will never be cleaned up as the daemon is still
// holding a reference to it. The daemon goroutine itself also does not go away and leaks.

// We can somewhat remedy this by giving our water factory a Destroy method that shuts down
// the daemon with a Context, but this forces the user to keep track of references to the water factory,
// and to properly clean it up when the last reference is about to be dropped. This means that such a
// water factory effectively needs manual memory management, despite Go being a garbage collected language!

type WaterFactoryWithDaemon struct {
	// Channels for atoms to send their arrival requests
	precomH chan chan struct{}
	precomO chan chan struct{}
}

func NewFactoryWithDaemon() WaterFactoryWithDaemon {
	wfd := WaterFactoryWithDaemon{
		precomH: make(chan chan struct{}),
		precomO: make(chan chan struct{}),
	}

	// Daemon goroutine
	go func() {
		for {
			// Step 1: (Precommit)
			//         Receive arrival requets from 2 hydrogen and 1 oxygen atoms
			h1 := <-wfd.precomH
			h2 := <-wfd.precomH
			o := <-wfd.precomO

			// Step 2: (Commit)
			//         Tell the 3 atoms to start bonding
			h1 <- struct{}{}
			h2 <- struct{}{}
			o <- struct{}{}

			// Step 3: (Postcommit)
			//         Wait until the 3 atoms have finished before looping
			// We re-use the same communication channel as (Commit)
			<-h1
			<-h2
			<-o
		}
	}()

	return wfd
}

func (wfd *WaterFactoryWithDaemon) hydrogen(bond func()) {
	commit := make(chan struct{}) // Step 1: Create private communication channel
	wfd.precomH <- commit         // Step 2: (Precommit)
	<-commit                      // Step 3: (Commit)
	bond()                        // Step 4: Bond
	commit <- struct{}{}          // Step 5: (Postcommit)
}

func (wfd *WaterFactoryWithDaemon) oxygen(bond func()) {
	commit := make(chan struct{}) // Step 1: Create private communication channel
	wfd.precomO <- commit         // Step 2: (Precommit)
	<-commit                      // Step 3: (Commit)
	bond()                        // Step 4: Bond
	commit <- struct{}{}          // Step 5: (Postcommit)
}

///////////////////////////////////////////////////////////////

func TestWaterFactoryWithDaemon() {
	oxygenBond := func() {
		fmt.Println("Bonding oxygen")
		time.Sleep(5 * time.Millisecond)
		fmt.Println("Done")
	}
	hydrogenBond := func() {
		fmt.Println("Bonding hydrogen")
		time.Sleep(5 * time.Millisecond)
		fmt.Println("Done")
	}

	wfd := NewFactoryWithDaemon()
	rand.Seed(time.Now().Unix())
	for i := 0; i < 33; i++ {
		if rand.Intn(3) == 2 {
			go wfd.oxygen(oxygenBond)
		} else {
			go wfd.hydrogen(hydrogenBond)
		}
	}
	time.Sleep(5 * time.Second)
}

func main() { TestWaterFactoryWithDaemon() }
