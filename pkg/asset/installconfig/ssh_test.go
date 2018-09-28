package installconfig

import (
	"testing"
)

func TestOpenSSHPublicKey(t *testing.T) {
	const invalidMsg = "invalid SSH public key"
	const multiLineMsg = "invalid SSH public key (should not contain any newline characters)"
	const privateKeyMsg = "invalid SSH public key (appears to be a private key)"
	tests := []struct {
		in       string
		expected string
	}{
		{"a", invalidMsg},
		{".", invalidMsg},
		{"日本語", invalidMsg},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL", ""},
		{"ssh-rsa \t AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL", ""},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL you@example.com", ""},
		{"\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL you@example.com", ""},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL you@example.com\n", ""},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL", multiLineMsg},
		{"ssh-rsa\nAAAAB3NzaC1yc2EAAAADAQABAAACAQDxL you@example.com", multiLineMsg},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL\nyou@example.com", multiLineMsg},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL", ""},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCt3BebCHqnSsgpLjo4kVvyfY/z2BS8t27r/7du+O2pb4xYkr7n+KFpbOz523vMTpQ+o1jY4u4TgexglyT9nqasWgLOvo1qjD1agHme8LlTPQSk07rXqOB85Uq5p7ig2zoOejF6qXhcc3n1c7+HkxHrgpBENjLVHOBpzPBIAHkAGaZcl07OCqbsG5yxqEmSGiAlh/IiUVOZgdDMaGjCRFy0wk0mQaGD66DmnFc1H5CzcPjsxr0qO65e7lTGsE930KkO1Vc+RHCVwvhdXs+c2NhJ2/3740Kpes9n1/YullaWZUzlCPDXtRuy6JRbFbvy39JUgHWGWzB3d+3f8oJ/N4qZ cardno:000603633110", ""},
		{"-----BEGIN CERTIFICATE-----abcd-----END CERTIFICATE-----", invalidMsg},
		{"-----BEGIN RSA PRIVATE KEY-----\nabc\n-----END RSA PRIVATE KEY-----", privateKeyMsg},
	}

	for _, test := range tests {
		err := validateOpenSSHPublicKey(test.in)
		if (err == nil && test.expected != "") || (err != nil && err.Error() != test.expected) {
			t.Errorf("For %q, expected %q, got %q", test.in, test.expected, err)
		}
	}
}
