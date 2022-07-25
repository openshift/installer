package ovirtclient

import (
	"fmt"

	"github.com/google/uuid"
	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateNIC(
	vmid VMID,
	vnicProfileID VNICProfileID,
	name string,
	_ OptionalNICParameters,
	retries ...RetryStrategy,
) (result NIC, err error) {
	if err := validateNICCreationParameters(vmid, name); err != nil {
		return nil, err
	}

	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("creating NIC for VM %s", vmid),
		o.logger,
		retries,
		func() error {
			nicBuilder := ovirtsdk.NewNicBuilder()
			nicBuilder.Name(name)
			nicBuilder.VnicProfile(ovirtsdk.NewVnicProfileBuilder().Id(string(vnicProfileID)).MustBuild())
			nic := nicBuilder.MustBuild()

			response, err := o.conn.SystemService().VmsService().VmService(string(vmid)).NicsService().Add().Nic(nic).Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Nic()
			if !ok {
				return newError(
					ENotFound,
					"no NIC returned creating NIC for VM ID %s",
					vmid,
				)
			}
			result, err = convertSDKNIC(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert newly created NIC for VM %s",
					vmid,
				)
			}
			return nil
		},
	)
	return result, err
}

func (m *mockClient) CreateNIC(
	vmid VMID,
	vnicProfileID VNICProfileID,
	name string,
	_ OptionalNICParameters,
	_ ...RetryStrategy,
) (NIC, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if err := validateNICCreationParameters(vmid, name); err != nil {
		return nil, err
	}
	if _, ok := m.vms[vmid]; !ok {
		return nil, newError(ENotFound, "VM with ID %s not found for NIC creation", vmid)
	}
	for _, n := range m.nics {
		if n.name == name {
			return nil, newError(ENotFound, "NIC with name %s is already in use", name)
		}
	}

	id := NICID(uuid.Must(uuid.NewUUID()).String())

	nic := &nic{
		client:        m,
		id:            id,
		name:          name,
		vmid:          vmid,
		vnicProfileID: vnicProfileID,
	}
	m.nics[id] = nic

	return nic, nil
}

func validateNICCreationParameters(vmid VMID, name string) error {
	if vmid == "" {
		return newError(EBadArgument, "VM ID cannot be empty")
	}
	if name == "" {
		return newError(EBadArgument, "NIC name cannot be empty")
	}
	return nil
}
