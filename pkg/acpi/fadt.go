package acpi

type FADTFeatureFlag uint32

const (
	WBINVD FADTFeatureFlag = 1<<0 + iota
	WBINVDFlush
	ProcC1
	PLvL2Up
	PwrButton
	SleepButton
	FixRTC
	RTCS4
	TmrValExt
	DCKCap
	ResetRegSup
	SealedCase
	Headless
	CPUSwSleep
	PCIExpWak
	UsePlatformClock
	S4RTCSTSValid
	RemotePowerOnCapable
	ForceAPICCluterModel
	ForceAPICPhysicalDestMode
	HwReducedACPI
	LowPowerS0IdleCapable
)

type FADT struct {
	Header
	FirmwareCTRL  uint32
	DSDTAddr      uint32
	_             uint8
	PrefPMProfile uint8
	SCIInt        uint16
	SMICmd        uint32
	ACPIEnable    uint8
	ACPIDisable   uint8
	S4BIOSReq     uint8
	PStateCnt     uint8
	PM1aEvtBlk    uint32
	PM1bEvtBlk    uint32
	PM1aCntBlk    uint32
	PM1bCntBlk    uint32
	PM2CntBlk     uint32
	PMTmrBlk      uint32
	GPE0Blk       uint32
	GPE1Blk       uint32
	PM1EvtLen     uint8
	PM1CntLen     uint8
	PM2CntLen     uint8
	PMTmrLen      uint8
	GPE0BlkLen    uint8
	GPE1BlkLen    uint8
	GPE1Base      uint8
	CstCnt        uint8
	PLvL2Lat      uint16
	PLvL3Lat      uint16
	FlushSize     uint16
	FlushStride   uint16
	DutyOffset    uint8
	DutyWidth     uint8
	DayALRM       uint8
	MonALRM       uint8
	Century       uint8
	IAPCBootArch  uint16
	_             uint8
	FADTFeatureFlag
	ResetReg      [12]uint8
	ResetValue    uint8
	ARMBootArch   uint16
	MinorVersion  uint8
	XFirmwareCntl uint64
	XDSDT         uint64
	XPM1aEvtBlk   [12]uint8
	XPM1bEvtBlk   [12]uint8
	XPM1aCntBlk   [12]uint8
	XPM1bCntBlk   [12]uint8
	XPM2CntBlk    [12]uint8
	XPMTmrBlk     [12]uint8
	XGPE0Blk      [12]uint8
	XGPE1Blk      [12]uint8
	SleepCtlReg   [12]uint8
	SleepStatReg  [12]uint8
	HyperVendorId [8]uint8
}
