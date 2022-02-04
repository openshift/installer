package validation

import (
	"testing"

	_ "github.com/stretchr/testify/assert"
	_ "k8s.io/apimachinery/pkg/util/validation/field"

	_ "github.com/openshift/installer/pkg/types/powervs"
)

func TestValidatePlatform(t *testing.T) {}
