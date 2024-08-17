package kimsufi

type KeyFunc func(Availability) string

func PlanCode(a Availability) string {
	return a.PlanCode
}
