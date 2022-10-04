package main

// This small example binary was inspired by the following blog post:
// https://blog.container-solutions.com/linux-capabilities-in-practice
// It is an implementation of the suggested "ambient" binary to set a child's ambient set.
// In order to use this binary, do the following:
// 1. Build it:
//      make ambient
// 2. Set the capabilities that you want to show up in the ambient set. All true bits from the permitted set will be
//    copied to the ambient set, so for example to have the binary set the cap_net_bind_service and cap_chown ambient
//    set flags to true for the child process, setcap the following:
//      sudo setcap "cap_net_bind_service+p cap_chown+p" bin/ambient
// 3. Call ambient with bash, then print the current capabilities:
//      bin/ambient /bin/bash
//      capsh --print | grep "Ambient set

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"syscall"

	"kernel.org/pub/linux/libs/security/libcap/cap"
)

// setInheritableFlags iterates over this thread's permitted capability set.
// If a given capability is enabled in the permitted capability set, make sure to enable the same
// capability in the inheritable set as well. This is a prerequisite for setting the same flag in the ambient
// set later on.
func setInheritableFlags() error {
	// Get this process's capability set.
	caps := cap.GetProc()

	// Iterate over all capabilities bits, set the inheritable flag to true
	// if the capability's permitted flag is true.
	for i := cap.CHOWN; i <= cap.CHECKPOINT_RESTORE; i++ {
		b, err := caps.GetFlag(cap.Permitted, i)
		if err != nil {
			return err
		}
		if b {
			if err := caps.SetFlag(cap.Inheritable, true, i); err != nil {
				return err
			}
		}
	}
	return caps.SetProc()
}

// setAmbientFlags iterates over this thread's permitted capability set.
// If a given capability is enabled in the permitted capability set, set the same flag in the ambient set.
// As a prerequisite, run setInheritableFlags first.
// man prctl
// (...)
//  PR_CAP_AMBIENT_RAISE
//  The capability specified in arg3 is added to the ambient set.  The specified capability  must
//  already be present in both the permitted and the inheritable sets of the process.
// (...)
func setAmbientFlags() error {
	// Get this process's capability set.
	caps := cap.GetProc()

	// Iterate over all capabilities bits, set the flag to true in the ambient set.
	for i := cap.CHOWN; i <= cap.CHECKPOINT_RESTORE; i++ {
		b, err := caps.GetFlag(cap.Permitted, i)
		if err != nil {
			return err
		}
		if b {
			if err := cap.SetAmbient(true, i); err != nil {
				return err
			}
		}
	}

	return nil
}

// printCap prints the current capabilities in a human friendly way.
func printCap() {
	caps := cap.GetProc()
	log.Printf("Current capabilities are: %s", caps.String())
}

// execCommand runs the provided command.
func execCommand() {
	flag.Parse()
	if len(os.Args) == 1 {
		return
	}
	c, err := exec.LookPath(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := syscall.Exec(c, flag.Args(), os.Environ()); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := setInheritableFlags(); err != nil {
		log.Fatal(err)
	}
	if err := setAmbientFlags(); err != nil {
		log.Fatal(err)
	}
	execCommand()
}
