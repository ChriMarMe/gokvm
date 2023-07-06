package memory

import (
	"syscall"
)

const (
	MAXPhysBits = 46
)

// Memory represents the a vm memory in an abstract manner and contains memory regions.
type Memory struct {
	regions []MemoryRegion
	Base    uint64
	Size    uint64
}

type MemoryRegion interface {
	ReadAt()
	WriteAt()
}

// MemoryRegion represents a specific part of memory of RegionType also in an abstract manner.
type RAMRegion struct {
	Start uint64
	Size  uint64
	mem   []byte
}

func (r *RAMRegion) ReadAt() {}

func (r *RAMRegion) WriteAt() {}

type SubRegion struct {
	Start uint64
	Size  uint64
}

func (s *SubRegion) ReadAt() {}

func (s *SubRegion) WriteAt() {}

type ReservedRegion struct {
	Start uint64
	Size  uint64
}

func (r *ReservedRegion) ReadAt() {}

func (r *ReservedRegion) WriteAt() {}

// NewMemory creates a Memory struct and sets up initial configuration in a fixed manner.
func NewMemory() (*Memory, error) {
	// Hardcoded max size to 4GiB starting from address 0
	mem := &Memory{
		Base: 0,
		Size: 1 << 32,
	}

	// Create RAM, Device and Reserved MemoryRegion
	if err := mem.SetupRegions(); err != nil {
		return nil, err
	}

	mem.Base = 0
	mem.Size = 1 << 32

	return mem, nil
}

// SetupRegions create RAM, Device and Reserved MemoryRegion
// RAM and Device MemoryRegions will be covered by the same memory allocated by mmap.
func (m *Memory) SetupRegions() error {
	var err error
	ram := &RAMRegion{
		Start: 0x0,
		Size:  Mem32BitReservedStart,
	}

	if ram.mem, err = syscall.Mmap(-1, 0, 1<<30,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_ANONYMOUS); err != nil {
		return err
	}

	m.regions = append(m.regions, ram)

	dev := &SubRegion{
		Start: Mem32BitReservedStart,
		Size:  Mem32BitDeviceSize,
	}

	m.regions = append(m.regions, dev)

	res := &ReservedRegion{
		Start: Mem32BitReservedStart + Mem32BitDeviceSize,
		Size:  Mem32BitReservedSize - Mem32BitDeviceSize,
	}

	m.regions = append(m.regions, res)

	return nil
}
