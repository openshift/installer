package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Sorter struct {
	Keys []string
	Vals []string
}

func newSorter(m map[string]string) *Sorter {
	hs := &Sorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]string, 0, len(m)),
	}

	for k, v := range m {
		hs.Keys = append(hs.Keys, k)
		hs.Vals = append(hs.Vals, v)
	}
	return hs
}

// Sort is an additional function for function SignHeader.
func (hs *Sorter) Sort() {
	sort.Sort(hs)
}

// Len is an additional function for function SignHeader.
func (hs *Sorter) Len() int {
	return len(hs.Vals)
}

// Less is an additional function for function SignHeader.
func (hs *Sorter) Less(i, j int) bool {
	return bytes.Compare([]byte(hs.Keys[i]), []byte(hs.Keys[j])) < 0
}

// Swap is an additional function for function SignHeader.
func (hs *Sorter) Swap(i, j int) {
	hs.Vals[i], hs.Vals[j] = hs.Vals[j], hs.Vals[i]
	hs.Keys[i], hs.Keys[j] = hs.Keys[j], hs.Keys[i]
}

func GetSignature(stringToSign *string, secret *string) *string {
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(tea.StringValue(secret)))
	io.WriteString(h, tea.StringValue(stringToSign))
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return tea.String(signedStr)
}

func GetStringToSign(request *tea.Request) *string {
	return tea.String(getStringToSign(request))
}

func ToForm(filter map[string]interface{}) *string {
	tmp := make(map[string]interface{})
	byt, _ := json.Marshal(filter)
	d := json.NewDecoder(bytes.NewReader(byt))
	d.UseNumber()
	_ = d.Decode(&tmp)

	result := make(map[string]*string)
	for key, value := range tmp {
		filterValue := reflect.ValueOf(value)
		flatRepeatedList(filterValue, result, key)
	}

	m := util.AnyifyMapValue(result)
	return util.ToFormString(m)
}

func Convert(input, output interface{}) {
	res := make(map[string]interface{})
	val := reflect.ValueOf(input).Elem()
	dataType := val.Type()
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		name, _ := field.Tag.Lookup("json")
		name = strings.Split(name, ",omitempty")[0]
		_, ok := val.Field(i).Interface().(io.Reader)
		if !ok {
			res[name] = val.Field(i).Interface()
		}
	}
	byt, _ := json.Marshal(res)
	json.Unmarshal(byt, output)
}

func flatRepeatedList(dataValue reflect.Value, result map[string]*string, prefix string) {
	if !dataValue.IsValid() {
		return
	}

	dataType := dataValue.Type()
	if dataType.Kind().String() == "slice" {
		handleRepeatedParams(dataValue, result, prefix)
	} else if dataType.Kind().String() == "map" {
		handleMap(dataValue, result, prefix)
	} else {
		result[prefix] = tea.String(fmt.Sprintf("%v", dataValue.Interface()))
	}
}

func handleRepeatedParams(repeatedFieldValue reflect.Value, result map[string]*string, prefix string) {
	if repeatedFieldValue.IsValid() && !repeatedFieldValue.IsNil() {
		for m := 0; m < repeatedFieldValue.Len(); m++ {
			elementValue := repeatedFieldValue.Index(m)
			key := prefix + "." + strconv.Itoa(m+1)
			fieldValue := reflect.ValueOf(elementValue.Interface())
			if fieldValue.Kind().String() == "map" {
				handleMap(fieldValue, result, key)
			} else {
				result[key] = tea.String(fmt.Sprintf("%v", fieldValue.Interface()))
			}
		}
	}
}

func handleMap(valueField reflect.Value, result map[string]*string, prefix string) {
	if valueField.IsValid() && valueField.String() != "" {
		valueFieldType := valueField.Type()
		if valueFieldType.Kind().String() == "map" {
			var byt []byte
			byt, _ = json.Marshal(valueField.Interface())
			cache := make(map[string]interface{})
			d := json.NewDecoder(bytes.NewReader(byt))
			d.UseNumber()
			_ = d.Decode(&cache)
			for key, value := range cache {
				pre := ""
				if prefix != "" {
					pre = prefix + "." + key
				} else {
					pre = key
				}
				fieldValue := reflect.ValueOf(value)
				flatRepeatedList(fieldValue, result, pre)
			}
		}
	}
}

func getStringToSign(request *tea.Request) string {
	resource := tea.StringValue(request.Pathname)
	queryParams := request.Query
	// sort QueryParams by key
	var queryKeys []string
	for key := range queryParams {
		queryKeys = append(queryKeys, key)
	}
	sort.Strings(queryKeys)
	tmp := ""
	for i := 0; i < len(queryKeys); i++ {
		queryKey := queryKeys[i]
		v := tea.StringValue(queryParams[queryKey])
		if v != "" {
			tmp = tmp + "&" + queryKey + "=" + v
		} else {
			tmp = tmp + "&" + queryKey
		}
	}
	if tmp != "" {
		tmp = strings.TrimLeft(tmp, "&")
		resource = resource + "?" + tmp
	}
	return getSignedStr(request, resource)
}

func getSignedStr(req *tea.Request, canonicalizedResource string) string {
	temp := make(map[string]string)

	for k, v := range req.Headers {
		if strings.HasPrefix(strings.ToLower(k), "x-acs-") {
			temp[strings.ToLower(k)] = tea.StringValue(v)
		}
	}
	hs := newSorter(temp)

	// Sort the temp by the ascending order
	hs.Sort()

	// Get the canonicalizedOSSHeaders
	canonicalizedOSSHeaders := ""
	for i := range hs.Keys {
		canonicalizedOSSHeaders += hs.Keys[i] + ":" + hs.Vals[i] + "\n"
	}

	// Give other parameters values
	// when sign URL, date is expires
	date := tea.StringValue(req.Headers["date"])
	accept := tea.StringValue(req.Headers["accept"])
	contentType := tea.StringValue(req.Headers["content-type"])
	contentMd5 := tea.StringValue(req.Headers["content-md5"])

	signStr := tea.StringValue(req.Method) + "\n" + accept + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	return signStr
}
