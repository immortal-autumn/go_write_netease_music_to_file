/**
Assumption 1: Get the title of the windows directly.
This failed because handlers are independent and not connected directly
May enhanced by further exploration,
*/

package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const TH32CS_SNAPPROCESS = 0x00000002

type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

// This allows search processes by name
func SearchProcessByName(name string) *WindowsProcess {
	procList, err := Processes()
	if err != nil {
		fmt.Println("Error getting the process list")
		os.Exit(-1)
	}
	for _, p := range procList {
		if strings.Contains(strings.ToLower(p.Exe), strings.ToLower(name)) {
			return &p
		}
	}
	return nil
}

func listAllProcesses() {
	procList, err := Processes()
	if err != nil {
		fmt.Println("Error getting the process list")
		os.Exit(-1)
	}
	for _, p := range procList {
		fmt.Println(p.Exe, p.ProcessID)
	}
}

func Processes() ([]WindowsProcess, error) {
	handle, err := syscall.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(handle)

	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	//get first process
	err = syscall.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}
	// iteratively get next entries
	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, NewWindowsProcess(&entry))
		fmt.Println(entry)

		err = syscall.Process32Next(handle, &entry)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func NewWindowsProcess(e *syscall.ProcessEntry32) WindowsProcess {
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}
	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}
