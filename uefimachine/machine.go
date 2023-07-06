package uefimachine

import (
	"os"

	"github.com/bobuhiro11/gokvm/kvm"
)

type UEFIMachine struct {
	kvmFd, vmFd uintptr
	CPUManager
	SysAllocator
	Config *vmConfig
}

func New() (*UEFIMachine, error) {
	kvmObj, err := os.OpenFile("/dev/kvm", os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	vmFd, err := kvm.CreateVM(kvmObj.Fd())
	if err != nil {
		return nil, err
	}

	return &UEFIMachine{
		kvmFd: kvmObj.Fd(),
		vmFd:  vmFd,
	}, nil
}
