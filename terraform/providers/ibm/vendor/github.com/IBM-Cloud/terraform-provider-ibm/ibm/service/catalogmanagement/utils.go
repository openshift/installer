package catalogmanagement

func SIToSS(i []interface{}) []string {
	var ss []string
	for _, iface := range i {
		ss = append(ss, iface.(string))
	}
	return ss
}
