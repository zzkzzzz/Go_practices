package main

import (
	"fmt"
	"math/rand"
	"time"
)

type WaterFactoryWithLeader struct {
	oxygenMutex chan struct{}
	precomH     chan chan struct{}
	commit      chan chan struct{}
}

//  Using oxygen atoms as leader goroutines
// In this case, it is easy to avoid creating a daemon goroutine by electing a “leader” that the
// ther atoms will communicate with. Since there is only one oxygen atom in a water molecule, we can make oxygen atoms leaders.

// To ensure that there’s never two leaders active at the same time, we can simply use a mutex.

func NewFactoryWithLeader() WaterFactoryWithLeader {
	wf := WaterFactoryWithLeader{
		oxygenMutex: make(chan struct{}, 1),
		precomH:     make(chan chan struct{}),
		commit:      make(chan chan struct{}),
	}
	wf.oxygenMutex <- struct{}{}
	return wf
}

func (wf *WaterFactoryWithLeader) hydrogen(bond func()) {
	commit := make(chan struct{}) // Step 1: Create private communication channel
	wf.precomH <- commit          // Step 2: (Precommit)
	<-commit                      // Step 3: (Commit)
	bond()                        // Step 4: Bond
	commit <- struct{}{}          // Step 5: (Postcommit)
}

func (wf *WaterFactoryWithLeader) oxygen(bond func()) {
	// Step 1: Become leader
	<-wf.oxygenMutex // For fun, we can use a channel as a mutex

	// Step 2: (Precommit)
	//         Receive arrival requets from 2 hydrogen atoms
	h1 := <-wf.precomH
	h2 := <-wf.precomH

	// Step 3: (Commit)
	//         Tell the 2 hydrogen atoms to start bonding
	h1 <- struct{}{}
	h2 <- struct{}{}

	// Step 4: Bond
	bond()

	// Step 5: (Postcommit)
	//         Wait until the 2 hydrogen atoms to finish
	// We re-use the same communication channel as (Commit)
	<-h1
	<-h2

	// Step 6: Step down from being leader
	wf.oxygenMutex <- struct{}{}
}

func TestWaterFactoryWithLeader() {
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

	wf := NewFactoryWithLeader()
	for i := 0; i < 33; i++ {
		if rand.Intn(3) == 2 {
			go wf.oxygen(oxygenBond)
		} else {
			go wf.hydrogen(hydrogenBond)
		}
	}
	time.Sleep(5 * time.Second)
}

func main() { TestWaterFactoryWithLeader() }
