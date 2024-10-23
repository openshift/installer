package ocm

import (
	"fmt"
	"strconv"
	"time"

	commonUtils "github.com/openshift-online/ocm-common/pkg/utils"
)

func Int32Validator(val interface{}) error {
	if val == "" { // if a value is not passed it should not throw an error (optional value)
		return nil
	}
	_, err := strconv.ParseInt(fmt.Sprintf("%v", val), 10, 32)
	if err != nil {
		return fmt.Errorf("Should provide an integer number between -2147483648 to 2147483647.")
	}
	return nil
}

func NonNegativeInt32Validator(val interface{}) error {
	if val == "" { // if a value is not passed it should not throw an error (optional value)
		return nil
	}
	number, err := strconv.ParseInt(fmt.Sprintf("%v", val), 10, 32)
	if err != nil {
		return fmt.Errorf("Should provide an integer number between 0 to 2147483647.")
	}

	if number < 0 {
		return fmt.Errorf("Number must be greater or equal to zero.")
	}

	return nil
}

func PositiveDurationStringValidator(val interface{}) error {
	if val == "" {
		return nil
	}
	input, ok := val.(string)

	if !ok {
		return fmt.Errorf("Can only validate strings, got %v", val)
	}

	duration, err := time.ParseDuration(input)

	if err != nil {
		return err
	}

	if duration < 0 {
		return fmt.Errorf("Only positive durations are allowed, got '%v'", val)
	}

	return nil
}

func PercentageValidator(val interface{}) error {
	if val == "" {
		return nil
	}

	number, err := strconv.ParseFloat(fmt.Sprintf("%v", val), commonUtils.MaxByteSize)
	if err != nil {
		return fmt.Errorf("Failed parsing '%v' into a floating-point number.", val)
	}

	if number > 1 || number < 0 {
		return fmt.Errorf("Expecting a floating-point number between 0 and 1.")
	}

	return nil
}
