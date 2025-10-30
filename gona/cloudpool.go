package gona

type CloudPool uint8

// The /pools API is currently undocumented. Hardcode some values and logic.

const (
	CloudPoolGeneralCompute CloudPool = 1
	CloudPoolAMDEPYC        CloudPool = 9
	CloudPoolDefault        CloudPool = 0
)

func (cp CloudPool) Name() string {
	switch cp {
	case CloudPoolGeneralCompute:
		return "General Compute"
	case CloudPoolAMDEPYC:
		return "AMD EPYC"
	case CloudPoolDefault:
		return "Default"
	default:
		return "Unknown"
	}
}

func CloudPoolFromName(s string) CloudPool {
	switch s {
	case CloudPoolGeneralCompute.Name():
		return CloudPoolGeneralCompute
	case CloudPoolAMDEPYC.Name():
		return CloudPoolAMDEPYC
	case CloudPoolDefault.Name():
		return CloudPoolDefault
	}
	return CloudPoolDefault
}
