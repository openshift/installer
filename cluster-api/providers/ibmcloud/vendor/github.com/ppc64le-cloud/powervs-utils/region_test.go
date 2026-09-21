package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetRegion(t *testing.T) {
	type args struct {
		zone string
	}
	tests := []struct {
		name       string
		args       args
		wantRegion string
		wantErr    bool
	}{
		{
			"London",
			args{"lon06"},
			"lon",
			false,
		},
		{
			"Dallas",
			args{"us-south"},
			"us-south",
			false,
		},
		{
			"Dallas 12",
			args{"dal12"},
			"dal",
			false,
		},
		{
			"Sao Paulo 01",
			args{"sao01"},
			"sao",
			false,
		},
		{
			"Sao Paulo 04",
			args{"sao04"},
			"sao",
			false,
		},
		{
			"Washington DC",
			args{"us-east"},
			"us-east",
			false,
		},
		{
			"Washington DC 06",
			args{"wdc06"},
			"wdc",
			false,
		},
		{
			"Washington DC 07",
			args{"wdc07"},
			"wdc",
			false,
		},
		{
			"Toronto",
			args{"tor01"},
			"tor",
			false,
		},
		{
			"Frankfurt",
			args{"eu-de-1"},
			"eu-de",
			false,
		},
		{
			"Sydney",
			args{"syd01"},
			"syd",
			false,
		},
		{
			"India",
			args{"blr01"},
			"",
			true,
		},
		{
			"Tokyo",
			args{"tok04"},
			"tok",
			false,
		},
		{
			"Montreal",
			args{"mon01"},
			"mon",
			false,
		},
		{
			"Osaka",
			args{"osa21"},
			"osa",
			false,
		},
		{
			"Madrid 02",
			args{"mad02"},
			"mad",
			false,
		},
		{
			"Madrid 04",
			args{"mad04"},
			"mad",
			false,
		},
		{
			"Chennai",
			args{"che01"},
			"che",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRegion, err := GetRegion(tt.args.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRegion != tt.wantRegion {
				t.Errorf("GetRegion() gotRegion = %v, want %v", gotRegion, tt.wantRegion)
			}
		})
	}
}

func TestCOSRegionForVPCRegion(t *testing.T) {
	type args struct {
		region string
	}
	tests := []struct {
		name       string
		args       args
		wantRegion string
		wantErr    bool
	}{
		{
			"Dallas",
			args{"us-south"},
			"us-south",
			false,
		},
		{
			"Osaka",
			args{"eu-de1"},
			"eu-de",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vpcRegion, err := COSRegionForVPCRegion(tt.args.region)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("COSRegionForVPCRegion() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if vpcRegion != tt.wantRegion {
				t.Errorf("COSRegionForVPCRegion() gotRegion = %v, want %v", vpcRegion, tt.wantRegion)
			}
		})
	}
}

func TestVPCRegionForPowerVSRegion(t *testing.T) {
	type args struct {
		region string
	}
	tests := []struct {
		name       string
		args       args
		wantRegion string
		wantErr    bool
	}{
		{
			"Dallas",
			args{"dal"},
			"us-south",
			false,
		},
		{
			"Osaka",
			args{"eu-de1"},
			"eu-de",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vpcRegion, err := VPCRegionForPowerVSRegion(tt.args.region)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("VPCRegionForPowerVSRegion() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if vpcRegion != tt.wantRegion {
				t.Errorf("VPCRegionForPowerVSRegion() gotRegion = %v, want %v", vpcRegion, tt.wantRegion)
			}
		})
	}
}

func TestCOSRegionForPowerVSRegion(t *testing.T) {
	type args struct {
		region string
	}
	tests := []struct {
		name       string
		args       args
		wantRegion string
		wantErr    bool
	}{
		{
			"Dallas",
			args{"dal"},
			"us-south",
			false,
		},
		{
			"Osaka",
			args{"eu-de1"},
			"eu-de",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			COSRegion, err := COSRegionForPowerVSRegion(tt.args.region)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("COSRegionForPowerVSRegion() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if COSRegion != tt.wantRegion {
				t.Errorf("COSRegionForPowerVSRegion() gotRegion = %v, want %v", COSRegion, tt.wantRegion)
			}
		})
	}
}
func TestValidateCOSRegion(t *testing.T) {
	tests := []struct {
		testcasename string
		COSRegion    string
		IsCosRegion  bool
	}{
		{
			testcasename: "Given COS region isn't tested",
			COSRegion:    "che",
			IsCosRegion:  false,
		},
		{
			testcasename: "Given COS region is tested",
			COSRegion:    "us-south",
			IsCosRegion:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			isKnown := ValidateCOSRegion(tt.COSRegion)
			if isKnown != tt.IsCosRegion {
				t.Errorf("ValidateCOSRegion(), expected: %t, returned: %t", tt.IsCosRegion, isKnown)
			}
		})
	}
}

func TestValidateVPCRegion(t *testing.T) {
	tests := []struct {
		testcasename string
		VPCRegion    string
		IsVPCRegion  bool
	}{
		{
			testcasename: "Given COS region isn't tested",
			VPCRegion:    "che",
			IsVPCRegion:  false,
		},
		{
			testcasename: "Given COS region is tested",
			VPCRegion:    "us-south",
			IsVPCRegion:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			isKnown := ValidateVPCRegion(tt.VPCRegion)
			if isKnown != tt.IsVPCRegion {
				t.Errorf("ValidateVPCRegion(), region: %s, expected: %t, returned: %t", tt.VPCRegion, tt.IsVPCRegion, isKnown)
			}
		})
	}
}

func TestValidateZone(t *testing.T) {
	tests := []struct {
		testcasename string
		zoneName     string
		isExists     bool
	}{
		{
			testcasename: "Given zone exists",
			zoneName:     "sao04",
			isExists:     true,
		},
		{
			testcasename: "Given zone doesn't exist",
			zoneName:     "wdc04",
			isExists:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			exists := ValidateZone(tt.zoneName)
			if exists != tt.isExists {
				t.Errorf("Zone %s expected to exist: %t, ValidateZone() returned: %t", tt.zoneName, tt.isExists, exists)
			}
		})
	}
}

func TestRegionFromZone(t *testing.T) {
	tests := []struct {
		testcasename   string
		zoneName       string
		expectedRegion string
	}{
		{
			testcasename:   "Region found for the given zone",
			zoneName:       "wdc06",
			expectedRegion: "wdc",
		},
		{
			testcasename:   "No region found for the given zone",
			zoneName:       "wdc04",
			expectedRegion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			returnedRegion := RegionFromZone(tt.zoneName)
			if returnedRegion != tt.expectedRegion {
				t.Errorf("For zone %s, expected region %s, returned region: %s", tt.zoneName, tt.expectedRegion, returnedRegion)
			}
		})
	}
}

func TestAvailableSysTypes(t *testing.T) {
	tests := []struct {
		testcasename  string
		region        string
		zoneName      string
		expectedError error
		systypes      []string
	}{
		{
			testcasename:  "Unknown region name",
			region:        "unknown",
			zoneName:      "unknown",
			expectedError: fmt.Errorf("unknown region name provided"),
		},
		{
			testcasename:  "Unknown zone name",
			region:        "us-south",
			zoneName:      "unknown",
			expectedError: fmt.Errorf("unknown zone name provided"),
		},
		{
			testcasename:  "Systypes available",
			region:        "wdc",
			zoneName:      "wdc07",
			systypes:      []string{"e1050", "e1080", "e980", "s1022", "s922"},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			returnedSystypes, returnedErr := AvailableSysTypes(tt.region, tt.zoneName)
			if tt.expectedError != nil {
				if returnedErr == nil || returnedErr.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error: %v, returned error %v", tt.expectedError, returnedErr)
				}
				return
			}
			if !reflect.DeepEqual(tt.systypes, returnedSystypes) {
				t.Errorf("Expected Systypes: %v, Returned systypes: %v", tt.systypes, returnedSystypes)
			}
		})
	}
}

func TestIsGlobalRoutingRequiredForTG(t *testing.T) {
	tests := []struct {
		testcasename            string
		powervsRegion           string
		vpcRegion               string
		isGlobalRoutingRequired bool
	}{
		{
			testcasename:            "Powervs and vpc regions are same.",
			powervsRegion:           "wdc",
			vpcRegion:               "us-east",
			isGlobalRoutingRequired: false,
		},
		{
			testcasename:            "Powervs and vpc regions are different.",
			powervsRegion:           "mon",
			vpcRegion:               "jp-osa",
			isGlobalRoutingRequired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			globalRoutingRequired := IsGlobalRoutingRequiredForTG(tt.powervsRegion, tt.vpcRegion)
			if globalRoutingRequired != tt.isGlobalRoutingRequired {
				t.Errorf("GlobalRoutingRequired? Expected: %t ,Received: %t", tt.isGlobalRoutingRequired, globalRoutingRequired)
			}
		})
	}
}

func TestVPCZonesForVPCRegion(t *testing.T) {
	tests := []struct {
		testcasename  string
		powervsRegion string
		vpcZones      []string
		expectedError error
	}{
		{
			testcasename:  "VPC Zones found for the corresponding region",
			powervsRegion: "jp-osa",
			vpcZones:      []string{"jp-osa-1", "jp-osa-2", "jp-osa-3"},
		}, {
			testcasename:  "VPC Zones not found for the corresponding region",
			powervsRegion: "unknown",
			expectedError: fmt.Errorf("VPC zones corresponding to the VPC region unknown is not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			returnedZones, err := VPCZonesForVPCRegion(tt.powervsRegion)
			if tt.expectedError != nil {
				if (err == nil) || tt.expectedError.Error() != err.Error() {
					t.Errorf("Expected error: %v, returned error: %v", tt.expectedError, err)
				}
				return
			}
			if !reflect.DeepEqual(tt.vpcZones, returnedZones) {
				t.Errorf("Expected vpcZones: %v, Returned vpcZones: %v", tt.vpcZones, returnedZones)
			}
		})
	}
}
