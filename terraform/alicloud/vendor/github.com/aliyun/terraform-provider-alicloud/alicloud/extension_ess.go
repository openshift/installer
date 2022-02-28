package alicloud

const UserId = "userId"
const ScalingGroup = "scaling_group"
const GroupId = "groupId"

type ScalingRuleType string

const (
	SimpleScalingRule         = ScalingRuleType("SimpleScalingRule")
	TargetTrackingScalingRule = ScalingRuleType("TargetTrackingScalingRule")
	StepScalingRule           = ScalingRuleType("StepScalingRule")
)

type BatchSize int

const (
	AttachDetachLoadbalancersBatchsize = BatchSize(5)
	AttachDetachDbinstancesBatchsize   = BatchSize(5)
)

type MaxItems int

const (
	MaxScalingConfigurationInstanceTypes = MaxItems(10)
)
