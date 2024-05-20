package ovirtclient

import (
	"fmt"
	"net"

	"github.com/google/uuid"
	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateNIC(
	vmid VMID,
	vnicProfileID VNICProfileID,
	name string,
	params OptionalNICParameters,
	retries ...RetryStrategy,
) (result NIC, err error) {
	if err := validateNICCreationParameters(vmid, name); err != nil {
		return nil, err
	}

	var mac string
	if params != nil {
		if err := validateNICCreationOptionalParameters(params); err != nil {
			return nil, err
		}
		mac = params.Mac()
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

			if mac != "" {
				nicBuilder.Mac(ovirtsdk.NewMacBuilder().Address(mac).MustBuild())
			}

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
	params OptionalNICParameters,
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

	if params != nil {
		if err := validateNICCreationOptionalParameters(params); err != nil {
			return nil, err
		}
		nic.mac = params.Mac()
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

func validateNICCreationOptionalParameters(params OptionalNICParameters) error {
	if mac := params.Mac(); mac != "" {
		if _, err := net.ParseMAC(mac); err != nil {
			return newError(EUnidentified, "Failed to parse MacAddress: %s", mac)
		}
	}
	return nil
}
