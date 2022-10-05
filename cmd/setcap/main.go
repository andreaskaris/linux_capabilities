package main

import (
	"log"

	"kernel.org/pub/linux/libs/security/libcap/cap"
)

// printCap prints the current capabilities in a human friendly way.
func printCap() {
	caps := cap.GetProc()
	log.Printf("Current capabilities are: %s", caps.String())
}

// setCap sets the capabilities flags of the effective and permitted sets
// if the flag in the inheritable set is set.
func setCap() error {
	caps := cap.GetProc()

	for i := cap.CHOWN; i <= cap.CHECKPOINT_RESTORE; i++ {
		b, err := caps.GetFlag(cap.Inheritable, i)
		if err != nil {
			return err
		}
		if b {
			if err := caps.SetFlag(cap.Permitted, true, i); err != nil {
				return err
			}
			if err := caps.SetFlag(cap.Effective, true, i); err != nil {
				return err
			}
		}
	}
	return caps.SetProc()
}

func main() {
	printCap()
	if err := setCap(); err != nil {
		log.Fatal(err)
	}
	printCap()
}
