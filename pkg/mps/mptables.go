package mps

import (
	"bytes"
	"encoding/binary"

	"github.com/bobuhiro11/gokvm/memory"
)

// We rely on the correct handling of overflow by the language itself.
func calcChecksum(data []byte) uint8 {
	var ret uint8
	for _, item := range data {
		ret += item
	}
	return ret
}

const (
	MPProcessor uint8 = 0 + iota
	MPBUS
	MPIOAPIC
	MPINTSRC
	MPLINTSRC
)

const (
	CPUEnable       = 1
	CPUBP           = 2
	MPCAPICUsable   = 1
	MPIRQDIRDefault = 0

	MPIRQSRCTypeMPInt    = 0
	MPIRQSRCTypeMPNMI    = 1
	MPIRQSRCTypeMPExtInt = 3

	MPCSpec        = uint8(4)
	APICVersion    = uint8(0x14)
	CPUStepping    = uint32(0x600)
	CPUFeatureAPIC = uint32(0x200)
	CPUFeatureFPU  = uint32(0x001)
)

var (
	SMPMagic   = [4]byte{0x5F, 0x4D, 0x50, 0x5F}                         // "_MP_"
	MPCSig     = [4]byte{0x50, 0x43, 0x4D, 0x50}                         // "PCMP"
	MPCOEM     = [8]byte{0x47, 0x4F, 0x4B, 0x56, 0x4D, 0x00, 0x00, 0x00} // "GOKVM   "
	BUSTypeISA = [6]byte{0x49, 0x53, 0x41, 0x00, 0x00, 0x00}             // "ISA   "
)

type MPFPointerStruct struct {
	Signature  [4]byte
	PhysAddr   uint32
	Length     uint8
	SpecRev    uint8
	Chksm      uint8
	MPFeature1 uint8
	MPFeature2 uint8
	_          [3]byte
}

func (mf *MPFPointerStruct) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, mf.Signature); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, mf.PhysAddr); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, mf.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, mf.SpecRev); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, mf.Chksm); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, mf.MPFeature1); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, mf.MPFeature2); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, [3]byte{0x0, 0x0, 0x0}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func newMPFPointerStruct(physAddr uint32) (MPFPointerStruct, error) {
	ret := MPFPointerStruct{
		Signature:  SMPMagic,
		PhysAddr:   physAddr,
		Length:     0x1,
		SpecRev:    4,
		Chksm:      0,
		MPFeature1: 0,
		MPFeature2: 0,
	}

	data, err := ret.Bytes()
	if err != nil {
		return MPFPointerStruct{}, err
	}

	chksm := calcChecksum(data)

	ret.Chksm = ^chksm + 1

	return ret, nil
}

const MPFPointerStructSize = 16

type mpcTable struct {
	Signature       [4]byte
	BaseTableLength uint16
	SpecRev         uint8
	Checksum        uint8
	OEMId           [8]byte
	ProdID          [12]byte
	TablePointer    uint32
	TableSize       uint16
	EntryCount      uint16
	LAPICAddr       uint32
	ExtTableLen     uint16
	ExtTableChksm   uint8
}

const mpcTableSize = 44

func newmpcTable(len uint16) mpcTable {
	return mpcTable{
		Signature:       MPCSig,
		BaseTableLength: len,
		SpecRev:         MPCSpec,
		Checksum:        0,
		OEMId:           MPCOEM,
		ProdID:          [12]byte{},
		LAPICAddr:       memory.APICStart,
	}

}

func (m *mpcTable) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, m.Signature); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.BaseTableLength); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.SpecRev); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.Checksum); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.OEMId); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.ProdID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.LAPICAddr); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type mpcCPU struct {
	EntryType uint8
	LAPICId   uint8
	LAPICVer  uint8
	CPUFlags  uint8
	CPUSig1   uint16
	CPUSig2   uint16
	FeatFlags uint32
	_         uint64
}

const mpcCPUSize = 20

func newmpcCPU(cpuid uint8) mpcCPU {
	f := uint8(CPUEnable)
	if cpuid == 0 {
		f |= CPUBP
	}

	return mpcCPU{
		EntryType: MPProcessor,
		LAPICId:   cpuid,
		LAPICVer:  APICVersion,
		CPUFlags:  f,
		CPUSig1:   uint16(CPUStepping),
		CPUSig2:   0,
		FeatFlags: CPUFeatureAPIC | CPUFeatureFPU,
	}
}

func (m *mpcCPU) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, m.EntryType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.LAPICId); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.LAPICVer); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.CPUFlags); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.CPUSig1); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.CPUSig2); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.FeatFlags); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type mpcBUS struct {
	EntryType  uint8
	BusId      uint8
	BusTypeStr [6]byte
}

const mpcBUSSize = 8

func newmpcBus() mpcBUS {
	return mpcBUS{
		EntryType:  MPBUS,
		BusId:      0,
		BusTypeStr: BUSTypeISA,
	}
}

func (m *mpcBUS) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, m.EntryType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.BusId); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.BusTypeStr); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type mpcIOAPCI struct {
	EntryType  uint8
	ID         uint8
	Version    uint8
	Flags      uint8
	IOAPICAddr uint32
}

const mpcIOAPCISize = 8

func newmpcIOAPIC(id uint8) mpcIOAPCI {
	return mpcIOAPCI{
		EntryType:  MPIOAPIC,
		ID:         id,
		Version:    APICVersion,
		Flags:      MPCAPICUsable,
		IOAPICAddr: memory.IOAPICStart,
	}
}

func (m *mpcIOAPCI) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, m.EntryType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.ID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.Flags); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.IOAPICAddr); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Combines Type 3 and Type 4 EntryTypes.
type mpcInterrupt struct {
	EntryType     uint8
	InterruptType uint8
	Flags         uint16
	SrcBusId      uint8
	SrcBusIRQ     uint8
	DestId        uint8
	DestPin       uint8
}

const mpcInterruptSize = 8

func newmpcInterrupt(ioapicid uint8) [16]mpcInterrupt {
	var ret [16]mpcInterrupt

	for i := 0; i < 16; i++ {
		ret[i].EntryType = MPINTSRC
		ret[i].InterruptType = MPIRQSRCTypeMPInt
		ret[i].Flags = MPIRQDIRDefault
		ret[i].SrcBusId = 0
		ret[i].SrcBusIRQ = uint8(i)
		ret[i].DestId = ioapicid
		ret[i].DestPin = uint8(i)
	}

	return ret
}

func (m *mpcInterrupt) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, m.EntryType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.InterruptType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.Flags); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.SrcBusId); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.SrcBusIRQ); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.DestId); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, m.DestPin); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func newmpcLocalInterrupt() [2]mpcInterrupt {
	var ret [2]mpcInterrupt

	for i := 0; i > 2; i++ {
		ret[i].EntryType = MPLINTSRC
		ret[i].InterruptType = MPIRQSRCTypeMPExtInt
		ret[i].Flags = MPIRQDIRDefault
		ret[i].SrcBusId = 0
		ret[i].SrcBusIRQ = 0
		ret[i].DestId = 0
		ret[i].DestPin = 0
	}

	ret[1].DestId = 0xFF // It is, what it is.

	return ret

}

func sizeOfMPTable(numCPU uint16) uint16 {
	ret := uint16(0)
	ret += MPFPointerStructSize
	ret += mpcTableSize
	ret += mpcCPUSize * numCPU // we just try 1 CPU for now.
	ret += mpcIOAPCISize
	ret += mpcBUSSize
	ret += mpcInterruptSize * 16
	ret += mpcInterruptSize * 2
	return ret
}

func CreateMPTables(physAddr uint32) ([]byte, error) {
	mpcTable, err := newMPFPointerStruct(physAddr)
	if err != nil {
		return nil, err
	}

	mpctablebytes, err := mpcTable.Bytes()
	if err != nil {
		return nil, err
	}

	datalen := sizeOfMPTable(1)

	table := newmpcTable(datalen)

	tablebytes, err := table.Bytes()
	if err != nil {
		return nil, err
	}

	cpu := newmpcCPU(1) // fixed one cpu for now
	cpubytes, err := cpu.Bytes()
	if err != nil {
		return nil, err
	}

	bus := newmpcBus()
	busbytes, err := bus.Bytes()
	if err != nil {
		return nil, err
	}

	ioapic := newmpcIOAPIC(2)
	ioapicbytes, err := ioapic.Bytes()
	if err != nil {
		return nil, err
	}

	intr := newmpcInterrupt(2) // NUM_CPU + 1
	intrbytes := make([]byte, 0)
	for _, entry := range intr {
		entrybytes, err := entry.Bytes()
		if err != nil {
			return nil, err
		}
		intrbytes = append(intrbytes, entrybytes...)
	}

	lint := newmpcLocalInterrupt()
	lintbytes := make([]byte, 0)
	for _, entry := range lint {
		entrybytes, err := entry.Bytes()
		if err != nil {
			return nil, err
		}
		lintbytes = append(lintbytes, entrybytes...)
	}

	data := make([]byte, 0)

	chskmsData := [][]byte{
		cpubytes,
		busbytes,
		ioapicbytes,
		intrbytes,
		lintbytes,
		tablebytes,
	}

	chksm := uint8(0)
	for _, item := range chskmsData {
		chksm += calcChecksum(item)
	}

	chksm = ^chksm + 1

	tablebytes[7] = byte(chksm)

	data = append(data, tablebytes...)
	data = append(data, cpubytes...)
	data = append(data, busbytes...)
	data = append(data, ioapicbytes...)
	data = append(data, intrbytes...)
	data = append(data, lintbytes...)

	mpctablebytes = append(mpctablebytes, data...)

	return mpctablebytes, nil
}
