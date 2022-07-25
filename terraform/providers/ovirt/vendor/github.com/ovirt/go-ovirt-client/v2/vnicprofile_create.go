package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateVNICProfile(
	name string,
	networkID NetworkID,
	params OptionalVNICProfileParameters,
	retries ...RetryStrategy,
) (result VNICProfile, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))

	if err := validateVNICProfileCreationParameters(name, networkID, params); err != nil {
		return nil, err
	}

	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("creating VNIC profile %s", name),
		o.logger,
		retries,
		func() error {
			profileBuilder := ovirtsdk.NewVnicProfileBuilder()
			profileBuilder.Name(name)
			profileBuilder.Network(ovirtsdk.NewNetworkBuilder().Id(string(networkID)).MustBuild())
			req := o.conn.SystemService().VnicProfilesService().Add()
			response, err := req.Profile(profileBuilder.MustBuild()).Send()
			if err != nil {
				return err
			}
			profile, ok := response.Profile()
			if !ok {
				return newFieldNotFound("response from VNIC profile creation", "profile")
			}
			result, err = convertSDKVNICProfile(profile, o)
			return err
		})
	return result, err
}

func (m *mockClient) CreateVNICProfile(
	name string,
	networkID NetworkID,
	params OptionalVNICProfileParameters,
	_ ...RetryStrategy,
) (VNICProfile, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if err := validateVNICProfileCreationParameters(name, networkID, params); err != nil {
		return nil, err
	}

	if _, ok := m.networks[networkID]; !ok {
		return nil, newError(ENotFound, "network not found")
	}

	for _, vnicProfile := range m.vnicProfiles {
		if vnicProfile.name == name {
			return nil, newError(EConflict, "VNIC profile name is already in use")
		}
	}

	id := VNICProfileID(m.GenerateUUID())
	m.vnicProfiles[id] = &vnicProfile{
		client: m,

		id:        id,
		networkID: networkID,
		name:      name,
	}

	return m.vnicProfiles[id], nil
}

func validateVNICProfileCreationParameters(name string, networkID NetworkID, _ OptionalVNICProfileParameters) error {
	if name == "" {
		return newError(EBadArgument, "name cannot be empty for VNIC profile creation")
	}
	if networkID == "" {
		return newError(EBadArgument, "network ID cannot be empty for VNIC profile creation")
	}
	return nil
}
