package acpi

type AMLOp uint8

const (
	OpZero = 0x00
	OpOne  = 0x01

	OpName            = 0x08
	OpBytePrefix      = 0x0A
	OpWordPrefix      = 0x0B
	OpFWordPrefix     = 0x0C
	OpString          = 0x0D
	QWordPrefix       = 0x0E
	OpScope           = 0x10
	OpBuffer          = 0x11
	OpPackage         = 0x12
	OpVarPackage      = 0x13
	OpMethod          = 0x14
	OpDualNamePrefix  = 0x2E
	OpMultiNamePrefix = 0x2F

	OpNameCharBase = 0x40

	OpExtPrefix   = 0x5b
	OpMutex       = 0x01
	OpCreateFile  = 0x13
	OpAcquire     = 0x23
	OpRelease     = 0x27
	OpRegionOp    = 0x80
	OpFile        = 0x81
	OpDevice      = 0x82
	OpPowerSource = 0x84

	OpLocal = 0x60
	OpArg   = 0x68
	OpStore = 0x70
	OpAdd   = 0x72
)
