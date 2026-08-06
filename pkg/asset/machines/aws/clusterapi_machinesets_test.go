package aws

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/core/v1beta2"

	icaws "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

func TestClusterAPIMachineSets(t *testing.T) {
	tests := []struct {
		name     string
		input    *MachineSetInput
		wantErr  string
		validate func(t *testing.T, templates []capa.AWSMachineTemplate, machineSets []capi.MachineSet)
	}{
		{
			name: "non-AWS pool returns error",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS = nil
				return in
			}(),
			wantErr: `non-AWS machine-pool: ""`,
		},
		{
			name: "missing subnet for zone in BYO VPC returns error",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.Subnets = icaws.SubnetsByZone{
					"us-east-1b": {ID: "subnet-b"},
				}
				return in
			}(),
			wantErr: `no subnet for zone us-east-1a`,
		},
		{
			name: "user tags conflict with reserved key returns error",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.InstallConfigPlatformAWS.UserTags = map[string]string{
					"kubernetes.io/cluster/test-cluster": "not-allowed",
				}
				return in
			}(),
			wantErr: "failed to create CAPA tags from user tags",
		},
		{
			name: "single zone managed VPC",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.Pool.Replicas = ptr.To(int64(2))
				return in
			}(),
			validate: validateSingleZoneManagedVPC,
		},
		{
			name:  "replicas distributed evenly across three zones",
			input: defaultMachineSetInput(),
			validate: func(t *testing.T, _ []capa.AWSMachineTemplate, machineSets []capi.MachineSet) {
				t.Helper()
				if len(machineSets) != 3 {
					t.Fatalf("expected 3 machinesets, got %d", len(machineSets))
				}
				for i, ms := range machineSets {
					if *ms.Spec.Replicas != 1 {
						t.Errorf("zone %d: replicas = %d, want 1", i, *ms.Spec.Replicas)
					}
				}
			},
		},
		{
			name: "1 replica across 3 zones assigns to first zone only",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Replicas = ptr.To(int64(1))
				return in
			}(),
			validate: func(t *testing.T, _ []capa.AWSMachineTemplate, machineSets []capi.MachineSet) {
				t.Helper()
				wantReplicas := []int32{1, 0, 0}
				for i, ms := range machineSets {
					if *ms.Spec.Replicas != wantReplicas[i] {
						t.Errorf("zone %d: replicas = %d, want %d", i, *ms.Spec.Replicas, wantReplicas[i])
					}
				}
			},
		},
		{
			name: "0 replicas",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Replicas = ptr.To(int64(0))
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a", "us-east-1b"}
				return in
			}(),
			validate: func(t *testing.T, _ []capa.AWSMachineTemplate, machineSets []capi.MachineSet) {
				t.Helper()
				for i, ms := range machineSets {
					if *ms.Spec.Replicas != 0 {
						t.Errorf("zone %d: replicas = %d, want 0", i, *ms.Spec.Replicas)
					}
				}
			},
		},
		{
			name: "BYO VPC with explicit subnets",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a", "us-east-1b"}
				in.Pool.Replicas = ptr.To(int64(2))
				in.Subnets = icaws.SubnetsByZone{
					"us-east-1a": {ID: "subnet-aaa", Public: false},
					"us-east-1b": {ID: "subnet-bbb", Public: true},
				}
				return in
			}(),
			validate: validateBYOVPCSubnets,
		},
		{
			name: "public subnet managed VPC uses public subnet filter",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.PublicSubnet = true
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a", "us-east-1b"}
				in.Pool.Replicas = ptr.To(int64(2))
				return in
			}(),
			validate: validatePublicSubnetFilter,
		},
		{
			name: "IMDS defaults to optional",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.Pool.Replicas = ptr.To(int64(1))
				return in
			}(),
			validate: func(t *testing.T, templates []capa.AWSMachineTemplate, _ []capi.MachineSet) {
				t.Helper()
				if len(templates) != 1 {
					t.Fatalf("expected 1 template, got %d", len(templates))
				}
				got := templates[0].Spec.Template.Spec.InstanceMetadataOptions.HTTPTokens
				if got != capa.HTTPTokensStateOptional {
					t.Errorf("IMDS HTTPTokens = %q, want %q", got, capa.HTTPTokensStateOptional)
				}
			},
		},
		{
			name: "IMDS required when authentication is Required",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.Pool.Replicas = ptr.To(int64(1))
				in.Pool.Platform.AWS.EC2Metadata = awstypes.EC2Metadata{
					Authentication: "Required",
				}
				return in
			}(),
			validate: func(t *testing.T, templates []capa.AWSMachineTemplate, _ []capi.MachineSet) {
				t.Helper()
				if len(templates) != 1 {
					t.Fatalf("expected 1 template, got %d", len(templates))
				}
				got := templates[0].Spec.Template.Spec.InstanceMetadataOptions.HTTPTokens
				if got != capa.HTTPTokensStateRequired {
					t.Errorf("IMDS HTTPTokens = %q, want %q", got, capa.HTTPTokensStateRequired)
				}
			},
		},
		{
			name: "default IAM profile uses cluster ID",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.Pool.Replicas = ptr.To(int64(1))
				return in
			}(),
			validate: func(t *testing.T, templates []capa.AWSMachineTemplate, _ []capi.MachineSet) {
				t.Helper()
				if len(templates) != 1 {
					t.Fatalf("expected 1 template, got %d", len(templates))
				}
				got := templates[0].Spec.Template.Spec.IAMInstanceProfile
				if got != "test-cluster-worker-profile" {
					t.Errorf("IAM profile = %q, want %q", got, "test-cluster-worker-profile")
				}
			},
		},
		{
			name: "custom IAM profile",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.Pool.Replicas = ptr.To(int64(1))
				in.Pool.Platform.AWS.IAMProfile = "my-custom-profile"
				return in
			}(),
			validate: func(t *testing.T, templates []capa.AWSMachineTemplate, _ []capi.MachineSet) {
				t.Helper()
				if len(templates) != 1 {
					t.Fatalf("expected 1 template, got %d", len(templates))
				}
				got := templates[0].Spec.Template.Spec.IAMInstanceProfile
				if got != "my-custom-profile" {
					t.Errorf("IAM profile = %q, want %q", got, "my-custom-profile")
				}
			},
		},
		{
			name: "user tags propagated to templates",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.Pool.Replicas = ptr.To(int64(1))
				in.InstallConfigPlatformAWS.UserTags = map[string]string{
					"env":   "prod",
					"owner": "team-a",
				}
				return in
			}(),
			validate: func(t *testing.T, templates []capa.AWSMachineTemplate, _ []capi.MachineSet) {
				t.Helper()
				if len(templates) != 1 {
					t.Fatalf("expected 1 template, got %d", len(templates))
				}
				got := templates[0].Spec.Template.Spec.AdditionalTags
				want := capa.Tags{
					"kubernetes.io/cluster/test-cluster": "owned",
					"env":                                "prod",
					"owner":                              "team-a",
				}
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("tags mismatch (-want +got):\n%s", diff)
				}
			},
		},
		{
			name: "additional security groups appended",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Platform.AWS.Zones = []string{"us-east-1a"}
				in.Pool.Replicas = ptr.To(int64(1))
				in.Pool.Platform.AWS.AdditionalSecurityGroupIDs = []string{"sg-extra-1", "sg-extra-2"}
				return in
			}(),
			validate: validateAdditionalSecurityGroups,
		},
		{
			name: "edge pool adds labels taints and preferred instance type",
			input: func() *MachineSetInput {
				in := defaultMachineSetInput()
				in.Pool.Name = types.MachinePoolEdgeRoleName
				in.Pool.Replicas = ptr.To(int64(1))
				in.Pool.Platform.AWS.Zones = []string{"us-east-1-bos-1a"}
				in.Zones = icaws.Zones{
					"us-east-1-bos-1a": &icaws.Zone{
						Name:                  "us-east-1-bos-1a",
						Type:                  "local-zone",
						GroupName:             "us-east-1-bos-1",
						ParentZoneName:        "us-east-1c",
						PreferredInstanceType: "m5.xlarge",
					},
				}
				return in
			}(),
			validate: validateEdgePool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates, machineSets, err := ClusterAPIMachineSets(tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error = %q, want substring %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.validate != nil {
				tt.validate(t, templates, machineSets)
			}
		})
	}
}

func defaultMachineSetInput() *MachineSetInput {
	return &MachineSetInput{
		ClusterID: "test-cluster",
		InstallConfigPlatformAWS: &awstypes.Platform{
			Region: "us-east-1",
		},
		Pool: &types.MachinePool{
			Name:     types.MachinePoolComputeRoleName,
			Replicas: ptr.To(int64(3)),
			Platform: types.MachinePoolPlatform{
				AWS: &awstypes.MachinePool{
					InstanceType: "m5.xlarge",
					Zones:        []string{"us-east-1a", "us-east-1b", "us-east-1c"},
					EC2RootVolume: awstypes.EC2RootVolume{
						Size: 120,
						Type: "gp3",
					},
				},
			},
		},
		UserDataSecret: "worker-user-data",
	}
}

func validateSingleZoneManagedVPC(t *testing.T, templates []capa.AWSMachineTemplate, machineSets []capi.MachineSet) {
	t.Helper()
	clusterID := "test-cluster"
	baseTags := capa.Tags{"kubernetes.io/cluster/test-cluster": "owned"}

	wantTemplates := []capa.AWSMachineTemplate{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
				Kind:       "AWSMachineTemplate",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-cluster-worker-us-east-1a",
				Namespace: "openshift-cluster-api",
				Labels:    map[string]string{"cluster.x-k8s.io/cluster-name": clusterID},
			},
			Spec: capa.AWSMachineTemplateSpec{
				Template: capa.AWSMachineTemplateResource{
					Spec: GenerateCAPIMachineSpec(&CAPIMachineSpecInput{
						InstanceType:       "m5.xlarge",
						IAMInstanceProfile: "test-cluster-worker-profile",
						Subnet: &capa.AWSResourceReference{
							Filters: []capa.Filter{{
								Name:   "tag:Name",
								Values: []string{"test-cluster-subnet-private-us-east-1a"},
							}},
						},
						Tags: baseTags,
						EC2RootVolume: awstypes.EC2RootVolume{
							Size: 120,
							Type: "gp3",
						},
						IMDS: capa.HTTPTokensStateOptional,
						SecurityGroups: []capa.AWSResourceReference{
							{Filters: []capa.Filter{{Name: "tag:Name", Values: []string{"test-cluster-node"}}}},
							{Filters: []capa.Filter{{Name: "tag:Name", Values: []string{"test-cluster-lb"}}}},
						},
						Ignition: &capa.Ignition{
							Version:     "3.2",
							StorageType: capa.IgnitionStorageTypeOptionUnencryptedUserData,
						},
					}),
				},
			},
		},
	}
	if diff := cmp.Diff(wantTemplates, templates); diff != "" {
		t.Errorf("AWSMachineTemplates mismatch (-want +got):\n%s", diff)
	}
	wantSets := []capi.MachineSet{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: capi.GroupVersion.String(),
				Kind:       "MachineSet",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-cluster-worker-us-east-1a",
				Namespace: "openshift-cluster-api",
				Labels:    map[string]string{"cluster.x-k8s.io/cluster-name": clusterID},
			},
			Spec: capi.MachineSetSpec{
				ClusterName: clusterID,
				Replicas:    ptr.To(int32(2)),
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"cluster.x-k8s.io/cluster-name": clusterID,
						"cluster.x-k8s.io/set-name":     "test-cluster-worker-us-east-1a",
					},
				},
				Template: capi.MachineTemplateSpec{
					ObjectMeta: capi.ObjectMeta{
						Labels: map[string]string{
							"cluster.x-k8s.io/cluster-name":  clusterID,
							"cluster.x-k8s.io/set-name":      "test-cluster-worker-us-east-1a",
							"node-role.kubernetes.io/worker": "",
						},
					},
					Spec: capi.MachineSpec{
						ClusterName: clusterID,
						Bootstrap: capi.Bootstrap{
							DataSecretName: ptr.To("worker-user-data"),
						},
						InfrastructureRef: capi.ContractVersionedObjectReference{
							APIGroup: "infrastructure.cluster.x-k8s.io",
							Kind:     "AWSMachineTemplate",
							Name:     "test-cluster-worker-us-east-1a",
						},
						FailureDomain: "us-east-1a",
						Taints:        []capi.MachineTaint{},
					},
				},
			},
		},
	}
	if diff := cmp.Diff(wantSets, machineSets); diff != "" {
		t.Errorf("MachineSets mismatch (-want +got):\n%s", diff)
	}
}

func validateBYOVPCSubnets(t *testing.T, templates []capa.AWSMachineTemplate, machineSets []capi.MachineSet) {
	t.Helper()
	wantSubnets := []struct {
		id       string
		publicIP bool
		zone     string
	}{
		{id: "subnet-aaa", publicIP: false, zone: "us-east-1a"},
		{id: "subnet-bbb", publicIP: true, zone: "us-east-1b"},
	}
	if len(templates) != len(wantSubnets) {
		t.Fatalf("expected %d templates, got %d", len(wantSubnets), len(templates))
	}
	for i, want := range wantSubnets {
		spec := templates[i].Spec.Template.Spec
		if spec.Subnet == nil || ptr.Deref(spec.Subnet.ID, "") != want.id {
			t.Errorf("zone %d: subnet ID = %+v, want %s", i, spec.Subnet, want.id)
		}
		if ptr.Deref(spec.PublicIP, false) != want.publicIP {
			t.Errorf("zone %d: PublicIP = %v, want %v", i, ptr.Deref(spec.PublicIP, false), want.publicIP)
		}
		if machineSets[i].Spec.Template.Spec.FailureDomain != want.zone {
			t.Errorf("zone %d: failure domain = %s, want %s", i, machineSets[i].Spec.Template.Spec.FailureDomain, want.zone)
		}
	}
}

func validatePublicSubnetFilter(t *testing.T, templates []capa.AWSMachineTemplate, _ []capi.MachineSet) {
	t.Helper()
	zones := []string{"us-east-1a", "us-east-1b"}
	if len(templates) != len(zones) {
		t.Fatalf("expected %d templates, got %d", len(zones), len(templates))
	}
	for i, az := range zones {
		spec := templates[i].Spec.Template.Spec
		if !ptr.Deref(spec.PublicIP, false) {
			t.Errorf("zone %s: PublicIP = false, want true", az)
		}
		if len(spec.Subnet.Filters) == 0 {
			t.Fatalf("zone %s: expected subnet filters", az)
		}
		got := spec.Subnet.Filters[0].Values[0]
		want := "test-cluster-subnet-public-" + az
		if got != want {
			t.Errorf("zone %s: subnet filter = %q, want %q", az, got, want)
		}
	}
}

func validateAdditionalSecurityGroups(t *testing.T, templates []capa.AWSMachineTemplate, _ []capi.MachineSet) {
	t.Helper()
	if len(templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(templates))
	}
	sgs := templates[0].Spec.Template.Spec.AdditionalSecurityGroups
	if len(sgs) != 4 {
		t.Fatalf("expected 4 security groups (2 default + 2 additional), got %d", len(sgs))
	}
	if sgs[0].Filters[0].Values[0] != "test-cluster-node" {
		t.Errorf("first SG filter = %v, want test-cluster-node", sgs[0].Filters[0].Values)
	}
	if sgs[1].Filters[0].Values[0] != "test-cluster-lb" {
		t.Errorf("second SG filter = %v, want test-cluster-lb", sgs[1].Filters[0].Values)
	}
	if ptr.Deref(sgs[2].ID, "") != "sg-extra-1" {
		t.Errorf("third SG ID = %q, want sg-extra-1", ptr.Deref(sgs[2].ID, ""))
	}
	if ptr.Deref(sgs[3].ID, "") != "sg-extra-2" {
		t.Errorf("fourth SG ID = %q, want sg-extra-2", ptr.Deref(sgs[3].ID, ""))
	}
}

func validateEdgePool(t *testing.T, templates []capa.AWSMachineTemplate, machineSets []capi.MachineSet) {
	t.Helper()
	if len(templates) != 1 || len(machineSets) != 1 {
		t.Fatalf("expected 1 template and 1 machineset, got %d and %d", len(templates), len(machineSets))
	}
	if templates[0].Spec.Template.Spec.InstanceType != "m5.xlarge" {
		t.Errorf("instance type = %q, want m5.xlarge", templates[0].Spec.Template.Spec.InstanceType)
	}
	msLabels := machineSets[0].Spec.Template.ObjectMeta.Labels
	wantLabels := map[string]string{
		"node-role.kubernetes.io/edge":          "",
		"node-role.kubernetes.io/worker":        "",
		"machine.openshift.io/zone-type":        "local-zone",
		"machine.openshift.io/zone-group":       "us-east-1-bos-1",
		"machine.openshift.io/parent-zone-name": "us-east-1c",
	}
	for k, v := range wantLabels {
		if got, ok := msLabels[k]; !ok || got != v {
			t.Errorf("label %s = %q (present=%v), want %q", k, got, ok, v)
		}
	}
	taints := machineSets[0].Spec.Template.Spec.Taints
	if len(taints) != 1 {
		t.Fatalf("expected 1 taint, got %d", len(taints))
	}
	if taints[0].Key != "node-role.kubernetes.io/edge" || taints[0].Effect != "NoSchedule" {
		t.Errorf("taint = %+v, want key=node-role.kubernetes.io/edge effect=NoSchedule", taints[0])
	}
	if taints[0].Propagation != capi.MachineTaintPropagationAlways {
		t.Errorf("taint propagation = %q, want Always", taints[0].Propagation)
	}
	if machineSets[0].Name != "test-cluster-edge-us-east-1-bos-1a" {
		t.Errorf("machineset name = %q, want test-cluster-edge-us-east-1-bos-1a", machineSets[0].Name)
	}
}
