package alicloud

var PvtzThrottlingUserCatcher = Catcher{ThrottlingUser, 30, 2}
var PvtzSystemBusyCatcher = Catcher{"System.Busy", 30, 5}

func PvtzInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(PvtzThrottlingUserCatcher)
	i.AddCatcher(ServiceBusyCatcher)
	i.AddCatcher(PvtzSystemBusyCatcher)
	return i
}
