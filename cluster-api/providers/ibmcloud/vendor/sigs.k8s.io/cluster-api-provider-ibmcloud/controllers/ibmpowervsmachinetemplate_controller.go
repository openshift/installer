/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
)

// defaultSMT is the default value of simultaneous multithreading.
const defaultSMT = 8

// IBMPowerVSMachineTemplateReconciler reconciles a IBMPowerVSMachineTemplate object.
type IBMPowerVSMachineTemplateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *IBMPowerVSMachineTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1beta2.IBMPowerVSMachineTemplate{}).
		Complete(r)
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachinetemplates,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachinetemplates/status,verbs=get;update;patch

func (r *IBMPowerVSMachineTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Reconciling IBMPowerVSMachineTemplate")

	var machineTemplate infrav1beta2.IBMPowerVSMachineTemplate
	if err := r.Get(ctx, req.NamespacedName, &machineTemplate); err != nil {
		log.Error(err, "Unable to fetch ibmpowervsmachinetemplate")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	helper, err := patch.NewHelper(&machineTemplate, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to init patch helper: %w", err)
	}

	capacity, err := getIBMPowerVSMachineCapacity(machineTemplate)
	if err != nil {
		log.Error(err, "Failed to get capacity from the ibmpowervsmachine template")
		return ctrl.Result{}, fmt.Errorf("failed to get capcity for machine template: %w", err)
	}
	log.V(3).Info("Calculated capacity for machine template", "capacity", capacity)
	if !reflect.DeepEqual(machineTemplate.Status.Capacity, capacity) {
		machineTemplate.Status.Capacity = capacity
		if err := helper.Patch(ctx, &machineTemplate); err != nil {
			if !apierrors.IsNotFound(err) {
				log.Error(err, "Failed to patch machineTemplate")
				return ctrl.Result{}, err
			}
		}
	}
	log.V(3).Info("Machine template status", "status", machineTemplate.Status.Capacity)
	return ctrl.Result{}, nil
}

func getIBMPowerVSMachineCapacity(machineTemplate infrav1beta2.IBMPowerVSMachineTemplate) (corev1.ResourceList, error) {
	capacity := make(corev1.ResourceList)
	memory := strconv.FormatInt(int64(machineTemplate.Spec.Template.Spec.MemoryGiB), 10)
	capacity[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%sG", memory))
	// There is a core-to-lCPU ratio of 1:1 for Dedicated processors. For shared processors, fractional cores round up to the nearest whole number. For example, 1.25 cores equals 2 lCPUs.
	// VM with 1 dedicated processor will see = 1 * SMT = 1 * 8 = 8 cpus in OS
	// VM with 1.5 shared processor will see = 2 * SMT = 2 * 8 = 16 cpus in OS
	// Here SMT: simultaneous multithreading which is default to 8
	// Here lCPU: number of online logical processors
	// example: on a Power VS machine with 0.5 cores
	// $ lparstat
	//	  System Configuration
	//	  type=Shared mode=Uncapped smt=8 lcpu=1 mem=33413760 kB cpus=20 ent=0.50
	var cores float64
	var err error
	switch machineTemplate.Spec.Template.Spec.Processors.Type {
	case intstr.Int:
		cores = float64(machineTemplate.Spec.Template.Spec.Processors.IntVal)
	case intstr.String:
		cores, err = strconv.ParseFloat(machineTemplate.Spec.Template.Spec.Processors.StrVal, 64)
		if err != nil {
			return nil, err
		}
	}

	virtualProcessors := fmt.Sprintf("%v", math.Ceil(cores)*defaultSMT)
	capacity[corev1.ResourceCPU] = resource.MustParse(virtualProcessors)
	return capacity, nil
}
