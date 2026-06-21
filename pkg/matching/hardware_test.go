package matching

import (
	"testing"
)

// TestWMIHardwareQuerier_SatisfiesInterface verifies that WMIHardwareQuerier
// implements the HardwareQuerier interface at compile time.
func TestWMIHardwareQuerier_SatisfiesInterface(t *testing.T) {
	var _ HardwareQuerier = WMIHardwareQuerier{}
}
