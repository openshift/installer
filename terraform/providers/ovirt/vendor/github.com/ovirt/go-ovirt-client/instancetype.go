package ovirtclient

import ovirtsdk "github.com/ovirt/go-ovirt"

// InstanceTypeID is a type alias for instance type IDs.
type InstanceTypeID string

// InstanceTypeClient lists the methods for working with instance types.
type InstanceTypeClient interface {
	GetInstanceType(id InstanceTypeID, retries ...RetryStrategy) (InstanceType, error)
	ListInstanceTypes(retries ...RetryStrategy) ([]InstanceType, error)
}

// InstanceTypeData is the data segment of the InstanceType type.
type InstanceTypeData interface {
	ID() InstanceTypeID
	Name() string
}

// InstanceType is a data structure that contains preconfigured instance parameters.
type InstanceType interface {
	InstanceTypeData
}

func convertSDKInstanceType(object *ovirtsdk.InstanceType, o *oVirtClient) (InstanceType, error) {
	name, ok := object.Name()
	if !ok {
		return nil, newFieldNotFound("instance type", "name")
	}

	return &instanceType{
		o,

		InstanceTypeID(object.MustId()),
		name,
	}, nil
}

type instanceType struct {
	client Client

	id   InstanceTypeID
	name string
}

func (i instanceType) Name() string {
	return i.name
}

func (i instanceType) ID() InstanceTypeID {
	return i.id
}
