package uefimachine

import (
	"github.com/bobuhiro11/gokvm/kvm"
	"github.com/bobuhiro11/gokvm/probe"
)

// CPUID feature bits
const (
	TSC_DEADLINE_TIMER_ECX_BIT uint32 = 24 // tsc deadline timer ecx bit.
	HYPERVISOR_ECX_BIT         uint32 = 31 // Hypervisor ecx bit.
	MTRR_EDX_BIT               uint32 = 12 // Hypervisor ecx bit.
	INVARIANT_TSC_EDX_BIT      uint32 = 8  // Invariant TSC bit on 0x8000_0007 EDX
)

const (
	// Default count of virtual cpus
	DefaultVCPUs uint32 = 1

	DefaultNUMPCISegments uint32 = 1

	DefaultRandomSrc = "/dev/urandom"
)

var (
	cpuidTSCPatch = kvm.CPUIDEntry2{
		Function: 1,
		Index:    0,
		Flags:    0,
		Eax:      0,
		Ebx:      0,
		Ecx:      TSC_DEADLINE_TIMER_ECX_BIT,
		Edx:      0,
	}

	cpuidHVPatch = kvm.CPUIDEntry2{
		Function: 1,
		Index:    0,
		Flags:    0,
		Eax:      0,
		Ebx:      0,
		Ecx:      HYPERVISOR_ECX_BIT,
		Edx:      0,
	}

	cpuidMTRRPatch = kvm.CPUIDEntry2{
		Function: 1,
		Index:    0,
		Flags:    0,
		Eax:      0,
		Ebx:      0,
		Ecx:      0,
		Edx:      MTRR_EDX_BIT,
	}

	patches = []kvm.CPUIDEntry2{
		cpuidTSCPatch,
		cpuidHVPatch,
		cpuidMTRRPatch,
	}
)

type MemCfg struct {
	RamSize           uint64
	MMIOAddrSpaceSize uint64
	StartPlatDevArea  uint64
	EndDeviceArea     uint64
}

const (
	// When booting with PVH boot the maximum physical addressable size
	// is a 46 bit address space even when the host supports with 5-level
	// paging.
	DefaultPhysBits uint32 = 46

	// Mem size constant for now
	DefaultMemSize = 1 << 30

	// Platform device area size is constant for now
	PlatformDeviceAreaSize = 1 << 20
)

func (m *MemCfg) Construct() {
	m.RamSize = DefaultMemSize

	m.MMIOAddrSpaceSize = uint64((1 << DefaultPhysBits) - (1 << 16))
	m.StartPlatDevArea = m.MMIOAddrSpaceSize - PlatformDeviceAreaSize

	m.EndDeviceArea = m.StartPlatDevArea - 1 // This doenst make any sense.....

}

type vmConfig struct {
	HasXSave2          bool
	XSaveSize          uint64
	MSRSupported       *kvm.MSRList
	SignalMsi          bool
	TscDeadlineTimer   bool
	SplitIRQChip       bool
	SetIdentityMapAddr bool
	SetTssAddr         bool
	ImmediatExit       bool
	GetTsckHz          bool
	SupportedCPUIDs    []kvm.CPUIDEntry2 // These CPUIDs will be patched.
	MemCfg
}

func (v *vmConfig) SetConfig(kvmfd uintptr) error {
	ret, err := kvm.CheckExtension(kvmfd, kvm.CapXSave2)
	if err != nil {
		return err
	}

	v.HasXSave2 = uint64(ret) != 0
	v.XSaveSize = uint64(ret)

	list := &kvm.MSRList{}

	if err := kvm.GetMSRIndexList(kvmfd, list); err != nil {
		return err
	}

	v.MSRSupported = list

	ret, err = kvm.CheckExtension(kvmfd, kvm.CapSignalMSI)
	if err != nil {
		return err
	}

	v.SignalMsi = ret != 0

	ret, err = kvm.CheckExtension(kvmfd, kvm.CapTSCDeadlineTimer)
	if err != nil {
		return err
	}

	v.TscDeadlineTimer = ret != 0

	ret, err = kvm.CheckExtension(kvmfd, kvm.CapSplitIRQChip)
	if err != nil {
		return err
	}

	v.SplitIRQChip = ret != 0

	ret, err = kvm.CheckExtension(kvmfd, kvm.CapSetIdentityMapAddr)
	if err != nil {
		return err
	}

	v.SetIdentityMapAddr = ret != 0

	ret, err = kvm.CheckExtension(kvmfd, kvm.CapSetTSSAddr)
	if err != nil {
		return err
	}

	v.SetTssAddr = ret != 0

	ret, err = kvm.CheckExtension(kvmfd, kvm.CapImmediateExit)
	if err != nil {
		return err
	}

	v.ImmediatExit = ret != 0

	ret, err = kvm.CheckExtension(kvmfd, kvm.CapGetTSCkHz)
	if err != nil {
		return err
	}

	v.GetTsckHz = ret != 0

	CPUID := kvm.CPUID{
		Nent:    100,
		Entries: make([]kvm.CPUIDEntry2, 100),
	}

	if err := kvm.GetSupportedCPUID(kvmfd, &CPUID); err != nil {
		return err
	}

	v.MemCfg.Construct()

	patchCPUID(&CPUID.Entries, patches)

	updateCPUIDwithHostValues(&CPUID.Entries)

	v.SupportedCPUIDs = CPUID.Entries

	return nil
}

// patchCPUID patches CPUIDs before vcpu generation
func patchCPUID(ids *[]kvm.CPUIDEntry2, patch []kvm.CPUIDEntry2) {
	for _, id := range *ids {
		for _, patch := range patches {
			if id.Function == patch.Function && id.Index == patch.Index {
				id.Flags |= 1 << patch.Flags
				id.Eax |= 1 << patch.Eax
				id.Ebx |= 1 << patch.Ebx
				id.Ecx |= 1 << patch.Ecx
				id.Edx |= 1 << patch.Edx
			}
		}
	}
}

func updateCPUIDwithHostValues(ids *[]kvm.CPUIDEntry2) {
	for _, id := range *ids {
		switch id.Function {
		case 0x8000_0006:
			// Copy host L2 cache details if not populated by KVM
			if id.Eax == 0 && id.Ebx == 0 && id.Ecx == 0 && id.Edx == 0 {
				if eax, _, _, _ := probe.CPUID(0x800_0000); eax >= 0x800_0006 {
					eax, ebx, ecx, edx := probe.CPUID(0x8000_0006)
					id.Eax = eax
					id.Ebx = ebx
					id.Ecx = ecx
					id.Edx = edx
				}
			}

		case 0x8000_0008:
			// Set CPU physical bits
			id.Eax = (id.Eax & 0xFFFF_FF00) | (uint32(DefaultPhysBits) & 0xFF)
		}
	}
}
