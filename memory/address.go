package memory

import (
	"errors"
)

var (
	errNoSize = errors.New("no size given")
)

// In PIO we dont have overlapping spaces.
// What about MMIO?
type AddressRange struct {
	FirstAddr uint64
	LastAddr  uint64
}

// ToDo: Change to a more sophisticated data structure for AddressRange for search efficiency
type AddressAllocator struct {
	Base   uint64
	End    uint64
	Ranges []*AddressRange
}

func NewAddressAllocator(base, size uint64) *AddressAllocator {
	if size <= 0 {
		return nil
	}

	ret := &AddressAllocator{
		Base: base,
		End:  base + (size - 1),
	}

	// We add the firstEntry for the size comparison
	firstEntry := &AddressRange{
		FirstAddr: base,
		LastAddr:  base + size,
	}

	ret.Ranges = append(ret.Ranges, firstEntry)

	return ret
}

func (a *AddressAllocator) Allocate(base, size uint64) error {
	if size == 0 {
		return errNoSize
	}

	add := &AddressRange{
		FirstAddr: base,
		LastAddr:  base + (size - 1),
	}

	// Simply add the entry. Optimisation is for professionals and we dont do this now.
	a.Ranges = append(a.Ranges, add)

	return nil
}

func (a *AddressAllocator) Find(firstAddr, lastAddr uint64) bool {
	for _, item := range a.Ranges {
		if item.FirstAddr == firstAddr && item.LastAddr == lastAddr {
			return true
		}
	}

	return false
}
