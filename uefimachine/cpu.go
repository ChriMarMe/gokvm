package uefimachine

import "github.com/bobuhiro11/gokvm/kvm"

// MTRR constants
const MTRR_ENABLE uint64 = 0x800 // IA32_MTRR_DEF_TYPE MSR: E (MTRRs enabled) flag, bit 11
const MTRR_MEM_TYPE_WB uint64 = 0x6

var (
	BootMSRfix = []kvm.MSR{
		kvm.MSRIA32SYSENTERCS,
		kvm.MSRIA32SYSENTERESP,
		kvm.MSRIA32SYSENTEREIP,
		kvm.MSRSTAR,
		kvm.MSRCSTAR,
		kvm.MSRLSTAR,
		kvm.MSRKERNELGSBASE,
		kvm.MSRFMASK,
		kvm.MSRIA32TSC,
	}

	firstmsr  = kvm.MSREntry{Index: uint32(kvm.MSRIA32MISCENABLE), Data: uint64(kvm.MSRIA32MISCENABLEFASTSTRING)}
	secondmsr = kvm.MSREntry{Index: uint32(kvm.MSRMTRRdefType), Data: uint64(MTRR_ENABLE | MTRR_MEM_TYPE_WB)}

	BootMSRVals = kvm.MSRS{
		NMSRs: 1,
		Entries: []kvm.MSREntry{
			firstmsr,
			secondmsr,
		},
	}
)

// CPU represents a virtual CPU
// Required if state should be restorable
type VCPU struct {
	FD         uintptr
	CPUID      []kvm.CPUIDEntry2
	MSRS       []kvm.MSRS
	Events     []kvm.VCPUEvents
	RunData    *kvm.RunData
	Regs       kvm.Regs
	SRegs      kvm.Sregs
	Fpu        kvm.FpuState
	Lapicstate kvm.LAPICState
	Xsave      kvm.XSave
	XCRS       kvm.XCRS
	MpState    kvm.MPState
}

type CPUManager struct {
	Vcpus []VCPU
}
