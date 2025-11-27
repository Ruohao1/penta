package parser

import (
	"net/netip"
	"testing"
)

func TestParseCIDR(t *testing.T) {
	got, err := ParseTargets("10.0.0.0/30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []netip.Addr{
		netip.MustParseAddr("10.0.0.0"),
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("10.0.0.2"),
		netip.MustParseAddr("10.0.0.3"),
	}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d, want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%s, want=%s", i, got[i], want[i])
		}
	}
}

func TestParseRange(t *testing.T) {
	got, err := ParseTargets("10.1-2.0.1-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []netip.Addr{
		netip.MustParseAddr("10.1.0.1"),
		netip.MustParseAddr("10.1.0.2"),
		netip.MustParseAddr("10.2.0.1"),
		netip.MustParseAddr("10.2.0.2"),
	}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d, want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%s, want=%s", i, got[i], want[i])
		}
	}
}

func TestParseSlice(t *testing.T) {
	got, err := ParseTargets("10.0.0.1,10.0.0.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []netip.Addr{
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("10.0.0.2"),
	}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d, want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%s, want=%s", i, got[i], want[i])
		}
	}
}
