package v1

import "fmt"

type Condition interface {
	ConditionType() ConditionType
}

type ConditionType interface {
	fmt.Stringer
}
