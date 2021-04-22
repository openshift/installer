package dns

import "net"

func lookupIP(host string) ([]string, []string, error) {
	records, err := net.LookupIP(host)
	if err != nil {
		return nil, nil, err
	}

	a := make([]string, 0)
	aaaa := make([]string, 0)
	for _, ip := range records {
		if ipv4 := ip.To4(); ipv4 != nil {
			a = append(a, ipv4.String())
		} else {
			aaaa = append(aaaa, ip.String())
		}
	}

	return a, aaaa, nil
}
