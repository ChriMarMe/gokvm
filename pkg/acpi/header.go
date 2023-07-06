package acpi

type Header struct {
	Signature  [4]byte
	Length     uint32
	Rev        uint8
	Checksum   uint8
	OEMId      [6]byte
	OEMTableID [8]byte
	OEMRev     uint32
	CreatorId  [4]byte
	CreatorRev [4]byte
}
