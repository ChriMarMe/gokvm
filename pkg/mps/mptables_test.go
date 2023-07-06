package mps

import "testing"

func TestCalcChecksum(t *testing.T) {
	for _, tt := range []struct {
		name   string
		input  []byte
		outexp uint8
	}{
		{
			name:   "SimpleUnderflow",
			input:  []byte{255},
			outexp: 255,
		},
		{
			name:   "SimpleOverflow",
			input:  []byte{255, 1},
			outexp: 0,
		},
		{
			name:   "SimpleOverflow",
			input:  []byte{255, 2},
			outexp: 1,
		},
		{
			name:   "SimpleOverflow",
			input:  []byte{255, 255, 255},
			outexp: 253,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			chksm := calcChecksum(tt.input)
			if chksm != tt.outexp {
				t.Fatalf("chksm: %d, exp: %d", chksm, tt.outexp)
			}
		})
	}
}

func TestMPFPointerStructure(t *testing.T) {
	mpfstr, err := newMPFPointerStruct(0xF_0000)
	if err != nil {
		t.Fatal(err)
	}

	data, err := mpfstr.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	chksm2 := ^(calcChecksum(data) - mpfstr.Chksm) + 1

	if chksm2 != mpfstr.Chksm {
		t.Fatalf("checksum algorithm does not work as intended: Have: %d", chksm2)
	}
}

func TestMPTablesCreate(t *testing.T) {
	tables, err := CreateMPTables(0xF_0000)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tables)
	t.Log(len(tables))
}
