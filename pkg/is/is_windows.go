// originally from github.com/inconshreveable/mousetrap (Apache license)

package is

import (
	"syscall"
	"unsafe"
)

func getProcessEntry(pid int) (*syscall.ProcessEntry32, error) {
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(snapshot)
	var procEntry syscall.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = syscall.Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}
	for {
		if procEntry.ProcessID == uint32(pid) {
			return &procEntry, nil
		}
		err = syscall.Process32Next(snapshot, &procEntry)
		if err != nil {
			return nil, err
		}
	}
}

func StartedByExplorer() bool {
	pe, err := getProcessEntry(syscall.Getppid())
	if err != nil {
		return false
	}
	return "explorer.exe" == syscall.UTF16ToString(pe.ExeFile[:])
}
