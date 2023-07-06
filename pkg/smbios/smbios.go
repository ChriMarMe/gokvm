package smbios

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/bobuhiro11/gokvm/memory"
)

type entry64 struct {
	Anchor             [5]uint8
	Checksum           uint8
	Length             uint8
	SMBIOSMajorVersion uint8
	SMBIOSMinorVersion uint8
	SMBIOSDocRev       uint8
	Revision           uint8
	Reserved           uint8
	StructMaxSize      uint32
	StructTableAddr    uint64
}

func calcChecksum(data []byte, skipIndex int) uint8 {
	var cs uint8
	for i, b := range data {
		if i == skipIndex {
			continue
		}
		cs += b
	}
	return uint8(0x100 - int(cs))
}

func newEntry64(maxsize uint32) entry64 {
	return entry64{
		Anchor:             [5]byte{0x5F, 0x53, 0x4D, 0x33, 0x5F}, // _SM3_ -string for indication of entry64 type
		Checksum:           0,
		Length:             0x18, // static length
		SMBIOSMajorVersion: 0x3,  // SMBIOS rev. 3.2.0
		SMBIOSMinorVersion: 0x2,
		SMBIOSDocRev:       0x0,
		Revision:           0x1, // SMBIOS 3.0
		Reserved:           0,
		StructMaxSize:      maxsize,
		StructTableAddr:    memory.SMBIOSStart + 0x18, // SMBIOS_START + 0x18 (SizeOf(Entry64{}))
	}
}

func (e *entry64) write(buf io.Writer) error {
	if err := binary.Write(buf, binary.LittleEndian, e.Anchor); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.Checksum); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.Length); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.SMBIOSMajorVersion); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.SMBIOSMinorVersion); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.SMBIOSDocRev); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.Revision); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.Reserved); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.StructMaxSize); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, e.StructTableAddr); err != nil {
		return err
	}
	return nil
}

// BIOSCharacteristics is defined in DSP0134 7.1.1.
type BIOSCharacteristics uint64

// BIOSCharacteristics fields are defined in DSP0134 7.1.1.
const (
	BIOSCharacteristicsReserved                                                           BIOSCharacteristics = 1 << 0  // Reserved.
	BIOSCharacteristicsReserved2                                                          BIOSCharacteristics = 1 << 1  // Reserved.
	BIOSCharacteristicsUnknown                                                            BIOSCharacteristics = 1 << 2  // Unknown.
	BIOSCharacteristicsBIOSCharacteristicsAreNotSupported                                 BIOSCharacteristics = 1 << 3  // BIOS Characteristics are not supported.
	BIOSCharacteristicsISAIsSupported                                                     BIOSCharacteristics = 1 << 4  // ISA is supported.
	BIOSCharacteristicsMCAIsSupported                                                     BIOSCharacteristics = 1 << 5  // MCA is supported.
	BIOSCharacteristicsEISAIsSupported                                                    BIOSCharacteristics = 1 << 6  // EISA is supported.
	BIOSCharacteristicsPCIIsSupported                                                     BIOSCharacteristics = 1 << 7  // PCI is supported.
	BIOSCharacteristicsPCCardPCMCIAIsSupported                                            BIOSCharacteristics = 1 << 8  // PC card (PCMCIA) is supported.
	BIOSCharacteristicsPlugAndPlayIsSupported                                             BIOSCharacteristics = 1 << 9  // Plug and Play is supported.
	BIOSCharacteristicsAPMIsSupported                                                     BIOSCharacteristics = 1 << 10 // APM is supported.
	BIOSCharacteristicsBIOSIsUpgradeableFlash                                             BIOSCharacteristics = 1 << 11 // BIOS is upgradeable (Flash).
	BIOSCharacteristicsBIOSShadowingIsAllowed                                             BIOSCharacteristics = 1 << 12 // BIOS shadowing is allowed.
	BIOSCharacteristicsVLVESAIsSupported                                                  BIOSCharacteristics = 1 << 13 // VL-VESA is supported.
	BIOSCharacteristicsESCDSupportIsAvailable                                             BIOSCharacteristics = 1 << 14 // ESCD support is available.
	BIOSCharacteristicsBootFromCDIsSupported                                              BIOSCharacteristics = 1 << 15 // Boot from CD is supported.
	BIOSCharacteristicsSelectableBootIsSupported                                          BIOSCharacteristics = 1 << 16 // Selectable boot is supported.
	BIOSCharacteristicsBIOSROMIsSocketed                                                  BIOSCharacteristics = 1 << 17 // BIOS ROM is socketed.
	BIOSCharacteristicsBootFromPCCardPCMCIAIsSupported                                    BIOSCharacteristics = 1 << 18 // Boot from PC card (PCMCIA) is supported.
	BIOSCharacteristicsEDDSpecificationIsSupported                                        BIOSCharacteristics = 1 << 19 // EDD specification is supported.
	BIOSCharacteristicsInt13hJapaneseFloppyForNEC980012MB351KBytessector360RPMIsSupported BIOSCharacteristics = 1 << 20 // Japanese floppy for NEC 9800 1.2 MB (3.5”, 1K bytes/sector, 360 RPM) is
	BIOSCharacteristicsInt13hJapaneseFloppyForToshiba12MB35360RPMIsSupported              BIOSCharacteristics = 1 << 21 // Japanese floppy for Toshiba 1.2 MB (3.5”, 360 RPM) is supported.
	BIOSCharacteristicsInt13h525360KBFloppyServicesAreSupported                           BIOSCharacteristics = 1 << 22 // 5.25” / 360 KB floppy services are supported.
	BIOSCharacteristicsInt13h52512MBFloppyServicesAreSupported                            BIOSCharacteristics = 1 << 23 // 5.25” /1.2 MB floppy services are supported.
	BIOSCharacteristicsInt13h35720KBFloppyServicesAreSupported                            BIOSCharacteristics = 1 << 24 // 3.5” / 720 KB floppy services are supported.
	BIOSCharacteristicsInt13h35288MBFloppyServicesAreSupported                            BIOSCharacteristics = 1 << 25 // 3.5” / 2.88 MB floppy services are supported.
	BIOSCharacteristicsInt5hPrintScreenServiceIsSupported                                 BIOSCharacteristics = 1 << 26 // Int 5h, print screen Service is supported.
	BIOSCharacteristicsInt9h8042KeyboardServicesAreSupported                              BIOSCharacteristics = 1 << 27 // Int 9h, 8042 keyboard services are supported.
	BIOSCharacteristicsInt14hSerialServicesAreSupported                                   BIOSCharacteristics = 1 << 28 // Int 14h, serial services are supported.
	BIOSCharacteristicsInt17hPrinterServicesAreSupported                                  BIOSCharacteristics = 1 << 29 // Int 17h, printer services are supported.
	BIOSCharacteristicsInt10hCGAMonoVideoServicesAreSupported                             BIOSCharacteristics = 1 << 30 // Int 10h, CGA/Mono Video Services are supported.
	BIOSCharacteristicsNECPC98                                                            BIOSCharacteristics = 1 << 31 // NEC PC-98.
)

// BIOSCharacteristicsExt1 is defined in DSP0134 7.1.2.1.
type BIOSCharacteristicsExt1 uint8

// BIOSCharacteristicsExt1 is defined in DSP0134 7.1.2.1.
const (
	BIOSCharacteristicsExt1ACPIIsSupported               BIOSCharacteristicsExt1 = 1 << 0 // ACPI is supported.
	BIOSCharacteristicsExt1USBLegacyIsSupported          BIOSCharacteristicsExt1 = 1 << 1 // USB Legacy is supported.
	BIOSCharacteristicsExt1AGPIsSupported                BIOSCharacteristicsExt1 = 1 << 2 // AGP is supported.
	BIOSCharacteristicsExt1I2OBootIsSupported            BIOSCharacteristicsExt1 = 1 << 3 // I2O boot is supported.
	BIOSCharacteristicsExt1LS120SuperDiskBootIsSupported BIOSCharacteristicsExt1 = 1 << 4 // LS-120 SuperDisk boot is supported.
	BIOSCharacteristicsExt1ATAPIZIPDriveBootIsSupported  BIOSCharacteristicsExt1 = 1 << 5 // ATAPI ZIP drive boot is supported.
	BIOSCharacteristicsExt11394BootIsSupported           BIOSCharacteristicsExt1 = 1 << 6 // 1394 boot is supported.
	BIOSCharacteristicsExt1SmartBatteryIsSupported       BIOSCharacteristicsExt1 = 1 << 7 // Smart battery is supported.
)

// BIOSCharacteristicsExt2 is defined in DSP0134 7.1.2.2.
type BIOSCharacteristicsExt2 uint8

// BIOSCharacteristicsExt1 is defined in DSP0134 7.1.2.2.
const (
	BIOSCharacteristicsExt2BIOSBootSpecificationIsSupported                  BIOSCharacteristicsExt2 = 1 << 0 // BIOS Boot Specification is supported.
	BIOSCharacteristicsExt2FunctionKeyinitiatedNetworkServiceBootIsSupported BIOSCharacteristicsExt2 = 1 << 1 // Function key-initiated network service boot is supported.
	BIOSCharacteristicsExt2TargetedContentDistributionIsSupported            BIOSCharacteristicsExt2 = 1 << 2 // Enable targeted content distribution.
	BIOSCharacteristicsExt2UEFISpecificationIsSupported                      BIOSCharacteristicsExt2 = 1 << 3 // UEFI Specification is supported.
	BIOSCharacteristicsExt2SMBIOSTableDescribesAVirtualMachine               BIOSCharacteristicsExt2 = 1 << 4 // SMBIOS table describes a virtual machine. (If this bit is not set, no inference can be made
)

// TableType specifies the DMI type of the table.
// Types are defined in DMTF DSP0134.
type TableType uint8

// Supported table types.
const (
	TableTypeBIOSInfo       TableType = 0
	TableTypeSystemInfo     TableType = 1
	TableTypeBaseboardInfo  TableType = 2
	TableTypeChassisInfo    TableType = 3
	TableTypeProcessorInfo  TableType = 4
	TableTypeCacheInfo      TableType = 7
	TableTypeSystemSlots    TableType = 9
	TableTypeOEMStrings     TableType = 11
	TableTypeMemoryDevice   TableType = 17
	TableTypeIPMIDeviceInfo TableType = 38
	TableTypeTPMDevice      TableType = 43
	TableTypeInactive       TableType = 126
	TableTypeEndOfTable     TableType = 127
)

type header struct {
	Type   TableType
	Length uint8
	Handle uint16
}

type biosInfo struct {
	header
	Vendor                                 byte                    // 04h
	Version                                byte                    // 05h
	StartingAddressSegment                 uint16                  // 06h
	ReleaseDate                            byte                    // 08h
	ROMSize                                uint8                   // 09h
	Characteristics                        BIOSCharacteristics     // 0Ah
	CharacteristicsExt1                    BIOSCharacteristicsExt1 // 12h
	CharacteristicsExt2                    BIOSCharacteristicsExt2 // 13h
	SystemBIOSMajorRelease                 uint8                   // 14h
	SystemBIOSMinorRelease                 uint8                   // 15h
	EmbeddedControllerFirmwareMajorRelease uint8                   // 16h
	EmbeddedControllerFirmwareMinorRelease uint8                   // 17h
	ExtendedROMSize                        uint16                  // 18h
}

func newBiosInfo() biosInfo {

	nheader := header{
		Type:   TableTypeBIOSInfo,
		Length: uint8(0x18),
		Handle: 1,
	}

	return biosInfo{
		header:              nheader,
		Vendor:              1,
		Version:             2,
		Characteristics:     BIOSCharacteristicsPCIIsSupported,
		CharacteristicsExt2: BIOSCharacteristicsExt2SMBIOSTableDescribesAVirtualMachine,
	}
}

func (bi *biosInfo) write(buf io.Writer, str ...string) error {
	if err := binary.Write(buf, binary.LittleEndian, bi.Type); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.Length); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.Handle); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.Vendor); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.Version); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.StartingAddressSegment); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.ReleaseDate); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.ROMSize); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.Characteristics); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.CharacteristicsExt1); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.CharacteristicsExt2); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.SystemBIOSMajorRelease); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.SystemBIOSMinorRelease); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.EmbeddedControllerFirmwareMajorRelease); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.EmbeddedControllerFirmwareMinorRelease); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, bi.ExtendedROMSize); err != nil {
		return err
	}

	if len(str) != 0 {
		for _, str := range str {
			if err := binary.Write(buf, binary.LittleEndian, []byte(str)); err != nil {
				return err
			}
		}
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(0x0)); err != nil {
		return err
	}

	return nil
}

type systemInfo struct {
	header
	Manufacturer uint8
	ProdName     uint8
	Version      uint8
	SerialNr     uint8
	UUID         [16]uint8
	WakeUp       uint8
	SKU          uint8
	Family       uint8
}

func newSystemInfo() systemInfo {
	nheader := header{
		Type:   TableTypeSystemInfo,
		Length: 27,
		Handle: 2,
	}

	return systemInfo{
		header:       nheader,
		Manufacturer: 1,
		ProdName:     1,
		SerialNr:     0,
		UUID:         [16]byte{},
	}

}

func (si *systemInfo) write(buf io.Writer, str ...string) error {
	if err := binary.Write(buf, binary.LittleEndian, si.Type); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.Length); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.Handle); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.Manufacturer); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.ProdName); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.Version); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.SerialNr); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.UUID); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.WakeUp); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.SKU); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, si.Family); err != nil {
		return err
	}

	if len(str) != 0 {
		for _, str := range str {
			if err := binary.Write(buf, binary.LittleEndian, []byte(str)); err != nil {
				return err
			}
		}
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(0x0)); err != nil {
		return err
	}

	return nil
}

type endOfTable struct {
	header
}

func newEndOfTable() endOfTable {
	nheader := header{
		Type:   TableTypeEndOfTable,
		Length: 4,
		Handle: 4,
	}
	return endOfTable{header: nheader}
}

func (et *endOfTable) write(buf io.Writer) error {
	if err := binary.Write(buf, binary.LittleEndian, et.Type); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, et.Length); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, et.Handle); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, byte(0x0)); err != nil {
		return err
	}

	return nil
}

func CreateSMBIOSData() ([]byte, error) {
	var bibuf, sibuf, eotbuf, entrybuf bytes.Buffer
	ret := make([]byte, 0)
	bi := newBiosInfo()
	si := newSystemInfo()
	eot := newEndOfTable()

	if err := bi.write(&bibuf, "cloud-hypervisor", "0"); err != nil {
		return nil, err
	}
	if err := si.write(&sibuf, "Cloud Hypervisor", "cloud-hypervisor"); err != nil {
		return nil, err
	}
	if err := eot.write(&eotbuf); err != nil {
		return nil, err
	}

	size := uint32(bibuf.Len() + sibuf.Len() + eotbuf.Len())

	entry := newEntry64(size)

	if err := entry.write(&entrybuf); err != nil {
		return nil, err
	}
	// calc checksum
	chksm := calcChecksum(entrybuf.Bytes(), 5)
	// Put the chksm into the entry bytes at position 5

	ret = append(ret, entrybuf.Bytes()...)
	copy(ret[5:], []byte{chksm})

	// now we append the rest
	ret = append(ret, bibuf.Bytes()...)
	ret = append(ret, sibuf.Bytes()...)
	ret = append(ret, eotbuf.Bytes()...)

	return ret, nil
}
