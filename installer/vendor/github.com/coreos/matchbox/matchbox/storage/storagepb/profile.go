package storagepb

import (
	"encoding/json"
	"errors"
)

var (
	ErrIdRequired = errors.New("Id is required")
)

// ParseProfile parses bytes into a Profile.
func ParseProfile(data []byte) (*Profile, error) {
	profile := new(Profile)
	err := json.Unmarshal(data, profile)
	return profile, err
}

// AssertValid validates a Profile. Returns nil if there are no validation
// errors.
func (p *Profile) AssertValid() error {
	// Id is required
	if p.Id == "" {
		return ErrIdRequired
	}
	return nil
}

func (p *Profile) Copy() *Profile {
	return &Profile{
		Id:         p.Id,
		Name:       p.Name,
		IgnitionId: p.IgnitionId,
		CloudId:    p.CloudId,
		GenericId:  p.GenericId,
		Boot:       p.Boot.Copy(),
	}
}

func (b *NetBoot) Copy() *NetBoot {
	initrd := make([]string, len(b.Initrd))
	copy(initrd, b.Initrd)
	args := make([]string, len(b.Args))
	copy(args, b.Args)
	cmdline := make(map[string]string)
	for k, v := range b.Cmdline {
		cmdline[k] = v
	}
	return &NetBoot{
		Kernel: b.Kernel,
		Initrd: initrd,
		Args:   args,
		// deprecated
		Cmdline: cmdline,
	}
}
