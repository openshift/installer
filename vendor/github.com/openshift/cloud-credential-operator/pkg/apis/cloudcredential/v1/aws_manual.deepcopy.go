package v1

// DeepCopy is a deepcopy function, copying the receiver, creating a new IAMPolicyCondition.
func (in *IAMPolicyCondition) DeepCopy() *IAMPolicyCondition {
	if in == nil {
		return nil
	}
	out := new(IAMPolicyCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IAMPolicyCondition) DeepCopyInto(out *IAMPolicyCondition) {
	if *in == nil {
		return
	}

	*out = make(IAMPolicyCondition, len(*in))
	tgt := *out

	for key, val := range *in {
		if val != nil {
			tgt[key] = make(IAMPolicyConditionKeyValue, len(val))
			for subKey, subVal := range val {
				tgt[key][subKey] = copyStringOrStringSlice(subVal)
			}
		}
	}
}

func copyStringOrStringSlice(from interface{}) interface{} {
	var to interface{}

	switch v := from.(type) {
	case string:
		// simple assignment creates copy
		to = from
	case []string:
		toSlice := make([]string, len(v))
		copy(toSlice, v)
		to = toSlice
	default:
		// unexpected type, this could mean we're not
		// doing a deepcopy
		to = from
	}

	return to
}
