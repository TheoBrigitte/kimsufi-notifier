package kimsufi

type KeyFunc func(Availability) string

func HardwareKey(a Availability) string {
	return a.Hardware
}

func RegionKey(a Availability) string {
	return a.Region
}
