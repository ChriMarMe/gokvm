package uefimachine

import (
	"github.com/bobuhiro11/gokvm/memory"
)

type SysAllocator struct {
	PIO          *memory.AddressAllocator
	PlatformMMIO *memory.AddressAllocator
	MMIOHole     *memory.AddressAllocator
	GSI          *GSIAllocator
}

func (u *UEFIMachine) InitSysAllocator() {
	pioalloc := memory.NewAddressAllocator(0, 1<<16)
	// Here is the place to allocate all required PIO Addresses and add them to the machine
	spaces := []*memory.AddressRange{
		{FirstAddr: 0x60, LastAddr: 0x70}, // PS/2 Keyboard (Always 8042 Chip)
		{FirstAddr: 0x70, LastAddr: 0x72}, // CMOS CLock
		{FirstAddr: 0x80, LastAddr: 0xA0}, // DMA Page register
		{FirstAddr: 0xED, LastAddr: 0xEE}, // Starndard Delay Port
	}

	for _, item := range spaces {
		pioalloc.Allocate(item.FirstAddr, item.LastAddr)
	}

	u.SysAllocator.PIO = pioalloc

	mmioalloc := memory.NewAddressAllocator(u.Config.MemCfg.StartPlatDevArea, memory.PlatformDeviceAreaSize)

	u.SysAllocator.PlatformMMIO = mmioalloc

	mmioHoleAlloc := memory.NewAddressAllocator(memory.Mem32BitDeviceStart, memory.Mem32BitDeviceSize)

	u.SysAllocator.MMIOHole = mmioHoleAlloc

}

func (u *UEFIMachine) AllocatePIO(port, size uint64) error {
	if err := u.SysAllocator.PIO.Allocate(port, size); err != nil {
		return err
	}

	return nil
}

func (u *UEFIMachine) AllocatePlaformMMIO(addr, size uint64) error {
	if err := u.SysAllocator.PlatformMMIO.Allocate(addr, size); err != nil {
		return err
	}

	return nil
}

func (u *UEFIMachine) InitMMIOHoleSpace(addr, size uint64) error {
	if err := u.SysAllocator.MMIOHole.Allocate(addr, size); err != nil {
		return err
	}

	return nil
}
