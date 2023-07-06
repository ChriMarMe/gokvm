package smbios

import (
	"bytes"
	"testing"
)

func TestBIOSTableWrite(t *testing.T) {
	biosTable := newBiosInfo()

	var buf bytes.Buffer

	if err := biosTable.write(&buf, "hallo", "world"); err != nil {
		t.Fatal(err)
	}

	t.Log(buf.Bytes())

}

func TestSystemInfoTableWrite(t *testing.T) {
	sysInfo := newSystemInfo()

	var buf bytes.Buffer

	if err := sysInfo.write(&buf, "hallo", "world"); err != nil {
		t.Fatal(err)
	}

	t.Log(buf.Bytes())

}
