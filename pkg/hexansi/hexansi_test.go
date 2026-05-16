package hexansi_test

import (
	"strings"
	"testing"

	"github.com/cumulus13/go-make_colors/pkg/hexansi"
)

func TestHexToRGB(t *testing.T) {
	cases := []struct {
		input    string
		wantR    uint8
		wantG    uint8
		wantB    uint8
		wantErr  bool
	}{
		{"#FF0000", 255, 0, 0, false},
		{"#00FF00", 0, 255, 0, false},
		{"#0000FF", 0, 0, 255, false},
		{"FF0000", 255, 0, 0, false},
		{"#F00", 255, 0, 0, false},   // 3-digit
		{"F00", 255, 0, 0, false},    // 3-digit no hash
		{"#FFFFFF", 255, 255, 255, false},
		{"#000000", 0, 0, 0, false},
		{"#808080", 128, 128, 128, false},
		{"invalid", 0, 0, 0, true},
		{"#GG0000", 0, 0, 0, true},
	}
	for _, tc := range cases {
		rgb, err := hexansi.HexToRGB(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("HexToRGB(%q) expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("HexToRGB(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if rgb[0] != tc.wantR || rgb[1] != tc.wantG || rgb[2] != tc.wantB {
			t.Errorf("HexToRGB(%q) = %v, want [%d,%d,%d]",
				tc.input, rgb, tc.wantR, tc.wantG, tc.wantB)
		}
	}
}

func TestConvert_TrueColor(t *testing.T) {
	res, err := hexansi.Convert("#FF0000", hexansi.ModeTrueColor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.FG, "38;2;255;0;0") {
		t.Errorf("FG = %q, expected 38;2;255;0;0", res.FG)
	}
	if !strings.Contains(res.BG, "48;2;255;0;0") {
		t.Errorf("BG = %q, expected 48;2;255;0;0", res.BG)
	}
}

func TestConvert_256(t *testing.T) {
	res, err := hexansi.Convert("#FF0000", hexansi.Mode256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.FG, "38;5;") {
		t.Errorf("FG = %q, expected 38;5;... format", res.FG)
	}
}

func TestConvert_16(t *testing.T) {
	res, err := hexansi.Convert("#FF0000", hexansi.Mode16)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(res.FG, "\x1b[") {
		t.Errorf("FG = %q, expected ANSI prefix", res.FG)
	}
}

func TestNameToHex(t *testing.T) {
	cases := []struct {
		name    string
		wantHex string
		wantErr bool
	}{
		{"red", "#FF0000", false},
		{"green", "#008000", false},
		{"blue", "#0000FF", false},
		{"black", "#000000", false},
		{"white", "#FFFFFF", false},
		{"nonexistentcolor12345", "", true},
	}
	for _, tc := range cases {
		hex, err := hexansi.NameToHex(tc.name)
		if tc.wantErr {
			if err == nil {
				t.Errorf("NameToHex(%q) expected error", tc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("NameToHex(%q) unexpected error: %v", tc.name, err)
			continue
		}
		if hex != tc.wantHex {
			t.Errorf("NameToHex(%q) = %q, want %q", tc.name, hex, tc.wantHex)
		}
	}
}

func TestHexToColorName(t *testing.T) {
	name, err := hexansi.HexToColorName("#FF0000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name == "" {
		t.Error("expected non-empty color name")
	}
}

func TestToANSI_ByName(t *testing.T) {
	res, err := hexansi.ToANSI("red", hexansi.ModeTrueColor)
	if err != nil {
		t.Fatalf("ToANSI(red) error: %v", err)
	}
	if !strings.Contains(res.FG, "38;2;") {
		t.Errorf("expected truecolor FG, got %q", res.FG)
	}
}

func TestIsHex(t *testing.T) {
	cases := map[string]bool{
		"FF0000": true,
		"F00":    true,
		"red":    false,
		"GGGGGG": false,
		"":       false,
		"12345":  false,
	}
	for in, want := range cases {
		got := hexansi.IsHex(in)
		if got != want {
			t.Errorf("IsHex(%q) = %v, want %v", in, got, want)
		}
	}
}
