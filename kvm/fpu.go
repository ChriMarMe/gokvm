package kvm

type FpuState struct {
	FPR    [8][16]uint8
	FCW    uint16
	FSW    uint16
	FTWX   uint8
	LastOp uint16
	LastIp uint64
	LastDp uint64
	XMM    [16][16]uint8
	MXCSR  uint32
}
