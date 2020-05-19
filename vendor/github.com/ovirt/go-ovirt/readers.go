//
// Copyright (c) 2017 Joey <majunjiev@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package ovirtsdk

import (
	"encoding/xml"
	"io"
	"strconv"
)

func XMLVmSummaryReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VmSummary, error) {
	builder := NewVmSummaryBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vm_summary"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Active(v)
			case "migrating":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Migrating(v)
			case "total":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Total(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVmSummaryReadMany(reader *XMLReader, start *xml.StartElement) (*VmSummarySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VmSummarySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vm_summary":
				one, err := XMLVmSummaryReadOne(reader, &t, "vm_summary")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNetworkFilterReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NetworkFilter, error) {
	builder := NewNetworkFilterBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "network_filter"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "version":
				v, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNetworkFilterReadMany(reader *XMLReader, start *xml.StartElement) (*NetworkFilterSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NetworkFilterSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "network_filter":
				one, err := XMLNetworkFilterReadOne(reader, &t, "network_filter")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLQuotaClusterLimitReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*QuotaClusterLimit, error) {
	builder := NewQuotaClusterLimitBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "quota_cluster_limit"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "memory_limit":
				v, err := reader.ReadFloat64(&t)
				if err != nil {
					return nil, err
				}
				builder.MemoryLimit(v)
			case "memory_usage":
				v, err := reader.ReadFloat64(&t)
				if err != nil {
					return nil, err
				}
				builder.MemoryUsage(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "vcpu_limit":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.VcpuLimit(v)
			case "vcpu_usage":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.VcpuUsage(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLQuotaClusterLimitReadMany(reader *XMLReader, start *xml.StartElement) (*QuotaClusterLimitSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result QuotaClusterLimitSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "quota_cluster_limit":
				one, err := XMLQuotaClusterLimitReadOne(reader, &t, "quota_cluster_limit")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBiosReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Bios, error) {
	builder := NewBiosBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "bios"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "boot_menu":
				v, err := XMLBootMenuReadOne(reader, &t, "boot_menu")
				if err != nil {
					return nil, err
				}
				builder.BootMenu(v)
			case "type":
				vp, err := XMLBiosTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBiosReadMany(reader *XMLReader, start *xml.StartElement) (*BiosSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BiosSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bios":
				one, err := XMLBiosReadOne(reader, &t, "bios")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLQuotaReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Quota, error) {
	builder := NewQuotaBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "quota"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster_hard_limit_pct":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ClusterHardLimitPct(v)
			case "cluster_soft_limit_pct":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ClusterSoftLimitPct(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disks":
				v, err := XMLDiskReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Disks(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "quota_cluster_limits":
				v, err := XMLQuotaClusterLimitReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.QuotaClusterLimits(v)
			case "quota_storage_limits":
				v, err := XMLQuotaStorageLimitReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.QuotaStorageLimits(v)
			case "storage_hard_limit_pct":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.StorageHardLimitPct(v)
			case "storage_soft_limit_pct":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.StorageSoftLimitPct(v)
			case "users":
				v, err := XMLUserReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Users(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "quotaclusterlimits":
			if one.quotaClusterLimits == nil {
				one.quotaClusterLimits = new(QuotaClusterLimitSlice)
			}
			one.quotaClusterLimits.href = link.href
		case "quotastoragelimits":
			if one.quotaStorageLimits == nil {
				one.quotaStorageLimits = new(QuotaStorageLimitSlice)
			}
			one.quotaStorageLimits.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLQuotaReadMany(reader *XMLReader, start *xml.StartElement) (*QuotaSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result QuotaSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "quota":
				one, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMigrationOptionsReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*MigrationOptions, error) {
	builder := NewMigrationOptionsBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "migration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "auto_converge":
				vp, err := XMLInheritableBooleanReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.AutoConverge(v)
			case "bandwidth":
				v, err := XMLMigrationBandwidthReadOne(reader, &t, "bandwidth")
				if err != nil {
					return nil, err
				}
				builder.Bandwidth(v)
			case "compressed":
				vp, err := XMLInheritableBooleanReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Compressed(v)
			case "encrypted":
				vp, err := XMLInheritableBooleanReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Encrypted(v)
			case "policy":
				v, err := XMLMigrationPolicyReadOne(reader, &t, "policy")
				if err != nil {
					return nil, err
				}
				builder.Policy(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMigrationOptionsReadMany(reader *XMLReader, start *xml.StartElement) (*MigrationOptionsSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MigrationOptionsSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "migration":
				one, err := XMLMigrationOptionsReadOne(reader, &t, "migration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRangeReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Range, error) {
	builder := NewRangeBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "range"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "from":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.From(v)
			case "to":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.To(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRangeReadMany(reader *XMLReader, start *xml.StartElement) (*RangeSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RangeSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "range":
				one, err := XMLRangeReadOne(reader, &t, "range")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVirtualNumaNodeReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VirtualNumaNode, error) {
	builder := NewVirtualNumaNodeBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vm_numa_node"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "index":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Index(v)
			case "memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Memory(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "node_distance":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.NodeDistance(v)
			case "numa_node_pins":
				v, err := XMLNumaNodePinReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NumaNodePins(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVirtualNumaNodeReadMany(reader *XMLReader, start *xml.StartElement) (*VirtualNumaNodeSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VirtualNumaNodeSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vm_numa_node":
				one, err := XMLVirtualNumaNodeReadOne(reader, &t, "vm_numa_node")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLUserReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*User, error) {
	builder := NewUserBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "user"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "department":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Department(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "domain":
				v, err := XMLDomainReadOne(reader, &t, "domain")
				if err != nil {
					return nil, err
				}
				builder.Domain(v)
			case "domain_entry_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.DomainEntryId(v)
			case "email":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Email(v)
			case "groups":
				v, err := XMLGroupReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Groups(v)
			case "last_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.LastName(v)
			case "logged_in":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.LoggedIn(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "namespace":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Namespace(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "principal":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Principal(v)
			case "roles":
				v, err := XMLRoleReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Roles(v)
			case "ssh_public_keys":
				v, err := XMLSshPublicKeyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.SshPublicKeys(v)
			case "tags":
				v, err := XMLTagReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Tags(v)
			case "user_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.UserName(v)
			case "user_options":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.UserOptions(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "groups":
			if one.groups == nil {
				one.groups = new(GroupSlice)
			}
			one.groups.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "roles":
			if one.roles == nil {
				one.roles = new(RoleSlice)
			}
			one.roles.href = link.href
		case "sshpublickeys":
			if one.sshPublicKeys == nil {
				one.sshPublicKeys = new(SshPublicKeySlice)
			}
			one.sshPublicKeys.href = link.href
		case "tags":
			if one.tags == nil {
				one.tags = new(TagSlice)
			}
			one.tags.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLUserReadMany(reader *XMLReader, start *xml.StartElement) (*UserSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result UserSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "user":
				one, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSshReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Ssh, error) {
	builder := NewSshBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "ssh"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authentication_method":
				vp, err := XMLSshAuthenticationMethodReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.AuthenticationMethod(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "fingerprint":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Fingerprint(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSshReadMany(reader *XMLReader, start *xml.StartElement) (*SshSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SshSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ssh":
				one, err := XMLSshReadOne(reader, &t, "ssh")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRoleReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Role, error) {
	builder := NewRoleBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "role"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "administrative":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Administrative(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "mutable":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Mutable(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "permits":
				v, err := XMLPermitReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permits(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permits":
			if one.permits == nil {
				one.permits = new(PermitSlice)
			}
			one.permits.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRoleReadMany(reader *XMLReader, start *xml.StartElement) (*RoleSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RoleSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "role":
				one, err := XMLRoleReadOne(reader, &t, "role")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLProductReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Product, error) {
	builder := NewProductBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "product"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLProductReadMany(reader *XMLReader, start *xml.StartElement) (*ProductSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ProductSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "product":
				one, err := XMLProductReadOne(reader, &t, "product")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenStackVolumeTypeReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenStackVolumeType, error) {
	builder := NewOpenStackVolumeTypeBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "open_stack_volume_type"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "openstack_volume_provider":
				v, err := XMLOpenStackVolumeProviderReadOne(reader, &t, "openstack_volume_provider")
				if err != nil {
					return nil, err
				}
				builder.OpenstackVolumeProvider(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenStackVolumeTypeReadMany(reader *XMLReader, start *xml.StartElement) (*OpenStackVolumeTypeSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenStackVolumeTypeSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "open_stack_volume_type":
				one, err := XMLOpenStackVolumeTypeReadOne(reader, &t, "open_stack_volume_type")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBootMenuReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*BootMenu, error) {
	builder := NewBootMenuBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "boot_menu"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBootMenuReadMany(reader *XMLReader, start *xml.StartElement) (*BootMenuSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BootMenuSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "boot_menu":
				one, err := XMLBootMenuReadOne(reader, &t, "boot_menu")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVersionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Version, error) {
	builder := NewVersionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "version"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "build":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Build_(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "full_version":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.FullVersion(v)
			case "major":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Major(v)
			case "minor":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Minor(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "revision":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Revision(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVersionReadMany(reader *XMLReader, start *xml.StartElement) (*VersionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VersionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "version":
				one, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLWatchdogReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Watchdog, error) {
	builder := NewWatchdogBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "watchdog"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "action":
				vp, err := XMLWatchdogActionReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Action(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "model":
				vp, err := XMLWatchdogModelReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Model(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLWatchdogReadMany(reader *XMLReader, start *xml.StartElement) (*WatchdogSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result WatchdogSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "watchdog":
				one, err := XMLWatchdogReadOne(reader, &t, "watchdog")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLStorageDomainReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*StorageDomain, error) {
	builder := NewStorageDomainBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "storage_domain"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "available":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Available(v)
			case "backup":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Backup(v)
			case "block_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.BlockSize(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "committed":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Committed(v)
			case "critical_space_action_blocker":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CriticalSpaceActionBlocker(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "data_centers":
				v, err := XMLDataCenterReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DataCenters(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "discard_after_delete":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DiscardAfterDelete(v)
			case "disk_profiles":
				v, err := XMLDiskProfileReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DiskProfiles(v)
			case "disk_snapshots":
				v, err := XMLDiskSnapshotReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DiskSnapshots(v)
			case "disks":
				v, err := XMLDiskReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Disks(v)
			case "external_status":
				vp, err := XMLExternalStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.ExternalStatus(v)
			case "files":
				v, err := XMLFileReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Files(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "images":
				v, err := XMLImageReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Images(v)
			case "import":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Import(v)
			case "master":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Master(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "status":
				vp, err := XMLStorageDomainStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "storage":
				v, err := XMLHostStorageReadOne(reader, &t, "storage")
				if err != nil {
					return nil, err
				}
				builder.Storage(v)
			case "storage_connections":
				v, err := XMLStorageConnectionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.StorageConnections(v)
			case "storage_format":
				vp, err := XMLStorageFormatReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageFormat(v)
			case "supports_discard":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SupportsDiscard(v)
			case "supports_discard_zeroes_data":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SupportsDiscardZeroesData(v)
			case "templates":
				v, err := XMLTemplateReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Templates(v)
			case "type":
				vp, err := XMLStorageDomainTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "used":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Used(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "warning_low_space_indicator":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.WarningLowSpaceIndicator(v)
			case "wipe_after_delete":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.WipeAfterDelete(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "datacenters":
			if one.dataCenters == nil {
				one.dataCenters = new(DataCenterSlice)
			}
			one.dataCenters.href = link.href
		case "diskprofiles":
			if one.diskProfiles == nil {
				one.diskProfiles = new(DiskProfileSlice)
			}
			one.diskProfiles.href = link.href
		case "disksnapshots":
			if one.diskSnapshots == nil {
				one.diskSnapshots = new(DiskSnapshotSlice)
			}
			one.diskSnapshots.href = link.href
		case "disks":
			if one.disks == nil {
				one.disks = new(DiskSlice)
			}
			one.disks.href = link.href
		case "files":
			if one.files == nil {
				one.files = new(FileSlice)
			}
			one.files.href = link.href
		case "images":
			if one.images == nil {
				one.images = new(ImageSlice)
			}
			one.images.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "storageconnections":
			if one.storageConnections == nil {
				one.storageConnections = new(StorageConnectionSlice)
			}
			one.storageConnections.href = link.href
		case "templates":
			if one.templates == nil {
				one.templates = new(TemplateSlice)
			}
			one.templates.href = link.href
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLStorageDomainReadMany(reader *XMLReader, start *xml.StartElement) (*StorageDomainSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result StorageDomainSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "storage_domain":
				one, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLKatelloErratumReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*KatelloErratum, error) {
	builder := NewKatelloErratumBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "katello_erratum"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "issued":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.Issued(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "packages":
				v, err := XMLPackageReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Packages(v)
			case "severity":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Severity(v)
			case "solution":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Solution(v)
			case "summary":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Summary(v)
			case "title":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Title(v)
			case "type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLKatelloErratumReadMany(reader *XMLReader, start *xml.StartElement) (*KatelloErratumSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result KatelloErratumSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "katello_erratum":
				one, err := XMLKatelloErratumReadOne(reader, &t, "katello_erratum")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMemoryOverCommitReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*MemoryOverCommit, error) {
	builder := NewMemoryOverCommitBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "memory_over_commit"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "percent":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Percent(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMemoryOverCommitReadMany(reader *XMLReader, start *xml.StartElement) (*MemoryOverCommitSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MemoryOverCommitSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "memory_over_commit":
				one, err := XMLMemoryOverCommitReadOne(reader, &t, "memory_over_commit")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMDevTypeReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*MDevType, error) {
	builder := NewMDevTypeBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "m_dev_type"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "available_instances":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.AvailableInstances(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMDevTypeReadMany(reader *XMLReader, start *xml.StartElement) (*MDevTypeSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MDevTypeSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "m_dev_type":
				one, err := XMLMDevTypeReadOne(reader, &t, "m_dev_type")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLTemplateVersionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*TemplateVersion, error) {
	builder := NewTemplateVersionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "template_version"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "base_template":
				v, err := XMLTemplateReadOne(reader, &t, "base_template")
				if err != nil {
					return nil, err
				}
				builder.BaseTemplate(v)
			case "version_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VersionName(v)
			case "version_number":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.VersionNumber(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLTemplateVersionReadMany(reader *XMLReader, start *xml.StartElement) (*TemplateVersionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result TemplateVersionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "template_version":
				one, err := XMLTemplateVersionReadOne(reader, &t, "template_version")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLStorageConnectionExtensionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*StorageConnectionExtension, error) {
	builder := NewStorageConnectionExtensionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "storage_connection_extension"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "target":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Target(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLStorageConnectionExtensionReadMany(reader *XMLReader, start *xml.StartElement) (*StorageConnectionExtensionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result StorageConnectionExtensionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "storage_connection_extension":
				one, err := XMLStorageConnectionExtensionReadOne(reader, &t, "storage_connection_extension")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCoreReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Core, error) {
	builder := NewCoreBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "core"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "index":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Index(v)
			case "socket":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Socket(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCoreReadMany(reader *XMLReader, start *xml.StartElement) (*CoreSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CoreSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "core":
				one, err := XMLCoreReadOne(reader, &t, "core")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVcpuPinReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VcpuPin, error) {
	builder := NewVcpuPinBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vcpu_pin"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cpu_set":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CpuSet(v)
			case "vcpu":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Vcpu(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVcpuPinReadMany(reader *XMLReader, start *xml.StartElement) (*VcpuPinSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VcpuPinSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vcpu_pin":
				one, err := XMLVcpuPinReadOne(reader, &t, "vcpu_pin")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLTemplateReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Template, error) {
	builder := NewTemplateBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "template"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bios":
				v, err := XMLBiosReadOne(reader, &t, "bios")
				if err != nil {
					return nil, err
				}
				builder.Bios(v)
			case "cdroms":
				v, err := XMLCdromReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Cdroms(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "console":
				v, err := XMLConsoleReadOne(reader, &t, "console")
				if err != nil {
					return nil, err
				}
				builder.Console(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "cpu_profile":
				v, err := XMLCpuProfileReadOne(reader, &t, "cpu_profile")
				if err != nil {
					return nil, err
				}
				builder.CpuProfile(v)
			case "cpu_shares":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CpuShares(v)
			case "creation_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.CreationTime(v)
			case "custom_compatibility_version":
				v, err := XMLVersionReadOne(reader, &t, "custom_compatibility_version")
				if err != nil {
					return nil, err
				}
				builder.CustomCompatibilityVersion(v)
			case "custom_cpu_model":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomCpuModel(v)
			case "custom_emulated_machine":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomEmulatedMachine(v)
			case "custom_properties":
				v, err := XMLCustomPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CustomProperties(v)
			case "delete_protected":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DeleteProtected(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk_attachments":
				v, err := XMLDiskAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DiskAttachments(v)
			case "display":
				v, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "domain":
				v, err := XMLDomainReadOne(reader, &t, "domain")
				if err != nil {
					return nil, err
				}
				builder.Domain(v)
			case "graphics_consoles":
				v, err := XMLGraphicsConsoleReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.GraphicsConsoles(v)
			case "high_availability":
				v, err := XMLHighAvailabilityReadOne(reader, &t, "high_availability")
				if err != nil {
					return nil, err
				}
				builder.HighAvailability(v)
			case "initialization":
				v, err := XMLInitializationReadOne(reader, &t, "initialization")
				if err != nil {
					return nil, err
				}
				builder.Initialization(v)
			case "io":
				v, err := XMLIoReadOne(reader, &t, "io")
				if err != nil {
					return nil, err
				}
				builder.Io(v)
			case "large_icon":
				v, err := XMLIconReadOne(reader, &t, "large_icon")
				if err != nil {
					return nil, err
				}
				builder.LargeIcon(v)
			case "lease":
				v, err := XMLStorageDomainLeaseReadOne(reader, &t, "lease")
				if err != nil {
					return nil, err
				}
				builder.Lease(v)
			case "memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Memory(v)
			case "memory_policy":
				v, err := XMLMemoryPolicyReadOne(reader, &t, "memory_policy")
				if err != nil {
					return nil, err
				}
				builder.MemoryPolicy(v)
			case "migration":
				v, err := XMLMigrationOptionsReadOne(reader, &t, "migration")
				if err != nil {
					return nil, err
				}
				builder.Migration(v)
			case "migration_downtime":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MigrationDowntime(v)
			case "multi_queues_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MultiQueuesEnabled(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "nics":
				v, err := XMLNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Nics(v)
			case "origin":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Origin(v)
			case "os":
				v, err := XMLOperatingSystemReadOne(reader, &t, "os")
				if err != nil {
					return nil, err
				}
				builder.Os(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "placement_policy":
				v, err := XMLVmPlacementPolicyReadOne(reader, &t, "placement_policy")
				if err != nil {
					return nil, err
				}
				builder.PlacementPolicy(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "rng_device":
				v, err := XMLRngDeviceReadOne(reader, &t, "rng_device")
				if err != nil {
					return nil, err
				}
				builder.RngDevice(v)
			case "serial_number":
				v, err := XMLSerialNumberReadOne(reader, &t, "serial_number")
				if err != nil {
					return nil, err
				}
				builder.SerialNumber(v)
			case "small_icon":
				v, err := XMLIconReadOne(reader, &t, "small_icon")
				if err != nil {
					return nil, err
				}
				builder.SmallIcon(v)
			case "soundcard_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SoundcardEnabled(v)
			case "sso":
				v, err := XMLSsoReadOne(reader, &t, "sso")
				if err != nil {
					return nil, err
				}
				builder.Sso(v)
			case "start_paused":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.StartPaused(v)
			case "stateless":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Stateless(v)
			case "status":
				vp, err := XMLTemplateStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "storage_error_resume_behaviour":
				vp, err := XMLVmStorageErrorResumeBehaviourReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageErrorResumeBehaviour(v)
			case "tags":
				v, err := XMLTagReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Tags(v)
			case "time_zone":
				v, err := XMLTimeZoneReadOne(reader, &t, "time_zone")
				if err != nil {
					return nil, err
				}
				builder.TimeZone(v)
			case "tunnel_migration":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.TunnelMigration(v)
			case "type":
				vp, err := XMLVmTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "usb":
				v, err := XMLUsbReadOne(reader, &t, "usb")
				if err != nil {
					return nil, err
				}
				builder.Usb(v)
			case "version":
				v, err := XMLTemplateVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "virtio_scsi":
				v, err := XMLVirtioScsiReadOne(reader, &t, "virtio_scsi")
				if err != nil {
					return nil, err
				}
				builder.VirtioScsi(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "watchdogs":
				v, err := XMLWatchdogReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Watchdogs(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "cdroms":
			if one.cdroms == nil {
				one.cdroms = new(CdromSlice)
			}
			one.cdroms.href = link.href
		case "diskattachments":
			if one.diskAttachments == nil {
				one.diskAttachments = new(DiskAttachmentSlice)
			}
			one.diskAttachments.href = link.href
		case "graphicsconsoles":
			if one.graphicsConsoles == nil {
				one.graphicsConsoles = new(GraphicsConsoleSlice)
			}
			one.graphicsConsoles.href = link.href
		case "nics":
			if one.nics == nil {
				one.nics = new(NicSlice)
			}
			one.nics.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "tags":
			if one.tags == nil {
				one.tags = new(TagSlice)
			}
			one.tags.href = link.href
		case "watchdogs":
			if one.watchdogs == nil {
				one.watchdogs = new(WatchdogSlice)
			}
			one.watchdogs.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLTemplateReadMany(reader *XMLReader, start *xml.StartElement) (*TemplateSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result TemplateSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "template":
				one, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNfsProfileDetailReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NfsProfileDetail, error) {
	builder := NewNfsProfileDetailBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "nfs_profile_detail"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "nfs_server_ip":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.NfsServerIp(v)
			case "profile_details":
				v, err := XMLProfileDetailReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ProfileDetails(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNfsProfileDetailReadMany(reader *XMLReader, start *xml.StartElement) (*NfsProfileDetailSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NfsProfileDetailSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "nfs_profile_detail":
				one, err := XMLNfsProfileDetailReadOne(reader, &t, "nfs_profile_detail")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBrickProfileDetailReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*BrickProfileDetail, error) {
	builder := NewBrickProfileDetailBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "brick_profile_detail"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "brick":
				v, err := XMLGlusterBrickReadOne(reader, &t, "brick")
				if err != nil {
					return nil, err
				}
				builder.Brick(v)
			case "profile_details":
				v, err := XMLProfileDetailReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ProfileDetails(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBrickProfileDetailReadMany(reader *XMLReader, start *xml.StartElement) (*BrickProfileDetailSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BrickProfileDetailSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "brick_profile_detail":
				one, err := XMLBrickProfileDetailReadOne(reader, &t, "brick_profile_detail")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLFloppyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Floppy, error) {
	builder := NewFloppyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "floppy"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "file":
				v, err := XMLFileReadOne(reader, &t, "file")
				if err != nil {
					return nil, err
				}
				builder.File(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLFloppyReadMany(reader *XMLReader, start *xml.StartElement) (*FloppySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result FloppySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "floppy":
				one, err := XMLFloppyReadOne(reader, &t, "floppy")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMacPoolReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*MacPool, error) {
	builder := NewMacPoolBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "mac_pool"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "allow_duplicates":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.AllowDuplicates(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "default_pool":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DefaultPool(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "ranges":
				v, err := XMLRangeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Ranges(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMacPoolReadMany(reader *XMLReader, start *xml.StartElement) (*MacPoolSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MacPoolSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "mac_pool":
				one, err := XMLMacPoolReadOne(reader, &t, "mac_pool")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVmPlacementPolicyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VmPlacementPolicy, error) {
	builder := NewVmPlacementPolicyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vm_placement_policy"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "affinity":
				vp, err := XMLVmAffinityReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Affinity(v)
			case "hosts":
				v, err := XMLHostReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Hosts(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "hosts":
			if one.hosts == nil {
				one.hosts = new(HostSlice)
			}
			one.hosts.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVmPlacementPolicyReadMany(reader *XMLReader, start *xml.StartElement) (*VmPlacementPolicySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VmPlacementPolicySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vm_placement_policy":
				one, err := XMLVmPlacementPolicyReadOne(reader, &t, "vm_placement_policy")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLProductInfoReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ProductInfo, error) {
	builder := NewProductInfoBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "product_info"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "vendor":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Vendor(v)
			case "version":
				v, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLProductInfoReadMany(reader *XMLReader, start *xml.StartElement) (*ProductInfoSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ProductInfoSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "product_info":
				one, err := XMLProductInfoReadOne(reader, &t, "product_info")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCpuTuneReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*CpuTune, error) {
	builder := NewCpuTuneBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cpu_tune"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vcpu_pins":
				v, err := XMLVcpuPinReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.VcpuPins(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCpuTuneReadMany(reader *XMLReader, start *xml.StartElement) (*CpuTuneSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CpuTuneSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cpu_tune":
				one, err := XMLCpuTuneReadOne(reader, &t, "cpu_tune")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLAffinityRuleReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*AffinityRule, error) {
	builder := NewAffinityRuleBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "affinity_rule"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "enforcing":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enforcing(v)
			case "positive":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Positive(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLAffinityRuleReadMany(reader *XMLReader, start *xml.StartElement) (*AffinityRuleSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result AffinityRuleSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "affinity_rule":
				one, err := XMLAffinityRuleReadOne(reader, &t, "affinity_rule")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHostNicReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*HostNic, error) {
	builder := NewHostNicBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "host_nic"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ad_aggregator_id":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.AdAggregatorId(v)
			case "base_interface":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.BaseInterface(v)
			case "bonding":
				v, err := XMLBondingReadOne(reader, &t, "bonding")
				if err != nil {
					return nil, err
				}
				builder.Bonding(v)
			case "boot_protocol":
				vp, err := XMLBootProtocolReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.BootProtocol(v)
			case "bridged":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Bridged(v)
			case "check_connectivity":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.CheckConnectivity(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "custom_configuration":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomConfiguration(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "ip":
				v, err := XMLIpReadOne(reader, &t, "ip")
				if err != nil {
					return nil, err
				}
				builder.Ip(v)
			case "ipv6":
				v, err := XMLIpReadOne(reader, &t, "ipv6")
				if err != nil {
					return nil, err
				}
				builder.Ipv6(v)
			case "ipv6_boot_protocol":
				vp, err := XMLBootProtocolReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Ipv6BootProtocol(v)
			case "mac":
				v, err := XMLMacReadOne(reader, &t, "mac")
				if err != nil {
					return nil, err
				}
				builder.Mac(v)
			case "mtu":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Mtu(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "network":
				v, err := XMLNetworkReadOne(reader, &t, "network")
				if err != nil {
					return nil, err
				}
				builder.Network(v)
			case "network_labels":
				v, err := XMLNetworkLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NetworkLabels(v)
			case "override_configuration":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.OverrideConfiguration(v)
			case "physical_function":
				v, err := XMLHostNicReadOne(reader, &t, "physical_function")
				if err != nil {
					return nil, err
				}
				builder.PhysicalFunction(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "qos":
				v, err := XMLQosReadOne(reader, &t, "qos")
				if err != nil {
					return nil, err
				}
				builder.Qos(v)
			case "speed":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Speed(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLNicStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "virtual_functions_configuration":
				v, err := XMLHostNicVirtualFunctionsConfigurationReadOne(reader, &t, "virtual_functions_configuration")
				if err != nil {
					return nil, err
				}
				builder.VirtualFunctionsConfiguration(v)
			case "vlan":
				v, err := XMLVlanReadOne(reader, &t, "vlan")
				if err != nil {
					return nil, err
				}
				builder.Vlan(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "networklabels":
			if one.networkLabels == nil {
				one.networkLabels = new(NetworkLabelSlice)
			}
			one.networkLabels.href = link.href
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHostNicReadMany(reader *XMLReader, start *xml.StartElement) (*HostNicSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HostNicSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "host_nic":
				one, err := XMLHostNicReadOne(reader, &t, "host_nic")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCustomPropertyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*CustomProperty, error) {
	builder := NewCustomPropertyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "custom_property"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "regexp":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Regexp(v)
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCustomPropertyReadMany(reader *XMLReader, start *xml.StartElement) (*CustomPropertySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CustomPropertySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "custom_property":
				one, err := XMLCustomPropertyReadOne(reader, &t, "custom_property")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLPortMirroringReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*PortMirroring, error) {
	builder := NewPortMirroringBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "port_mirroring"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	reader.Skip()
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLPortMirroringReadMany(reader *XMLReader, start *xml.StartElement) (*PortMirroringSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result PortMirroringSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "port_mirroring":
				one, err := XMLPortMirroringReadOne(reader, &t, "port_mirroring")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBookmarkReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Bookmark, error) {
	builder := NewBookmarkBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "bookmark"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBookmarkReadMany(reader *XMLReader, start *xml.StartElement) (*BookmarkSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BookmarkSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bookmark":
				one, err := XMLBookmarkReadOne(reader, &t, "bookmark")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRegistrationAffinityGroupMappingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RegistrationAffinityGroupMapping, error) {
	builder := NewRegistrationAffinityGroupMappingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "registration_affinity_group_mapping"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "from":
				v, err := XMLAffinityGroupReadOne(reader, &t, "from")
				if err != nil {
					return nil, err
				}
				builder.From(v)
			case "to":
				v, err := XMLAffinityGroupReadOne(reader, &t, "to")
				if err != nil {
					return nil, err
				}
				builder.To(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRegistrationAffinityGroupMappingReadMany(reader *XMLReader, start *xml.StartElement) (*RegistrationAffinityGroupMappingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RegistrationAffinityGroupMappingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "registration_affinity_group_mapping":
				one, err := XMLRegistrationAffinityGroupMappingReadOne(reader, &t, "registration_affinity_group_mapping")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLProfileDetailReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ProfileDetail, error) {
	builder := NewProfileDetailBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "profile_detail"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "block_statistics":
				v, err := XMLBlockStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.BlockStatistics(v)
			case "duration":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Duration(v)
			case "fop_statistics":
				v, err := XMLFopStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.FopStatistics(v)
			case "profile_type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ProfileType(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLProfileDetailReadMany(reader *XMLReader, start *xml.StartElement) (*ProfileDetailSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ProfileDetailSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "profile_detail":
				one, err := XMLProfileDetailReadOne(reader, &t, "profile_detail")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLTimeZoneReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*TimeZone, error) {
	builder := NewTimeZoneBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "time_zone"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "utc_offset":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.UtcOffset(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLTimeZoneReadMany(reader *XMLReader, start *xml.StartElement) (*TimeZoneSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result TimeZoneSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "time_zone":
				one, err := XMLTimeZoneReadOne(reader, &t, "time_zone")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLEventReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Event, error) {
	builder := NewEventBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "event"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "code":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Code(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "correlation_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CorrelationId(v)
			case "custom_data":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomData(v)
			case "custom_id":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomId(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "flood_rate":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.FloodRate(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "index":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Index(v)
			case "log_on_host":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.LogOnHost(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "origin":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Origin(v)
			case "severity":
				vp, err := XMLLogSeverityReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Severity(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.Time(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLEventReadMany(reader *XMLReader, start *xml.StartElement) (*EventSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result EventSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "event":
				one, err := XMLEventReadOne(reader, &t, "event")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLFileReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*File, error) {
	builder := NewFileBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "file"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "content":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Content(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLFileReadMany(reader *XMLReader, start *xml.StartElement) (*FileSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result FileSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "file":
				one, err := XMLFileReadOne(reader, &t, "file")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLLogicalUnitReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*LogicalUnit, error) {
	builder := NewLogicalUnitBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "logical_unit"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "discard_max_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.DiscardMaxSize(v)
			case "discard_zeroes_data":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DiscardZeroesData(v)
			case "disk_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.DiskId(v)
			case "lun_mapping":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.LunMapping(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "paths":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Paths(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "portal":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Portal(v)
			case "product_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ProductId(v)
			case "serial":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Serial(v)
			case "size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Size(v)
			case "status":
				vp, err := XMLLunStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "storage_domain_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.StorageDomainId(v)
			case "target":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Target(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "vendor_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VendorId(v)
			case "volume_group_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VolumeGroupId(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLLogicalUnitReadMany(reader *XMLReader, start *xml.StartElement) (*LogicalUnitSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result LogicalUnitSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "logical_unit":
				one, err := XMLLogicalUnitReadOne(reader, &t, "logical_unit")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLApplicationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Application, error) {
	builder := NewApplicationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "application"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLApplicationReadMany(reader *XMLReader, start *xml.StartElement) (*ApplicationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ApplicationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "application":
				one, err := XMLApplicationReadOne(reader, &t, "application")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLImageTransferReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ImageTransfer, error) {
	builder := NewImageTransferBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "image_transfer"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Active(v)
			case "backup":
				v, err := XMLBackupReadOne(reader, &t, "backup")
				if err != nil {
					return nil, err
				}
				builder.Backup(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "direction":
				vp, err := XMLImageTransferDirectionReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Direction(v)
			case "disk":
				v, err := XMLDiskReadOne(reader, &t, "disk")
				if err != nil {
					return nil, err
				}
				builder.Disk(v)
			case "format":
				vp, err := XMLDiskFormatReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Format(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "image":
				v, err := XMLImageReadOne(reader, &t, "image")
				if err != nil {
					return nil, err
				}
				builder.Image(v)
			case "inactivity_timeout":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.InactivityTimeout(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "phase":
				vp, err := XMLImageTransferPhaseReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Phase(v)
			case "proxy_url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ProxyUrl(v)
			case "signed_ticket":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.SignedTicket(v)
			case "snapshot":
				v, err := XMLDiskSnapshotReadOne(reader, &t, "snapshot")
				if err != nil {
					return nil, err
				}
				builder.Snapshot(v)
			case "transfer_url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.TransferUrl(v)
			case "transferred":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Transferred(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLImageTransferReadMany(reader *XMLReader, start *xml.StartElement) (*ImageTransferSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ImageTransferSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "image_transfer":
				one, err := XMLImageTransferReadOne(reader, &t, "image_transfer")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLIpReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Ip, error) {
	builder := NewIpBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "ip"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "gateway":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Gateway(v)
			case "netmask":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Netmask(v)
			case "version":
				vp, err := XMLIpVersionReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLIpReadMany(reader *XMLReader, start *xml.StartElement) (*IpSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result IpSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ip":
				one, err := XMLIpReadOne(reader, &t, "ip")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVnicProfileMappingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VnicProfileMapping, error) {
	builder := NewVnicProfileMappingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vnic_profile_mapping"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "source_network_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.SourceNetworkName(v)
			case "source_network_profile_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.SourceNetworkProfileName(v)
			case "target_vnic_profile":
				v, err := XMLVnicProfileReadOne(reader, &t, "target_vnic_profile")
				if err != nil {
					return nil, err
				}
				builder.TargetVnicProfile(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVnicProfileMappingReadMany(reader *XMLReader, start *xml.StartElement) (*VnicProfileMappingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VnicProfileMappingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vnic_profile_mapping":
				one, err := XMLVnicProfileMappingReadOne(reader, &t, "vnic_profile_mapping")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSystemOptionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SystemOption, error) {
	builder := NewSystemOptionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "system_option"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "values":
				v, err := XMLSystemOptionValueReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Values(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSystemOptionReadMany(reader *XMLReader, start *xml.StartElement) (*SystemOptionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SystemOptionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "system_option":
				one, err := XMLSystemOptionReadOne(reader, &t, "system_option")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterBrickReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterBrick, error) {
	builder := NewGlusterBrickBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "brick"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "brick_dir":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.BrickDir(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "device":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Device(v)
			case "fs_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.FsName(v)
			case "gluster_clients":
				v, err := XMLGlusterClientReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.GlusterClients(v)
			case "gluster_volume":
				v, err := XMLGlusterVolumeReadOne(reader, &t, "gluster_volume")
				if err != nil {
					return nil, err
				}
				builder.GlusterVolume(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "memory_pools":
				v, err := XMLGlusterMemoryPoolReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.MemoryPools(v)
			case "mnt_options":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.MntOptions(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "pid":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Pid(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "server_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ServerId(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLGlusterBrickStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterBrickReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterBrickSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterBrickSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "brick":
				one, err := XMLGlusterBrickReadOne(reader, &t, "brick")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenStackProviderReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenStackProvider, error) {
	builder := NewOpenStackProviderBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "open_stack_provider"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authentication_url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.AuthenticationUrl(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "requires_authentication":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RequiresAuthentication(v)
			case "tenant_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.TenantName(v)
			case "url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Url(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenStackProviderReadMany(reader *XMLReader, start *xml.StartElement) (*OpenStackProviderSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenStackProviderSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "open_stack_provider":
				one, err := XMLOpenStackProviderReadOne(reader, &t, "open_stack_provider")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRegistrationClusterMappingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RegistrationClusterMapping, error) {
	builder := NewRegistrationClusterMappingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "registration_cluster_mapping"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "from":
				v, err := XMLClusterReadOne(reader, &t, "from")
				if err != nil {
					return nil, err
				}
				builder.From(v)
			case "to":
				v, err := XMLClusterReadOne(reader, &t, "to")
				if err != nil {
					return nil, err
				}
				builder.To(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRegistrationClusterMappingReadMany(reader *XMLReader, start *xml.StartElement) (*RegistrationClusterMappingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RegistrationClusterMappingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "registration_cluster_mapping":
				one, err := XMLRegistrationClusterMappingReadOne(reader, &t, "registration_cluster_mapping")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSystemOptionValueReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SystemOptionValue, error) {
	builder := NewSystemOptionValueBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "system_option_value"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "version":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSystemOptionValueReadMany(reader *XMLReader, start *xml.StartElement) (*SystemOptionValueSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SystemOptionValueSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "system_option_value":
				one, err := XMLSystemOptionValueReadOne(reader, &t, "system_option_value")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBondingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Bonding, error) {
	builder := NewBondingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "bonding"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active_slave":
				v, err := XMLHostNicReadOne(reader, &t, "active_slave")
				if err != nil {
					return nil, err
				}
				builder.ActiveSlave(v)
			case "ad_partner_mac":
				v, err := XMLMacReadOne(reader, &t, "ad_partner_mac")
				if err != nil {
					return nil, err
				}
				builder.AdPartnerMac(v)
			case "options":
				v, err := XMLOptionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Options(v)
			case "slaves":
				v, err := XMLHostNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Slaves(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBondingReadMany(reader *XMLReader, start *xml.StartElement) (*BondingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BondingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bonding":
				one, err := XMLBondingReadOne(reader, &t, "bonding")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSeLinuxReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SeLinux, error) {
	builder := NewSeLinuxBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "se_linux"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "mode":
				vp, err := XMLSeLinuxModeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Mode(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSeLinuxReadMany(reader *XMLReader, start *xml.StartElement) (*SeLinuxSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SeLinuxSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "se_linux":
				one, err := XMLSeLinuxReadOne(reader, &t, "se_linux")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLApiReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Api, error) {
	builder := NewApiBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "api"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authenticated_user":
				v, err := XMLUserReadOne(reader, &t, "authenticated_user")
				if err != nil {
					return nil, err
				}
				builder.AuthenticatedUser(v)
			case "effective_user":
				v, err := XMLUserReadOne(reader, &t, "effective_user")
				if err != nil {
					return nil, err
				}
				builder.EffectiveUser(v)
			case "product_info":
				v, err := XMLProductInfoReadOne(reader, &t, "product_info")
				if err != nil {
					return nil, err
				}
				builder.ProductInfo(v)
			case "special_objects":
				v, err := XMLSpecialObjectsReadOne(reader, &t, "special_objects")
				if err != nil {
					return nil, err
				}
				builder.SpecialObjects(v)
			case "summary":
				v, err := XMLApiSummaryReadOne(reader, &t, "summary")
				if err != nil {
					return nil, err
				}
				builder.Summary(v)
			case "time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.Time(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLApiReadMany(reader *XMLReader, start *xml.StartElement) (*ApiSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ApiSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "api":
				one, err := XMLApiReadOne(reader, &t, "api")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMemoryPolicyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*MemoryPolicy, error) {
	builder := NewMemoryPolicyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "memory_policy"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ballooning":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Ballooning(v)
			case "guaranteed":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Guaranteed(v)
			case "max":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Max(v)
			case "over_commit":
				v, err := XMLMemoryOverCommitReadOne(reader, &t, "over_commit")
				if err != nil {
					return nil, err
				}
				builder.OverCommit(v)
			case "transparent_hugepages":
				v, err := XMLTransparentHugePagesReadOne(reader, &t, "transparent_hugepages")
				if err != nil {
					return nil, err
				}
				builder.TransparentHugePages(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMemoryPolicyReadMany(reader *XMLReader, start *xml.StartElement) (*MemoryPolicySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MemoryPolicySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "memory_policy":
				one, err := XMLMemoryPolicyReadOne(reader, &t, "memory_policy")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenStackNetworkReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenStackNetwork, error) {
	builder := NewOpenStackNetworkBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "openstack_network"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "openstack_network_provider":
				v, err := XMLOpenStackNetworkProviderReadOne(reader, &t, "openstack_network_provider")
				if err != nil {
					return nil, err
				}
				builder.OpenstackNetworkProvider(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenStackNetworkReadMany(reader *XMLReader, start *xml.StartElement) (*OpenStackNetworkSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenStackNetworkSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "openstack_network":
				one, err := XMLOpenStackNetworkReadOne(reader, &t, "openstack_network")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVmReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Vm, error) {
	builder := NewVmBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vm"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "affinity_labels":
				v, err := XMLAffinityLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.AffinityLabels(v)
			case "applications":
				v, err := XMLApplicationReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Applications(v)
			case "bios":
				v, err := XMLBiosReadOne(reader, &t, "bios")
				if err != nil {
					return nil, err
				}
				builder.Bios(v)
			case "cdroms":
				v, err := XMLCdromReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Cdroms(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "console":
				v, err := XMLConsoleReadOne(reader, &t, "console")
				if err != nil {
					return nil, err
				}
				builder.Console(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "cpu_profile":
				v, err := XMLCpuProfileReadOne(reader, &t, "cpu_profile")
				if err != nil {
					return nil, err
				}
				builder.CpuProfile(v)
			case "cpu_shares":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CpuShares(v)
			case "creation_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.CreationTime(v)
			case "custom_compatibility_version":
				v, err := XMLVersionReadOne(reader, &t, "custom_compatibility_version")
				if err != nil {
					return nil, err
				}
				builder.CustomCompatibilityVersion(v)
			case "custom_cpu_model":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomCpuModel(v)
			case "custom_emulated_machine":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomEmulatedMachine(v)
			case "custom_properties":
				v, err := XMLCustomPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CustomProperties(v)
			case "delete_protected":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DeleteProtected(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk_attachments":
				v, err := XMLDiskAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DiskAttachments(v)
			case "display":
				v, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "domain":
				v, err := XMLDomainReadOne(reader, &t, "domain")
				if err != nil {
					return nil, err
				}
				builder.Domain(v)
			case "external_host_provider":
				v, err := XMLExternalHostProviderReadOne(reader, &t, "external_host_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalHostProvider(v)
			case "floppies":
				v, err := XMLFloppyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Floppies(v)
			case "fqdn":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Fqdn(v)
			case "graphics_consoles":
				v, err := XMLGraphicsConsoleReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.GraphicsConsoles(v)
			case "guest_operating_system":
				v, err := XMLGuestOperatingSystemReadOne(reader, &t, "guest_operating_system")
				if err != nil {
					return nil, err
				}
				builder.GuestOperatingSystem(v)
			case "guest_time_zone":
				v, err := XMLTimeZoneReadOne(reader, &t, "guest_time_zone")
				if err != nil {
					return nil, err
				}
				builder.GuestTimeZone(v)
			case "has_illegal_images":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.HasIllegalImages(v)
			case "high_availability":
				v, err := XMLHighAvailabilityReadOne(reader, &t, "high_availability")
				if err != nil {
					return nil, err
				}
				builder.HighAvailability(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "host_devices":
				v, err := XMLHostDeviceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.HostDevices(v)
			case "initialization":
				v, err := XMLInitializationReadOne(reader, &t, "initialization")
				if err != nil {
					return nil, err
				}
				builder.Initialization(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "io":
				v, err := XMLIoReadOne(reader, &t, "io")
				if err != nil {
					return nil, err
				}
				builder.Io(v)
			case "katello_errata":
				v, err := XMLKatelloErratumReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.KatelloErrata(v)
			case "large_icon":
				v, err := XMLIconReadOne(reader, &t, "large_icon")
				if err != nil {
					return nil, err
				}
				builder.LargeIcon(v)
			case "lease":
				v, err := XMLStorageDomainLeaseReadOne(reader, &t, "lease")
				if err != nil {
					return nil, err
				}
				builder.Lease(v)
			case "memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Memory(v)
			case "memory_policy":
				v, err := XMLMemoryPolicyReadOne(reader, &t, "memory_policy")
				if err != nil {
					return nil, err
				}
				builder.MemoryPolicy(v)
			case "migration":
				v, err := XMLMigrationOptionsReadOne(reader, &t, "migration")
				if err != nil {
					return nil, err
				}
				builder.Migration(v)
			case "migration_downtime":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MigrationDowntime(v)
			case "multi_queues_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MultiQueuesEnabled(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "next_run_configuration_exists":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.NextRunConfigurationExists(v)
			case "nics":
				v, err := XMLNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Nics(v)
			case "host_numa_nodes":
				v, err := XMLNumaNodeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NumaNodes(v)
			case "numa_tune_mode":
				vp, err := XMLNumaTuneModeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.NumaTuneMode(v)
			case "origin":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Origin(v)
			case "original_template":
				v, err := XMLTemplateReadOne(reader, &t, "original_template")
				if err != nil {
					return nil, err
				}
				builder.OriginalTemplate(v)
			case "os":
				v, err := XMLOperatingSystemReadOne(reader, &t, "os")
				if err != nil {
					return nil, err
				}
				builder.Os(v)
			case "payloads":
				v, err := XMLPayloadReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Payloads(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "placement_policy":
				v, err := XMLVmPlacementPolicyReadOne(reader, &t, "placement_policy")
				if err != nil {
					return nil, err
				}
				builder.PlacementPolicy(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "reported_devices":
				v, err := XMLReportedDeviceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ReportedDevices(v)
			case "rng_device":
				v, err := XMLRngDeviceReadOne(reader, &t, "rng_device")
				if err != nil {
					return nil, err
				}
				builder.RngDevice(v)
			case "run_once":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RunOnce(v)
			case "serial_number":
				v, err := XMLSerialNumberReadOne(reader, &t, "serial_number")
				if err != nil {
					return nil, err
				}
				builder.SerialNumber(v)
			case "sessions":
				v, err := XMLSessionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Sessions(v)
			case "small_icon":
				v, err := XMLIconReadOne(reader, &t, "small_icon")
				if err != nil {
					return nil, err
				}
				builder.SmallIcon(v)
			case "snapshots":
				v, err := XMLSnapshotReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Snapshots(v)
			case "soundcard_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SoundcardEnabled(v)
			case "sso":
				v, err := XMLSsoReadOne(reader, &t, "sso")
				if err != nil {
					return nil, err
				}
				builder.Sso(v)
			case "start_paused":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.StartPaused(v)
			case "start_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.StartTime(v)
			case "stateless":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Stateless(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLVmStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "status_detail":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.StatusDetail(v)
			case "stop_reason":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.StopReason(v)
			case "stop_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.StopTime(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "storage_error_resume_behaviour":
				vp, err := XMLVmStorageErrorResumeBehaviourReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageErrorResumeBehaviour(v)
			case "tags":
				v, err := XMLTagReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Tags(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "time_zone":
				v, err := XMLTimeZoneReadOne(reader, &t, "time_zone")
				if err != nil {
					return nil, err
				}
				builder.TimeZone(v)
			case "tunnel_migration":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.TunnelMigration(v)
			case "type":
				vp, err := XMLVmTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "usb":
				v, err := XMLUsbReadOne(reader, &t, "usb")
				if err != nil {
					return nil, err
				}
				builder.Usb(v)
			case "use_latest_template_version":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UseLatestTemplateVersion(v)
			case "virtio_scsi":
				v, err := XMLVirtioScsiReadOne(reader, &t, "virtio_scsi")
				if err != nil {
					return nil, err
				}
				builder.VirtioScsi(v)
			case "vm_pool":
				v, err := XMLVmPoolReadOne(reader, &t, "vm_pool")
				if err != nil {
					return nil, err
				}
				builder.VmPool(v)
			case "watchdogs":
				v, err := XMLWatchdogReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Watchdogs(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "affinitylabels":
			if one.affinityLabels == nil {
				one.affinityLabels = new(AffinityLabelSlice)
			}
			one.affinityLabels.href = link.href
		case "applications":
			if one.applications == nil {
				one.applications = new(ApplicationSlice)
			}
			one.applications.href = link.href
		case "cdroms":
			if one.cdroms == nil {
				one.cdroms = new(CdromSlice)
			}
			one.cdroms.href = link.href
		case "diskattachments":
			if one.diskAttachments == nil {
				one.diskAttachments = new(DiskAttachmentSlice)
			}
			one.diskAttachments.href = link.href
		case "floppies":
			if one.floppies == nil {
				one.floppies = new(FloppySlice)
			}
			one.floppies.href = link.href
		case "graphicsconsoles":
			if one.graphicsConsoles == nil {
				one.graphicsConsoles = new(GraphicsConsoleSlice)
			}
			one.graphicsConsoles.href = link.href
		case "hostdevices":
			if one.hostDevices == nil {
				one.hostDevices = new(HostDeviceSlice)
			}
			one.hostDevices.href = link.href
		case "katelloerrata":
			if one.katelloErrata == nil {
				one.katelloErrata = new(KatelloErratumSlice)
			}
			one.katelloErrata.href = link.href
		case "nics":
			if one.nics == nil {
				one.nics = new(NicSlice)
			}
			one.nics.href = link.href
		case "numanodes":
			if one.numaNodes == nil {
				one.numaNodes = new(NumaNodeSlice)
			}
			one.numaNodes.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "reporteddevices":
			if one.reportedDevices == nil {
				one.reportedDevices = new(ReportedDeviceSlice)
			}
			one.reportedDevices.href = link.href
		case "sessions":
			if one.sessions == nil {
				one.sessions = new(SessionSlice)
			}
			one.sessions.href = link.href
		case "snapshots":
			if one.snapshots == nil {
				one.snapshots = new(SnapshotSlice)
			}
			one.snapshots.href = link.href
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		case "tags":
			if one.tags == nil {
				one.tags = new(TagSlice)
			}
			one.tags.href = link.href
		case "watchdogs":
			if one.watchdogs == nil {
				one.watchdogs = new(WatchdogSlice)
			}
			one.watchdogs.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVmReadMany(reader *XMLReader, start *xml.StartElement) (*VmSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VmSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vm":
				one, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLExternalHostGroupReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ExternalHostGroup, error) {
	builder := NewExternalHostGroupBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "external_host_group"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "architecture_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ArchitectureName(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "domain_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.DomainName(v)
			case "external_host_provider":
				v, err := XMLExternalHostProviderReadOne(reader, &t, "external_host_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalHostProvider(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "operating_system_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.OperatingSystemName(v)
			case "subnet_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.SubnetName(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLExternalHostGroupReadMany(reader *XMLReader, start *xml.StartElement) (*ExternalHostGroupSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ExternalHostGroupSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "external_host_group":
				one, err := XMLExternalHostGroupReadOne(reader, &t, "external_host_group")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLExternalComputeResourceReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ExternalComputeResource, error) {
	builder := NewExternalComputeResourceBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "external_compute_resource"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "external_host_provider":
				v, err := XMLExternalHostProviderReadOne(reader, &t, "external_host_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalHostProvider(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "provider":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Provider(v)
			case "url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Url(v)
			case "user":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLExternalComputeResourceReadMany(reader *XMLReader, start *xml.StartElement) (*ExternalComputeResourceSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ExternalComputeResourceSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "external_compute_resource":
				one, err := XMLExternalComputeResourceReadOne(reader, &t, "external_compute_resource")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenStackVolumeProviderReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenStackVolumeProvider, error) {
	builder := NewOpenStackVolumeProviderBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "openstack_volume_provider"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authentication_keys":
				v, err := XMLOpenstackVolumeAuthenticationKeyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.AuthenticationKeys(v)
			case "authentication_url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.AuthenticationUrl(v)
			case "certificates":
				v, err := XMLCertificateReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Certificates(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "requires_authentication":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RequiresAuthentication(v)
			case "tenant_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.TenantName(v)
			case "url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Url(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "volume_types":
				v, err := XMLOpenStackVolumeTypeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.VolumeTypes(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "authenticationkeys":
			if one.authenticationKeys == nil {
				one.authenticationKeys = new(OpenstackVolumeAuthenticationKeySlice)
			}
			one.authenticationKeys.href = link.href
		case "certificates":
			if one.certificates == nil {
				one.certificates = new(CertificateSlice)
			}
			one.certificates.href = link.href
		case "volumetypes":
			if one.volumeTypes == nil {
				one.volumeTypes = new(OpenStackVolumeTypeSlice)
			}
			one.volumeTypes.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenStackVolumeProviderReadMany(reader *XMLReader, start *xml.StartElement) (*OpenStackVolumeProviderSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenStackVolumeProviderSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "openstack_volume_provider":
				one, err := XMLOpenStackVolumeProviderReadOne(reader, &t, "openstack_volume_provider")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLStepReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Step, error) {
	builder := NewStepBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "step"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "end_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.EndTime(v)
			case "execution_host":
				v, err := XMLHostReadOne(reader, &t, "execution_host")
				if err != nil {
					return nil, err
				}
				builder.ExecutionHost(v)
			case "external":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.External(v)
			case "external_type":
				vp, err := XMLExternalSystemTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.ExternalType(v)
			case "job":
				v, err := XMLJobReadOne(reader, &t, "job")
				if err != nil {
					return nil, err
				}
				builder.Job(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "number":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Number(v)
			case "parent_step":
				v, err := XMLStepReadOne(reader, &t, "parent_step")
				if err != nil {
					return nil, err
				}
				builder.ParentStep(v)
			case "progress":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Progress(v)
			case "start_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.StartTime(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLStepStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "type":
				vp, err := XMLStepEnumReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLStepReadMany(reader *XMLReader, start *xml.StartElement) (*StepSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result StepSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "step":
				one, err := XMLStepReadOne(reader, &t, "step")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenStackNetworkProviderReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenStackNetworkProvider, error) {
	builder := NewOpenStackNetworkProviderBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "openstack_network_provider"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "agent_configuration":
				v, err := XMLAgentConfigurationReadOne(reader, &t, "agent_configuration")
				if err != nil {
					return nil, err
				}
				builder.AgentConfiguration(v)
			case "authentication_url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.AuthenticationUrl(v)
			case "auto_sync":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.AutoSync(v)
			case "certificates":
				v, err := XMLCertificateReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Certificates(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "external_plugin_type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ExternalPluginType(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "networks":
				v, err := XMLOpenStackNetworkReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Networks(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "plugin_type":
				vp, err := XMLNetworkPluginTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.PluginType(v)
			case "project_domain_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ProjectDomainName(v)
			case "project_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ProjectName(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "read_only":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ReadOnly(v)
			case "requires_authentication":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RequiresAuthentication(v)
			case "subnets":
				v, err := XMLOpenStackSubnetReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Subnets(v)
			case "tenant_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.TenantName(v)
			case "type":
				vp, err := XMLOpenStackNetworkProviderTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "unmanaged":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Unmanaged(v)
			case "url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Url(v)
			case "user_domain_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.UserDomainName(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "certificates":
			if one.certificates == nil {
				one.certificates = new(CertificateSlice)
			}
			one.certificates.href = link.href
		case "networks":
			if one.networks == nil {
				one.networks = new(OpenStackNetworkSlice)
			}
			one.networks.href = link.href
		case "subnets":
			if one.subnets == nil {
				one.subnets = new(OpenStackSubnetSlice)
			}
			one.subnets.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenStackNetworkProviderReadMany(reader *XMLReader, start *xml.StartElement) (*OpenStackNetworkProviderSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenStackNetworkProviderSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "openstack_network_provider":
				one, err := XMLOpenStackNetworkProviderReadOne(reader, &t, "openstack_network_provider")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHookReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Hook, error) {
	builder := NewHookBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "hook"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "event_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.EventName(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "md5":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Md5(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHookReadMany(reader *XMLReader, start *xml.StartElement) (*HookSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HookSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "hook":
				one, err := XMLHookReadOne(reader, &t, "hook")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Configuration, error) {
	builder := NewConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "data":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Data(v)
			case "type":
				vp, err := XMLConfigurationTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*ConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "configuration":
				one, err := XMLConfigurationReadOne(reader, &t, "configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOperatingSystemInfoReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OperatingSystemInfo, error) {
	builder := NewOperatingSystemInfoBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "operating_system"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "architecture":
				vp, err := XMLArchitectureReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Architecture(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "large_icon":
				v, err := XMLIconReadOne(reader, &t, "large_icon")
				if err != nil {
					return nil, err
				}
				builder.LargeIcon(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "small_icon":
				v, err := XMLIconReadOne(reader, &t, "small_icon")
				if err != nil {
					return nil, err
				}
				builder.SmallIcon(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOperatingSystemInfoReadMany(reader *XMLReader, start *xml.StartElement) (*OperatingSystemInfoSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OperatingSystemInfoSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "operating_system":
				one, err := XMLOperatingSystemInfoReadOne(reader, &t, "operating_system")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRegistrationLunMappingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RegistrationLunMapping, error) {
	builder := NewRegistrationLunMappingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "registration_lun_mapping"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "from":
				v, err := XMLDiskReadOne(reader, &t, "from")
				if err != nil {
					return nil, err
				}
				builder.From(v)
			case "to":
				v, err := XMLDiskReadOne(reader, &t, "to")
				if err != nil {
					return nil, err
				}
				builder.To(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRegistrationLunMappingReadMany(reader *XMLReader, start *xml.StartElement) (*RegistrationLunMappingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RegistrationLunMappingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "registration_lun_mapping":
				one, err := XMLRegistrationLunMappingReadOne(reader, &t, "registration_lun_mapping")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDeviceReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Device, error) {
	builder := NewDeviceBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "device"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDeviceReadMany(reader *XMLReader, start *xml.StartElement) (*DeviceSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DeviceSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "device":
				one, err := XMLDeviceReadOne(reader, &t, "device")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNumaNodePinReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NumaNodePin, error) {
	builder := NewNumaNodePinBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "numa_node_pin"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "host_numa_node":
				v, err := XMLNumaNodeReadOne(reader, &t, "host_numa_node")
				if err != nil {
					return nil, err
				}
				builder.HostNumaNode(v)
			case "index":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Index(v)
			case "pinned":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Pinned(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNumaNodePinReadMany(reader *XMLReader, start *xml.StartElement) (*NumaNodePinSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NumaNodePinSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "numa_node_pin":
				one, err := XMLNumaNodePinReadOne(reader, &t, "numa_node_pin")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLReportedDeviceReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ReportedDevice, error) {
	builder := NewReportedDeviceBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "reported_device"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "ips":
				v, err := XMLIpReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Ips(v)
			case "mac":
				v, err := XMLMacReadOne(reader, &t, "mac")
				if err != nil {
					return nil, err
				}
				builder.Mac(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "type":
				vp, err := XMLReportedDeviceTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLReportedDeviceReadMany(reader *XMLReader, start *xml.StartElement) (*ReportedDeviceSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ReportedDeviceSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "reported_device":
				one, err := XMLReportedDeviceReadOne(reader, &t, "reported_device")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHostDevicePassthroughReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*HostDevicePassthrough, error) {
	builder := NewHostDevicePassthroughBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "host_device_passthrough"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHostDevicePassthroughReadMany(reader *XMLReader, start *xml.StartElement) (*HostDevicePassthroughSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HostDevicePassthroughSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "host_device_passthrough":
				one, err := XMLHostDevicePassthroughReadOne(reader, &t, "host_device_passthrough")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLAgentReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Agent, error) {
	builder := NewAgentBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "agent"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "concurrent":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Concurrent(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "encrypt_options":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.EncryptOptions(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "options":
				v, err := XMLOptionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Options(v)
			case "order":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Order(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLAgentReadMany(reader *XMLReader, start *xml.StartElement) (*AgentSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result AgentSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "agent":
				one, err := XMLAgentReadOne(reader, &t, "agent")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSerialNumberReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SerialNumber, error) {
	builder := NewSerialNumberBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "serial_number"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "policy":
				vp, err := XMLSerialNumberPolicyReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Policy(v)
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSerialNumberReadMany(reader *XMLReader, start *xml.StartElement) (*SerialNumberSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SerialNumberSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "serial_number":
				one, err := XMLSerialNumberReadOne(reader, &t, "serial_number")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLUsbReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Usb, error) {
	builder := NewUsbBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "usb"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "type":
				vp, err := XMLUsbTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLUsbReadMany(reader *XMLReader, start *xml.StartElement) (*UsbSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result UsbSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "usb":
				one, err := XMLUsbReadOne(reader, &t, "usb")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRegistrationDomainMappingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RegistrationDomainMapping, error) {
	builder := NewRegistrationDomainMappingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "registration_domain_mapping"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "from":
				v, err := XMLDomainReadOne(reader, &t, "from")
				if err != nil {
					return nil, err
				}
				builder.From(v)
			case "to":
				v, err := XMLDomainReadOne(reader, &t, "to")
				if err != nil {
					return nil, err
				}
				builder.To(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRegistrationDomainMappingReadMany(reader *XMLReader, start *xml.StartElement) (*RegistrationDomainMappingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RegistrationDomainMappingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "registration_domain_mapping":
				one, err := XMLRegistrationDomainMappingReadOne(reader, &t, "registration_domain_mapping")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGroupReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Group, error) {
	builder := NewGroupBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "group"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "domain":
				v, err := XMLDomainReadOne(reader, &t, "domain")
				if err != nil {
					return nil, err
				}
				builder.Domain(v)
			case "domain_entry_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.DomainEntryId(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "namespace":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Namespace(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "roles":
				v, err := XMLRoleReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Roles(v)
			case "tags":
				v, err := XMLTagReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Tags(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "roles":
			if one.roles == nil {
				one.roles = new(RoleSlice)
			}
			one.roles.href = link.href
		case "tags":
			if one.tags == nil {
				one.tags = new(TagSlice)
			}
			one.tags.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGroupReadMany(reader *XMLReader, start *xml.StartElement) (*GroupSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GroupSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "group":
				one, err := XMLGroupReadOne(reader, &t, "group")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRateReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Rate, error) {
	builder := NewRateBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "rate"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bytes":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Bytes(v)
			case "period":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Period(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRateReadMany(reader *XMLReader, start *xml.StartElement) (*RateSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RateSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "rate":
				one, err := XMLRateReadOne(reader, &t, "rate")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLAffinityLabelReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*AffinityLabel, error) {
	builder := NewAffinityLabelBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "affinity_label"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "has_implicit_affinity_group":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.HasImplicitAffinityGroup(v)
			case "hosts":
				v, err := XMLHostReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Hosts(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "read_only":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ReadOnly(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "hosts":
			if one.hosts == nil {
				one.hosts = new(HostSlice)
			}
			one.hosts.href = link.href
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLAffinityLabelReadMany(reader *XMLReader, start *xml.StartElement) (*AffinityLabelSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result AffinityLabelSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "affinity_label":
				one, err := XMLAffinityLabelReadOne(reader, &t, "affinity_label")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHostNicVirtualFunctionsConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*HostNicVirtualFunctionsConfiguration, error) {
	builder := NewHostNicVirtualFunctionsConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "host_nic_virtual_functions_configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "all_networks_allowed":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.AllNetworksAllowed(v)
			case "max_number_of_virtual_functions":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxNumberOfVirtualFunctions(v)
			case "number_of_virtual_functions":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.NumberOfVirtualFunctions(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHostNicVirtualFunctionsConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*HostNicVirtualFunctionsConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HostNicVirtualFunctionsConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "host_nic_virtual_functions_configuration":
				one, err := XMLHostNicVirtualFunctionsConfigurationReadOne(reader, &t, "host_nic_virtual_functions_configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenStackImageReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenStackImage, error) {
	builder := NewOpenStackImageBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "openstack_image"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "openstack_image_provider":
				v, err := XMLOpenStackImageProviderReadOne(reader, &t, "openstack_image_provider")
				if err != nil {
					return nil, err
				}
				builder.OpenstackImageProvider(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenStackImageReadMany(reader *XMLReader, start *xml.StartElement) (*OpenStackImageSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenStackImageSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "openstack_image":
				one, err := XMLOpenStackImageReadOne(reader, &t, "openstack_image")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLKernelReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Kernel, error) {
	builder := NewKernelBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "kernel"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "version":
				v, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLKernelReadMany(reader *XMLReader, start *xml.StartElement) (*KernelSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result KernelSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "kernel":
				one, err := XMLKernelReadOne(reader, &t, "kernel")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLExternalHostProviderReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ExternalHostProvider, error) {
	builder := NewExternalHostProviderBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "external_host_provider"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authentication_url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.AuthenticationUrl(v)
			case "certificates":
				v, err := XMLCertificateReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Certificates(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "compute_resources":
				v, err := XMLExternalComputeResourceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ComputeResources(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "discovered_hosts":
				v, err := XMLExternalDiscoveredHostReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DiscoveredHosts(v)
			case "host_groups":
				v, err := XMLExternalHostGroupReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.HostGroups(v)
			case "hosts":
				v, err := XMLHostReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Hosts(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "requires_authentication":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RequiresAuthentication(v)
			case "url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Url(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "certificates":
			if one.certificates == nil {
				one.certificates = new(CertificateSlice)
			}
			one.certificates.href = link.href
		case "computeresources":
			if one.computeResources == nil {
				one.computeResources = new(ExternalComputeResourceSlice)
			}
			one.computeResources.href = link.href
		case "discoveredhosts":
			if one.discoveredHosts == nil {
				one.discoveredHosts = new(ExternalDiscoveredHostSlice)
			}
			one.discoveredHosts.href = link.href
		case "hostgroups":
			if one.hostGroups == nil {
				one.hostGroups = new(ExternalHostGroupSlice)
			}
			one.hostGroups.href = link.href
		case "hosts":
			if one.hosts == nil {
				one.hosts = new(HostSlice)
			}
			one.hosts.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLExternalHostProviderReadMany(reader *XMLReader, start *xml.StartElement) (*ExternalHostProviderSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ExternalHostProviderSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "external_host_provider":
				one, err := XMLExternalHostProviderReadOne(reader, &t, "external_host_provider")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSsoReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Sso, error) {
	builder := NewSsoBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "sso"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "methods":
				v, err := XMLMethodReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Methods(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSsoReadMany(reader *XMLReader, start *xml.StartElement) (*SsoSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SsoSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "sso":
				one, err := XMLSsoReadOne(reader, &t, "sso")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLIscsiDetailsReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*IscsiDetails, error) {
	builder := NewIscsiDetailsBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "iscsi_details"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "disk_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.DiskId(v)
			case "initiator":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Initiator(v)
			case "lun_mapping":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.LunMapping(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "paths":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Paths(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "portal":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Portal(v)
			case "product_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ProductId(v)
			case "serial":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Serial(v)
			case "size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Size(v)
			case "status":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "storage_domain_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.StorageDomainId(v)
			case "target":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Target(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "vendor_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VendorId(v)
			case "volume_group_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VolumeGroupId(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLIscsiDetailsReadMany(reader *XMLReader, start *xml.StartElement) (*IscsiDetailsSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result IscsiDetailsSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "iscsi_details":
				one, err := XMLIscsiDetailsReadOne(reader, &t, "iscsi_details")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDnsResolverConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*DnsResolverConfiguration, error) {
	builder := NewDnsResolverConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "dns_resolver_configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "name_servers":
				v, err := reader.ReadStrings(&t)
				if err != nil {
					return nil, err
				}
				builder.NameServers(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDnsResolverConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*DnsResolverConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DnsResolverConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "dns_resolver_configuration":
				one, err := XMLDnsResolverConfigurationReadOne(reader, &t, "dns_resolver_configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMethodReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Method, error) {
	builder := NewMethodBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "method"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(SsoMethod(value))
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	reader.Skip()
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMethodReadMany(reader *XMLReader, start *xml.StartElement) (*MethodSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MethodSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "method":
				one, err := XMLMethodReadOne(reader, &t, "method")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDataCenterReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*DataCenter, error) {
	builder := NewDataCenterBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "data_center"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "clusters":
				v, err := XMLClusterReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Clusters(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "iscsi_bonds":
				v, err := XMLIscsiBondReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.IscsiBonds(v)
			case "local":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Local(v)
			case "mac_pool":
				v, err := XMLMacPoolReadOne(reader, &t, "mac_pool")
				if err != nil {
					return nil, err
				}
				builder.MacPool(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "networks":
				v, err := XMLNetworkReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Networks(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "qoss":
				v, err := XMLQosReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Qoss(v)
			case "quota_mode":
				vp, err := XMLQuotaModeTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.QuotaMode(v)
			case "quotas":
				v, err := XMLQuotaReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Quotas(v)
			case "status":
				vp, err := XMLDataCenterStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "storage_domains":
				v, err := XMLStorageDomainReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.StorageDomains(v)
			case "storage_format":
				vp, err := XMLStorageFormatReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageFormat(v)
			case "supported_versions":
				v, err := XMLVersionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.SupportedVersions(v)
			case "version":
				v, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "clusters":
			if one.clusters == nil {
				one.clusters = new(ClusterSlice)
			}
			one.clusters.href = link.href
		case "iscsibonds":
			if one.iscsiBonds == nil {
				one.iscsiBonds = new(IscsiBondSlice)
			}
			one.iscsiBonds.href = link.href
		case "networks":
			if one.networks == nil {
				one.networks = new(NetworkSlice)
			}
			one.networks.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "qoss":
			if one.qoss == nil {
				one.qoss = new(QosSlice)
			}
			one.qoss.href = link.href
		case "quotas":
			if one.quotas == nil {
				one.quotas = new(QuotaSlice)
			}
			one.quotas.href = link.href
		case "storagedomains":
			if one.storageDomains == nil {
				one.storageDomains = new(StorageDomainSlice)
			}
			one.storageDomains.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDataCenterReadMany(reader *XMLReader, start *xml.StartElement) (*DataCenterSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DataCenterSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "data_center":
				one, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRegistrationRoleMappingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RegistrationRoleMapping, error) {
	builder := NewRegistrationRoleMappingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "registration_role_mapping"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "from":
				v, err := XMLRoleReadOne(reader, &t, "from")
				if err != nil {
					return nil, err
				}
				builder.From(v)
			case "to":
				v, err := XMLRoleReadOne(reader, &t, "to")
				if err != nil {
					return nil, err
				}
				builder.To(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRegistrationRoleMappingReadMany(reader *XMLReader, start *xml.StartElement) (*RegistrationRoleMappingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RegistrationRoleMappingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "registration_role_mapping":
				one, err := XMLRegistrationRoleMappingReadOne(reader, &t, "registration_role_mapping")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLFopStatisticReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*FopStatistic, error) {
	builder := NewFopStatisticBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "fop_statistic"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLFopStatisticReadMany(reader *XMLReader, start *xml.StartElement) (*FopStatisticSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result FopStatisticSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "fop_statistic":
				one, err := XMLFopStatisticReadOne(reader, &t, "fop_statistic")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLIdentifiedReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Identified, error) {
	builder := NewIdentifiedBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "identified"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLIdentifiedReadMany(reader *XMLReader, start *xml.StartElement) (*IdentifiedSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result IdentifiedSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "identified":
				one, err := XMLIdentifiedReadOne(reader, &t, "identified")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLEntityProfileDetailReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*EntityProfileDetail, error) {
	builder := NewEntityProfileDetailBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "entity_profile_detail"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "profile_details":
				v, err := XMLProfileDetailReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ProfileDetails(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLEntityProfileDetailReadMany(reader *XMLReader, start *xml.StartElement) (*EntityProfileDetailSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result EntityProfileDetailSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "entity_profile_detail":
				one, err := XMLEntityProfileDetailReadOne(reader, &t, "entity_profile_detail")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLWeightReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Weight, error) {
	builder := NewWeightBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "weight"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "factor":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Factor(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "scheduling_policy":
				v, err := XMLSchedulingPolicyReadOne(reader, &t, "scheduling_policy")
				if err != nil {
					return nil, err
				}
				builder.SchedulingPolicy(v)
			case "scheduling_policy_unit":
				v, err := XMLSchedulingPolicyUnitReadOne(reader, &t, "scheduling_policy_unit")
				if err != nil {
					return nil, err
				}
				builder.SchedulingPolicyUnit(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLWeightReadMany(reader *XMLReader, start *xml.StartElement) (*WeightSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result WeightSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "weight":
				one, err := XMLWeightReadOne(reader, &t, "weight")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterVolumeReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterVolume, error) {
	builder := NewGlusterVolumeBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "gluster_volume"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bricks":
				v, err := XMLGlusterBrickReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Bricks(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disperse_count":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.DisperseCount(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "options":
				v, err := XMLOptionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Options(v)
			case "redundancy_count":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.RedundancyCount(v)
			case "replica_count":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ReplicaCount(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLGlusterVolumeStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "stripe_count":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.StripeCount(v)
			case "transport_types":
				v, err := XMLTransportTypeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.TransportTypes(v)
			case "volume_type":
				vp, err := XMLGlusterVolumeTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.VolumeType(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "bricks":
			if one.bricks == nil {
				one.bricks = new(GlusterBrickSlice)
			}
			one.bricks.href = link.href
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterVolumeReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterVolumeSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterVolumeSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "gluster_volume":
				one, err := XMLGlusterVolumeReadOne(reader, &t, "gluster_volume")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLStorageConnectionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*StorageConnection, error) {
	builder := NewStorageConnectionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "storage_connection"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "gluster_volume":
				v, err := XMLGlusterVolumeReadOne(reader, &t, "gluster_volume")
				if err != nil {
					return nil, err
				}
				builder.GlusterVolume(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "mount_options":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.MountOptions(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "nfs_retrans":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.NfsRetrans(v)
			case "nfs_timeo":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.NfsTimeo(v)
			case "nfs_version":
				vp, err := XMLNfsVersionReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.NfsVersion(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "path":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Path(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "portal":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Portal(v)
			case "target":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Target(v)
			case "type":
				vp, err := XMLStorageTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "vfs_type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VfsType(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLStorageConnectionReadMany(reader *XMLReader, start *xml.StartElement) (*StorageConnectionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result StorageConnectionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "storage_connection":
				one, err := XMLStorageConnectionReadOne(reader, &t, "storage_connection")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNetworkFilterParameterReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NetworkFilterParameter, error) {
	builder := NewNetworkFilterParameterBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "network_filter_parameter"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "nic":
				v, err := XMLNicReadOne(reader, &t, "nic")
				if err != nil {
					return nil, err
				}
				builder.Nic(v)
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNetworkFilterParameterReadMany(reader *XMLReader, start *xml.StartElement) (*NetworkFilterParameterSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NetworkFilterParameterSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "network_filter_parameter":
				one, err := XMLNetworkFilterParameterReadOne(reader, &t, "network_filter_parameter")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenStackSubnetReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenStackSubnet, error) {
	builder := NewOpenStackSubnetBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "openstack_subnet"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cidr":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Cidr(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "dns_servers":
				v, err := reader.ReadStrings(&t)
				if err != nil {
					return nil, err
				}
				builder.DnsServers(v)
			case "gateway":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Gateway(v)
			case "ip_version":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.IpVersion(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "openstack_network":
				v, err := XMLOpenStackNetworkReadOne(reader, &t, "openstack_network")
				if err != nil {
					return nil, err
				}
				builder.OpenstackNetwork(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenStackSubnetReadMany(reader *XMLReader, start *xml.StartElement) (*OpenStackSubnetSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenStackSubnetSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "openstack_subnet":
				one, err := XMLOpenStackSubnetReadOne(reader, &t, "openstack_subnet")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLTagReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Tag, error) {
	builder := NewTagBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "tag"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "group":
				v, err := XMLGroupReadOne(reader, &t, "group")
				if err != nil {
					return nil, err
				}
				builder.Group(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "parent":
				v, err := XMLTagReadOne(reader, &t, "parent")
				if err != nil {
					return nil, err
				}
				builder.Parent(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLTagReadMany(reader *XMLReader, start *xml.StartElement) (*TagSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result TagSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "tag":
				one, err := XMLTagReadOne(reader, &t, "tag")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLAuthorizedKeyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*AuthorizedKey, error) {
	builder := NewAuthorizedKeyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "authorized_key"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "key":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Key(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLAuthorizedKeyReadMany(reader *XMLReader, start *xml.StartElement) (*AuthorizedKeySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result AuthorizedKeySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authorized_key":
				one, err := XMLAuthorizedKeyReadOne(reader, &t, "authorized_key")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVnicProfileReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VnicProfile, error) {
	builder := NewVnicProfileBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vnic_profile"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "custom_properties":
				v, err := XMLCustomPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CustomProperties(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "migratable":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Migratable(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "network":
				v, err := XMLNetworkReadOne(reader, &t, "network")
				if err != nil {
					return nil, err
				}
				builder.Network(v)
			case "network_filter":
				v, err := XMLNetworkFilterReadOne(reader, &t, "network_filter")
				if err != nil {
					return nil, err
				}
				builder.NetworkFilter(v)
			case "pass_through":
				v, err := XMLVnicPassThroughReadOne(reader, &t, "pass_through")
				if err != nil {
					return nil, err
				}
				builder.PassThrough(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "port_mirroring":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.PortMirroring(v)
			case "qos":
				v, err := XMLQosReadOne(reader, &t, "qos")
				if err != nil {
					return nil, err
				}
				builder.Qos(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVnicProfileReadMany(reader *XMLReader, start *xml.StartElement) (*VnicProfileSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VnicProfileSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vnic_profile":
				one, err := XMLVnicProfileReadOne(reader, &t, "vnic_profile")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGraphicsConsoleReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GraphicsConsole, error) {
	builder := NewGraphicsConsoleBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "graphics_console"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "protocol":
				vp, err := XMLGraphicsTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Protocol(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "tls_port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.TlsPort(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGraphicsConsoleReadMany(reader *XMLReader, start *xml.StartElement) (*GraphicsConsoleSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GraphicsConsoleSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "graphics_console":
				one, err := XMLGraphicsConsoleReadOne(reader, &t, "graphics_console")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLIconReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Icon, error) {
	builder := NewIconBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "icon"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "data":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Data(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "media_type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.MediaType(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLIconReadMany(reader *XMLReader, start *xml.StartElement) (*IconSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result IconSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "icon":
				one, err := XMLIconReadOne(reader, &t, "icon")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDiskProfileReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*DiskProfile, error) {
	builder := NewDiskProfileBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "disk_profile"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "qos":
				v, err := XMLQosReadOne(reader, &t, "qos")
				if err != nil {
					return nil, err
				}
				builder.Qos(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDiskProfileReadMany(reader *XMLReader, start *xml.StartElement) (*DiskProfileSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DiskProfileSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "disk_profile":
				one, err := XMLDiskProfileReadOne(reader, &t, "disk_profile")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLImageReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Image, error) {
	builder := NewImageBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "image"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Size(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "type":
				vp, err := XMLImageFileTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLImageReadMany(reader *XMLReader, start *xml.StartElement) (*ImageSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ImageSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "image":
				one, err := XMLImageReadOne(reader, &t, "image")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLEventSubscriptionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*EventSubscription, error) {
	builder := NewEventSubscriptionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "event_subscription"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "event":
				vp, err := XMLNotifiableEventReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Event(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "notification_method":
				vp, err := XMLNotificationMethodReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.NotificationMethod(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLEventSubscriptionReadMany(reader *XMLReader, start *xml.StartElement) (*EventSubscriptionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result EventSubscriptionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "event_subscription":
				one, err := XMLEventSubscriptionReadOne(reader, &t, "event_subscription")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLQuotaStorageLimitReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*QuotaStorageLimit, error) {
	builder := NewQuotaStorageLimitBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "quota_storage_limit"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "limit":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Limit(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "usage":
				v, err := reader.ReadFloat64(&t)
				if err != nil {
					return nil, err
				}
				builder.Usage(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLQuotaStorageLimitReadMany(reader *XMLReader, start *xml.StartElement) (*QuotaStorageLimitSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result QuotaStorageLimitSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "quota_storage_limit":
				one, err := XMLQuotaStorageLimitReadOne(reader, &t, "quota_storage_limit")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOperatingSystemReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OperatingSystem, error) {
	builder := NewOperatingSystemBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "os"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "boot":
				v, err := XMLBootReadOne(reader, &t, "boot")
				if err != nil {
					return nil, err
				}
				builder.Boot(v)
			case "cmdline":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Cmdline(v)
			case "custom_kernel_cmdline":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomKernelCmdline(v)
			case "initrd":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Initrd(v)
			case "kernel":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Kernel(v)
			case "reported_kernel_cmdline":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ReportedKernelCmdline(v)
			case "type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "version":
				v, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOperatingSystemReadMany(reader *XMLReader, start *xml.StartElement) (*OperatingSystemSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OperatingSystemSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "os":
				one, err := XMLOperatingSystemReadOne(reader, &t, "os")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVirtioScsiReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VirtioScsi, error) {
	builder := NewVirtioScsiBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "virtio_scsi"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVirtioScsiReadMany(reader *XMLReader, start *xml.StartElement) (*VirtioScsiSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VirtioScsiSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "virtio_scsi":
				one, err := XMLVirtioScsiReadOne(reader, &t, "virtio_scsi")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVendorReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Vendor, error) {
	builder := NewVendorBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vendor"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVendorReadMany(reader *XMLReader, start *xml.StartElement) (*VendorSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VendorSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vendor":
				one, err := XMLVendorReadOne(reader, &t, "vendor")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLAgentConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*AgentConfiguration, error) {
	builder := NewAgentConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "agent_configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "broker_type":
				vp, err := XMLMessageBrokerTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.BrokerType(v)
			case "network_mappings":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.NetworkMappings(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLAgentConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*AgentConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result AgentConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "agent_configuration":
				one, err := XMLAgentConfigurationReadOne(reader, &t, "agent_configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenStackImageProviderReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenStackImageProvider, error) {
	builder := NewOpenStackImageProviderBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "openstack_image_provider"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authentication_url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.AuthenticationUrl(v)
			case "certificates":
				v, err := XMLCertificateReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Certificates(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "images":
				v, err := XMLOpenStackImageReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Images(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "requires_authentication":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RequiresAuthentication(v)
			case "tenant_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.TenantName(v)
			case "url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Url(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "certificates":
			if one.certificates == nil {
				one.certificates = new(CertificateSlice)
			}
			one.certificates.href = link.href
		case "images":
			if one.images == nil {
				one.images = new(OpenStackImageSlice)
			}
			one.images.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenStackImageProviderReadMany(reader *XMLReader, start *xml.StartElement) (*OpenStackImageProviderSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenStackImageProviderSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "openstack_image_provider":
				one, err := XMLOpenStackImageProviderReadOne(reader, &t, "openstack_image_provider")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMigrationBandwidthReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*MigrationBandwidth, error) {
	builder := NewMigrationBandwidthBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "migration_bandwidth"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "assignment_method":
				vp, err := XMLMigrationBandwidthAssignmentMethodReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.AssignmentMethod(v)
			case "custom_value":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomValue(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMigrationBandwidthReadMany(reader *XMLReader, start *xml.StartElement) (*MigrationBandwidthSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MigrationBandwidthSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "migration_bandwidth":
				one, err := XMLMigrationBandwidthReadOne(reader, &t, "migration_bandwidth")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLUnmanagedNetworkReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*UnmanagedNetwork, error) {
	builder := NewUnmanagedNetworkBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "unmanaged_network"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "host_nic":
				v, err := XMLHostNicReadOne(reader, &t, "host_nic")
				if err != nil {
					return nil, err
				}
				builder.HostNic(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLUnmanagedNetworkReadMany(reader *XMLReader, start *xml.StartElement) (*UnmanagedNetworkSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result UnmanagedNetworkSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "unmanaged_network":
				one, err := XMLUnmanagedNetworkReadOne(reader, &t, "unmanaged_network")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHostStorageReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*HostStorage, error) {
	builder := NewHostStorageBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "host_storage"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "driver_options":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DriverOptions(v)
			case "driver_sensitive_options":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DriverSensitiveOptions(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "logical_units":
				v, err := XMLLogicalUnitReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.LogicalUnits(v)
			case "mount_options":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.MountOptions(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "nfs_retrans":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.NfsRetrans(v)
			case "nfs_timeo":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.NfsTimeo(v)
			case "nfs_version":
				vp, err := XMLNfsVersionReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.NfsVersion(v)
			case "override_luns":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.OverrideLuns(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "path":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Path(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "portal":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Portal(v)
			case "target":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Target(v)
			case "type":
				vp, err := XMLStorageTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "vfs_type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VfsType(v)
			case "volume_group":
				v, err := XMLVolumeGroupReadOne(reader, &t, "volume_group")
				if err != nil {
					return nil, err
				}
				builder.VolumeGroup(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHostStorageReadMany(reader *XMLReader, start *xml.StartElement) (*HostStorageSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HostStorageSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "host_storage":
				one, err := XMLHostStorageReadOne(reader, &t, "host_storage")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLErrorHandlingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ErrorHandling, error) {
	builder := NewErrorHandlingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "error_handling"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "on_error":
				vp, err := XMLMigrateOnErrorReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.OnError(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLErrorHandlingReadMany(reader *XMLReader, start *xml.StartElement) (*ErrorHandlingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ErrorHandlingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "error_handling":
				one, err := XMLErrorHandlingReadOne(reader, &t, "error_handling")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterHookReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterHook, error) {
	builder := NewGlusterHookBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "gluster_hook"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "checksum":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Checksum(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "conflict_status":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ConflictStatus(v)
			case "conflicts":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Conflicts(v)
			case "content":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Content(v)
			case "content_type":
				vp, err := XMLHookContentTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.ContentType(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "gluster_command":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.GlusterCommand(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "server_hooks":
				v, err := XMLGlusterServerHookReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ServerHooks(v)
			case "stage":
				vp, err := XMLHookStageReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Stage(v)
			case "status":
				vp, err := XMLGlusterHookStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "serverhooks":
			if one.serverHooks == nil {
				one.serverHooks = new(GlusterServerHookSlice)
			}
			one.serverHooks.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterHookReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterHookSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterHookSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "gluster_hook":
				one, err := XMLGlusterHookReadOne(reader, &t, "gluster_hook")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVolumeGroupReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VolumeGroup, error) {
	builder := NewVolumeGroupBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "volume_group"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "logical_units":
				v, err := XMLLogicalUnitReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.LogicalUnits(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVolumeGroupReadMany(reader *XMLReader, start *xml.StartElement) (*VolumeGroupSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VolumeGroupSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "volume_group":
				one, err := XMLVolumeGroupReadOne(reader, &t, "volume_group")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDomainReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Domain, error) {
	builder := NewDomainBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "domain"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "groups":
				v, err := XMLGroupReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Groups(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "users":
				v, err := XMLUserReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Users(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "groups":
			if one.groups == nil {
				one.groups = new(GroupSlice)
			}
			one.groups.href = link.href
		case "users":
			if one.users == nil {
				one.users = new(UserSlice)
			}
			one.users.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDomainReadMany(reader *XMLReader, start *xml.StartElement) (*DomainSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DomainSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "domain":
				one, err := XMLDomainReadOne(reader, &t, "domain")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRegistrationConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RegistrationConfiguration, error) {
	builder := NewRegistrationConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "registration_configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "affinity_group_mappings":
				v, err := XMLRegistrationAffinityGroupMappingReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.AffinityGroupMappings(v)
			case "affinity_label_mappings":
				v, err := XMLRegistrationAffinityLabelMappingReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.AffinityLabelMappings(v)
			case "cluster_mappings":
				v, err := XMLRegistrationClusterMappingReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ClusterMappings(v)
			case "domain_mappings":
				v, err := XMLRegistrationDomainMappingReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DomainMappings(v)
			case "lun_mappings":
				v, err := XMLRegistrationLunMappingReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.LunMappings(v)
			case "role_mappings":
				v, err := XMLRegistrationRoleMappingReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.RoleMappings(v)
			case "vnic_profile_mappings":
				v, err := XMLRegistrationVnicProfileMappingReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.VnicProfileMappings(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRegistrationConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*RegistrationConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RegistrationConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "registration_configuration":
				one, err := XMLRegistrationConfigurationReadOne(reader, &t, "registration_configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLInstanceTypeReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*InstanceType, error) {
	builder := NewInstanceTypeBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "instance_type"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bios":
				v, err := XMLBiosReadOne(reader, &t, "bios")
				if err != nil {
					return nil, err
				}
				builder.Bios(v)
			case "cdroms":
				v, err := XMLCdromReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Cdroms(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "console":
				v, err := XMLConsoleReadOne(reader, &t, "console")
				if err != nil {
					return nil, err
				}
				builder.Console(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "cpu_profile":
				v, err := XMLCpuProfileReadOne(reader, &t, "cpu_profile")
				if err != nil {
					return nil, err
				}
				builder.CpuProfile(v)
			case "cpu_shares":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CpuShares(v)
			case "creation_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.CreationTime(v)
			case "custom_compatibility_version":
				v, err := XMLVersionReadOne(reader, &t, "custom_compatibility_version")
				if err != nil {
					return nil, err
				}
				builder.CustomCompatibilityVersion(v)
			case "custom_cpu_model":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomCpuModel(v)
			case "custom_emulated_machine":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomEmulatedMachine(v)
			case "custom_properties":
				v, err := XMLCustomPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CustomProperties(v)
			case "delete_protected":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DeleteProtected(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk_attachments":
				v, err := XMLDiskAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DiskAttachments(v)
			case "display":
				v, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "domain":
				v, err := XMLDomainReadOne(reader, &t, "domain")
				if err != nil {
					return nil, err
				}
				builder.Domain(v)
			case "graphics_consoles":
				v, err := XMLGraphicsConsoleReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.GraphicsConsoles(v)
			case "high_availability":
				v, err := XMLHighAvailabilityReadOne(reader, &t, "high_availability")
				if err != nil {
					return nil, err
				}
				builder.HighAvailability(v)
			case "initialization":
				v, err := XMLInitializationReadOne(reader, &t, "initialization")
				if err != nil {
					return nil, err
				}
				builder.Initialization(v)
			case "io":
				v, err := XMLIoReadOne(reader, &t, "io")
				if err != nil {
					return nil, err
				}
				builder.Io(v)
			case "large_icon":
				v, err := XMLIconReadOne(reader, &t, "large_icon")
				if err != nil {
					return nil, err
				}
				builder.LargeIcon(v)
			case "lease":
				v, err := XMLStorageDomainLeaseReadOne(reader, &t, "lease")
				if err != nil {
					return nil, err
				}
				builder.Lease(v)
			case "memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Memory(v)
			case "memory_policy":
				v, err := XMLMemoryPolicyReadOne(reader, &t, "memory_policy")
				if err != nil {
					return nil, err
				}
				builder.MemoryPolicy(v)
			case "migration":
				v, err := XMLMigrationOptionsReadOne(reader, &t, "migration")
				if err != nil {
					return nil, err
				}
				builder.Migration(v)
			case "migration_downtime":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MigrationDowntime(v)
			case "multi_queues_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MultiQueuesEnabled(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "nics":
				v, err := XMLNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Nics(v)
			case "origin":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Origin(v)
			case "os":
				v, err := XMLOperatingSystemReadOne(reader, &t, "os")
				if err != nil {
					return nil, err
				}
				builder.Os(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "placement_policy":
				v, err := XMLVmPlacementPolicyReadOne(reader, &t, "placement_policy")
				if err != nil {
					return nil, err
				}
				builder.PlacementPolicy(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "rng_device":
				v, err := XMLRngDeviceReadOne(reader, &t, "rng_device")
				if err != nil {
					return nil, err
				}
				builder.RngDevice(v)
			case "serial_number":
				v, err := XMLSerialNumberReadOne(reader, &t, "serial_number")
				if err != nil {
					return nil, err
				}
				builder.SerialNumber(v)
			case "small_icon":
				v, err := XMLIconReadOne(reader, &t, "small_icon")
				if err != nil {
					return nil, err
				}
				builder.SmallIcon(v)
			case "soundcard_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SoundcardEnabled(v)
			case "sso":
				v, err := XMLSsoReadOne(reader, &t, "sso")
				if err != nil {
					return nil, err
				}
				builder.Sso(v)
			case "start_paused":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.StartPaused(v)
			case "stateless":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Stateless(v)
			case "status":
				vp, err := XMLTemplateStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "storage_error_resume_behaviour":
				vp, err := XMLVmStorageErrorResumeBehaviourReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageErrorResumeBehaviour(v)
			case "tags":
				v, err := XMLTagReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Tags(v)
			case "time_zone":
				v, err := XMLTimeZoneReadOne(reader, &t, "time_zone")
				if err != nil {
					return nil, err
				}
				builder.TimeZone(v)
			case "tunnel_migration":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.TunnelMigration(v)
			case "type":
				vp, err := XMLVmTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "usb":
				v, err := XMLUsbReadOne(reader, &t, "usb")
				if err != nil {
					return nil, err
				}
				builder.Usb(v)
			case "version":
				v, err := XMLTemplateVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "virtio_scsi":
				v, err := XMLVirtioScsiReadOne(reader, &t, "virtio_scsi")
				if err != nil {
					return nil, err
				}
				builder.VirtioScsi(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "watchdogs":
				v, err := XMLWatchdogReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Watchdogs(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "cdroms":
			if one.cdroms == nil {
				one.cdroms = new(CdromSlice)
			}
			one.cdroms.href = link.href
		case "diskattachments":
			if one.diskAttachments == nil {
				one.diskAttachments = new(DiskAttachmentSlice)
			}
			one.diskAttachments.href = link.href
		case "graphicsconsoles":
			if one.graphicsConsoles == nil {
				one.graphicsConsoles = new(GraphicsConsoleSlice)
			}
			one.graphicsConsoles.href = link.href
		case "nics":
			if one.nics == nil {
				one.nics = new(NicSlice)
			}
			one.nics.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "tags":
			if one.tags == nil {
				one.tags = new(TagSlice)
			}
			one.tags.href = link.href
		case "watchdogs":
			if one.watchdogs == nil {
				one.watchdogs = new(WatchdogSlice)
			}
			one.watchdogs.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLInstanceTypeReadMany(reader *XMLReader, start *xml.StartElement) (*InstanceTypeSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result InstanceTypeSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "instance_type":
				one, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOptionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Option, error) {
	builder := NewOptionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "option"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOptionReadMany(reader *XMLReader, start *xml.StartElement) (*OptionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OptionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "option":
				one, err := XMLOptionReadOne(reader, &t, "option")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLPropertyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Property, error) {
	builder := NewPropertyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "property"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLPropertyReadMany(reader *XMLReader, start *xml.StartElement) (*PropertySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result PropertySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "property":
				one, err := XMLPropertyReadOne(reader, &t, "property")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLPermissionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Permission, error) {
	builder := NewPermissionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "permission"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk":
				v, err := XMLDiskReadOne(reader, &t, "disk")
				if err != nil {
					return nil, err
				}
				builder.Disk(v)
			case "group":
				v, err := XMLGroupReadOne(reader, &t, "group")
				if err != nil {
					return nil, err
				}
				builder.Group(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "role":
				v, err := XMLRoleReadOne(reader, &t, "role")
				if err != nil {
					return nil, err
				}
				builder.Role(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vm_pool":
				v, err := XMLVmPoolReadOne(reader, &t, "vm_pool")
				if err != nil {
					return nil, err
				}
				builder.VmPool(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLPermissionReadMany(reader *XMLReader, start *xml.StartElement) (*PermissionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result PermissionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "permission":
				one, err := XMLPermissionReadOne(reader, &t, "permission")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRngDeviceReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RngDevice, error) {
	builder := NewRngDeviceBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "rng_device"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "rate":
				v, err := XMLRateReadOne(reader, &t, "rate")
				if err != nil {
					return nil, err
				}
				builder.Rate(v)
			case "source":
				vp, err := XMLRngSourceReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Source(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRngDeviceReadMany(reader *XMLReader, start *xml.StartElement) (*RngDeviceSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RngDeviceSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "rng_device":
				one, err := XMLRngDeviceReadOne(reader, &t, "rng_device")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLTicketReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Ticket, error) {
	builder := NewTicketBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "ticket"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "expiry":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Expiry(v)
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLTicketReadMany(reader *XMLReader, start *xml.StartElement) (*TicketSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result TicketSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ticket":
				one, err := XMLTicketReadOne(reader, &t, "ticket")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLJobReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Job, error) {
	builder := NewJobBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "job"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "auto_cleared":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.AutoCleared(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "end_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.EndTime(v)
			case "external":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.External(v)
			case "last_updated":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.LastUpdated(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "owner":
				v, err := XMLUserReadOne(reader, &t, "owner")
				if err != nil {
					return nil, err
				}
				builder.Owner(v)
			case "start_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.StartTime(v)
			case "status":
				vp, err := XMLJobStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "steps":
				v, err := XMLStepReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Steps(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "steps":
			if one.steps == nil {
				one.steps = new(StepSlice)
			}
			one.steps.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLJobReadMany(reader *XMLReader, start *xml.StartElement) (*JobSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result JobSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "job":
				one, err := XMLJobReadOne(reader, &t, "job")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLClusterFeatureReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ClusterFeature, error) {
	builder := NewClusterFeatureBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cluster_feature"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster_level":
				v, err := XMLClusterLevelReadOne(reader, &t, "cluster_level")
				if err != nil {
					return nil, err
				}
				builder.ClusterLevel(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLClusterFeatureReadMany(reader *XMLReader, start *xml.StartElement) (*ClusterFeatureSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ClusterFeatureSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster_feature":
				one, err := XMLClusterFeatureReadOne(reader, &t, "cluster_feature")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBootReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Boot, error) {
	builder := NewBootBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "boot"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "devices":
				v, err := XMLBootDeviceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Devices(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBootReadMany(reader *XMLReader, start *xml.StartElement) (*BootSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BootSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "boot":
				one, err := XMLBootReadOne(reader, &t, "boot")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLTransparentHugePagesReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*TransparentHugePages, error) {
	builder := NewTransparentHugePagesBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "transparent_hugepages"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLTransparentHugePagesReadMany(reader *XMLReader, start *xml.StartElement) (*TransparentHugePagesSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result TransparentHugePagesSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "transparent_hugepages":
				one, err := XMLTransparentHugePagesReadOne(reader, &t, "transparent_hugepages")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLExternalNetworkProviderConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ExternalNetworkProviderConfiguration, error) {
	builder := NewExternalNetworkProviderConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "external_network_provider_configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "external_network_provider":
				v, err := XMLExternalProviderReadOne(reader, &t, "external_network_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalNetworkProvider(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLExternalNetworkProviderConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*ExternalNetworkProviderConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ExternalNetworkProviderConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "external_network_provider_configuration":
				one, err := XMLExternalNetworkProviderConfigurationReadOne(reader, &t, "external_network_provider_configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVlanReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Vlan, error) {
	builder := NewVlanBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vlan"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			builder.Id(v)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	reader.Skip()
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVlanReadMany(reader *XMLReader, start *xml.StartElement) (*VlanSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VlanSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vlan":
				one, err := XMLVlanReadOne(reader, &t, "vlan")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterVolumeProfileDetailsReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterVolumeProfileDetails, error) {
	builder := NewGlusterVolumeProfileDetailsBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "gluster_volume_profile_details"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "brick_profile_details":
				v, err := XMLBrickProfileDetailReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.BrickProfileDetails(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "nfs_profile_details":
				v, err := XMLNfsProfileDetailReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NfsProfileDetails(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterVolumeProfileDetailsReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterVolumeProfileDetailsSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterVolumeProfileDetailsSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "gluster_volume_profile_details":
				one, err := XMLGlusterVolumeProfileDetailsReadOne(reader, &t, "gluster_volume_profile_details")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRegistrationVnicProfileMappingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RegistrationVnicProfileMapping, error) {
	builder := NewRegistrationVnicProfileMappingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "registration_vnic_profile_mapping"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "from":
				v, err := XMLVnicProfileReadOne(reader, &t, "from")
				if err != nil {
					return nil, err
				}
				builder.From(v)
			case "to":
				v, err := XMLVnicProfileReadOne(reader, &t, "to")
				if err != nil {
					return nil, err
				}
				builder.To(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRegistrationVnicProfileMappingReadMany(reader *XMLReader, start *xml.StartElement) (*RegistrationVnicProfileMappingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RegistrationVnicProfileMappingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "registration_vnic_profile_mapping":
				one, err := XMLRegistrationVnicProfileMappingReadOne(reader, &t, "registration_vnic_profile_mapping")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNetworkReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Network, error) {
	builder := NewNetworkBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "network"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "display":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "dns_resolver_configuration":
				v, err := XMLDnsResolverConfigurationReadOne(reader, &t, "dns_resolver_configuration")
				if err != nil {
					return nil, err
				}
				builder.DnsResolverConfiguration(v)
			case "external_provider":
				v, err := XMLOpenStackNetworkProviderReadOne(reader, &t, "external_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalProvider(v)
			case "external_provider_physical_network":
				v, err := XMLNetworkReadOne(reader, &t, "external_provider_physical_network")
				if err != nil {
					return nil, err
				}
				builder.ExternalProviderPhysicalNetwork(v)
			case "ip":
				v, err := XMLIpReadOne(reader, &t, "ip")
				if err != nil {
					return nil, err
				}
				builder.Ip(v)
			case "mtu":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Mtu(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "network_labels":
				v, err := XMLNetworkLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NetworkLabels(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "profile_required":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ProfileRequired(v)
			case "qos":
				v, err := XMLQosReadOne(reader, &t, "qos")
				if err != nil {
					return nil, err
				}
				builder.Qos(v)
			case "required":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Required(v)
			case "status":
				vp, err := XMLNetworkStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "stp":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Stp(v)
			case "usages":
				v, err := XMLNetworkUsageReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Usages(v)
			case "vdsm_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VdsmName(v)
			case "vlan":
				v, err := XMLVlanReadOne(reader, &t, "vlan")
				if err != nil {
					return nil, err
				}
				builder.Vlan(v)
			case "vnic_profiles":
				v, err := XMLVnicProfileReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.VnicProfiles(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "networklabels":
			if one.networkLabels == nil {
				one.networkLabels = new(NetworkLabelSlice)
			}
			one.networkLabels.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "vnicprofiles":
			if one.vnicProfiles == nil {
				one.vnicProfiles = new(VnicProfileSlice)
			}
			one.vnicProfiles.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNetworkReadMany(reader *XMLReader, start *xml.StartElement) (*NetworkSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NetworkSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "network":
				one, err := XMLNetworkReadOne(reader, &t, "network")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLClusterLevelReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ClusterLevel, error) {
	builder := NewClusterLevelBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cluster_level"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster_features":
				v, err := XMLClusterFeatureReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ClusterFeatures(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "cpu_types":
				v, err := XMLCpuTypeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CpuTypes(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "permits":
				v, err := XMLPermitReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permits(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "clusterfeatures":
			if one.clusterFeatures == nil {
				one.clusterFeatures = new(ClusterFeatureSlice)
			}
			one.clusterFeatures.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLClusterLevelReadMany(reader *XMLReader, start *xml.StartElement) (*ClusterLevelSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ClusterLevelSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster_level":
				one, err := XMLClusterLevelReadOne(reader, &t, "cluster_level")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLPowerManagementReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*PowerManagement, error) {
	builder := NewPowerManagementBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "power_management"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "agents":
				v, err := XMLAgentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Agents(v)
			case "automatic_pm_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.AutomaticPmEnabled(v)
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "kdump_detection":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.KdumpDetection(v)
			case "options":
				v, err := XMLOptionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Options(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "pm_proxies":
				v, err := XMLPmProxyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.PmProxies(v)
			case "status":
				vp, err := XMLPowerManagementStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLPowerManagementReadMany(reader *XMLReader, start *xml.StartElement) (*PowerManagementSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result PowerManagementSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "power_management":
				one, err := XMLPowerManagementReadOne(reader, &t, "power_management")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLAffinityGroupReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*AffinityGroup, error) {
	builder := NewAffinityGroupBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "affinity_group"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "enforcing":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enforcing(v)
			case "host_labels":
				v, err := XMLAffinityLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.HostLabels(v)
			case "hosts":
				v, err := XMLHostReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Hosts(v)
			case "hosts_rule":
				v, err := XMLAffinityRuleReadOne(reader, &t, "hosts_rule")
				if err != nil {
					return nil, err
				}
				builder.HostsRule(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "positive":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Positive(v)
			case "priority":
				v, err := reader.ReadFloat64(&t)
				if err != nil {
					return nil, err
				}
				builder.Priority(v)
			case "vm_labels":
				v, err := XMLAffinityLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.VmLabels(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "vms_rule":
				v, err := XMLAffinityRuleReadOne(reader, &t, "vms_rule")
				if err != nil {
					return nil, err
				}
				builder.VmsRule(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "hostlabels":
			if one.hostLabels == nil {
				one.hostLabels = new(AffinityLabelSlice)
			}
			one.hostLabels.href = link.href
		case "hosts":
			if one.hosts == nil {
				one.hosts = new(HostSlice)
			}
			one.hosts.href = link.href
		case "vmlabels":
			if one.vmLabels == nil {
				one.vmLabels = new(AffinityLabelSlice)
			}
			one.vmLabels.href = link.href
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLAffinityGroupReadMany(reader *XMLReader, start *xml.StartElement) (*AffinityGroupSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result AffinityGroupSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "affinity_group":
				one, err := XMLAffinityGroupReadOne(reader, &t, "affinity_group")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHardwareInformationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*HardwareInformation, error) {
	builder := NewHardwareInformationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "hardware_information"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "family":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Family(v)
			case "manufacturer":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Manufacturer(v)
			case "product_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ProductName(v)
			case "serial_number":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.SerialNumber(v)
			case "supported_rng_sources":
				v, err := XMLRngSourceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.SupportedRngSources(v)
			case "uuid":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Uuid(v)
			case "version":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHardwareInformationReadMany(reader *XMLReader, start *xml.StartElement) (*HardwareInformationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HardwareInformationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "hardware_information":
				one, err := XMLHardwareInformationReadOne(reader, &t, "hardware_information")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBalanceReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Balance, error) {
	builder := NewBalanceBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "balance"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "scheduling_policy":
				v, err := XMLSchedulingPolicyReadOne(reader, &t, "scheduling_policy")
				if err != nil {
					return nil, err
				}
				builder.SchedulingPolicy(v)
			case "scheduling_policy_unit":
				v, err := XMLSchedulingPolicyUnitReadOne(reader, &t, "scheduling_policy_unit")
				if err != nil {
					return nil, err
				}
				builder.SchedulingPolicyUnit(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBalanceReadMany(reader *XMLReader, start *xml.StartElement) (*BalanceSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BalanceSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "balance":
				one, err := XMLBalanceReadOne(reader, &t, "balance")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLIoReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Io, error) {
	builder := NewIoBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "io"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "threads":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Threads(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLIoReadMany(reader *XMLReader, start *xml.StartElement) (*IoSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result IoSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "io":
				one, err := XMLIoReadOne(reader, &t, "io")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterClientReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterClient, error) {
	builder := NewGlusterClientBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "gluster_client"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bytes_read":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.BytesRead(v)
			case "bytes_written":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.BytesWritten(v)
			case "client_port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ClientPort(v)
			case "host_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.HostName(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterClientReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterClientSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterClientSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "gluster_client":
				one, err := XMLGlusterClientReadOne(reader, &t, "gluster_client")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterMemoryPoolReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterMemoryPool, error) {
	builder := NewGlusterMemoryPoolBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "memory_pool"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "alloc_count":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.AllocCount(v)
			case "cold_count":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ColdCount(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "hot_count":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.HotCount(v)
			case "max_alloc":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxAlloc(v)
			case "max_stdalloc":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxStdalloc(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "padded_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.PaddedSize(v)
			case "pool_misses":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.PoolMisses(v)
			case "type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterMemoryPoolReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterMemoryPoolSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterMemoryPoolSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "memory_pool":
				one, err := XMLGlusterMemoryPoolReadOne(reader, &t, "memory_pool")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLInitializationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Initialization, error) {
	builder := NewInitializationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "initialization"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active_directory_ou":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ActiveDirectoryOu(v)
			case "authorized_ssh_keys":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.AuthorizedSshKeys(v)
			case "cloud_init":
				v, err := XMLCloudInitReadOne(reader, &t, "cloud_init")
				if err != nil {
					return nil, err
				}
				builder.CloudInit(v)
			case "cloud_init_network_protocol":
				vp, err := XMLCloudInitNetworkProtocolReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.CloudInitNetworkProtocol(v)
			case "configuration":
				v, err := XMLConfigurationReadOne(reader, &t, "configuration")
				if err != nil {
					return nil, err
				}
				builder.Configuration(v)
			case "custom_script":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomScript(v)
			case "dns_search":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.DnsSearch(v)
			case "dns_servers":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.DnsServers(v)
			case "domain":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Domain(v)
			case "host_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.HostName(v)
			case "input_locale":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.InputLocale(v)
			case "nic_configurations":
				v, err := XMLNicConfigurationReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NicConfigurations(v)
			case "org_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.OrgName(v)
			case "regenerate_ids":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RegenerateIds(v)
			case "regenerate_ssh_keys":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RegenerateSshKeys(v)
			case "root_password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.RootPassword(v)
			case "system_locale":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.SystemLocale(v)
			case "timezone":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Timezone(v)
			case "ui_language":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.UiLanguage(v)
			case "user_locale":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.UserLocale(v)
			case "user_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.UserName(v)
			case "windows_license_key":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.WindowsLicenseKey(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLInitializationReadMany(reader *XMLReader, start *xml.StartElement) (*InitializationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result InitializationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "initialization":
				one, err := XMLInitializationReadOne(reader, &t, "initialization")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNumaNodeReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NumaNode, error) {
	builder := NewNumaNodeBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "host_numa_node"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "index":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Index(v)
			case "memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Memory(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "node_distance":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.NodeDistance(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNumaNodeReadMany(reader *XMLReader, start *xml.StartElement) (*NumaNodeSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NumaNodeSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "host_numa_node":
				one, err := XMLNumaNodeReadOne(reader, &t, "host_numa_node")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMacReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Mac, error) {
	builder := NewMacBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "mac"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMacReadMany(reader *XMLReader, start *xml.StartElement) (*MacSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MacSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "mac":
				one, err := XMLMacReadOne(reader, &t, "mac")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNetworkConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NetworkConfiguration, error) {
	builder := NewNetworkConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "network_configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "dns":
				v, err := XMLDnsReadOne(reader, &t, "dns")
				if err != nil {
					return nil, err
				}
				builder.Dns(v)
			case "nics":
				v, err := XMLNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Nics(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNetworkConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*NetworkConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NetworkConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "network_configuration":
				one, err := XMLNetworkConfigurationReadOne(reader, &t, "network_configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLExternalDiscoveredHostReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ExternalDiscoveredHost, error) {
	builder := NewExternalDiscoveredHostBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "external_discovered_host"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "external_host_provider":
				v, err := XMLExternalHostProviderReadOne(reader, &t, "external_host_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalHostProvider(v)
			case "ip":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Ip(v)
			case "last_report":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.LastReport(v)
			case "mac":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Mac(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "subnet_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.SubnetName(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLExternalDiscoveredHostReadMany(reader *XMLReader, start *xml.StartElement) (*ExternalDiscoveredHostSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ExternalDiscoveredHostSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "external_discovered_host":
				one, err := XMLExternalDiscoveredHostReadOne(reader, &t, "external_discovered_host")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVmBaseReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VmBase, error) {
	builder := NewVmBaseBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vm_base"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "bios":
				v, err := XMLBiosReadOne(reader, &t, "bios")
				if err != nil {
					return nil, err
				}
				builder.Bios(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "console":
				v, err := XMLConsoleReadOne(reader, &t, "console")
				if err != nil {
					return nil, err
				}
				builder.Console(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "cpu_profile":
				v, err := XMLCpuProfileReadOne(reader, &t, "cpu_profile")
				if err != nil {
					return nil, err
				}
				builder.CpuProfile(v)
			case "cpu_shares":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CpuShares(v)
			case "creation_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.CreationTime(v)
			case "custom_compatibility_version":
				v, err := XMLVersionReadOne(reader, &t, "custom_compatibility_version")
				if err != nil {
					return nil, err
				}
				builder.CustomCompatibilityVersion(v)
			case "custom_cpu_model":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomCpuModel(v)
			case "custom_emulated_machine":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomEmulatedMachine(v)
			case "custom_properties":
				v, err := XMLCustomPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CustomProperties(v)
			case "delete_protected":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DeleteProtected(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "display":
				v, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "domain":
				v, err := XMLDomainReadOne(reader, &t, "domain")
				if err != nil {
					return nil, err
				}
				builder.Domain(v)
			case "high_availability":
				v, err := XMLHighAvailabilityReadOne(reader, &t, "high_availability")
				if err != nil {
					return nil, err
				}
				builder.HighAvailability(v)
			case "initialization":
				v, err := XMLInitializationReadOne(reader, &t, "initialization")
				if err != nil {
					return nil, err
				}
				builder.Initialization(v)
			case "io":
				v, err := XMLIoReadOne(reader, &t, "io")
				if err != nil {
					return nil, err
				}
				builder.Io(v)
			case "large_icon":
				v, err := XMLIconReadOne(reader, &t, "large_icon")
				if err != nil {
					return nil, err
				}
				builder.LargeIcon(v)
			case "lease":
				v, err := XMLStorageDomainLeaseReadOne(reader, &t, "lease")
				if err != nil {
					return nil, err
				}
				builder.Lease(v)
			case "memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Memory(v)
			case "memory_policy":
				v, err := XMLMemoryPolicyReadOne(reader, &t, "memory_policy")
				if err != nil {
					return nil, err
				}
				builder.MemoryPolicy(v)
			case "migration":
				v, err := XMLMigrationOptionsReadOne(reader, &t, "migration")
				if err != nil {
					return nil, err
				}
				builder.Migration(v)
			case "migration_downtime":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MigrationDowntime(v)
			case "multi_queues_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MultiQueuesEnabled(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "origin":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Origin(v)
			case "os":
				v, err := XMLOperatingSystemReadOne(reader, &t, "os")
				if err != nil {
					return nil, err
				}
				builder.Os(v)
			case "placement_policy":
				v, err := XMLVmPlacementPolicyReadOne(reader, &t, "placement_policy")
				if err != nil {
					return nil, err
				}
				builder.PlacementPolicy(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "rng_device":
				v, err := XMLRngDeviceReadOne(reader, &t, "rng_device")
				if err != nil {
					return nil, err
				}
				builder.RngDevice(v)
			case "serial_number":
				v, err := XMLSerialNumberReadOne(reader, &t, "serial_number")
				if err != nil {
					return nil, err
				}
				builder.SerialNumber(v)
			case "small_icon":
				v, err := XMLIconReadOne(reader, &t, "small_icon")
				if err != nil {
					return nil, err
				}
				builder.SmallIcon(v)
			case "soundcard_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SoundcardEnabled(v)
			case "sso":
				v, err := XMLSsoReadOne(reader, &t, "sso")
				if err != nil {
					return nil, err
				}
				builder.Sso(v)
			case "start_paused":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.StartPaused(v)
			case "stateless":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Stateless(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "storage_error_resume_behaviour":
				vp, err := XMLVmStorageErrorResumeBehaviourReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageErrorResumeBehaviour(v)
			case "time_zone":
				v, err := XMLTimeZoneReadOne(reader, &t, "time_zone")
				if err != nil {
					return nil, err
				}
				builder.TimeZone(v)
			case "tunnel_migration":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.TunnelMigration(v)
			case "type":
				vp, err := XMLVmTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "usb":
				v, err := XMLUsbReadOne(reader, &t, "usb")
				if err != nil {
					return nil, err
				}
				builder.Usb(v)
			case "virtio_scsi":
				v, err := XMLVirtioScsiReadOne(reader, &t, "virtio_scsi")
				if err != nil {
					return nil, err
				}
				builder.VirtioScsi(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVmBaseReadMany(reader *XMLReader, start *xml.StartElement) (*VmBaseSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VmBaseSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vm_base":
				one, err := XMLVmBaseReadOne(reader, &t, "vm_base")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDnsReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Dns, error) {
	builder := NewDnsBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "dns"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "search_domains":
				v, err := XMLHostReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.SearchDomains(v)
			case "servers":
				v, err := XMLHostReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Servers(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDnsReadMany(reader *XMLReader, start *xml.StartElement) (*DnsSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DnsSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "dns":
				one, err := XMLDnsReadOne(reader, &t, "dns")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCpuProfileReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*CpuProfile, error) {
	builder := NewCpuProfileBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cpu_profile"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "qos":
				v, err := XMLQosReadOne(reader, &t, "qos")
				if err != nil {
					return nil, err
				}
				builder.Qos(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCpuProfileReadMany(reader *XMLReader, start *xml.StartElement) (*CpuProfileSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CpuProfileSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cpu_profile":
				one, err := XMLCpuProfileReadOne(reader, &t, "cpu_profile")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLExternalHostReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ExternalHost, error) {
	builder := NewExternalHostBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "external_host"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "external_host_provider":
				v, err := XMLExternalHostProviderReadOne(reader, &t, "external_host_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalHostProvider(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLExternalHostReadMany(reader *XMLReader, start *xml.StartElement) (*ExternalHostSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ExternalHostSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "external_host":
				one, err := XMLExternalHostReadOne(reader, &t, "external_host")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNicReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Nic, error) {
	builder := NewNicBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "nic"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "boot_protocol":
				vp, err := XMLBootProtocolReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.BootProtocol(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "interface":
				vp, err := XMLNicInterfaceReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Interface(v)
			case "linked":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Linked(v)
			case "mac":
				v, err := XMLMacReadOne(reader, &t, "mac")
				if err != nil {
					return nil, err
				}
				builder.Mac(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "network":
				v, err := XMLNetworkReadOne(reader, &t, "network")
				if err != nil {
					return nil, err
				}
				builder.Network(v)
			case "network_attachments":
				v, err := XMLNetworkAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NetworkAttachments(v)
			case "network_filter_parameters":
				v, err := XMLNetworkFilterParameterReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NetworkFilterParameters(v)
			case "network_labels":
				v, err := XMLNetworkLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NetworkLabels(v)
			case "on_boot":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.OnBoot(v)
			case "plugged":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Plugged(v)
			case "reported_devices":
				v, err := XMLReportedDeviceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ReportedDevices(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "virtual_function_allowed_labels":
				v, err := XMLNetworkLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.VirtualFunctionAllowedLabels(v)
			case "virtual_function_allowed_networks":
				v, err := XMLNetworkReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.VirtualFunctionAllowedNetworks(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "vnic_profile":
				v, err := XMLVnicProfileReadOne(reader, &t, "vnic_profile")
				if err != nil {
					return nil, err
				}
				builder.VnicProfile(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "networkattachments":
			if one.networkAttachments == nil {
				one.networkAttachments = new(NetworkAttachmentSlice)
			}
			one.networkAttachments.href = link.href
		case "networkfilterparameters":
			if one.networkFilterParameters == nil {
				one.networkFilterParameters = new(NetworkFilterParameterSlice)
			}
			one.networkFilterParameters.href = link.href
		case "networklabels":
			if one.networkLabels == nil {
				one.networkLabels = new(NetworkLabelSlice)
			}
			one.networkLabels.href = link.href
		case "reporteddevices":
			if one.reportedDevices == nil {
				one.reportedDevices = new(ReportedDeviceSlice)
			}
			one.reportedDevices.href = link.href
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		case "virtualfunctionallowedlabels":
			if one.virtualFunctionAllowedLabels == nil {
				one.virtualFunctionAllowedLabels = new(NetworkLabelSlice)
			}
			one.virtualFunctionAllowedLabels.href = link.href
		case "virtualfunctionallowednetworks":
			if one.virtualFunctionAllowedNetworks == nil {
				one.virtualFunctionAllowedNetworks = new(NetworkSlice)
			}
			one.virtualFunctionAllowedNetworks.href = link.href
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNicReadMany(reader *XMLReader, start *xml.StartElement) (*NicSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NicSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "nic":
				one, err := XMLNicReadOne(reader, &t, "nic")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLOpenstackVolumeAuthenticationKeyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*OpenstackVolumeAuthenticationKey, error) {
	builder := NewOpenstackVolumeAuthenticationKeyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "openstack_volume_authentication_key"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "creation_date":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.CreationDate(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "openstack_volume_provider":
				v, err := XMLOpenStackVolumeProviderReadOne(reader, &t, "openstack_volume_provider")
				if err != nil {
					return nil, err
				}
				builder.OpenstackVolumeProvider(v)
			case "usage_type":
				vp, err := XMLOpenstackVolumeAuthenticationKeyUsageTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.UsageType(v)
			case "uuid":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Uuid(v)
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLOpenstackVolumeAuthenticationKeyReadMany(reader *XMLReader, start *xml.StartElement) (*OpenstackVolumeAuthenticationKeySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result OpenstackVolumeAuthenticationKeySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "openstack_volume_authentication_key":
				one, err := XMLOpenstackVolumeAuthenticationKeyReadOne(reader, &t, "openstack_volume_authentication_key")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCpuTopologyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*CpuTopology, error) {
	builder := NewCpuTopologyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cpu_topology"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cores":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Cores(v)
			case "sockets":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Sockets(v)
			case "threads":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Threads(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCpuTopologyReadMany(reader *XMLReader, start *xml.StartElement) (*CpuTopologySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CpuTopologySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cpu_topology":
				one, err := XMLCpuTopologyReadOne(reader, &t, "cpu_topology")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLExternalProviderReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ExternalProvider, error) {
	builder := NewExternalProviderBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "external_provider"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authentication_url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.AuthenticationUrl(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "requires_authentication":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RequiresAuthentication(v)
			case "url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Url(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLExternalProviderReadMany(reader *XMLReader, start *xml.StartElement) (*ExternalProviderSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ExternalProviderSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "external_provider":
				one, err := XMLExternalProviderReadOne(reader, &t, "external_provider")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLExternalVmImportReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ExternalVmImport, error) {
	builder := NewExternalVmImportBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "external_vm_import"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "cpu_profile":
				v, err := XMLCpuProfileReadOne(reader, &t, "cpu_profile")
				if err != nil {
					return nil, err
				}
				builder.CpuProfile(v)
			case "drivers_iso":
				v, err := XMLFileReadOne(reader, &t, "drivers_iso")
				if err != nil {
					return nil, err
				}
				builder.DriversIso(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Password(v)
			case "provider":
				vp, err := XMLExternalVmProviderTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Provider(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "sparse":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Sparse(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "url":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Url(v)
			case "username":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Username(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLExternalVmImportReadMany(reader *XMLReader, start *xml.StartElement) (*ExternalVmImportSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ExternalVmImportSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "external_vm_import":
				one, err := XMLExternalVmImportReadOne(reader, &t, "external_vm_import")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLPackageReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Package, error) {
	builder := NewPackageBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "package"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLPackageReadMany(reader *XMLReader, start *xml.StartElement) (*PackageSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result PackageSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "package":
				one, err := XMLPackageReadOne(reader, &t, "package")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLQosReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Qos, error) {
	builder := NewQosBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "qos"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "cpu_limit":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CpuLimit(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "inbound_average":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.InboundAverage(v)
			case "inbound_burst":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.InboundBurst(v)
			case "inbound_peak":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.InboundPeak(v)
			case "max_iops":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxIops(v)
			case "max_read_iops":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxReadIops(v)
			case "max_read_throughput":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxReadThroughput(v)
			case "max_throughput":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxThroughput(v)
			case "max_write_iops":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxWriteIops(v)
			case "max_write_throughput":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxWriteThroughput(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "outbound_average":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.OutboundAverage(v)
			case "outbound_average_linkshare":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.OutboundAverageLinkshare(v)
			case "outbound_average_realtime":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.OutboundAverageRealtime(v)
			case "outbound_average_upperlimit":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.OutboundAverageUpperlimit(v)
			case "outbound_burst":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.OutboundBurst(v)
			case "outbound_peak":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.OutboundPeak(v)
			case "type":
				vp, err := XMLQosTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLQosReadMany(reader *XMLReader, start *xml.StartElement) (*QosSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result QosSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "qos":
				one, err := XMLQosReadOne(reader, &t, "qos")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNicConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NicConfiguration, error) {
	builder := NewNicConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "nic_configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "boot_protocol":
				vp, err := XMLBootProtocolReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.BootProtocol(v)
			case "ip":
				v, err := XMLIpReadOne(reader, &t, "ip")
				if err != nil {
					return nil, err
				}
				builder.Ip(v)
			case "ipv6":
				v, err := XMLIpReadOne(reader, &t, "ipv6")
				if err != nil {
					return nil, err
				}
				builder.Ipv6(v)
			case "ipv6_boot_protocol":
				vp, err := XMLBootProtocolReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Ipv6BootProtocol(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "on_boot":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.OnBoot(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNicConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*NicConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NicConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "nic_configuration":
				one, err := XMLNicConfigurationReadOne(reader, &t, "nic_configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBlockStatisticReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*BlockStatistic, error) {
	builder := NewBlockStatisticBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "block_statistic"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBlockStatisticReadMany(reader *XMLReader, start *xml.StartElement) (*BlockStatisticSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BlockStatisticSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "block_statistic":
				one, err := XMLBlockStatisticReadOne(reader, &t, "block_statistic")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHighAvailabilityReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*HighAvailability, error) {
	builder := NewHighAvailabilityBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "high_availability"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "priority":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Priority(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHighAvailabilityReadMany(reader *XMLReader, start *xml.StartElement) (*HighAvailabilitySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HighAvailabilitySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "high_availability":
				one, err := XMLHighAvailabilityReadOne(reader, &t, "high_availability")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDiskSnapshotReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*DiskSnapshot, error) {
	builder := NewDiskSnapshotBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "disk_snapshot"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Active(v)
			case "actual_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ActualSize(v)
			case "alias":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Alias(v)
			case "backup":
				vp, err := XMLDiskBackupReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Backup(v)
			case "bootable":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Bootable(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "content_type":
				vp, err := XMLDiskContentTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.ContentType(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk":
				v, err := XMLDiskReadOne(reader, &t, "disk")
				if err != nil {
					return nil, err
				}
				builder.Disk(v)
			case "disk_profile":
				v, err := XMLDiskProfileReadOne(reader, &t, "disk_profile")
				if err != nil {
					return nil, err
				}
				builder.DiskProfile(v)
			case "format":
				vp, err := XMLDiskFormatReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Format(v)
			case "image_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ImageId(v)
			case "initial_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.InitialSize(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "interface":
				vp, err := XMLDiskInterfaceReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Interface(v)
			case "logical_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.LogicalName(v)
			case "lun_storage":
				v, err := XMLHostStorageReadOne(reader, &t, "lun_storage")
				if err != nil {
					return nil, err
				}
				builder.LunStorage(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "openstack_volume_type":
				v, err := XMLOpenStackVolumeTypeReadOne(reader, &t, "openstack_volume_type")
				if err != nil {
					return nil, err
				}
				builder.OpenstackVolumeType(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "propagate_errors":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.PropagateErrors(v)
			case "provisioned_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ProvisionedSize(v)
			case "qcow_version":
				vp, err := XMLQcowVersionReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.QcowVersion(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "read_only":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ReadOnly(v)
			case "sgio":
				vp, err := XMLScsiGenericIOReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Sgio(v)
			case "shareable":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Shareable(v)
			case "snapshot":
				v, err := XMLSnapshotReadOne(reader, &t, "snapshot")
				if err != nil {
					return nil, err
				}
				builder.Snapshot(v)
			case "sparse":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Sparse(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLDiskStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "storage_domains":
				v, err := XMLStorageDomainReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.StorageDomains(v)
			case "storage_type":
				vp, err := XMLDiskStorageTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageType(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "total_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.TotalSize(v)
			case "uses_scsi_reservation":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UsesScsiReservation(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "wipe_after_delete":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.WipeAfterDelete(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		case "storagedomains":
			if one.storageDomains == nil {
				one.storageDomains = new(StorageDomainSlice)
			}
			one.storageDomains.href = link.href
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDiskSnapshotReadMany(reader *XMLReader, start *xml.StartElement) (*DiskSnapshotSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DiskSnapshotSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "disk_snapshot":
				one, err := XMLDiskSnapshotReadOne(reader, &t, "disk_snapshot")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCdromReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Cdrom, error) {
	builder := NewCdromBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cdrom"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "file":
				v, err := XMLFileReadOne(reader, &t, "file")
				if err != nil {
					return nil, err
				}
				builder.File(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCdromReadMany(reader *XMLReader, start *xml.StartElement) (*CdromSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CdromSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cdrom":
				one, err := XMLCdromReadOne(reader, &t, "cdrom")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLConsoleReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Console, error) {
	builder := NewConsoleBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "console"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLConsoleReadMany(reader *XMLReader, start *xml.StartElement) (*ConsoleSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ConsoleSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "console":
				one, err := XMLConsoleReadOne(reader, &t, "console")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLFencingPolicyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*FencingPolicy, error) {
	builder := NewFencingPolicyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "fencing_policy"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "skip_if_connectivity_broken":
				v, err := XMLSkipIfConnectivityBrokenReadOne(reader, &t, "skip_if_connectivity_broken")
				if err != nil {
					return nil, err
				}
				builder.SkipIfConnectivityBroken(v)
			case "skip_if_gluster_bricks_up":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SkipIfGlusterBricksUp(v)
			case "skip_if_gluster_quorum_not_met":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SkipIfGlusterQuorumNotMet(v)
			case "skip_if_sd_active":
				v, err := XMLSkipIfSdActiveReadOne(reader, &t, "skip_if_sd_active")
				if err != nil {
					return nil, err
				}
				builder.SkipIfSdActive(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLFencingPolicyReadMany(reader *XMLReader, start *xml.StartElement) (*FencingPolicySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result FencingPolicySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "fencing_policy":
				one, err := XMLFencingPolicyReadOne(reader, &t, "fencing_policy")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSkipIfSdActiveReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SkipIfSdActive, error) {
	builder := NewSkipIfSdActiveBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "skip_if_sd_active"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSkipIfSdActiveReadMany(reader *XMLReader, start *xml.StartElement) (*SkipIfSdActiveSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SkipIfSdActiveSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "skip_if_sd_active":
				one, err := XMLSkipIfSdActiveReadOne(reader, &t, "skip_if_sd_active")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLPayloadReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Payload, error) {
	builder := NewPayloadBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "payload"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "files":
				v, err := XMLFileReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Files(v)
			case "type":
				vp, err := XMLVmDeviceTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "volume_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.VolumeId(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLPayloadReadMany(reader *XMLReader, start *xml.StartElement) (*PayloadSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result PayloadSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "payload":
				one, err := XMLPayloadReadOne(reader, &t, "payload")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSshPublicKeyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SshPublicKey, error) {
	builder := NewSshPublicKeyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "ssh_public_key"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "content":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Content(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSshPublicKeyReadMany(reader *XMLReader, start *xml.StartElement) (*SshPublicKeySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SshPublicKeySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ssh_public_key":
				one, err := XMLSshPublicKeyReadOne(reader, &t, "ssh_public_key")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLApiSummaryItemReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ApiSummaryItem, error) {
	builder := NewApiSummaryItemBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "api_summary_item"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Active(v)
			case "total":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Total(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLApiSummaryItemReadMany(reader *XMLReader, start *xml.StartElement) (*ApiSummaryItemSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ApiSummaryItemSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "api_summary_item":
				one, err := XMLApiSummaryItemReadOne(reader, &t, "api_summary_item")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDiskReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Disk, error) {
	builder := NewDiskBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "disk"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Active(v)
			case "actual_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ActualSize(v)
			case "alias":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Alias(v)
			case "backup":
				vp, err := XMLDiskBackupReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Backup(v)
			case "bootable":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Bootable(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "content_type":
				vp, err := XMLDiskContentTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.ContentType(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk_profile":
				v, err := XMLDiskProfileReadOne(reader, &t, "disk_profile")
				if err != nil {
					return nil, err
				}
				builder.DiskProfile(v)
			case "format":
				vp, err := XMLDiskFormatReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Format(v)
			case "image_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ImageId(v)
			case "initial_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.InitialSize(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "interface":
				vp, err := XMLDiskInterfaceReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Interface(v)
			case "logical_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.LogicalName(v)
			case "lun_storage":
				v, err := XMLHostStorageReadOne(reader, &t, "lun_storage")
				if err != nil {
					return nil, err
				}
				builder.LunStorage(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "openstack_volume_type":
				v, err := XMLOpenStackVolumeTypeReadOne(reader, &t, "openstack_volume_type")
				if err != nil {
					return nil, err
				}
				builder.OpenstackVolumeType(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "propagate_errors":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.PropagateErrors(v)
			case "provisioned_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ProvisionedSize(v)
			case "qcow_version":
				vp, err := XMLQcowVersionReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.QcowVersion(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "read_only":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ReadOnly(v)
			case "sgio":
				vp, err := XMLScsiGenericIOReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Sgio(v)
			case "shareable":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Shareable(v)
			case "snapshot":
				v, err := XMLSnapshotReadOne(reader, &t, "snapshot")
				if err != nil {
					return nil, err
				}
				builder.Snapshot(v)
			case "sparse":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Sparse(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLDiskStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "storage_domains":
				v, err := XMLStorageDomainReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.StorageDomains(v)
			case "storage_type":
				vp, err := XMLDiskStorageTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageType(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "total_size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.TotalSize(v)
			case "uses_scsi_reservation":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UsesScsiReservation(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "wipe_after_delete":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.WipeAfterDelete(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		case "storagedomains":
			if one.storageDomains == nil {
				one.storageDomains = new(StorageDomainSlice)
			}
			one.storageDomains.href = link.href
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDiskReadMany(reader *XMLReader, start *xml.StartElement) (*DiskSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DiskSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "disk":
				one, err := XMLDiskReadOne(reader, &t, "disk")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNetworkLabelReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NetworkLabel, error) {
	builder := NewNetworkLabelBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "network_label"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "host_nic":
				v, err := XMLHostNicReadOne(reader, &t, "host_nic")
				if err != nil {
					return nil, err
				}
				builder.HostNic(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "network":
				v, err := XMLNetworkReadOne(reader, &t, "network")
				if err != nil {
					return nil, err
				}
				builder.Network(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNetworkLabelReadMany(reader *XMLReader, start *xml.StartElement) (*NetworkLabelSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NetworkLabelSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "network_label":
				one, err := XMLNetworkLabelReadOne(reader, &t, "network_label")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSpmReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Spm, error) {
	builder := NewSpmBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "spm"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "priority":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Priority(v)
			case "status":
				vp, err := XMLSpmStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSpmReadMany(reader *XMLReader, start *xml.StartElement) (*SpmSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SpmSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "spm":
				one, err := XMLSpmReadOne(reader, &t, "spm")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHostDeviceReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*HostDevice, error) {
	builder := NewHostDeviceBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "host_device"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "capability":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Capability(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "driver":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Driver(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "iommu_group":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.IommuGroup(v)
			case "m_dev_types":
				v, err := XMLMDevTypeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.MDevTypes(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "parent_device":
				v, err := XMLHostDeviceReadOne(reader, &t, "parent_device")
				if err != nil {
					return nil, err
				}
				builder.ParentDevice(v)
			case "physical_function":
				v, err := XMLHostDeviceReadOne(reader, &t, "physical_function")
				if err != nil {
					return nil, err
				}
				builder.PhysicalFunction(v)
			case "placeholder":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Placeholder(v)
			case "product":
				v, err := XMLProductReadOne(reader, &t, "product")
				if err != nil {
					return nil, err
				}
				builder.Product(v)
			case "vendor":
				v, err := XMLVendorReadOne(reader, &t, "vendor")
				if err != nil {
					return nil, err
				}
				builder.Vendor(v)
			case "virtual_functions":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.VirtualFunctions(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHostDeviceReadMany(reader *XMLReader, start *xml.StartElement) (*HostDeviceSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HostDeviceSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "host_device":
				one, err := XMLHostDeviceReadOne(reader, &t, "host_device")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLPmProxyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*PmProxy, error) {
	builder := NewPmProxyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "pm_proxy"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "type":
				vp, err := XMLPmProxyTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLPmProxyReadMany(reader *XMLReader, start *xml.StartElement) (*PmProxySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result PmProxySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "pm_proxy":
				one, err := XMLPmProxyReadOne(reader, &t, "pm_proxy")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSkipIfConnectivityBrokenReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SkipIfConnectivityBroken, error) {
	builder := NewSkipIfConnectivityBrokenBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "skip_if_connectivity_broken"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "threshold":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Threshold(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSkipIfConnectivityBrokenReadMany(reader *XMLReader, start *xml.StartElement) (*SkipIfConnectivityBrokenSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SkipIfConnectivityBrokenSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "skip_if_connectivity_broken":
				one, err := XMLSkipIfConnectivityBrokenReadOne(reader, &t, "skip_if_connectivity_broken")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSchedulingPolicyUnitReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SchedulingPolicyUnit, error) {
	builder := NewSchedulingPolicyUnitBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "scheduling_policy_unit"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "internal":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Internal(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "type":
				vp, err := XMLPolicyUnitTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSchedulingPolicyUnitReadMany(reader *XMLReader, start *xml.StartElement) (*SchedulingPolicyUnitSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SchedulingPolicyUnitSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "scheduling_policy_unit":
				one, err := XMLSchedulingPolicyUnitReadOne(reader, &t, "scheduling_policy_unit")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVnicPassThroughReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VnicPassThrough, error) {
	builder := NewVnicPassThroughBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vnic_pass_through"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "mode":
				vp, err := XMLVnicPassThroughModeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Mode(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVnicPassThroughReadMany(reader *XMLReader, start *xml.StartElement) (*VnicPassThroughSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VnicPassThroughSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vnic_pass_through":
				one, err := XMLVnicPassThroughReadOne(reader, &t, "vnic_pass_through")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLPermitReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Permit, error) {
	builder := NewPermitBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "permit"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "administrative":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Administrative(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "role":
				v, err := XMLRoleReadOne(reader, &t, "role")
				if err != nil {
					return nil, err
				}
				builder.Role(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLPermitReadMany(reader *XMLReader, start *xml.StartElement) (*PermitSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result PermitSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "permit":
				one, err := XMLPermitReadOne(reader, &t, "permit")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLClusterReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Cluster, error) {
	builder := NewClusterBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cluster"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "affinity_groups":
				v, err := XMLAffinityGroupReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.AffinityGroups(v)
			case "ballooning_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.BallooningEnabled(v)
			case "bios_type":
				vp, err := XMLBiosTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.BiosType(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "cpu_profiles":
				v, err := XMLCpuProfileReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CpuProfiles(v)
			case "custom_scheduling_policy_properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CustomSchedulingPolicyProperties(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "display":
				v, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "enabled_features":
				v, err := XMLClusterFeatureReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.EnabledFeatures(v)
			case "error_handling":
				v, err := XMLErrorHandlingReadOne(reader, &t, "error_handling")
				if err != nil {
					return nil, err
				}
				builder.ErrorHandling(v)
			case "external_network_providers":
				v, err := XMLExternalProviderReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ExternalNetworkProviders(v)
			case "fencing_policy":
				v, err := XMLFencingPolicyReadOne(reader, &t, "fencing_policy")
				if err != nil {
					return nil, err
				}
				builder.FencingPolicy(v)
			case "firewall_type":
				vp, err := XMLFirewallTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.FirewallType(v)
			case "gluster_hooks":
				v, err := XMLGlusterHookReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.GlusterHooks(v)
			case "gluster_service":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.GlusterService(v)
			case "gluster_tuned_profile":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.GlusterTunedProfile(v)
			case "gluster_volumes":
				v, err := XMLGlusterVolumeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.GlusterVolumes(v)
			case "ha_reservation":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.HaReservation(v)
			case "ksm":
				v, err := XMLKsmReadOne(reader, &t, "ksm")
				if err != nil {
					return nil, err
				}
				builder.Ksm(v)
			case "log_max_memory_used_threshold":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.LogMaxMemoryUsedThreshold(v)
			case "log_max_memory_used_threshold_type":
				vp, err := XMLLogMaxMemoryUsedThresholdTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.LogMaxMemoryUsedThresholdType(v)
			case "mac_pool":
				v, err := XMLMacPoolReadOne(reader, &t, "mac_pool")
				if err != nil {
					return nil, err
				}
				builder.MacPool(v)
			case "maintenance_reason_required":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MaintenanceReasonRequired(v)
			case "management_network":
				v, err := XMLNetworkReadOne(reader, &t, "management_network")
				if err != nil {
					return nil, err
				}
				builder.ManagementNetwork(v)
			case "memory_policy":
				v, err := XMLMemoryPolicyReadOne(reader, &t, "memory_policy")
				if err != nil {
					return nil, err
				}
				builder.MemoryPolicy(v)
			case "migration":
				v, err := XMLMigrationOptionsReadOne(reader, &t, "migration")
				if err != nil {
					return nil, err
				}
				builder.Migration(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "network_filters":
				v, err := XMLNetworkFilterReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NetworkFilters(v)
			case "networks":
				v, err := XMLNetworkReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Networks(v)
			case "optional_reason":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.OptionalReason(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "required_rng_sources":
				v, err := XMLRngSourceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.RequiredRngSources(v)
			case "scheduling_policy":
				v, err := XMLSchedulingPolicyReadOne(reader, &t, "scheduling_policy")
				if err != nil {
					return nil, err
				}
				builder.SchedulingPolicy(v)
			case "serial_number":
				v, err := XMLSerialNumberReadOne(reader, &t, "serial_number")
				if err != nil {
					return nil, err
				}
				builder.SerialNumber(v)
			case "supported_versions":
				v, err := XMLVersionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.SupportedVersions(v)
			case "switch_type":
				vp, err := XMLSwitchTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.SwitchType(v)
			case "threads_as_cores":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ThreadsAsCores(v)
			case "trusted_service":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.TrustedService(v)
			case "tunnel_migration":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.TunnelMigration(v)
			case "version":
				v, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "virt_service":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.VirtService(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "affinitygroups":
			if one.affinityGroups == nil {
				one.affinityGroups = new(AffinityGroupSlice)
			}
			one.affinityGroups.href = link.href
		case "cpuprofiles":
			if one.cpuProfiles == nil {
				one.cpuProfiles = new(CpuProfileSlice)
			}
			one.cpuProfiles.href = link.href
		case "enabledfeatures":
			if one.enabledFeatures == nil {
				one.enabledFeatures = new(ClusterFeatureSlice)
			}
			one.enabledFeatures.href = link.href
		case "externalnetworkproviders":
			if one.externalNetworkProviders == nil {
				one.externalNetworkProviders = new(ExternalProviderSlice)
			}
			one.externalNetworkProviders.href = link.href
		case "glusterhooks":
			if one.glusterHooks == nil {
				one.glusterHooks = new(GlusterHookSlice)
			}
			one.glusterHooks.href = link.href
		case "glustervolumes":
			if one.glusterVolumes == nil {
				one.glusterVolumes = new(GlusterVolumeSlice)
			}
			one.glusterVolumes.href = link.href
		case "networkfilters":
			if one.networkFilters == nil {
				one.networkFilters = new(NetworkFilterSlice)
			}
			one.networkFilters.href = link.href
		case "networks":
			if one.networks == nil {
				one.networks = new(NetworkSlice)
			}
			one.networks.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLClusterReadMany(reader *XMLReader, start *xml.StartElement) (*ClusterSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ClusterSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cluster":
				one, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterBrickMemoryInfoReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterBrickMemoryInfo, error) {
	builder := NewGlusterBrickMemoryInfoBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "brick_memoryinfo"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "memory_pools":
				v, err := XMLGlusterMemoryPoolReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.MemoryPools(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterBrickMemoryInfoReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterBrickMemoryInfoSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterBrickMemoryInfoSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "brick_memoryinfo":
				one, err := XMLGlusterBrickMemoryInfoReadOne(reader, &t, "brick_memoryinfo")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLVmPoolReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*VmPool, error) {
	builder := NewVmPoolBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "vm_pool"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "auto_storage_select":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.AutoStorageSelect(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "display":
				v, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "max_user_vms":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxUserVms(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "prestarted_vms":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.PrestartedVms(v)
			case "rng_device":
				v, err := XMLRngDeviceReadOne(reader, &t, "rng_device")
				if err != nil {
					return nil, err
				}
				builder.RngDevice(v)
			case "size":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Size(v)
			case "soundcard_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SoundcardEnabled(v)
			case "stateful":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Stateful(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "type":
				vp, err := XMLVmPoolTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "use_latest_template_version":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UseLatestTemplateVersion(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLVmPoolReadMany(reader *XMLReader, start *xml.StartElement) (*VmPoolSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result VmPoolSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "vm_pool":
				one, err := XMLVmPoolReadOne(reader, &t, "vm_pool")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLProxyTicketReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ProxyTicket, error) {
	builder := NewProxyTicketBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "proxy_ticket"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Value(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLProxyTicketReadMany(reader *XMLReader, start *xml.StartElement) (*ProxyTicketSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ProxyTicketSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "proxy_ticket":
				one, err := XMLProxyTicketReadOne(reader, &t, "proxy_ticket")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterServerHookReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterServerHook, error) {
	builder := NewGlusterServerHookBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "server_hook"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "checksum":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Checksum(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "content_type":
				vp, err := XMLHookContentTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.ContentType(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "status":
				vp, err := XMLGlusterHookStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterServerHookReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterServerHookSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterServerHookSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "server_hook":
				one, err := XMLGlusterServerHookReadOne(reader, &t, "server_hook")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLStatisticReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Statistic, error) {
	builder := NewStatisticBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "statistic"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "brick":
				v, err := XMLGlusterBrickReadOne(reader, &t, "brick")
				if err != nil {
					return nil, err
				}
				builder.Brick(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk":
				v, err := XMLDiskReadOne(reader, &t, "disk")
				if err != nil {
					return nil, err
				}
				builder.Disk(v)
			case "gluster_volume":
				v, err := XMLGlusterVolumeReadOne(reader, &t, "gluster_volume")
				if err != nil {
					return nil, err
				}
				builder.GlusterVolume(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "host_nic":
				v, err := XMLHostNicReadOne(reader, &t, "host_nic")
				if err != nil {
					return nil, err
				}
				builder.HostNic(v)
			case "host_numa_node":
				v, err := XMLNumaNodeReadOne(reader, &t, "host_numa_node")
				if err != nil {
					return nil, err
				}
				builder.HostNumaNode(v)
			case "kind":
				vp, err := XMLStatisticKindReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Kind(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "nic":
				v, err := XMLNicReadOne(reader, &t, "nic")
				if err != nil {
					return nil, err
				}
				builder.Nic(v)
			case "step":
				v, err := XMLStepReadOne(reader, &t, "step")
				if err != nil {
					return nil, err
				}
				builder.Step(v)
			case "type":
				vp, err := XMLValueTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "unit":
				vp, err := XMLStatisticUnitReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Unit(v)
			case "values":
				v, err := XMLValueReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Values(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLStatisticReadMany(reader *XMLReader, start *xml.StartElement) (*StatisticSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result StatisticSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "statistic":
				one, err := XMLStatisticReadOne(reader, &t, "statistic")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSnapshotReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Snapshot, error) {
	builder := NewSnapshotBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "snapshot"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "affinity_labels":
				v, err := XMLAffinityLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.AffinityLabels(v)
			case "applications":
				v, err := XMLApplicationReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Applications(v)
			case "bios":
				v, err := XMLBiosReadOne(reader, &t, "bios")
				if err != nil {
					return nil, err
				}
				builder.Bios(v)
			case "cdroms":
				v, err := XMLCdromReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Cdroms(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "console":
				v, err := XMLConsoleReadOne(reader, &t, "console")
				if err != nil {
					return nil, err
				}
				builder.Console(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "cpu_profile":
				v, err := XMLCpuProfileReadOne(reader, &t, "cpu_profile")
				if err != nil {
					return nil, err
				}
				builder.CpuProfile(v)
			case "cpu_shares":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.CpuShares(v)
			case "creation_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.CreationTime(v)
			case "custom_compatibility_version":
				v, err := XMLVersionReadOne(reader, &t, "custom_compatibility_version")
				if err != nil {
					return nil, err
				}
				builder.CustomCompatibilityVersion(v)
			case "custom_cpu_model":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomCpuModel(v)
			case "custom_emulated_machine":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.CustomEmulatedMachine(v)
			case "custom_properties":
				v, err := XMLCustomPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.CustomProperties(v)
			case "date":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.Date(v)
			case "delete_protected":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DeleteProtected(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk_attachments":
				v, err := XMLDiskAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DiskAttachments(v)
			case "disks":
				v, err := XMLDiskReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Disks(v)
			case "display":
				v, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "domain":
				v, err := XMLDomainReadOne(reader, &t, "domain")
				if err != nil {
					return nil, err
				}
				builder.Domain(v)
			case "external_host_provider":
				v, err := XMLExternalHostProviderReadOne(reader, &t, "external_host_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalHostProvider(v)
			case "floppies":
				v, err := XMLFloppyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Floppies(v)
			case "fqdn":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Fqdn(v)
			case "graphics_consoles":
				v, err := XMLGraphicsConsoleReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.GraphicsConsoles(v)
			case "guest_operating_system":
				v, err := XMLGuestOperatingSystemReadOne(reader, &t, "guest_operating_system")
				if err != nil {
					return nil, err
				}
				builder.GuestOperatingSystem(v)
			case "guest_time_zone":
				v, err := XMLTimeZoneReadOne(reader, &t, "guest_time_zone")
				if err != nil {
					return nil, err
				}
				builder.GuestTimeZone(v)
			case "has_illegal_images":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.HasIllegalImages(v)
			case "high_availability":
				v, err := XMLHighAvailabilityReadOne(reader, &t, "high_availability")
				if err != nil {
					return nil, err
				}
				builder.HighAvailability(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "host_devices":
				v, err := XMLHostDeviceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.HostDevices(v)
			case "initialization":
				v, err := XMLInitializationReadOne(reader, &t, "initialization")
				if err != nil {
					return nil, err
				}
				builder.Initialization(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "io":
				v, err := XMLIoReadOne(reader, &t, "io")
				if err != nil {
					return nil, err
				}
				builder.Io(v)
			case "katello_errata":
				v, err := XMLKatelloErratumReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.KatelloErrata(v)
			case "large_icon":
				v, err := XMLIconReadOne(reader, &t, "large_icon")
				if err != nil {
					return nil, err
				}
				builder.LargeIcon(v)
			case "lease":
				v, err := XMLStorageDomainLeaseReadOne(reader, &t, "lease")
				if err != nil {
					return nil, err
				}
				builder.Lease(v)
			case "memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Memory(v)
			case "memory_policy":
				v, err := XMLMemoryPolicyReadOne(reader, &t, "memory_policy")
				if err != nil {
					return nil, err
				}
				builder.MemoryPolicy(v)
			case "migration":
				v, err := XMLMigrationOptionsReadOne(reader, &t, "migration")
				if err != nil {
					return nil, err
				}
				builder.Migration(v)
			case "migration_downtime":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MigrationDowntime(v)
			case "multi_queues_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MultiQueuesEnabled(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "next_run_configuration_exists":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.NextRunConfigurationExists(v)
			case "nics":
				v, err := XMLNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Nics(v)
			case "host_numa_nodes":
				v, err := XMLNumaNodeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NumaNodes(v)
			case "numa_tune_mode":
				vp, err := XMLNumaTuneModeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.NumaTuneMode(v)
			case "origin":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Origin(v)
			case "original_template":
				v, err := XMLTemplateReadOne(reader, &t, "original_template")
				if err != nil {
					return nil, err
				}
				builder.OriginalTemplate(v)
			case "os":
				v, err := XMLOperatingSystemReadOne(reader, &t, "os")
				if err != nil {
					return nil, err
				}
				builder.Os(v)
			case "payloads":
				v, err := XMLPayloadReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Payloads(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "persist_memorystate":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.PersistMemorystate(v)
			case "placement_policy":
				v, err := XMLVmPlacementPolicyReadOne(reader, &t, "placement_policy")
				if err != nil {
					return nil, err
				}
				builder.PlacementPolicy(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "reported_devices":
				v, err := XMLReportedDeviceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ReportedDevices(v)
			case "rng_device":
				v, err := XMLRngDeviceReadOne(reader, &t, "rng_device")
				if err != nil {
					return nil, err
				}
				builder.RngDevice(v)
			case "run_once":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RunOnce(v)
			case "serial_number":
				v, err := XMLSerialNumberReadOne(reader, &t, "serial_number")
				if err != nil {
					return nil, err
				}
				builder.SerialNumber(v)
			case "sessions":
				v, err := XMLSessionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Sessions(v)
			case "small_icon":
				v, err := XMLIconReadOne(reader, &t, "small_icon")
				if err != nil {
					return nil, err
				}
				builder.SmallIcon(v)
			case "snapshot_status":
				vp, err := XMLSnapshotStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.SnapshotStatus(v)
			case "snapshot_type":
				vp, err := XMLSnapshotTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.SnapshotType(v)
			case "snapshots":
				v, err := XMLSnapshotReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Snapshots(v)
			case "soundcard_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SoundcardEnabled(v)
			case "sso":
				v, err := XMLSsoReadOne(reader, &t, "sso")
				if err != nil {
					return nil, err
				}
				builder.Sso(v)
			case "start_paused":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.StartPaused(v)
			case "start_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.StartTime(v)
			case "stateless":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Stateless(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLVmStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "status_detail":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.StatusDetail(v)
			case "stop_reason":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.StopReason(v)
			case "stop_time":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.StopTime(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "storage_error_resume_behaviour":
				vp, err := XMLVmStorageErrorResumeBehaviourReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.StorageErrorResumeBehaviour(v)
			case "tags":
				v, err := XMLTagReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Tags(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "time_zone":
				v, err := XMLTimeZoneReadOne(reader, &t, "time_zone")
				if err != nil {
					return nil, err
				}
				builder.TimeZone(v)
			case "tunnel_migration":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.TunnelMigration(v)
			case "type":
				vp, err := XMLVmTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "usb":
				v, err := XMLUsbReadOne(reader, &t, "usb")
				if err != nil {
					return nil, err
				}
				builder.Usb(v)
			case "use_latest_template_version":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UseLatestTemplateVersion(v)
			case "virtio_scsi":
				v, err := XMLVirtioScsiReadOne(reader, &t, "virtio_scsi")
				if err != nil {
					return nil, err
				}
				builder.VirtioScsi(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vm_pool":
				v, err := XMLVmPoolReadOne(reader, &t, "vm_pool")
				if err != nil {
					return nil, err
				}
				builder.VmPool(v)
			case "watchdogs":
				v, err := XMLWatchdogReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Watchdogs(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "affinitylabels":
			if one.affinityLabels == nil {
				one.affinityLabels = new(AffinityLabelSlice)
			}
			one.affinityLabels.href = link.href
		case "applications":
			if one.applications == nil {
				one.applications = new(ApplicationSlice)
			}
			one.applications.href = link.href
		case "cdroms":
			if one.cdroms == nil {
				one.cdroms = new(CdromSlice)
			}
			one.cdroms.href = link.href
		case "diskattachments":
			if one.diskAttachments == nil {
				one.diskAttachments = new(DiskAttachmentSlice)
			}
			one.diskAttachments.href = link.href
		case "disks":
			if one.disks == nil {
				one.disks = new(DiskSlice)
			}
			one.disks.href = link.href
		case "floppies":
			if one.floppies == nil {
				one.floppies = new(FloppySlice)
			}
			one.floppies.href = link.href
		case "graphicsconsoles":
			if one.graphicsConsoles == nil {
				one.graphicsConsoles = new(GraphicsConsoleSlice)
			}
			one.graphicsConsoles.href = link.href
		case "hostdevices":
			if one.hostDevices == nil {
				one.hostDevices = new(HostDeviceSlice)
			}
			one.hostDevices.href = link.href
		case "katelloerrata":
			if one.katelloErrata == nil {
				one.katelloErrata = new(KatelloErratumSlice)
			}
			one.katelloErrata.href = link.href
		case "nics":
			if one.nics == nil {
				one.nics = new(NicSlice)
			}
			one.nics.href = link.href
		case "numanodes":
			if one.numaNodes == nil {
				one.numaNodes = new(NumaNodeSlice)
			}
			one.numaNodes.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "reporteddevices":
			if one.reportedDevices == nil {
				one.reportedDevices = new(ReportedDeviceSlice)
			}
			one.reportedDevices.href = link.href
		case "sessions":
			if one.sessions == nil {
				one.sessions = new(SessionSlice)
			}
			one.sessions.href = link.href
		case "snapshots":
			if one.snapshots == nil {
				one.snapshots = new(SnapshotSlice)
			}
			one.snapshots.href = link.href
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		case "tags":
			if one.tags == nil {
				one.tags = new(TagSlice)
			}
			one.tags.href = link.href
		case "watchdogs":
			if one.watchdogs == nil {
				one.watchdogs = new(WatchdogSlice)
			}
			one.watchdogs.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSnapshotReadMany(reader *XMLReader, start *xml.StartElement) (*SnapshotSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SnapshotSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "snapshot":
				one, err := XMLSnapshotReadOne(reader, &t, "snapshot")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHostedEngineReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*HostedEngine, error) {
	builder := NewHostedEngineBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "hosted_engine"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Active(v)
			case "configured":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Configured(v)
			case "global_maintenance":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.GlobalMaintenance(v)
			case "local_maintenance":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.LocalMaintenance(v)
			case "score":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Score(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHostedEngineReadMany(reader *XMLReader, start *xml.StartElement) (*HostedEngineSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HostedEngineSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "hosted_engine":
				one, err := XMLHostedEngineReadOne(reader, &t, "hosted_engine")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDiskAttachmentReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*DiskAttachment, error) {
	builder := NewDiskAttachmentBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "disk_attachment"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "active":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Active(v)
			case "bootable":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Bootable(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disk":
				v, err := XMLDiskReadOne(reader, &t, "disk")
				if err != nil {
					return nil, err
				}
				builder.Disk(v)
			case "interface":
				vp, err := XMLDiskInterfaceReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Interface(v)
			case "logical_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.LogicalName(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "pass_discard":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.PassDiscard(v)
			case "read_only":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ReadOnly(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "uses_scsi_reservation":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UsesScsiReservation(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDiskAttachmentReadMany(reader *XMLReader, start *xml.StartElement) (*DiskAttachmentSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DiskAttachmentSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "disk_attachment":
				one, err := XMLDiskAttachmentReadOne(reader, &t, "disk_attachment")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCertificateReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Certificate, error) {
	builder := NewCertificateBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "certificate"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "content":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Content(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "organization":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Organization(v)
			case "subject":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Subject(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCertificateReadMany(reader *XMLReader, start *xml.StartElement) (*CertificateSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CertificateSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "certificate":
				one, err := XMLCertificateReadOne(reader, &t, "certificate")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLReportedConfigurationReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ReportedConfiguration, error) {
	builder := NewReportedConfigurationBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "reported_configuration"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "actual_value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ActualValue(v)
			case "expected_value":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ExpectedValue(v)
			case "in_sync":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.InSync(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLReportedConfigurationReadMany(reader *XMLReader, start *xml.StartElement) (*ReportedConfigurationSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ReportedConfigurationSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "reported_configuration":
				one, err := XMLReportedConfigurationReadOne(reader, &t, "reported_configuration")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLBackupReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Backup, error) {
	builder := NewBackupBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "backup"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "creation_date":
				v, err := reader.ReadTime(&t)
				if err != nil {
					return nil, err
				}
				builder.CreationDate(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "disks":
				v, err := XMLDiskReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Disks(v)
			case "from_checkpoint_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.FromCheckpointId(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "phase":
				vp, err := XMLBackupPhaseReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Phase(v)
			case "to_checkpoint_id":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ToCheckpointId(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "disks":
			if one.disks == nil {
				one.disks = new(DiskSlice)
			}
			one.disks.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLBackupReadMany(reader *XMLReader, start *xml.StartElement) (*BackupSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result BackupSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "backup":
				one, err := XMLBackupReadOne(reader, &t, "backup")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSchedulingPolicyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SchedulingPolicy, error) {
	builder := NewSchedulingPolicyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "scheduling_policy"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "balances":
				v, err := XMLBalanceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Balances(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "default_policy":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DefaultPolicy(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "filters":
				v, err := XMLFilterReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Filters(v)
			case "locked":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Locked(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "weight":
				v, err := XMLWeightReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Weight(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "balances":
			if one.balances == nil {
				one.balances = new(BalanceSlice)
			}
			one.balances.href = link.href
		case "filters":
			if one.filters == nil {
				one.filters = new(FilterSlice)
			}
			one.filters.href = link.href
		case "weight":
			if one.weight == nil {
				one.weight = new(WeightSlice)
			}
			one.weight.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSchedulingPolicyReadMany(reader *XMLReader, start *xml.StartElement) (*SchedulingPolicySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SchedulingPolicySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "scheduling_policy":
				one, err := XMLSchedulingPolicyReadOne(reader, &t, "scheduling_policy")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCpuReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Cpu, error) {
	builder := NewCpuBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cpu"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "architecture":
				vp, err := XMLArchitectureReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Architecture(v)
			case "cores":
				v, err := XMLCoreReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Cores(v)
			case "cpu_tune":
				v, err := XMLCpuTuneReadOne(reader, &t, "cpu_tune")
				if err != nil {
					return nil, err
				}
				builder.CpuTune(v)
			case "level":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Level(v)
			case "mode":
				vp, err := XMLCpuModeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Mode(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "speed":
				v, err := reader.ReadFloat64(&t)
				if err != nil {
					return nil, err
				}
				builder.Speed(v)
			case "topology":
				v, err := XMLCpuTopologyReadOne(reader, &t, "topology")
				if err != nil {
					return nil, err
				}
				builder.Topology(v)
			case "type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCpuReadMany(reader *XMLReader, start *xml.StartElement) (*CpuSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CpuSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cpu":
				one, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLRegistrationAffinityLabelMappingReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*RegistrationAffinityLabelMapping, error) {
	builder := NewRegistrationAffinityLabelMappingBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "registration_affinity_label_mapping"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "from":
				v, err := XMLAffinityLabelReadOne(reader, &t, "from")
				if err != nil {
					return nil, err
				}
				builder.From(v)
			case "to":
				v, err := XMLAffinityLabelReadOne(reader, &t, "to")
				if err != nil {
					return nil, err
				}
				builder.To(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLRegistrationAffinityLabelMappingReadMany(reader *XMLReader, start *xml.StartElement) (*RegistrationAffinityLabelMappingSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result RegistrationAffinityLabelMappingSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "registration_affinity_label_mapping":
				one, err := XMLRegistrationAffinityLabelMappingReadOne(reader, &t, "registration_affinity_label_mapping")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLIpAddressAssignmentReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*IpAddressAssignment, error) {
	builder := NewIpAddressAssignmentBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "ip_address_assignment"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "assignment_method":
				vp, err := XMLBootProtocolReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.AssignmentMethod(v)
			case "ip":
				v, err := XMLIpReadOne(reader, &t, "ip")
				if err != nil {
					return nil, err
				}
				builder.Ip(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLIpAddressAssignmentReadMany(reader *XMLReader, start *xml.StartElement) (*IpAddressAssignmentSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result IpAddressAssignmentSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ip_address_assignment":
				one, err := XMLIpAddressAssignmentReadOne(reader, &t, "ip_address_assignment")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGuestOperatingSystemReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GuestOperatingSystem, error) {
	builder := NewGuestOperatingSystemBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "guest_operating_system"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "architecture":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Architecture(v)
			case "codename":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Codename(v)
			case "distribution":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Distribution(v)
			case "family":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Family(v)
			case "kernel":
				v, err := XMLKernelReadOne(reader, &t, "kernel")
				if err != nil {
					return nil, err
				}
				builder.Kernel(v)
			case "version":
				v, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGuestOperatingSystemReadMany(reader *XMLReader, start *xml.StartElement) (*GuestOperatingSystemSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GuestOperatingSystemSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "guest_operating_system":
				one, err := XMLGuestOperatingSystemReadOne(reader, &t, "guest_operating_system")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCpuTypeReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*CpuType, error) {
	builder := NewCpuTypeBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cpu_type"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "architecture":
				vp, err := XMLArchitectureReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Architecture(v)
			case "level":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Level(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCpuTypeReadMany(reader *XMLReader, start *xml.StartElement) (*CpuTypeSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CpuTypeSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cpu_type":
				one, err := XMLCpuTypeReadOne(reader, &t, "cpu_type")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSpecialObjectsReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*SpecialObjects, error) {
	builder := NewSpecialObjectsBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "special_objects"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "blank_template":
				v, err := XMLTemplateReadOne(reader, &t, "blank_template")
				if err != nil {
					return nil, err
				}
				builder.BlankTemplate(v)
			case "root_tag":
				v, err := XMLTagReadOne(reader, &t, "root_tag")
				if err != nil {
					return nil, err
				}
				builder.RootTag(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSpecialObjectsReadMany(reader *XMLReader, start *xml.StartElement) (*SpecialObjectsSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SpecialObjectsSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "special_objects":
				one, err := XMLSpecialObjectsReadOne(reader, &t, "special_objects")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLSessionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Session, error) {
	builder := NewSessionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "session"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "console_user":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ConsoleUser(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "ip":
				v, err := XMLIpReadOne(reader, &t, "ip")
				if err != nil {
					return nil, err
				}
				builder.Ip(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "protocol":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Protocol(v)
			case "user":
				v, err := XMLUserReadOne(reader, &t, "user")
				if err != nil {
					return nil, err
				}
				builder.User(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLSessionReadMany(reader *XMLReader, start *xml.StartElement) (*SessionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result SessionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "session":
				one, err := XMLSessionReadOne(reader, &t, "session")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterBrickAdvancedDetailsReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GlusterBrickAdvancedDetails, error) {
	builder := NewGlusterBrickAdvancedDetailsBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "gluster_brick_advanced_details"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "device":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Device(v)
			case "fs_name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.FsName(v)
			case "gluster_clients":
				v, err := XMLGlusterClientReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.GlusterClients(v)
			case "instance_type":
				v, err := XMLInstanceTypeReadOne(reader, &t, "instance_type")
				if err != nil {
					return nil, err
				}
				builder.InstanceType(v)
			case "memory_pools":
				v, err := XMLGlusterMemoryPoolReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.MemoryPools(v)
			case "mnt_options":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.MntOptions(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "pid":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Pid(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vms":
				v, err := XMLVmReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "vms":
			if one.vms == nil {
				one.vms = new(VmSlice)
			}
			one.vms.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGlusterBrickAdvancedDetailsReadMany(reader *XMLReader, start *xml.StartElement) (*GlusterBrickAdvancedDetailsSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GlusterBrickAdvancedDetailsSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "gluster_brick_advanced_details":
				one, err := XMLGlusterBrickAdvancedDetailsReadOne(reader, &t, "gluster_brick_advanced_details")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLKsmReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Ksm, error) {
	builder := NewKsmBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "ksm"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Enabled(v)
			case "merge_across_nodes":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MergeAcrossNodes(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLKsmReadMany(reader *XMLReader, start *xml.StartElement) (*KsmSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result KsmSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ksm":
				one, err := XMLKsmReadOne(reader, &t, "ksm")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLValueReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Value, error) {
	builder := NewValueBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "value"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "datum":
				v, err := reader.ReadFloat64(&t)
				if err != nil {
					return nil, err
				}
				builder.Datum(v)
			case "detail":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Detail(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLValueReadMany(reader *XMLReader, start *xml.StartElement) (*ValueSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ValueSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "value":
				one, err := XMLValueReadOne(reader, &t, "value")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLCloudInitReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*CloudInit, error) {
	builder := NewCloudInitBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "cloud_init"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "authorized_keys":
				v, err := XMLAuthorizedKeyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.AuthorizedKeys(v)
			case "files":
				v, err := XMLFileReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Files(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "network_configuration":
				v, err := XMLNetworkConfigurationReadOne(reader, &t, "network_configuration")
				if err != nil {
					return nil, err
				}
				builder.NetworkConfiguration(v)
			case "regenerate_ssh_keys":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RegenerateSshKeys(v)
			case "timezone":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Timezone(v)
			case "users":
				v, err := XMLUserReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Users(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLCloudInitReadMany(reader *XMLReader, start *xml.StartElement) (*CloudInitSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result CloudInitSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "cloud_init":
				one, err := XMLCloudInitReadOne(reader, &t, "cloud_init")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLDisplayReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Display, error) {
	builder := NewDisplayBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "display"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "allow_override":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.AllowOverride(v)
			case "certificate":
				v, err := XMLCertificateReadOne(reader, &t, "certificate")
				if err != nil {
					return nil, err
				}
				builder.Certificate(v)
			case "copy_paste_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.CopyPasteEnabled(v)
			case "disconnect_action":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.DisconnectAction(v)
			case "file_transfer_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.FileTransferEnabled(v)
			case "keyboard_layout":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.KeyboardLayout(v)
			case "monitors":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Monitors(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "proxy":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Proxy(v)
			case "secure_port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.SecurePort(v)
			case "single_qxl_pci":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SingleQxlPci(v)
			case "smartcard_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.SmartcardEnabled(v)
			case "type":
				vp, err := XMLDisplayTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLDisplayReadMany(reader *XMLReader, start *xml.StartElement) (*DisplaySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result DisplaySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "display":
				one, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLStorageDomainLeaseReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*StorageDomainLease, error) {
	builder := NewStorageDomainLeaseBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "storage_domain_lease"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLStorageDomainLeaseReadMany(reader *XMLReader, start *xml.StartElement) (*StorageDomainLeaseSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result StorageDomainLeaseSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "storage_domain_lease":
				one, err := XMLStorageDomainLeaseReadOne(reader, &t, "storage_domain_lease")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLMigrationPolicyReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*MigrationPolicy, error) {
	builder := NewMigrationPolicyBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "migration_policy"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLMigrationPolicyReadMany(reader *XMLReader, start *xml.StartElement) (*MigrationPolicySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result MigrationPolicySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "migration_policy":
				one, err := XMLMigrationPolicyReadOne(reader, &t, "migration_policy")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLHostReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Host, error) {
	builder := NewHostBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "host"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "address":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Address(v)
			case "affinity_labels":
				v, err := XMLAffinityLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.AffinityLabels(v)
			case "agents":
				v, err := XMLAgentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Agents(v)
			case "auto_numa_status":
				vp, err := XMLAutoNumaStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.AutoNumaStatus(v)
			case "certificate":
				v, err := XMLCertificateReadOne(reader, &t, "certificate")
				if err != nil {
					return nil, err
				}
				builder.Certificate(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "cpu":
				v, err := XMLCpuReadOne(reader, &t, "cpu")
				if err != nil {
					return nil, err
				}
				builder.Cpu(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "device_passthrough":
				v, err := XMLHostDevicePassthroughReadOne(reader, &t, "device_passthrough")
				if err != nil {
					return nil, err
				}
				builder.DevicePassthrough(v)
			case "devices":
				v, err := XMLHostDeviceReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Devices(v)
			case "display":
				v, err := XMLDisplayReadOne(reader, &t, "display")
				if err != nil {
					return nil, err
				}
				builder.Display(v)
			case "external_host_provider":
				v, err := XMLExternalHostProviderReadOne(reader, &t, "external_host_provider")
				if err != nil {
					return nil, err
				}
				builder.ExternalHostProvider(v)
			case "external_network_provider_configurations":
				v, err := XMLExternalNetworkProviderConfigurationReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ExternalNetworkProviderConfigurations(v)
			case "external_status":
				vp, err := XMLExternalStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.ExternalStatus(v)
			case "hardware_information":
				v, err := XMLHardwareInformationReadOne(reader, &t, "hardware_information")
				if err != nil {
					return nil, err
				}
				builder.HardwareInformation(v)
			case "hooks":
				v, err := XMLHookReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Hooks(v)
			case "hosted_engine":
				v, err := XMLHostedEngineReadOne(reader, &t, "hosted_engine")
				if err != nil {
					return nil, err
				}
				builder.HostedEngine(v)
			case "iscsi":
				v, err := XMLIscsiDetailsReadOne(reader, &t, "iscsi")
				if err != nil {
					return nil, err
				}
				builder.Iscsi(v)
			case "katello_errata":
				v, err := XMLKatelloErratumReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.KatelloErrata(v)
			case "kdump_status":
				vp, err := XMLKdumpStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.KdumpStatus(v)
			case "ksm":
				v, err := XMLKsmReadOne(reader, &t, "ksm")
				if err != nil {
					return nil, err
				}
				builder.Ksm(v)
			case "libvirt_version":
				v, err := XMLVersionReadOne(reader, &t, "libvirt_version")
				if err != nil {
					return nil, err
				}
				builder.LibvirtVersion(v)
			case "max_scheduling_memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.MaxSchedulingMemory(v)
			case "memory":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Memory(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "network_attachments":
				v, err := XMLNetworkAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NetworkAttachments(v)
			case "network_operation_in_progress":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.NetworkOperationInProgress(v)
			case "nics":
				v, err := XMLHostNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Nics(v)
			case "host_numa_nodes":
				v, err := XMLNumaNodeReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.NumaNodes(v)
			case "numa_supported":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.NumaSupported(v)
			case "os":
				v, err := XMLOperatingSystemReadOne(reader, &t, "os")
				if err != nil {
					return nil, err
				}
				builder.Os(v)
			case "override_iptables":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.OverrideIptables(v)
			case "permissions":
				v, err := XMLPermissionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Permissions(v)
			case "port":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Port(v)
			case "power_management":
				v, err := XMLPowerManagementReadOne(reader, &t, "power_management")
				if err != nil {
					return nil, err
				}
				builder.PowerManagement(v)
			case "protocol":
				vp, err := XMLHostProtocolReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Protocol(v)
			case "root_password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.RootPassword(v)
			case "se_linux":
				v, err := XMLSeLinuxReadOne(reader, &t, "se_linux")
				if err != nil {
					return nil, err
				}
				builder.SeLinux(v)
			case "spm":
				v, err := XMLSpmReadOne(reader, &t, "spm")
				if err != nil {
					return nil, err
				}
				builder.Spm(v)
			case "ssh":
				v, err := XMLSshReadOne(reader, &t, "ssh")
				if err != nil {
					return nil, err
				}
				builder.Ssh(v)
			case "statistics":
				v, err := XMLStatisticReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Statistics(v)
			case "status":
				vp, err := XMLHostStatusReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "status_detail":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.StatusDetail(v)
			case "storage_connection_extensions":
				v, err := XMLStorageConnectionExtensionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.StorageConnectionExtensions(v)
			case "storages":
				v, err := XMLHostStorageReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Storages(v)
			case "summary":
				v, err := XMLVmSummaryReadOne(reader, &t, "summary")
				if err != nil {
					return nil, err
				}
				builder.Summary(v)
			case "tags":
				v, err := XMLTagReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Tags(v)
			case "transparent_hugepages":
				v, err := XMLTransparentHugePagesReadOne(reader, &t, "transparent_hugepages")
				if err != nil {
					return nil, err
				}
				builder.TransparentHugePages(v)
			case "type":
				vp, err := XMLHostTypeReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "unmanaged_networks":
				v, err := XMLUnmanagedNetworkReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.UnmanagedNetworks(v)
			case "update_available":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UpdateAvailable(v)
			case "version":
				v, err := XMLVersionReadOne(reader, &t, "version")
				if err != nil {
					return nil, err
				}
				builder.Version(v)
			case "vgpu_placement":
				vp, err := XMLVgpuPlacementReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.VgpuPlacement(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "affinitylabels":
			if one.affinityLabels == nil {
				one.affinityLabels = new(AffinityLabelSlice)
			}
			one.affinityLabels.href = link.href
		case "agents":
			if one.agents == nil {
				one.agents = new(AgentSlice)
			}
			one.agents.href = link.href
		case "devices":
			if one.devices == nil {
				one.devices = new(HostDeviceSlice)
			}
			one.devices.href = link.href
		case "externalnetworkproviderconfigurations":
			if one.externalNetworkProviderConfigurations == nil {
				one.externalNetworkProviderConfigurations = new(ExternalNetworkProviderConfigurationSlice)
			}
			one.externalNetworkProviderConfigurations.href = link.href
		case "hooks":
			if one.hooks == nil {
				one.hooks = new(HookSlice)
			}
			one.hooks.href = link.href
		case "katelloerrata":
			if one.katelloErrata == nil {
				one.katelloErrata = new(KatelloErratumSlice)
			}
			one.katelloErrata.href = link.href
		case "networkattachments":
			if one.networkAttachments == nil {
				one.networkAttachments = new(NetworkAttachmentSlice)
			}
			one.networkAttachments.href = link.href
		case "nics":
			if one.nics == nil {
				one.nics = new(HostNicSlice)
			}
			one.nics.href = link.href
		case "numanodes":
			if one.numaNodes == nil {
				one.numaNodes = new(NumaNodeSlice)
			}
			one.numaNodes.href = link.href
		case "permissions":
			if one.permissions == nil {
				one.permissions = new(PermissionSlice)
			}
			one.permissions.href = link.href
		case "statistics":
			if one.statistics == nil {
				one.statistics = new(StatisticSlice)
			}
			one.statistics.href = link.href
		case "storageconnectionextensions":
			if one.storageConnectionExtensions == nil {
				one.storageConnectionExtensions = new(StorageConnectionExtensionSlice)
			}
			one.storageConnectionExtensions.href = link.href
		case "storages":
			if one.storages == nil {
				one.storages = new(HostStorageSlice)
			}
			one.storages.href = link.href
		case "tags":
			if one.tags == nil {
				one.tags = new(TagSlice)
			}
			one.tags.href = link.href
		case "unmanagednetworks":
			if one.unmanagedNetworks == nil {
				one.unmanagedNetworks = new(UnmanagedNetworkSlice)
			}
			one.unmanagedNetworks.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLHostReadMany(reader *XMLReader, start *xml.StartElement) (*HostSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result HostSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "host":
				one, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLNetworkAttachmentReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*NetworkAttachment, error) {
	builder := NewNetworkAttachmentBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "network_attachment"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "dns_resolver_configuration":
				v, err := XMLDnsResolverConfigurationReadOne(reader, &t, "dns_resolver_configuration")
				if err != nil {
					return nil, err
				}
				builder.DnsResolverConfiguration(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "host_nic":
				v, err := XMLHostNicReadOne(reader, &t, "host_nic")
				if err != nil {
					return nil, err
				}
				builder.HostNic(v)
			case "in_sync":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.InSync(v)
			case "ip_address_assignments":
				v, err := XMLIpAddressAssignmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.IpAddressAssignments(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "network":
				v, err := XMLNetworkReadOne(reader, &t, "network")
				if err != nil {
					return nil, err
				}
				builder.Network(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "qos":
				v, err := XMLQosReadOne(reader, &t, "qos")
				if err != nil {
					return nil, err
				}
				builder.Qos(v)
			case "reported_configurations":
				v, err := XMLReportedConfigurationReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ReportedConfigurations(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLNetworkAttachmentReadMany(reader *XMLReader, start *xml.StartElement) (*NetworkAttachmentSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result NetworkAttachmentSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "network_attachment":
				one, err := XMLNetworkAttachmentReadOne(reader, &t, "network_attachment")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLFilterReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Filter, error) {
	builder := NewFilterBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "filter"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "position":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Position(v)
			case "scheduling_policy_unit":
				v, err := XMLSchedulingPolicyUnitReadOne(reader, &t, "scheduling_policy_unit")
				if err != nil {
					return nil, err
				}
				builder.SchedulingPolicyUnit(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLFilterReadMany(reader *XMLReader, start *xml.StartElement) (*FilterSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result FilterSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "filter":
				one, err := XMLFilterReadOne(reader, &t, "filter")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLIscsiBondReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*IscsiBond, error) {
	builder := NewIscsiBondBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "iscsi_bond"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "networks":
				v, err := XMLNetworkReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Networks(v)
			case "storage_connections":
				v, err := XMLStorageConnectionReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.StorageConnections(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		case "networks":
			if one.networks == nil {
				one.networks = new(NetworkSlice)
			}
			one.networks.href = link.href
		case "storageconnections":
			if one.storageConnections == nil {
				one.storageConnections = new(StorageConnectionSlice)
			}
			one.storageConnections.href = link.href
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLIscsiBondReadMany(reader *XMLReader, start *xml.StartElement) (*IscsiBondSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result IscsiBondSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "iscsi_bond":
				one, err := XMLIscsiBondReadOne(reader, &t, "iscsi_bond")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLLinkLayerDiscoveryProtocolElementReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*LinkLayerDiscoveryProtocolElement, error) {
	builder := NewLinkLayerDiscoveryProtocolElementBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "link_layer_discovery_protocol_element"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "oui":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Oui(v)
			case "properties":
				v, err := XMLPropertyReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Properties(v)
			case "subtype":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Subtype(v)
			case "type":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Type(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLLinkLayerDiscoveryProtocolElementReadMany(reader *XMLReader, start *xml.StartElement) (*LinkLayerDiscoveryProtocolElementSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result LinkLayerDiscoveryProtocolElementSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "link_layer_discovery_protocol_element":
				one, err := XMLLinkLayerDiscoveryProtocolElementReadOne(reader, &t, "link_layer_discovery_protocol_element")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLApiSummaryReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*ApiSummary, error) {
	builder := NewApiSummaryBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "api_summary"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "hosts":
				v, err := XMLApiSummaryItemReadOne(reader, &t, "hosts")
				if err != nil {
					return nil, err
				}
				builder.Hosts(v)
			case "storage_domains":
				v, err := XMLApiSummaryItemReadOne(reader, &t, "storage_domains")
				if err != nil {
					return nil, err
				}
				builder.StorageDomains(v)
			case "users":
				v, err := XMLApiSummaryItemReadOne(reader, &t, "users")
				if err != nil {
					return nil, err
				}
				builder.Users(v)
			case "vms":
				v, err := XMLApiSummaryItemReadOne(reader, &t, "vms")
				if err != nil {
					return nil, err
				}
				builder.Vms(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLApiSummaryReadMany(reader *XMLReader, start *xml.StartElement) (*ApiSummarySlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ApiSummarySlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "api_summary":
				one, err := XMLApiSummaryReadOne(reader, &t, "api_summary")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLFaultReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Fault, error) {
	builder := NewFaultBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "fault"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "detail":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Detail(v)
			case "reason":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Reason(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLFaultReadMany(reader *XMLReader, start *xml.StartElement) (*FaultSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result FaultSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "fault":
				one, err := XMLFaultReadOne(reader, &t, "fault")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGracePeriodReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*GracePeriod, error) {
	builder := NewGracePeriodBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "grace_period"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "expiry":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Expiry(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLGracePeriodReadMany(reader *XMLReader, start *xml.StartElement) (*GracePeriodSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result GracePeriodSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "grace_period":
				one, err := XMLGracePeriodReadOne(reader, &t, "grace_period")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLActionReadOne(reader *XMLReader, start *xml.StartElement, expectedTag string) (*Action, error) {
	builder := NewActionBuilder()
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	if expectedTag == "" {
		expectedTag = "action"
	}
	if start.Name.Local != expectedTag {
		return nil, XMLTagNotMatchError{start.Name.Local, expectedTag}
	}
	// Process the attributes
	for _, attr := range start.Attr {
		name := attr.Name.Local
		value := attr.Value
		switch name {
		case "id":
			builder.Id(value)
		case "href":
			builder.Href(value)
		}
	}
	var links []Link
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "activate":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Activate(v)
			case "allow_partial_import":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.AllowPartialImport(v)
			case "async":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Async(v)
			case "attachment":
				v, err := XMLDiskAttachmentReadOne(reader, &t, "attachment")
				if err != nil {
					return nil, err
				}
				builder.Attachment(v)
			case "authorized_key":
				v, err := XMLAuthorizedKeyReadOne(reader, &t, "authorized_key")
				if err != nil {
					return nil, err
				}
				builder.AuthorizedKey(v)
			case "bricks":
				v, err := XMLGlusterBrickReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Bricks(v)
			case "certificates":
				v, err := XMLCertificateReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Certificates(v)
			case "check_connectivity":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.CheckConnectivity(v)
			case "clone":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Clone(v)
			case "clone_permissions":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ClonePermissions(v)
			case "cluster":
				v, err := XMLClusterReadOne(reader, &t, "cluster")
				if err != nil {
					return nil, err
				}
				builder.Cluster(v)
			case "collapse_snapshots":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.CollapseSnapshots(v)
			case "comment":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Comment(v)
			case "commit_on_success":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.CommitOnSuccess(v)
			case "connection":
				v, err := XMLStorageConnectionReadOne(reader, &t, "connection")
				if err != nil {
					return nil, err
				}
				builder.Connection(v)
			case "connectivity_timeout":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.ConnectivityTimeout(v)
			case "data_center":
				v, err := XMLDataCenterReadOne(reader, &t, "data_center")
				if err != nil {
					return nil, err
				}
				builder.DataCenter(v)
			case "deploy_hosted_engine":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DeployHostedEngine(v)
			case "description":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Description(v)
			case "details":
				v, err := XMLGlusterVolumeProfileDetailsReadOne(reader, &t, "details")
				if err != nil {
					return nil, err
				}
				builder.Details(v)
			case "directory":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Directory(v)
			case "discard_snapshots":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.DiscardSnapshots(v)
			case "discovered_targets":
				v, err := XMLIscsiDetailsReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.DiscoveredTargets(v)
			case "disk":
				v, err := XMLDiskReadOne(reader, &t, "disk")
				if err != nil {
					return nil, err
				}
				builder.Disk(v)
			case "disk_profile":
				v, err := XMLDiskProfileReadOne(reader, &t, "disk_profile")
				if err != nil {
					return nil, err
				}
				builder.DiskProfile(v)
			case "disks":
				v, err := XMLDiskReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.Disks(v)
			case "exclusive":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Exclusive(v)
			case "fault":
				v, err := XMLFaultReadOne(reader, &t, "fault")
				if err != nil {
					return nil, err
				}
				builder.Fault(v)
			case "fence_type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.FenceType(v)
			case "filename":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Filename(v)
			case "filter":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Filter(v)
			case "fix_layout":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.FixLayout(v)
			case "force":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Force(v)
			case "grace_period":
				v, err := XMLGracePeriodReadOne(reader, &t, "grace_period")
				if err != nil {
					return nil, err
				}
				builder.GracePeriod(v)
			case "host":
				v, err := XMLHostReadOne(reader, &t, "host")
				if err != nil {
					return nil, err
				}
				builder.Host(v)
			case "image":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Image(v)
			case "image_transfer":
				v, err := XMLImageTransferReadOne(reader, &t, "image_transfer")
				if err != nil {
					return nil, err
				}
				builder.ImageTransfer(v)
			case "import_as_template":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ImportAsTemplate(v)
			case "is_attached":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.IsAttached(v)
			case "iscsi":
				v, err := XMLIscsiDetailsReadOne(reader, &t, "iscsi")
				if err != nil {
					return nil, err
				}
				builder.Iscsi(v)
			case "iscsi_targets":
				v, err := reader.ReadStrings(&t)
				if err != nil {
					return nil, err
				}
				builder.IscsiTargets(v)
			case "job":
				v, err := XMLJobReadOne(reader, &t, "job")
				if err != nil {
					return nil, err
				}
				builder.Job(v)
			case "lease":
				v, err := XMLStorageDomainLeaseReadOne(reader, &t, "lease")
				if err != nil {
					return nil, err
				}
				builder.Lease(v)
			case "logical_units":
				v, err := XMLLogicalUnitReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.LogicalUnits(v)
			case "maintenance_after_restart":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MaintenanceAfterRestart(v)
			case "maintenance_enabled":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MaintenanceEnabled(v)
			case "migrate_vms_in_affinity_closure":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.MigrateVmsInAffinityClosure(v)
			case "modified_bonds":
				v, err := XMLHostNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ModifiedBonds(v)
			case "modified_labels":
				v, err := XMLNetworkLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ModifiedLabels(v)
			case "modified_network_attachments":
				v, err := XMLNetworkAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.ModifiedNetworkAttachments(v)
			case "name":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Name(v)
			case "option":
				v, err := XMLOptionReadOne(reader, &t, "option")
				if err != nil {
					return nil, err
				}
				builder.Option(v)
			case "pause":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Pause(v)
			case "permission":
				v, err := XMLPermissionReadOne(reader, &t, "permission")
				if err != nil {
					return nil, err
				}
				builder.Permission(v)
			case "power_management":
				v, err := XMLPowerManagementReadOne(reader, &t, "power_management")
				if err != nil {
					return nil, err
				}
				builder.PowerManagement(v)
			case "proxy_ticket":
				v, err := XMLProxyTicketReadOne(reader, &t, "proxy_ticket")
				if err != nil {
					return nil, err
				}
				builder.ProxyTicket(v)
			case "quota":
				v, err := XMLQuotaReadOne(reader, &t, "quota")
				if err != nil {
					return nil, err
				}
				builder.Quota(v)
			case "reason":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Reason(v)
			case "reassign_bad_macs":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.ReassignBadMacs(v)
			case "reboot":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Reboot(v)
			case "registration_configuration":
				v, err := XMLRegistrationConfigurationReadOne(reader, &t, "registration_configuration")
				if err != nil {
					return nil, err
				}
				builder.RegistrationConfiguration(v)
			case "remote_viewer_connection_file":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.RemoteViewerConnectionFile(v)
			case "removed_bonds":
				v, err := XMLHostNicReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.RemovedBonds(v)
			case "removed_labels":
				v, err := XMLNetworkLabelReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.RemovedLabels(v)
			case "removed_network_attachments":
				v, err := XMLNetworkAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.RemovedNetworkAttachments(v)
			case "resolution_type":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.ResolutionType(v)
			case "restore_memory":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.RestoreMemory(v)
			case "root_password":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.RootPassword(v)
			case "seal":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Seal(v)
			case "snapshot":
				v, err := XMLSnapshotReadOne(reader, &t, "snapshot")
				if err != nil {
					return nil, err
				}
				builder.Snapshot(v)
			case "source_host":
				v, err := XMLHostReadOne(reader, &t, "source_host")
				if err != nil {
					return nil, err
				}
				builder.SourceHost(v)
			case "ssh":
				v, err := XMLSshReadOne(reader, &t, "ssh")
				if err != nil {
					return nil, err
				}
				builder.Ssh(v)
			case "status":
				v, err := reader.ReadString(&t)
				if err != nil {
					return nil, err
				}
				builder.Status(v)
			case "stop_gluster_service":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.StopGlusterService(v)
			case "storage_domain":
				v, err := XMLStorageDomainReadOne(reader, &t, "storage_domain")
				if err != nil {
					return nil, err
				}
				builder.StorageDomain(v)
			case "storage_domains":
				v, err := XMLStorageDomainReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.StorageDomains(v)
			case "succeeded":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Succeeded(v)
			case "synchronized_network_attachments":
				v, err := XMLNetworkAttachmentReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.SynchronizedNetworkAttachments(v)
			case "template":
				v, err := XMLTemplateReadOne(reader, &t, "template")
				if err != nil {
					return nil, err
				}
				builder.Template(v)
			case "ticket":
				v, err := XMLTicketReadOne(reader, &t, "ticket")
				if err != nil {
					return nil, err
				}
				builder.Ticket(v)
			case "timeout":
				v, err := reader.ReadInt64(&t)
				if err != nil {
					return nil, err
				}
				builder.Timeout(v)
			case "undeploy_hosted_engine":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UndeployHostedEngine(v)
			case "upgrade_action":
				vp, err := XMLClusterUpgradeActionReadOne(reader, &t)
				v := *vp
				if err != nil {
					return nil, err
				}
				builder.UpgradeAction(v)
			case "use_cloud_init":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UseCloudInit(v)
			case "use_ignition":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UseIgnition(v)
			case "use_initialization":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UseInitialization(v)
			case "use_sysprep":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.UseSysprep(v)
			case "virtual_functions_configuration":
				v, err := XMLHostNicVirtualFunctionsConfigurationReadOne(reader, &t, "virtual_functions_configuration")
				if err != nil {
					return nil, err
				}
				builder.VirtualFunctionsConfiguration(v)
			case "vm":
				v, err := XMLVmReadOne(reader, &t, "vm")
				if err != nil {
					return nil, err
				}
				builder.Vm(v)
			case "vnic_profile_mappings":
				v, err := XMLVnicProfileMappingReadMany(reader, &t)
				if err != nil {
					return nil, err
				}
				builder.VnicProfileMappings(v)
			case "volatile":
				v, err := reader.ReadBool(&t)
				if err != nil {
					return nil, err
				}
				builder.Volatile(v)
			case "link":
				var rel, href string
				for _, attr := range t.Attr {
					name := attr.Name.Local
					value := attr.Value
					switch name {
					case "href":
						href = value
					case "rel":
						rel = value
					}
				}
				if rel != "" && href != "" {
					links = append(links, Link{&href, &rel})
				}
				// <link> just has attributes, so must skip manually
				reader.Skip()
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	one, err := builder.Build()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		switch *link.rel {
		} // end of switch
	} // end of for-links
	return one, nil
}

func XMLActionReadMany(reader *XMLReader, start *xml.StartElement) (*ActionSlice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var result ActionSlice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "action":
				one, err := XMLActionReadOne(reader, &t, "action")
				if err != nil {
					return nil, err
				}
				if one != nil {
					result.slice = append(result.slice, one)
				}
			default:
				reader.Skip()
			}
		case xml.EndElement:
			depth--
		}
	}
	return &result, nil
}

func XMLGlusterHookStatusReadOne(reader *XMLReader, start *xml.StartElement) (*GlusterHookStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(GlusterHookStatus)
	*result = GlusterHookStatus(s)
	return result, nil
}

func XMLGlusterHookStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]GlusterHookStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []GlusterHookStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, GlusterHookStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLGlusterBrickStatusReadOne(reader *XMLReader, start *xml.StartElement) (*GlusterBrickStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(GlusterBrickStatus)
	*result = GlusterBrickStatus(s)
	return result, nil
}

func XMLGlusterBrickStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]GlusterBrickStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []GlusterBrickStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, GlusterBrickStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLGlusterVolumeStatusReadOne(reader *XMLReader, start *xml.StartElement) (*GlusterVolumeStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(GlusterVolumeStatus)
	*result = GlusterVolumeStatus(s)
	return result, nil
}

func XMLGlusterVolumeStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]GlusterVolumeStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []GlusterVolumeStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, GlusterVolumeStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNicInterfaceReadOne(reader *XMLReader, start *xml.StartElement) (*NicInterface, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NicInterface)
	*result = NicInterface(s)
	return result, nil
}

func XMLNicInterfaceReadMany(reader *XMLReader, start *xml.StartElement) ([]NicInterface, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NicInterface
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NicInterface(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLGraphicsTypeReadOne(reader *XMLReader, start *xml.StartElement) (*GraphicsType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(GraphicsType)
	*result = GraphicsType(s)
	return result, nil
}

func XMLGraphicsTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]GraphicsType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []GraphicsType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, GraphicsType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLSnapshotStatusReadOne(reader *XMLReader, start *xml.StartElement) (*SnapshotStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(SnapshotStatus)
	*result = SnapshotStatus(s)
	return result, nil
}

func XMLSnapshotStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]SnapshotStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []SnapshotStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, SnapshotStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLImageFileTypeReadOne(reader *XMLReader, start *xml.StartElement) (*ImageFileType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ImageFileType)
	*result = ImageFileType(s)
	return result, nil
}

func XMLImageFileTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]ImageFileType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ImageFileType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ImageFileType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDiskStorageTypeReadOne(reader *XMLReader, start *xml.StartElement) (*DiskStorageType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DiskStorageType)
	*result = DiskStorageType(s)
	return result, nil
}

func XMLDiskStorageTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]DiskStorageType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DiskStorageType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DiskStorageType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLSpmStatusReadOne(reader *XMLReader, start *xml.StartElement) (*SpmStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(SpmStatus)
	*result = SpmStatus(s)
	return result, nil
}

func XMLSpmStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]SpmStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []SpmStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, SpmStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLImageTransferDirectionReadOne(reader *XMLReader, start *xml.StartElement) (*ImageTransferDirection, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ImageTransferDirection)
	*result = ImageTransferDirection(s)
	return result, nil
}

func XMLImageTransferDirectionReadMany(reader *XMLReader, start *xml.StartElement) ([]ImageTransferDirection, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ImageTransferDirection
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ImageTransferDirection(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLAccessProtocolReadOne(reader *XMLReader, start *xml.StartElement) (*AccessProtocol, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(AccessProtocol)
	*result = AccessProtocol(s)
	return result, nil
}

func XMLAccessProtocolReadMany(reader *XMLReader, start *xml.StartElement) ([]AccessProtocol, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []AccessProtocol
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, AccessProtocol(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLStorageFormatReadOne(reader *XMLReader, start *xml.StartElement) (*StorageFormat, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(StorageFormat)
	*result = StorageFormat(s)
	return result, nil
}

func XMLStorageFormatReadMany(reader *XMLReader, start *xml.StartElement) ([]StorageFormat, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []StorageFormat
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, StorageFormat(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLVnicPassThroughModeReadOne(reader *XMLReader, start *xml.StartElement) (*VnicPassThroughMode, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(VnicPassThroughMode)
	*result = VnicPassThroughMode(s)
	return result, nil
}

func XMLVnicPassThroughModeReadMany(reader *XMLReader, start *xml.StartElement) ([]VnicPassThroughMode, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []VnicPassThroughMode
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, VnicPassThroughMode(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNicStatusReadOne(reader *XMLReader, start *xml.StartElement) (*NicStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NicStatus)
	*result = NicStatus(s)
	return result, nil
}

func XMLNicStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]NicStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NicStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NicStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLLogMaxMemoryUsedThresholdTypeReadOne(reader *XMLReader, start *xml.StartElement) (*LogMaxMemoryUsedThresholdType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(LogMaxMemoryUsedThresholdType)
	*result = LogMaxMemoryUsedThresholdType(s)
	return result, nil
}

func XMLLogMaxMemoryUsedThresholdTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]LogMaxMemoryUsedThresholdType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []LogMaxMemoryUsedThresholdType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, LogMaxMemoryUsedThresholdType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNetworkStatusReadOne(reader *XMLReader, start *xml.StartElement) (*NetworkStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NetworkStatus)
	*result = NetworkStatus(s)
	return result, nil
}

func XMLNetworkStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]NetworkStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NetworkStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NetworkStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDiskStatusReadOne(reader *XMLReader, start *xml.StartElement) (*DiskStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DiskStatus)
	*result = DiskStatus(s)
	return result, nil
}

func XMLDiskStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]DiskStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DiskStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DiskStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLPolicyUnitTypeReadOne(reader *XMLReader, start *xml.StartElement) (*PolicyUnitType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(PolicyUnitType)
	*result = PolicyUnitType(s)
	return result, nil
}

func XMLPolicyUnitTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]PolicyUnitType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []PolicyUnitType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, PolicyUnitType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLHookStageReadOne(reader *XMLReader, start *xml.StartElement) (*HookStage, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(HookStage)
	*result = HookStage(s)
	return result, nil
}

func XMLHookStageReadMany(reader *XMLReader, start *xml.StartElement) ([]HookStage, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []HookStage
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, HookStage(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNotifiableEventReadOne(reader *XMLReader, start *xml.StartElement) (*NotifiableEvent, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NotifiableEvent)
	*result = NotifiableEvent(s)
	return result, nil
}

func XMLNotifiableEventReadMany(reader *XMLReader, start *xml.StartElement) ([]NotifiableEvent, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NotifiableEvent
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NotifiableEvent(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLCreationStatusReadOne(reader *XMLReader, start *xml.StartElement) (*CreationStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(CreationStatus)
	*result = CreationStatus(s)
	return result, nil
}

func XMLCreationStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]CreationStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []CreationStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, CreationStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLConfigurationTypeReadOne(reader *XMLReader, start *xml.StartElement) (*ConfigurationType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ConfigurationType)
	*result = ConfigurationType(s)
	return result, nil
}

func XMLConfigurationTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]ConfigurationType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ConfigurationType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ConfigurationType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLVgpuPlacementReadOne(reader *XMLReader, start *xml.StartElement) (*VgpuPlacement, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(VgpuPlacement)
	*result = VgpuPlacement(s)
	return result, nil
}

func XMLVgpuPlacementReadMany(reader *XMLReader, start *xml.StartElement) ([]VgpuPlacement, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []VgpuPlacement
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, VgpuPlacement(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLOpenstackVolumeAuthenticationKeyUsageTypeReadOne(reader *XMLReader, start *xml.StartElement) (*OpenstackVolumeAuthenticationKeyUsageType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(OpenstackVolumeAuthenticationKeyUsageType)
	*result = OpenstackVolumeAuthenticationKeyUsageType(s)
	return result, nil
}

func XMLOpenstackVolumeAuthenticationKeyUsageTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]OpenstackVolumeAuthenticationKeyUsageType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []OpenstackVolumeAuthenticationKeyUsageType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, OpenstackVolumeAuthenticationKeyUsageType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLTransportTypeReadOne(reader *XMLReader, start *xml.StartElement) (*TransportType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(TransportType)
	*result = TransportType(s)
	return result, nil
}

func XMLTransportTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]TransportType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []TransportType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, TransportType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLExternalVmProviderTypeReadOne(reader *XMLReader, start *xml.StartElement) (*ExternalVmProviderType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ExternalVmProviderType)
	*result = ExternalVmProviderType(s)
	return result, nil
}

func XMLExternalVmProviderTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]ExternalVmProviderType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ExternalVmProviderType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ExternalVmProviderType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLMessageBrokerTypeReadOne(reader *XMLReader, start *xml.StartElement) (*MessageBrokerType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(MessageBrokerType)
	*result = MessageBrokerType(s)
	return result, nil
}

func XMLMessageBrokerTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]MessageBrokerType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []MessageBrokerType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, MessageBrokerType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLBootDeviceReadOne(reader *XMLReader, start *xml.StartElement) (*BootDevice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(BootDevice)
	*result = BootDevice(s)
	return result, nil
}

func XMLBootDeviceReadMany(reader *XMLReader, start *xml.StartElement) ([]BootDevice, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []BootDevice
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, BootDevice(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLStatisticKindReadOne(reader *XMLReader, start *xml.StartElement) (*StatisticKind, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(StatisticKind)
	*result = StatisticKind(s)
	return result, nil
}

func XMLStatisticKindReadMany(reader *XMLReader, start *xml.StartElement) ([]StatisticKind, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []StatisticKind
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, StatisticKind(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLJobStatusReadOne(reader *XMLReader, start *xml.StartElement) (*JobStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(JobStatus)
	*result = JobStatus(s)
	return result, nil
}

func XMLJobStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]JobStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []JobStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, JobStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDiskBackupReadOne(reader *XMLReader, start *xml.StartElement) (*DiskBackup, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DiskBackup)
	*result = DiskBackup(s)
	return result, nil
}

func XMLDiskBackupReadMany(reader *XMLReader, start *xml.StartElement) ([]DiskBackup, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DiskBackup
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DiskBackup(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLQcowVersionReadOne(reader *XMLReader, start *xml.StartElement) (*QcowVersion, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(QcowVersion)
	*result = QcowVersion(s)
	return result, nil
}

func XMLQcowVersionReadMany(reader *XMLReader, start *xml.StartElement) ([]QcowVersion, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []QcowVersion
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, QcowVersion(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLFenceTypeReadOne(reader *XMLReader, start *xml.StartElement) (*FenceType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(FenceType)
	*result = FenceType(s)
	return result, nil
}

func XMLFenceTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]FenceType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []FenceType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, FenceType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLPowerManagementStatusReadOne(reader *XMLReader, start *xml.StartElement) (*PowerManagementStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(PowerManagementStatus)
	*result = PowerManagementStatus(s)
	return result, nil
}

func XMLPowerManagementStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]PowerManagementStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []PowerManagementStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, PowerManagementStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLQosTypeReadOne(reader *XMLReader, start *xml.StartElement) (*QosType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(QosType)
	*result = QosType(s)
	return result, nil
}

func XMLQosTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]QosType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []QosType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, QosType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLStepEnumReadOne(reader *XMLReader, start *xml.StartElement) (*StepEnum, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(StepEnum)
	*result = StepEnum(s)
	return result, nil
}

func XMLStepEnumReadMany(reader *XMLReader, start *xml.StartElement) ([]StepEnum, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []StepEnum
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, StepEnum(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLSsoMethodReadOne(reader *XMLReader, start *xml.StartElement) (*SsoMethod, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(SsoMethod)
	*result = SsoMethod(s)
	return result, nil
}

func XMLSsoMethodReadMany(reader *XMLReader, start *xml.StartElement) ([]SsoMethod, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []SsoMethod
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, SsoMethod(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLSnapshotTypeReadOne(reader *XMLReader, start *xml.StartElement) (*SnapshotType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(SnapshotType)
	*result = SnapshotType(s)
	return result, nil
}

func XMLSnapshotTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]SnapshotType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []SnapshotType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, SnapshotType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLGlusterStateReadOne(reader *XMLReader, start *xml.StartElement) (*GlusterState, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(GlusterState)
	*result = GlusterState(s)
	return result, nil
}

func XMLGlusterStateReadMany(reader *XMLReader, start *xml.StartElement) ([]GlusterState, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []GlusterState
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, GlusterState(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNumaTuneModeReadOne(reader *XMLReader, start *xml.StartElement) (*NumaTuneMode, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NumaTuneMode)
	*result = NumaTuneMode(s)
	return result, nil
}

func XMLNumaTuneModeReadMany(reader *XMLReader, start *xml.StartElement) ([]NumaTuneMode, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NumaTuneMode
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NumaTuneMode(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLMigrationBandwidthAssignmentMethodReadOne(reader *XMLReader, start *xml.StartElement) (*MigrationBandwidthAssignmentMethod, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(MigrationBandwidthAssignmentMethod)
	*result = MigrationBandwidthAssignmentMethod(s)
	return result, nil
}

func XMLMigrationBandwidthAssignmentMethodReadMany(reader *XMLReader, start *xml.StartElement) ([]MigrationBandwidthAssignmentMethod, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []MigrationBandwidthAssignmentMethod
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, MigrationBandwidthAssignmentMethod(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLBackupPhaseReadOne(reader *XMLReader, start *xml.StartElement) (*BackupPhase, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(BackupPhase)
	*result = BackupPhase(s)
	return result, nil
}

func XMLBackupPhaseReadMany(reader *XMLReader, start *xml.StartElement) ([]BackupPhase, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []BackupPhase
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, BackupPhase(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLResolutionTypeReadOne(reader *XMLReader, start *xml.StartElement) (*ResolutionType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ResolutionType)
	*result = ResolutionType(s)
	return result, nil
}

func XMLResolutionTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]ResolutionType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ResolutionType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ResolutionType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLVmDeviceTypeReadOne(reader *XMLReader, start *xml.StartElement) (*VmDeviceType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(VmDeviceType)
	*result = VmDeviceType(s)
	return result, nil
}

func XMLVmDeviceTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]VmDeviceType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []VmDeviceType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, VmDeviceType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLIpVersionReadOne(reader *XMLReader, start *xml.StartElement) (*IpVersion, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(IpVersion)
	*result = IpVersion(s)
	return result, nil
}

func XMLIpVersionReadMany(reader *XMLReader, start *xml.StartElement) ([]IpVersion, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []IpVersion
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, IpVersion(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLValueTypeReadOne(reader *XMLReader, start *xml.StartElement) (*ValueType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ValueType)
	*result = ValueType(s)
	return result, nil
}

func XMLValueTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]ValueType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ValueType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ValueType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLSshAuthenticationMethodReadOne(reader *XMLReader, start *xml.StartElement) (*SshAuthenticationMethod, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(SshAuthenticationMethod)
	*result = SshAuthenticationMethod(s)
	return result, nil
}

func XMLSshAuthenticationMethodReadMany(reader *XMLReader, start *xml.StartElement) ([]SshAuthenticationMethod, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []SshAuthenticationMethod
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, SshAuthenticationMethod(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLKdumpStatusReadOne(reader *XMLReader, start *xml.StartElement) (*KdumpStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(KdumpStatus)
	*result = KdumpStatus(s)
	return result, nil
}

func XMLKdumpStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]KdumpStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []KdumpStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, KdumpStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLTemplateStatusReadOne(reader *XMLReader, start *xml.StartElement) (*TemplateStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(TemplateStatus)
	*result = TemplateStatus(s)
	return result, nil
}

func XMLTemplateStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]TemplateStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []TemplateStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, TemplateStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNetworkUsageReadOne(reader *XMLReader, start *xml.StartElement) (*NetworkUsage, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NetworkUsage)
	*result = NetworkUsage(s)
	return result, nil
}

func XMLNetworkUsageReadMany(reader *XMLReader, start *xml.StartElement) ([]NetworkUsage, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NetworkUsage
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NetworkUsage(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLStepStatusReadOne(reader *XMLReader, start *xml.StartElement) (*StepStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(StepStatus)
	*result = StepStatus(s)
	return result, nil
}

func XMLStepStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]StepStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []StepStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, StepStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLLunStatusReadOne(reader *XMLReader, start *xml.StartElement) (*LunStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(LunStatus)
	*result = LunStatus(s)
	return result, nil
}

func XMLLunStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]LunStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []LunStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, LunStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNetworkPluginTypeReadOne(reader *XMLReader, start *xml.StartElement) (*NetworkPluginType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NetworkPluginType)
	*result = NetworkPluginType(s)
	return result, nil
}

func XMLNetworkPluginTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]NetworkPluginType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NetworkPluginType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NetworkPluginType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLHostProtocolReadOne(reader *XMLReader, start *xml.StartElement) (*HostProtocol, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(HostProtocol)
	*result = HostProtocol(s)
	return result, nil
}

func XMLHostProtocolReadMany(reader *XMLReader, start *xml.StartElement) ([]HostProtocol, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []HostProtocol
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, HostProtocol(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLFirewallTypeReadOne(reader *XMLReader, start *xml.StartElement) (*FirewallType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(FirewallType)
	*result = FirewallType(s)
	return result, nil
}

func XMLFirewallTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]FirewallType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []FirewallType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, FirewallType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLVmStorageErrorResumeBehaviourReadOne(reader *XMLReader, start *xml.StartElement) (*VmStorageErrorResumeBehaviour, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(VmStorageErrorResumeBehaviour)
	*result = VmStorageErrorResumeBehaviour(s)
	return result, nil
}

func XMLVmStorageErrorResumeBehaviourReadMany(reader *XMLReader, start *xml.StartElement) ([]VmStorageErrorResumeBehaviour, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []VmStorageErrorResumeBehaviour
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, VmStorageErrorResumeBehaviour(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLSeLinuxModeReadOne(reader *XMLReader, start *xml.StartElement) (*SeLinuxMode, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(SeLinuxMode)
	*result = SeLinuxMode(s)
	return result, nil
}

func XMLSeLinuxModeReadMany(reader *XMLReader, start *xml.StartElement) ([]SeLinuxMode, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []SeLinuxMode
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, SeLinuxMode(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLStorageTypeReadOne(reader *XMLReader, start *xml.StartElement) (*StorageType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(StorageType)
	*result = StorageType(s)
	return result, nil
}

func XMLStorageTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]StorageType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []StorageType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, StorageType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLOpenStackNetworkProviderTypeReadOne(reader *XMLReader, start *xml.StartElement) (*OpenStackNetworkProviderType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(OpenStackNetworkProviderType)
	*result = OpenStackNetworkProviderType(s)
	return result, nil
}

func XMLOpenStackNetworkProviderTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]OpenStackNetworkProviderType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []OpenStackNetworkProviderType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, OpenStackNetworkProviderType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLSerialNumberPolicyReadOne(reader *XMLReader, start *xml.StartElement) (*SerialNumberPolicy, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(SerialNumberPolicy)
	*result = SerialNumberPolicy(s)
	return result, nil
}

func XMLSerialNumberPolicyReadMany(reader *XMLReader, start *xml.StartElement) ([]SerialNumberPolicy, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []SerialNumberPolicy
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, SerialNumberPolicy(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLEntityExternalStatusReadOne(reader *XMLReader, start *xml.StartElement) (*EntityExternalStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(EntityExternalStatus)
	*result = EntityExternalStatus(s)
	return result, nil
}

func XMLEntityExternalStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]EntityExternalStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []EntityExternalStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, EntityExternalStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLScsiGenericIOReadOne(reader *XMLReader, start *xml.StartElement) (*ScsiGenericIO, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ScsiGenericIO)
	*result = ScsiGenericIO(s)
	return result, nil
}

func XMLScsiGenericIOReadMany(reader *XMLReader, start *xml.StartElement) ([]ScsiGenericIO, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ScsiGenericIO
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ScsiGenericIO(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLCloudInitNetworkProtocolReadOne(reader *XMLReader, start *xml.StartElement) (*CloudInitNetworkProtocol, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(CloudInitNetworkProtocol)
	*result = CloudInitNetworkProtocol(s)
	return result, nil
}

func XMLCloudInitNetworkProtocolReadMany(reader *XMLReader, start *xml.StartElement) ([]CloudInitNetworkProtocol, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []CloudInitNetworkProtocol
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, CloudInitNetworkProtocol(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLSwitchTypeReadOne(reader *XMLReader, start *xml.StartElement) (*SwitchType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(SwitchType)
	*result = SwitchType(s)
	return result, nil
}

func XMLSwitchTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]SwitchType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []SwitchType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, SwitchType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDiskTypeReadOne(reader *XMLReader, start *xml.StartElement) (*DiskType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DiskType)
	*result = DiskType(s)
	return result, nil
}

func XMLDiskTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]DiskType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DiskType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DiskType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDiskInterfaceReadOne(reader *XMLReader, start *xml.StartElement) (*DiskInterface, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DiskInterface)
	*result = DiskInterface(s)
	return result, nil
}

func XMLDiskInterfaceReadMany(reader *XMLReader, start *xml.StartElement) ([]DiskInterface, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DiskInterface
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DiskInterface(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDisplayTypeReadOne(reader *XMLReader, start *xml.StartElement) (*DisplayType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DisplayType)
	*result = DisplayType(s)
	return result, nil
}

func XMLDisplayTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]DisplayType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DisplayType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DisplayType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLInheritableBooleanReadOne(reader *XMLReader, start *xml.StartElement) (*InheritableBoolean, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(InheritableBoolean)
	*result = InheritableBoolean(s)
	return result, nil
}

func XMLInheritableBooleanReadMany(reader *XMLReader, start *xml.StartElement) ([]InheritableBoolean, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []InheritableBoolean
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, InheritableBoolean(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLOsTypeReadOne(reader *XMLReader, start *xml.StartElement) (*OsType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(OsType)
	*result = OsType(s)
	return result, nil
}

func XMLOsTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]OsType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []OsType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, OsType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLImageTransferPhaseReadOne(reader *XMLReader, start *xml.StartElement) (*ImageTransferPhase, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ImageTransferPhase)
	*result = ImageTransferPhase(s)
	return result, nil
}

func XMLImageTransferPhaseReadMany(reader *XMLReader, start *xml.StartElement) ([]ImageTransferPhase, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ImageTransferPhase
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ImageTransferPhase(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLCpuModeReadOne(reader *XMLReader, start *xml.StartElement) (*CpuMode, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(CpuMode)
	*result = CpuMode(s)
	return result, nil
}

func XMLCpuModeReadMany(reader *XMLReader, start *xml.StartElement) ([]CpuMode, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []CpuMode
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, CpuMode(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLHostTypeReadOne(reader *XMLReader, start *xml.StartElement) (*HostType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(HostType)
	*result = HostType(s)
	return result, nil
}

func XMLHostTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]HostType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []HostType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, HostType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLBootProtocolReadOne(reader *XMLReader, start *xml.StartElement) (*BootProtocol, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(BootProtocol)
	*result = BootProtocol(s)
	return result, nil
}

func XMLBootProtocolReadMany(reader *XMLReader, start *xml.StartElement) ([]BootProtocol, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []BootProtocol
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, BootProtocol(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLReportedDeviceTypeReadOne(reader *XMLReader, start *xml.StartElement) (*ReportedDeviceType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ReportedDeviceType)
	*result = ReportedDeviceType(s)
	return result, nil
}

func XMLReportedDeviceTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]ReportedDeviceType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ReportedDeviceType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ReportedDeviceType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLHookContentTypeReadOne(reader *XMLReader, start *xml.StartElement) (*HookContentType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(HookContentType)
	*result = HookContentType(s)
	return result, nil
}

func XMLHookContentTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]HookContentType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []HookContentType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, HookContentType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLStorageDomainStatusReadOne(reader *XMLReader, start *xml.StartElement) (*StorageDomainStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(StorageDomainStatus)
	*result = StorageDomainStatus(s)
	return result, nil
}

func XMLStorageDomainStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]StorageDomainStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []StorageDomainStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, StorageDomainStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLWatchdogActionReadOne(reader *XMLReader, start *xml.StartElement) (*WatchdogAction, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(WatchdogAction)
	*result = WatchdogAction(s)
	return result, nil
}

func XMLWatchdogActionReadMany(reader *XMLReader, start *xml.StartElement) ([]WatchdogAction, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []WatchdogAction
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, WatchdogAction(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLHookStatusReadOne(reader *XMLReader, start *xml.StartElement) (*HookStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(HookStatus)
	*result = HookStatus(s)
	return result, nil
}

func XMLHookStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]HookStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []HookStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, HookStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLStatisticUnitReadOne(reader *XMLReader, start *xml.StartElement) (*StatisticUnit, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(StatisticUnit)
	*result = StatisticUnit(s)
	return result, nil
}

func XMLStatisticUnitReadMany(reader *XMLReader, start *xml.StartElement) ([]StatisticUnit, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []StatisticUnit
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, StatisticUnit(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLArchitectureReadOne(reader *XMLReader, start *xml.StartElement) (*Architecture, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(Architecture)
	*result = Architecture(s)
	return result, nil
}

func XMLArchitectureReadMany(reader *XMLReader, start *xml.StartElement) ([]Architecture, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []Architecture
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, Architecture(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNfsVersionReadOne(reader *XMLReader, start *xml.StartElement) (*NfsVersion, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NfsVersion)
	*result = NfsVersion(s)
	return result, nil
}

func XMLNfsVersionReadMany(reader *XMLReader, start *xml.StartElement) ([]NfsVersion, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NfsVersion
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NfsVersion(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLGlusterVolumeTypeReadOne(reader *XMLReader, start *xml.StartElement) (*GlusterVolumeType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(GlusterVolumeType)
	*result = GlusterVolumeType(s)
	return result, nil
}

func XMLGlusterVolumeTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]GlusterVolumeType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []GlusterVolumeType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, GlusterVolumeType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDiskContentTypeReadOne(reader *XMLReader, start *xml.StartElement) (*DiskContentType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DiskContentType)
	*result = DiskContentType(s)
	return result, nil
}

func XMLDiskContentTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]DiskContentType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DiskContentType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DiskContentType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLVmAffinityReadOne(reader *XMLReader, start *xml.StartElement) (*VmAffinity, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(VmAffinity)
	*result = VmAffinity(s)
	return result, nil
}

func XMLVmAffinityReadMany(reader *XMLReader, start *xml.StartElement) ([]VmAffinity, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []VmAffinity
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, VmAffinity(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLAutoNumaStatusReadOne(reader *XMLReader, start *xml.StartElement) (*AutoNumaStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(AutoNumaStatus)
	*result = AutoNumaStatus(s)
	return result, nil
}

func XMLAutoNumaStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]AutoNumaStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []AutoNumaStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, AutoNumaStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLClusterUpgradeActionReadOne(reader *XMLReader, start *xml.StartElement) (*ClusterUpgradeAction, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ClusterUpgradeAction)
	*result = ClusterUpgradeAction(s)
	return result, nil
}

func XMLClusterUpgradeActionReadMany(reader *XMLReader, start *xml.StartElement) ([]ClusterUpgradeAction, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ClusterUpgradeAction
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ClusterUpgradeAction(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLExternalSystemTypeReadOne(reader *XMLReader, start *xml.StartElement) (*ExternalSystemType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ExternalSystemType)
	*result = ExternalSystemType(s)
	return result, nil
}

func XMLExternalSystemTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]ExternalSystemType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ExternalSystemType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ExternalSystemType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLRngSourceReadOne(reader *XMLReader, start *xml.StartElement) (*RngSource, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(RngSource)
	*result = RngSource(s)
	return result, nil
}

func XMLRngSourceReadMany(reader *XMLReader, start *xml.StartElement) ([]RngSource, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []RngSource
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, RngSource(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLBiosTypeReadOne(reader *XMLReader, start *xml.StartElement) (*BiosType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(BiosType)
	*result = BiosType(s)
	return result, nil
}

func XMLBiosTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]BiosType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []BiosType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, BiosType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDiskFormatReadOne(reader *XMLReader, start *xml.StartElement) (*DiskFormat, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DiskFormat)
	*result = DiskFormat(s)
	return result, nil
}

func XMLDiskFormatReadMany(reader *XMLReader, start *xml.StartElement) ([]DiskFormat, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DiskFormat
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DiskFormat(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLRoleTypeReadOne(reader *XMLReader, start *xml.StartElement) (*RoleType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(RoleType)
	*result = RoleType(s)
	return result, nil
}

func XMLRoleTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]RoleType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []RoleType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, RoleType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLDataCenterStatusReadOne(reader *XMLReader, start *xml.StartElement) (*DataCenterStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(DataCenterStatus)
	*result = DataCenterStatus(s)
	return result, nil
}

func XMLDataCenterStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]DataCenterStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []DataCenterStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, DataCenterStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLVmPoolTypeReadOne(reader *XMLReader, start *xml.StartElement) (*VmPoolType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(VmPoolType)
	*result = VmPoolType(s)
	return result, nil
}

func XMLVmPoolTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]VmPoolType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []VmPoolType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, VmPoolType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLPayloadEncodingReadOne(reader *XMLReader, start *xml.StartElement) (*PayloadEncoding, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(PayloadEncoding)
	*result = PayloadEncoding(s)
	return result, nil
}

func XMLPayloadEncodingReadMany(reader *XMLReader, start *xml.StartElement) ([]PayloadEncoding, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []PayloadEncoding
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, PayloadEncoding(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLStorageDomainTypeReadOne(reader *XMLReader, start *xml.StartElement) (*StorageDomainType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(StorageDomainType)
	*result = StorageDomainType(s)
	return result, nil
}

func XMLStorageDomainTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]StorageDomainType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []StorageDomainType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, StorageDomainType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLNotificationMethodReadOne(reader *XMLReader, start *xml.StartElement) (*NotificationMethod, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(NotificationMethod)
	*result = NotificationMethod(s)
	return result, nil
}

func XMLNotificationMethodReadMany(reader *XMLReader, start *xml.StartElement) ([]NotificationMethod, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []NotificationMethod
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, NotificationMethod(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLQuotaModeTypeReadOne(reader *XMLReader, start *xml.StartElement) (*QuotaModeType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(QuotaModeType)
	*result = QuotaModeType(s)
	return result, nil
}

func XMLQuotaModeTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]QuotaModeType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []QuotaModeType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, QuotaModeType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLWatchdogModelReadOne(reader *XMLReader, start *xml.StartElement) (*WatchdogModel, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(WatchdogModel)
	*result = WatchdogModel(s)
	return result, nil
}

func XMLWatchdogModelReadMany(reader *XMLReader, start *xml.StartElement) ([]WatchdogModel, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []WatchdogModel
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, WatchdogModel(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLPmProxyTypeReadOne(reader *XMLReader, start *xml.StartElement) (*PmProxyType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(PmProxyType)
	*result = PmProxyType(s)
	return result, nil
}

func XMLPmProxyTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]PmProxyType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []PmProxyType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, PmProxyType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLHostStatusReadOne(reader *XMLReader, start *xml.StartElement) (*HostStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(HostStatus)
	*result = HostStatus(s)
	return result, nil
}

func XMLHostStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]HostStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []HostStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, HostStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLLogSeverityReadOne(reader *XMLReader, start *xml.StartElement) (*LogSeverity, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(LogSeverity)
	*result = LogSeverity(s)
	return result, nil
}

func XMLLogSeverityReadMany(reader *XMLReader, start *xml.StartElement) ([]LogSeverity, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []LogSeverity
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, LogSeverity(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLMigrateOnErrorReadOne(reader *XMLReader, start *xml.StartElement) (*MigrateOnError, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(MigrateOnError)
	*result = MigrateOnError(s)
	return result, nil
}

func XMLMigrateOnErrorReadMany(reader *XMLReader, start *xml.StartElement) ([]MigrateOnError, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []MigrateOnError
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, MigrateOnError(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLVmStatusReadOne(reader *XMLReader, start *xml.StartElement) (*VmStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(VmStatus)
	*result = VmStatus(s)
	return result, nil
}

func XMLVmStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]VmStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []VmStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, VmStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLExternalStatusReadOne(reader *XMLReader, start *xml.StartElement) (*ExternalStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(ExternalStatus)
	*result = ExternalStatus(s)
	return result, nil
}

func XMLExternalStatusReadMany(reader *XMLReader, start *xml.StartElement) ([]ExternalStatus, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []ExternalStatus
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, ExternalStatus(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLVmTypeReadOne(reader *XMLReader, start *xml.StartElement) (*VmType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(VmType)
	*result = VmType(s)
	return result, nil
}

func XMLVmTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]VmType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []VmType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, VmType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}

func XMLUsbTypeReadOne(reader *XMLReader, start *xml.StartElement) (*UsbType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	s, err := reader.ReadString(start)
	if err != nil {
		return nil, err
	}
	result := new(UsbType)
	*result = UsbType(s)
	return result, nil
}

func XMLUsbTypeReadMany(reader *XMLReader, start *xml.StartElement) ([]UsbType, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		start = st
	}
	var results []UsbType
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		t = xml.CopyToken(t)
		switch t := t.(type) {
		case xml.StartElement:
			one, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			results = append(results, UsbType(one))
		case xml.EndElement:
			depth--
		}
	}
	return results, nil
}
