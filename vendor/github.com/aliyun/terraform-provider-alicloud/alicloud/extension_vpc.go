package alicloud

const (
	EcsInstance = "EcsInstance"
	SlbInstance = "SlbInstance"
	Nat         = "Nat"
	HaVip       = "HaVip"
)

type RouterType string
type Role string
type Spec string

const (
	VRouter = RouterType("VRouter")
	VBR     = RouterType("VBR")

	InitiatingSide = Role("InitiatingSide")
	AcceptingSide  = Role("AcceptingSide")

	Mini2   = Spec("Mini.2")
	Mini5   = Spec("Mini.5")
	Small1  = Spec("Small.1")
	Small2  = Spec("Small.2")
	Small5  = Spec("Small.5")
	Middle1 = Spec("Middle.1")
	Middle2 = Spec("Middle.2")
	Middle5 = Spec("Middle.5")
	Large1  = Spec("Large.1")
	Large2  = Spec("Large.2")
	Large5  = Spec("Large.5")
	Xlarge1 = Spec("Xlarge.1")

	Negative = Spec(("Negative"))
)

func GetAllRouterInterfaceSpec() (specifications []string) {
	specifications = append(specifications, string(Mini2), string(Mini5),
		string(Small1), string(Small2), string(Small5),
		string(Middle1), string(Middle2), string(Middle5),
		string(Large1), string(Large2), string(Large5), string(Xlarge1),
		string(Negative))
	return
}
