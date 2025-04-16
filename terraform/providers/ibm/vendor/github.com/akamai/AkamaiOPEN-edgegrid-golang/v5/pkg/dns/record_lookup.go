package dns

import (
	"context"
	"fmt"
	"net/http"

	"encoding/hex"
	"net"
	"strconv"
	"strings"
)

func (p *dns) FullIPv6(ctx context.Context, ip net.IP) string {

	logger := p.Log(ctx)
	logger.Debug("FullIPv6")

	dst := make([]byte, hex.EncodedLen(len(ip)))
	_ = hex.Encode(dst, ip)
	return string(dst[0:4]) + ":" +
		string(dst[4:8]) + ":" +
		string(dst[8:12]) + ":" +
		string(dst[12:16]) + ":" +
		string(dst[16:20]) + ":" +
		string(dst[20:24]) + ":" +
		string(dst[24:28]) + ":" +
		string(dst[28:])
}

func padvalue(str string) string {
	vstr := strings.Replace(str, "m", "", -1)
	vfloat, err := strconv.ParseFloat(vstr, 32)
	if err != nil {
		return "FAIL"
	}

	return fmt.Sprintf("%.2f", vfloat)
}

func (p *dns) PadCoordinates(ctx context.Context, str string) string {

	logger := p.Log(ctx)
	logger.Debug("PadCoordinates")

	s := strings.Split(str, " ")
	if len(s) < 12 {
		return ""
	}

	latd, latm, lats, latDir, longd, longm, longs, longDir, altitude, size, horizPrecision, vertPrecision := s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11]

	return latd + " " + latm + " " + lats + " " + latDir + " " + longd + " " + longm + " " + longs + " " + longDir + " " + padvalue(altitude) + "m " + padvalue(size) + "m " + padvalue(horizPrecision) + "m " + padvalue(vertPrecision) + "m"

}

func (p *dns) GetRecord(ctx context.Context, zone string, name string, recordType string) (*RecordBody, error) {

	logger := p.Log(ctx)
	logger.Debug("GetRecord")

	var rec RecordBody
	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", zone, name, recordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRecord request: %w", err)
	}

	resp, err := p.Exec(req, &rec)
	if err != nil {
		return nil, fmt.Errorf("GetRecord request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rec, nil
}

func (p *dns) GetRecordList(ctx context.Context, zone string, _ string, recordType string) (*RecordSetResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetRecordList")

	var records RecordSetResponse
	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets?types=%s&showAll=true", zone, recordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRecordList request: %w", err)
	}

	resp, err := p.Exec(req, &records)
	if err != nil {
		return nil, fmt.Errorf("GetRecordList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &records, nil
}

func (p *dns) GetRdata(ctx context.Context, zone string, name string, recordType string) ([]string, error) {

	logger := p.Log(ctx)
	logger.Debug("GetRdata")

	records, err := p.GetRecordList(ctx, zone, name, recordType)
	if err != nil {
		return nil, err
	}

	var arrLength int
	for _, c := range records.Recordsets {
		if c.Name == name {
			arrLength = len(c.Rdata)
		}
	}

	rdata := make([]string, 0, arrLength)

	for _, r := range records.Recordsets {
		if r.Name == name {
			for _, i := range r.Rdata {
				str := i

				if recordType == "AAAA" {
					addr := net.ParseIP(str)
					result := p.FullIPv6(ctx, addr)
					str = result
				} else if recordType == "LOC" {
					str = p.PadCoordinates(ctx, str)
				}
				rdata = append(rdata, str)
			}
		}
	}
	return rdata, nil
}

func (p *dns) ProcessRdata(ctx context.Context, rdata []string, rtype string) []string {

	logger := p.Log(ctx)
	logger.Debug("ProcessRdata")

	newrdata := make([]string, 0, len(rdata))
	for _, i := range rdata {
		str := i
		if rtype == "AAAA" {
			addr := net.ParseIP(str)
			result := p.FullIPv6(ctx, addr)
			str = result
		} else if rtype == "LOC" {
			str = p.PadCoordinates(ctx, str)
		}
		newrdata = append(newrdata, str)
	}
	return newrdata

}

func (p *dns) ParseRData(ctx context.Context, rtype string, rdata []string) map[string]interface{} {

	logger := p.Log(ctx)
	logger.Debug("ParseRData")

	fieldMap := make(map[string]interface{}, 0)
	if len(rdata) == 0 {
		return fieldMap
	}
	newrdata := make([]string, 0, len(rdata))
	fieldMap["target"] = newrdata

	switch rtype {
	case "AFSDB":
		resolveAFSDBType(rdata, newrdata, fieldMap)

	case "DNSKEY":
		resolveDNSKEYType(rdata, fieldMap)

	case "DS":
		resolveDSType(rdata, fieldMap)

	case "HINFO":
		resolveHINFOType(rdata, fieldMap)
	/*
		// too many variations to calculate pri and increment
		case "MX":
			sort.Strings(rdata)
			parts := strings.Split(rdata[0], " ")
			fieldMap["priority"], _ = strconv.Atoi(parts[0])
			if len(rdata) > 1 {
				parts = strings.Split(rdata[1], " ")
				tpri, _ := strconv.Atoi(parts[0])
				fieldMap["priority_increment"] = tpri - fieldMap["priority"].(int)
			}
			for _, rcontent := range rdata {
				parts := strings.Split(rcontent, " ")
				newrdata = append(newrdata, parts[1])
			}
			fieldMap["target"] = newrdata
	*/

	case "NAPTR":
		resolveNAPTRType(rdata, fieldMap)

	case "NSEC3":
		resolveNSEC3Type(rdata, fieldMap)

	case "NSEC3PARAM":
		resolveNSEC3PARAMType(rdata, fieldMap)

	case "RP":
		resolveRPType(rdata, fieldMap)

	case "RRSIG":
		resolveRRSIGType(rdata, fieldMap)

	case "SRV":
		resolveSRVType(rdata, newrdata, fieldMap)

	case "SSHFP":
		resolveSSHFPType(rdata, fieldMap)

	case "SOA":
		resolveSOAType(rdata, fieldMap)

	case "AKAMAITLC":
		resolveAKAMAITLCType(rdata, fieldMap)

	case "SPF":
		resolveSPFType(rdata, newrdata, fieldMap)

	case "TXT":
		resolveTXTType(rdata, newrdata, fieldMap)

	case "AAAA":
		resolveAAAAType(ctx, p, rdata, newrdata, fieldMap)

	case "LOC":
		resolveLOCType(ctx, p, rdata, newrdata, fieldMap)

	case "CERT":
		resolveCERTType(rdata, fieldMap)

	case "TLSA":
		resolveTLSAType(rdata, fieldMap)

	case "SVCB":
		resolveSVCBType(rdata, fieldMap)

	case "HTTPS":
		resolveHTTPSType(rdata, fieldMap)

	default:
		for _, rcontent := range rdata {
			newrdata = append(newrdata, rcontent)
		}
		fieldMap["target"] = newrdata
	}

	return fieldMap
}

func resolveAFSDBType(rdata, newrdata []string, fieldMap map[string]interface{}) {
	parts := strings.Split(rdata[0], " ")
	fieldMap["subtype"], _ = strconv.Atoi(parts[0])
	for _, rcontent := range rdata {
		parts = strings.Split(rcontent, " ")
		newrdata = append(newrdata, parts[1])
	}
	fieldMap["target"] = newrdata
}

func resolveDNSKEYType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["flags"], _ = strconv.Atoi(parts[0])
		fieldMap["protocol"], _ = strconv.Atoi(parts[1])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[2])
		key := parts[3]
		// key can have whitespace
		if len(parts) > 4 {
			i := 4
			for i < len(parts) {
				key += " " + parts[i]
			}
		}
		fieldMap["key"] = key
		break
	}
}

func resolveSVCBType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.SplitN(rcontent, " ", 3)
		// has to be at least two fields.
		if len(parts) < 2 {
			break
		}
		fieldMap["svc_priority"], _ = strconv.Atoi(parts[0])
		fieldMap["target_name"] = parts[1]
		if len(parts) > 2 {
			fieldMap["svc_params"] = parts[2]
		}
		break
	}
}

func resolveDSType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["keytag"], _ = strconv.Atoi(parts[0])
		fieldMap["digest_type"], _ = strconv.Atoi(parts[2])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[1])
		dig := parts[3]
		// digest can have whitespace
		if len(parts) > 4 {
			i := 4
			for i < len(parts) {
				dig += " " + parts[i]
			}
		}
		fieldMap["digest"] = dig
		break
	}
}

func resolveHINFOType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["hardware"] = parts[0]
		fieldMap["software"] = parts[1]
		break
	}
}

func resolveNAPTRType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["order"], _ = strconv.Atoi(parts[0])
		fieldMap["preference"], _ = strconv.Atoi(parts[1])
		fieldMap["flagsnaptr"] = parts[2]
		fieldMap["service"] = parts[3]
		fieldMap["regexp"] = parts[4]
		fieldMap["replacement"] = parts[5]
		break
	}
}

func resolveNSEC3Type(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["flags"], _ = strconv.Atoi(parts[1])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[0])
		fieldMap["iterations"], _ = strconv.Atoi(parts[2])
		fieldMap["salt"] = parts[3]
		fieldMap["next_hashed_owner_name"] = parts[4]
		fieldMap["type_bitmaps"] = parts[5]
		break
	}
}

func resolveNSEC3PARAMType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["flags"], _ = strconv.Atoi(parts[1])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[0])
		fieldMap["iterations"], _ = strconv.Atoi(parts[2])
		fieldMap["salt"] = parts[3]
		break
	}
}

func resolveRPType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["mailbox"] = parts[0]
		fieldMap["txt"] = parts[1]
		break
	}
}

func resolveRRSIGType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["type_covered"] = parts[0]
		fieldMap["algorithm"], _ = strconv.Atoi(parts[1])
		fieldMap["labels"], _ = strconv.Atoi(parts[2])
		fieldMap["original_ttl"], _ = strconv.Atoi(parts[3])
		fieldMap["expiration"] = parts[4]
		fieldMap["inception"] = parts[5]
		fieldMap["signer"] = parts[7]
		fieldMap["keytag"], _ = strconv.Atoi(parts[6])
		sig := parts[8]
		// sig can have whitespace
		if len(parts) > 9 {
			i := 9
			for i < len(parts) {
				sig += " " + parts[i]
			}
		}
		fieldMap["signature"] = sig
		break
	}
}

func resolveSRVType(rdata, newrdata []string, fieldMap map[string]interface{}) {
	// pull out some fields
	parts := strings.Split(rdata[0], " ")
	fieldMap["priority"], _ = strconv.Atoi(parts[0])
	fieldMap["weight"], _ = strconv.Atoi(parts[1])
	fieldMap["port"], _ = strconv.Atoi(parts[2])
	// populate target
	for _, rcontent := range rdata {
		parts = strings.Split(rcontent, " ")
		newrdata = append(newrdata, parts[3])
	}
	fieldMap["target"] = newrdata
}

func resolveSSHFPType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["algorithm"], _ = strconv.Atoi(parts[0])
		fieldMap["fingerprint_type"], _ = strconv.Atoi(parts[1])
		fieldMap["fingerprint"] = parts[2]
		break
	}
}

func resolveSOAType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["name_server"] = parts[0]
		fieldMap["email_address"] = parts[1]
		fieldMap["serial"], _ = strconv.Atoi(parts[2])
		fieldMap["refresh"], _ = strconv.Atoi(parts[3])
		fieldMap["retry"], _ = strconv.Atoi(parts[4])
		fieldMap["expiry"], _ = strconv.Atoi(parts[5])
		fieldMap["nxdomain_ttl"], _ = strconv.Atoi(parts[6])
		break
	}
}

func resolveAKAMAITLCType(rdata []string, fieldMap map[string]interface{}) {
	parts := strings.Split(rdata[0], " ")
	fieldMap["answer_type"] = parts[0]
	fieldMap["dns_name"] = parts[1]
}

func resolveSPFType(rdata, newrdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		newrdata = append(newrdata, rcontent)
	}
	fieldMap["target"] = newrdata
}

func resolveTXTType(rdata, newrdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		newrdata = append(newrdata, rcontent)
	}
	fieldMap["target"] = newrdata
}

func resolveAAAAType(ctx context.Context, p *dns, rdata, newrdata []string, fieldMap map[string]interface{}) {
	for _, i := range rdata {
		str := i
		addr := net.ParseIP(str)
		result := p.FullIPv6(ctx, addr)
		str = result
		newrdata = append(newrdata, str)
	}
	fieldMap["target"] = newrdata
}

func resolveLOCType(ctx context.Context, p *dns, rdata, newrdata []string, fieldMap map[string]interface{}) {
	for _, i := range rdata {
		str := i
		str = p.PadCoordinates(ctx, str)
		newrdata = append(newrdata, str)
	}
	fieldMap["target"] = newrdata
}

func resolveCERTType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		val, err := strconv.Atoi(parts[0])
		if err == nil {
			fieldMap["type_value"] = val
		} else {
			fieldMap["type_mnemonic"] = parts[0]
		}
		fieldMap["keytag"], _ = strconv.Atoi(parts[1])
		fieldMap["algorithm"], _ = strconv.Atoi(parts[2])
		fieldMap["certificate"] = parts[3]
		break
	}
}

func resolveTLSAType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.Split(rcontent, " ")
		fieldMap["usage"], _ = strconv.Atoi(parts[0])
		fieldMap["selector"], _ = strconv.Atoi(parts[1])
		fieldMap["match_type"], _ = strconv.Atoi(parts[2])
		fieldMap["certificate"] = parts[3]
		break
	}
}

func resolveHTTPSType(rdata []string, fieldMap map[string]interface{}) {
	for _, rcontent := range rdata {
		parts := strings.SplitN(rcontent, " ", 3)
		// has to be at least two fields.
		if len(parts) < 2 {
			break
		}
		fieldMap["svc_priority"], _ = strconv.Atoi(parts[0])
		fieldMap["target_name"] = parts[1]
		if len(parts) > 2 {
			fieldMap["svc_params"] = parts[2]
		}
		break
	}
}
