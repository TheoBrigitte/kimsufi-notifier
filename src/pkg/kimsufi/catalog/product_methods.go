package catalog

import "fmt"

// Format returns a formatted string representation of the bandwidth.
// e.g. 100 Mbit/s
func (b ProductBlobsTechnicalBandwidth) Format() string {
	return fmt.Sprintf("%.0f Mbit/s", b.Level)
}

// Format returns a formatted string representation of the CPU.
// e.g. Intel Xeon E3 1245v6 3.7 GHz
func (c ProductBlobsTechnicalCPU) Format() string {
	return fmt.Sprintf("%s %s %.2f GHz",
		c.Brand,
		c.Model,
		c.Frequency,
	)
}

// Format returns a formatted string representation of the memory.
// e.g. 64 Go DDR4 ECC 2133 MHz
func (m ProductBlobsTechnicalMemory) Format() string {
	return fmt.Sprintf("%d Go %s",
		m.Size,
		m.RamType,
	)
}

// Format returns a slice of formatted string representation of the storage disks.
// Calls Format() on each disk
func (s ProductBlobsTechnicalStorage) Format() []string {
	var disks []string
	for _, disk := range s.Disks {
		disks = append(disks, disk.Format())
	}

	return disks
}

// FormatFirst returns a formatted string representation of the first disk.
func (s ProductBlobsTechnicalStorage) FormatFirst() string {
	if len(s.Disks) <= 0 {
		return "-"
	}

	disk := s.Disks[0]

	return disk.Format()
}

// Format returns a formatted string representation of a storage disk.
// e.g. 2 x 240 Go SSD
func (d ProductBlobsTechnicalStorageDisk) Format() string {
	return fmt.Sprintf("%d x %d Go %s",
		d.Number,
		d.Capacity,
		d.Technology,
	)
}
