package uefimachine

const (
	NumIOAPICPins uint32 = 24
)

type APIC struct {
	Base uint32
	irqs uint32
}

type GSIAllocator struct {
	APICS   *[]APIC
	nextIRQ uint32
	nextGSI uint32
}

func (u *UEFIMachine) NewGSIAllocator(apics []APIC) {
	a := make([]APIC, 0)
	gsi := &GSIAllocator{
		nextIRQ: 0xFFFF_FFFF,
		nextGSI: 0,
		APICS:   &a,
	}

	for _, apic := range apics {
		if apic.Base < u.GSI.nextGSI {
			gsi.nextIRQ = apic.Base
		}

		if apic.Base+apic.irqs > u.GSI.nextGSI {
			u.GSI.nextGSI = apic.Base + apic.irqs
		}

		a = append(a, apic)
	}

	u.GSI = gsi
}
