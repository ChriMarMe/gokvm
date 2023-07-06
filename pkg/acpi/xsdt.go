package acpi

type XSDTEntry uint64

type XSDT struct {
	Header
	Entries []XSDTEntry
}
