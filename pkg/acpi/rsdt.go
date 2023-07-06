package acpi

type RSDTEntry uint32

type RSDT struct {
	Header
	Entries []RSDTEntry
}
