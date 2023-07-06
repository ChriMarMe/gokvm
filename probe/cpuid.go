package probe

func cpuid_low(arg1, arg2 uint32) (eax, ebx, ecx, edx uint32) // implemented in cpuid.s

func CPUID(leaf uint32) (uint32, uint32, uint32, uint32) {
	return cpuid_low(leaf, 0)
}
