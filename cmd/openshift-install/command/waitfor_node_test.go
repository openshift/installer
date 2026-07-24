package command

import (
	"context"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func newNode(name, role string, ready bool) *corev1.Node {
	status := corev1.ConditionFalse
	if ready {
		status = corev1.ConditionTrue
	}
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"node-role.kubernetes.io/" + role: "",
			},
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeReady,
					Status: status,
				},
			},
		},
	}
}

func TestVerifyExpectedNodesWithClient(t *testing.T) {
	tests := []struct {
		name            string
		nodes           []*corev1.Node
		expectedMasters int
		expectedWorkers int
		expectErr       bool
		errContains     string
	}{
		{
			name:            "all masters ready",
			nodes:           []*corev1.Node{newNode("master-0", "master", true), newNode("master-1", "master", true), newNode("master-2", "master", true)},
			expectedMasters: 3,
			expectedWorkers: 0,
			expectErr:       false,
		},
		{
			name:            "missing master node",
			nodes:           []*corev1.Node{newNode("master-0", "master", true), newNode("master-1", "master", true)},
			expectedMasters: 3,
			expectedWorkers: 0,
			expectErr:       true,
			errContains:     "expected 3 master node(s) to be Ready but only 2 found",
		},
		{
			name:            "master node not ready",
			nodes:           []*corev1.Node{newNode("master-0", "master", true), newNode("master-1", "master", true), newNode("master-2", "master", false)},
			expectedMasters: 3,
			expectedWorkers: 0,
			expectErr:       true,
			errContains:     "expected 3 master node(s) to be Ready but only 2 found",
		},
		{
			name:            "all workers ready",
			nodes:           []*corev1.Node{newNode("worker-0", "worker", true), newNode("worker-1", "worker", true)},
			expectedMasters: 0,
			expectedWorkers: 2,
			expectErr:       false,
		},
		{
			name:            "missing worker node",
			nodes:           []*corev1.Node{newNode("worker-0", "worker", true)},
			expectedMasters: 0,
			expectedWorkers: 2,
			expectErr:       true,
			errContains:     "expected 2 worker node(s) to be Ready but only 1 found",
		},
		{
			name: "masters and workers all ready",
			nodes: []*corev1.Node{
				newNode("master-0", "master", true), newNode("master-1", "master", true), newNode("master-2", "master", true),
				newNode("worker-0", "worker", true), newNode("worker-1", "worker", true),
			},
			expectedMasters: 3,
			expectedWorkers: 2,
			expectErr:       false,
		},
		{
			name: "masters ready but workers missing",
			nodes: []*corev1.Node{
				newNode("master-0", "master", true), newNode("master-1", "master", true), newNode("master-2", "master", true),
				newNode("worker-0", "worker", true),
			},
			expectedMasters: 3,
			expectedWorkers: 2,
			expectErr:       true,
			errContains:     "expected 2 worker node(s) to be Ready but only 1 found",
		},
		{
			name:            "no nodes exist at all",
			nodes:           []*corev1.Node{},
			expectedMasters: 3,
			expectedWorkers: 0,
			expectErr:       true,
			errContains:     "expected 3 master node(s) to be Ready but only 0 found",
		},
		{
			name:            "zero expectations succeeds with no nodes",
			nodes:           []*corev1.Node{},
			expectedMasters: 0,
			expectedWorkers: 0,
			expectErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtimeObjs := make([]runtime.Object, len(tt.nodes))
			for i, n := range tt.nodes {
				runtimeObjs[i] = n
			}
			client := fake.NewSimpleClientset(runtimeObjs...)

			err := verifyExpectedNodesWithClient(context.Background(), client, tt.expectedMasters, tt.expectedWorkers)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if tt.expectErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q but got: %v", tt.errContains, err)
				}
			}
		})
	}
}
