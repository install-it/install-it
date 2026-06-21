package matching

import (
	"fmt"
	"math"
	"strings"

	"install-it/pkg/storage"
	"install-it/pkg/sysinfo"
)

// WMIHardwareQuerier queries hardware via WMI and formats results
// into strings matching the frontend's getHardware() formatting.
type WMIHardwareQuerier struct{}

// HardwareMap queries WMI and returns formatted strings per rule source.
// Errors from individual hardware queries are silently ignored —
// this is intentional degradation for systems missing WMI providers.
func (WMIHardwareQuerier) HardwareMap() (map[storage.RuleSource][]string, error) {
	si := sysinfo.SysInfo{}
	hw := make(map[storage.RuleSource][]string)

	if cpus, err := si.CpuInfo(); err == nil {
		for _, v := range cpus {
			hw[storage.Cpu] = append(hw[storage.Cpu], v.Name)
		}
	}

	if gpus, err := si.GpuInfo(); err == nil {
		for _, v := range gpus {
			gb := int(math.Round(float64(v.AdapterRAM) / math.Pow(1024, 3)))
			hw[storage.Gpu] = append(hw[storage.Gpu], fmt.Sprintf("%s (%dGB)", v.Name, gb))
		}
	}

	if mems, err := si.MemoryInfo(); err == nil {
		for _, v := range mems {
			gb := float64(v.Capacity) / math.Pow(1024, 3)
			gbStr := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.10g", gb), "0"), ".")
			hw[storage.Memory] = append(hw[storage.Memory], fmt.Sprintf("%s %s %sGB %dMHz",
				v.Manufacturer, strings.TrimSpace(v.PartNumber), gbStr, v.Speed))
		}
	}

	if boards, err := si.MotherboardInfo(); err == nil {
		for _, v := range boards {
			hw[storage.Motherboard] = append(hw[storage.Motherboard], fmt.Sprintf("%s %s", v.Manufacturer, v.Product))
		}
	}

	if nics, err := si.NicInfo(); err == nil {
		for _, v := range nics {
			hw[storage.Nic] = append(hw[storage.Nic], v.Name)
		}
	}

	if disks, err := si.DiskInfo(); err == nil {
		for _, v := range disks {
			gb := int(math.Round(float64(v.Size) / math.Pow(1024, 3)))
			hw[storage.Storage] = append(hw[storage.Storage], fmt.Sprintf("%s (%dGB)", v.Model, gb))
		}
	}

	return hw, nil
}