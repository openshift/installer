package alicloud

type PrimaryKeyTypeString string

const (
	IntegerType = PrimaryKeyTypeString("Integer")
	StringType  = PrimaryKeyTypeString("String")
	BinaryType  = PrimaryKeyTypeString("Binary")
)

type InstanceAccessedByType string

const (
	AnyNetwork   = InstanceAccessedByType("Any")
	VpcOnly      = InstanceAccessedByType("Vpc")
	VpcOrConsole = InstanceAccessedByType("ConsoleOrVpc")
)

type OtsInstanceType string

const (
	OtsCapacity        = OtsInstanceType("Capacity")
	OtsHighPerformance = OtsInstanceType("HighPerformance")
)

func convertInstanceAccessedBy(accessed InstanceAccessedByType) string {
	switch accessed {
	case VpcOnly:
		return "VPC"
	case VpcOrConsole:
		return "VPC_CONSOLE"
	default:
		return "NORMAL"
	}
}

func convertInstanceAccessedByRevert(network string) InstanceAccessedByType {
	switch network {
	case "VPC":
		return VpcOnly
	case "VPC_CONSOLE":
		return VpcOrConsole
	default:
		return AnyNetwork
	}
}

func convertInstanceType(instanceType OtsInstanceType) string {
	switch instanceType {
	case OtsHighPerformance:
		return "SSD"
	default:
		return "HYBRID"
	}
}

func convertInstanceTypeRevert(instanceType string) OtsInstanceType {
	switch instanceType {
	case "SSD":
		return OtsHighPerformance
	default:
		return OtsCapacity
	}
}

// OTS instance total status: S_RUNNING = 1, S_DISABLED = 2, S_DELETING = 3
func convertOtsInstanceStatus(status Status) int {
	switch status {
	case Running:
		return 1
	case DisabledStatus:
		return 2
	case Deleting:
		return 3
	default:
		return -1
	}
}

func convertOtsInstanceStatusConvert(status int) Status {
	switch status {
	case 1:
		return Running
	case 2:
		return DisabledStatus
	case 3:
		return Deleting
	default:
		return ""
	}
}
