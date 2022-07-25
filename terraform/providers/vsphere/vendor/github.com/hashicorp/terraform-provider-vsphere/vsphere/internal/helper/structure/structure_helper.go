package structure

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/govmomi/vim25/types"
)

// ResourceIDStringer is a small interface that can be used to supply
// ResourceData and ResourceDiff to functions that need to print the ID of a
// resource, namely used by logging.
type ResourceIDStringer interface {
	Id() string
}

// ResourceIDString prints a friendly string for a resource, supplied by name.
func ResourceIDString(d ResourceIDStringer, name string) string {
	id := d.Id()
	if id == "" {
		id = "<new resource>"
	}
	return fmt.Sprintf("%s (ID = %s)", name, id)
}

// SliceInterfacesToStrings converts an interface slice to a string slice. The
// function does not attempt to do any sanity checking and will panic if one of
// the items in the slice is not a string.
func SliceInterfacesToStrings(s []interface{}) []string {
	var d []string
	for _, v := range s {
		if o, ok := v.(string); ok {
			d = append(d, o)
		}
	}
	return d
}

// SliceStringsToInterfaces converts a string slice to an interface slice.
func SliceStringsToInterfaces(s []string) []interface{} {
	var d []interface{}
	for _, v := range s {
		d = append(d, v)
	}
	return d
}

// SliceInterfacesToManagedObjectReferences converts an interface slice into a
// slice of ManagedObjectReferences with the type of t.
func SliceInterfacesToManagedObjectReferences(s []interface{}, t string) []types.ManagedObjectReference {
	var d []types.ManagedObjectReference
	for _, v := range s {
		d = append(d, types.ManagedObjectReference{
			Type:  t,
			Value: v.(string),
		})
	}
	return d
}

// SliceStringsToManagedObjectReferences converts a string slice into a slice
// of ManagedObjectReferences with the type of t.
func SliceStringsToManagedObjectReferences(s []string, t string) []types.ManagedObjectReference {
	var d []types.ManagedObjectReference
	for _, v := range s {
		d = append(d, types.ManagedObjectReference{
			Type:  t,
			Value: v,
		})
	}
	return d
}

// MergeSchema merges the map[string]*schema.Schema from src into dst. Safety
// against conflicts is enforced by panicing.
func MergeSchema(dst, src map[string]*schema.Schema) {
	for k, v := range src {
		if _, ok := dst[k]; ok {
			panic(fmt.Errorf("conflicting schema key: %s", k))
		}
		dst[k] = v
	}
}

// BoolPtr makes a *bool out of the value passed in through v.
//
// vSphere uses nil values in bools to omit values in the SOAP XML request, and
// helps denote inheritance in certain cases.
func BoolPtr(v bool) *bool {
	return &v
}

// GetBoolPtr reads a ResourceData and returns an appropriate *bool for the
// state of the definition. nil is returned if it does not exist.
func GetBoolPtr(d *schema.ResourceData, key string) *bool {
	v, e := d.GetOkExists(key)
	if e {
		return BoolPtr(v.(bool))
	}
	return nil
}

// GetBool reads a ResourceData and returns a *bool. This differs from
// GetBoolPtr in that a nil value is never returned.
func GetBool(d *schema.ResourceData, key string) *bool {
	return BoolPtr(d.Get(key).(bool))
}

// SetBoolPtr sets a ResourceData field depending on if a *bool exists or not.
// The field is not set if it's nil.
func SetBoolPtr(d *schema.ResourceData, key string, val *bool) error {
	if val == nil {
		return nil
	}
	err := d.Set(key, val)
	return err
}

// GetBoolStringPtr reads a ResourceData *string* field. This field is handled
// in the following way:
//
// * If it's empty, nil is returned.
// * The string is then sent through ParseBool. This will return a valid value
// for anything ParseBool returns a value for.
// * If it's anything else, an error is returned.
//
// This is designed to address the current lack of HCL and Terraform to be able
// to distinguish between nil states and zero values properly. This is a
// systemic issue that affects reading, writing, and diffing of these values.
// These issues will eventually be addressed in HCL2.
func GetBoolStringPtr(d *schema.ResourceData, key string) (*bool, error) {
	v, ok := d.GetOk(key)
	if !ok {
		return nil, nil
	}
	b, err := strconv.ParseBool(v.(string))
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// SetBoolStringPtr sets a stringified ResoruceData bool field. This is a field
// that is supposed to behave like a bool (true/false), but needs to be a
// string to represent a nil state as well.
//
// This is designed to address the current lack of HCL and Terraform to be able
// to distinguish between nil states and zero values properly. This is a
// systemic issue that affects reading, writing, and diffing of these values.
// These issues will eventually be addressed in HCL2.
func SetBoolStringPtr(d *schema.ResourceData, key string, val *bool) error {
	var s string
	if val != nil {
		s = strconv.FormatBool(*val)
	}
	return d.Set(key, s)
}

// BoolStringPtrState is a state normalization function for stringified 3-state
// bool pointers.
//
// The function silently drops any result that can't be parsed with ParseBool,
// and will return an empty string for these cases.
//
// This is designed to address the current lack of HCL and Terraform to be able
// to distinguish between nil states and zero values properly. This is a
// systemic issue that affects reading, writing, and diffing of these values.
// These issues will eventually be addressed in HCL2.
func BoolStringPtrState(v interface{}) string {
	b, err := strconv.ParseBool(v.(string))
	if err != nil {
		return ""
	}
	return strconv.FormatBool(b)
}

// ValidateBoolStringPtr validates that the input value can be parsed by
// ParseBool. It also succeeds on empty strings.
//
// This is designed to address the current lack of HCL and Terraform to be able
// to distinguish between nil states and zero values properly. This is a
// systemic issue that affects reading, writing, and diffing of these values.
// These issues will eventually be addressed in HCL2.
func ValidateBoolStringPtr() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v := i.(string)
		if v == "" {
			return
		}
		if _, err := strconv.ParseBool(v); err != nil {
			es = append(es, err)
		}
		return
	}
}

// Int64Ptr makes an *int64 out of the value passed in through v.
func Int64Ptr(v int64) *int64 {
	return &v
}

// Int32Ptr makes an *int32 out of the value passed in through v.
func Int32Ptr(v int32) *int32 {
	return &v
}

// GetInt64Ptr reads a ResourceData and returns an appropriate *int64 for the
// state of the definition. nil is returned if it does not exist.
func GetInt64Ptr(d *schema.ResourceData, key string) *int64 {
	v, e := d.GetOkExists(key)
	if e {
		return Int64Ptr(int64(v.(int)))
	}
	return nil
}

// GetInt64PtrEmptyZero reads a ResourceData and returns an appropriate *int64
// for the state of the definition. 0 is returned if it does not exist.
func GetInt64PtrEmptyZero(d *schema.ResourceData, key string) *int64 {
	i := GetInt64Ptr(d, key)
	if i != nil {
		return i
	}
	return Int64Ptr(int64(0))
}

// SetInt64Ptr sets a ResourceData field depending on if an *int64 exists or
// not.  The field is not set if it's nil.
func SetInt64Ptr(d *schema.ResourceData, key string, val *int64) error {
	if val == nil {
		return nil
	}
	err := d.Set(key, val)
	return err
}

// ByteToMB returns n/1000000. The input must be an integer that can be divisible
// by 1000000.
func ByteToMB(n interface{}) interface{} {
	switch v := n.(type) {
	case int:
		return v / 1000000
	case int32:
		return v / 1000000
	case int64:
		return v / 1000000
	}
	panic(fmt.Errorf("non-integer type %T for value", n))
}

// ByteToGiB returns n/1024^3, *rounded up*.
//
// Standard integer division results in fractional GiB being discarded,
// resulting in errors errors cloning virtual machines having disk size
// in non-integer GiB. The result is rounded up to avoid this edge case.
func ByteToGiB(n interface{}) int {
	switch n.(type) {
	case int, int32, int64:
		return int(math.Ceil(float64(n.(int64)) / math.Pow(1024, 3)))
	}
	panic(fmt.Errorf("non-integer type %T for value", n))
}

// GiBToByte returns n*1024^3.
//
// The output is returned as int64 - if another type is needed, it needs to be
// cast. Remember that int32 overflows at around 2GiB and uint32 will overflow at 4GiB.
func GiBToByte(n interface{}) int64 {
	switch v := n.(type) {
	case int:
		return int64(v * int(math.Pow(1024, 3)))
	case int32:
		return int64(v * int32(math.Pow(1024, 3)))
	case int64:
		return v * int64(math.Pow(1024, 3))
	}
	panic(fmt.Errorf("non-integer type %T for value", n))
}

// BoolPolicy converts a bool into a VMware BoolPolicy value.
func BoolPolicy(b bool) *types.BoolPolicy {
	bp := &types.BoolPolicy{
		Value: BoolPtr(b),
	}
	return bp
}

// GetBoolPolicy reads a ResourceData and returns an appropriate BoolPolicy for
// the state of the definition. nil is returned if it does not exist.
func GetBoolPolicy(d *schema.ResourceData, key string) *types.BoolPolicy {
	v, e := d.GetOkExists(key)
	if e {
		return BoolPolicy(v.(bool))
	}
	return nil
}

// SetBoolPolicy sets a ResourceData field depending on if a BoolPolicy exists
// or not. The field is not set if it's nil.
func SetBoolPolicy(d *schema.ResourceData, key string, val *types.BoolPolicy) error {
	if val == nil {
		return nil
	}
	err := d.Set(key, val.Value)
	return err
}

// GetBoolPolicyReverse acts like GetBoolPolicy, but the value is inverted.
func GetBoolPolicyReverse(d *schema.ResourceData, key string) *types.BoolPolicy {
	v, e := d.GetOkExists(key)
	if e {
		return BoolPolicy(!v.(bool))
	}
	return nil
}

// SetBoolPolicyReverse acts like SetBoolPolicy, but the value is inverted.
func SetBoolPolicyReverse(d *schema.ResourceData, key string, val *types.BoolPolicy) error {
	if val == nil {
		return nil
	}
	err := d.Set(key, !*val.Value)
	return err
}

// StringPolicy converts a string into a VMware StringPolicy value.
func StringPolicy(s string) *types.StringPolicy {
	sp := &types.StringPolicy{
		Value: s,
	}
	return sp
}

// GetStringPolicy reads a ResourceData and returns an appropriate StringPolicy
// for the state of the definition. nil is returned if it does not exist.
func GetStringPolicy(d *schema.ResourceData, key string) *types.StringPolicy {
	v, e := d.GetOkExists(key)
	if e {
		return StringPolicy(v.(string))
	}
	return nil
}

// SetStringPolicy sets a ResourceData field depending on if a StringPolicy
// exists or not. The field is not set if it's nil.
func SetStringPolicy(d *schema.ResourceData, key string, val *types.StringPolicy) error {
	if val == nil {
		return nil
	}
	err := d.Set(key, val.Value)
	return err
}

// LongPolicy converts a supported number into a VMware LongPolicy value. This
// will panic if there is no implicit conversion of the value into an int64.
func LongPolicy(n interface{}) *types.LongPolicy {
	lp := &types.LongPolicy{}
	switch v := n.(type) {
	case int:
		lp.Value = int64(v)
	case int8:
		lp.Value = int64(v)
	case int16:
		lp.Value = int64(v)
	case int32:
		lp.Value = int64(v)
	case uint:
		lp.Value = int64(v)
	case uint8:
		lp.Value = int64(v)
	case uint16:
		lp.Value = int64(v)
	case uint32:
		lp.Value = int64(v)
	case int64:
		lp.Value = v
	default:
		panic(fmt.Errorf("non-convertible type %T for value", n))
	}
	return lp
}

// GetLongPolicy reads a ResourceData and returns an appropriate LongPolicy
// for the state of the definition. nil is returned if it does not exist.
func GetLongPolicy(d *schema.ResourceData, key string) *types.LongPolicy {
	v, e := d.GetOkExists(key)
	if e {
		return LongPolicy(v)
	}
	return nil
}

// SetLongPolicy sets a ResourceData field depending on if a LongPolicy
// exists or not. The field is not set if it's nil.
func SetLongPolicy(d *schema.ResourceData, key string, val *types.LongPolicy) error {
	if val == nil {
		return nil
	}
	err := d.Set(key, val.Value)
	return err
}

// AllFieldsEmpty checks to see if all fields in a given struct are zero
// values. It does not recurse, so finer-grained checking should be done for
// deep accuracy when necessary. It also does not dereference pointers, except
// if the value itself is a pointer and is not nil.
func AllFieldsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct && (t.Kind() == reflect.Ptr && t.Elem().Kind() != reflect.Struct) {
		return reflect.Zero(t).Interface() == reflect.ValueOf(v).Interface()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		fv := reflect.ValueOf(v).Elem().Field(i)
		ft := t.Field(i).Type
		fz := reflect.Zero(ft)
		switch ft.Kind() {
		case reflect.Map, reflect.Slice:
			if fv.Len() > 0 {
				return false
			}
		default:
			if fz.Interface() != fv.Interface() {
				return false
			}
		}
	}

	return true
}

// DeRef returns the value pointed to by the interface if the interface is a
// pointer and is not nil, otherwise returns nil, or the direct value if it's
// not a pointer.
func DeRef(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	k := reflect.TypeOf(v).Kind()
	if k != reflect.Ptr {
		return v
	}
	if reflect.ValueOf(v) == reflect.Zero(reflect.TypeOf(v)) {
		// All zero-value pointers are nil
		return nil
	}
	return reflect.ValueOf(v).Elem().Interface()
}

// NormalizeValue converts a value to something that is suitable to be set in a
// ResourceData and can be useful in situations where there is not access to
// normal helper/schema functionality, but you still need saved fields to
// behave in the same way.
//
// Specifically, this will run the value through DeRef to dereference any
// pointers first, and then convert numeric primitives, if necessary.
func NormalizeValue(v interface{}) interface{} {
	v = DeRef(v)
	if v == nil {
		return nil
	}
	k := reflect.TypeOf(v).Kind()
	switch {
	case k >= reflect.Int8 && k <= reflect.Uint64:
		v = reflect.ValueOf(v).Convert(reflect.TypeOf(int(0))).Interface()
	case k == reflect.Float32:
		v = reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0))).Interface()
	}
	return v
}

// LogCond takes a boolean (which can be passed in as a bool or as a
// conditional), and returns either the first value if true, and the second if
// false. It's designed to be a "ternary" operator of sorts for logging.
func LogCond(c bool, t, f interface{}) interface{} {
	if c {
		return t
	}
	return f
}

// SetBatch takes a map of values and sets the appropriate top-level attributes
// for each item.
//
// attrs is a map[string]interface{} that follows a pattern in the example
// below:
//
//   err := SetBatch(d, map[string]interface{}{
//  	"foo": obj.Foo,
//  	"bar": obj.Bar,
//   })
//   if err != nil {
//  	return err
//   }
//
// For best results, supplied values should be or have concrete values that map
// to the correct values for the respective type in helper/schema. This is
// enforced by way of checking each Set call for errors. If there is an error
// setting a particular key, processing stops immediately.
func SetBatch(d *schema.ResourceData, attrs map[string]interface{}) error {
	for k, v := range attrs {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("error setting attribute %q: %s", k, err)
		}
	}

	return nil
}

// MoRefSorter is a sorting wrapper for a slice of MangedObjectReference.
type MoRefSorter []types.ManagedObjectReference

// Len implements sort.Interface for MoRefSorter.
func (s MoRefSorter) Len() int {
	return len(s)
}

// Less helps implement sort.Interface for MoRefSorter.
func (s MoRefSorter) Less(i, j int) bool {
	return s[i].Value < s[j].Value
}

// Swap helps implement sort.Interface for MoRefSorter.
func (s MoRefSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// ValuesAvailable takes a subresource path and a list of keys and checks that
// the value for each key is available at CustomizeDiff time. This function
// will return false if any of they values are based on computed values from
// other new or updated resources.
func ValuesAvailable(base string, keys []string, d *schema.ResourceDiff) bool {
	for _, k := range keys {
		if !d.NewValueKnown(fmt.Sprintf("%s%s", base, k)) {
			return false
		}
	}
	return true
}

func DiffSlice(a, b []interface{}) []interface{} {
	var c []interface{}
	for _, aa := range a {
		found := false
		for _, bb := range b {
			if reflect.DeepEqual(aa, bb) {
				found = true
			}
		}
		if !found {
			c = append(c, aa)
		}
	}
	return c
}

func DropSliceItem(a []interface{}, n int) []interface{} {
	var b []interface{}
	if n > 0 {
		b = a[:n-1]
	}
	if n < len(a)-1 {
		b = append(b, a[n+1:]...)
	}
	return b
}
