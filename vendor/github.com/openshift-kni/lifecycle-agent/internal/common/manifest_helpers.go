package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func SortAndGroupByApplyWave[T metav1.Object](objs []T) ([][]T, error) {
	groups := [][]T{}
	if len(objs) == 0 {
		return groups, nil
	}
	sort.SliceStable(objs, func(i, j int) bool {
		iStr, iExist := objs[i].GetAnnotations()[ApplyWaveAnn]
		jStr, jExist := objs[j].GetAnnotations()[ApplyWaveAnn]
		// ignoring errors here since they will be handled when the slice is traversed
		iVal, _ := strconv.Atoi(iStr)
		jVal, _ := strconv.Atoi(jStr)
		if !jExist || jStr == "" {
			jVal = defaultApplyWave
		}
		if !iExist || iStr == "" {
			iVal = defaultApplyWave
		}
		if iVal != jVal {
			return iVal < jVal
		}
		return objs[i].GetName() < objs[j].GetName()
	})
	previousGroupValue := 0
	var groupValue int
	var err error
	for i, obj := range objs {
		groupStr, ok := objs[i].GetAnnotations()[ApplyWaveAnn]
		if !ok || groupStr == "" {
			groupValue = defaultApplyWave
		} else {
			groupValue, err = strconv.Atoi(groupStr)
			if err != nil {
				return groups, fmt.Errorf("error convert: %w", err)
			}
		}
		if i == 0 || groupValue != previousGroupValue {
			previousGroupValue = groupValue
			groups = append(groups, []T{})
		}
		index := len(groups) - 1
		groups[index] = append(groups[index], obj)
	}
	return groups, nil
}

func ExtractResourcesFromConfigmaps[T metav1.Object](configmaps []corev1.ConfigMap, gvk schema.GroupVersionKind) ([]T, error) {
	var objs []T

	for _, cm := range configmaps {
		for _, value := range cm.Data {
			decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(value), 4096)
			for {
				resource := unstructured.Unstructured{}
				err := decoder.Decode(&resource)
				if err != nil {
					if errors.Is(err, io.EOF) {
						// Reach the end of the data, exit the loop
						break
					}
					return nil, fmt.Errorf("failed to decode file: %w", err)
				}

				if resource.GroupVersionKind() != gvk {
					continue
				}

				obj := new(T)
				if err := runtime.DefaultUnstructuredConverter.FromUnstructured(resource.Object, obj); err != nil {
					return nil, fmt.Errorf("failed to convert to unstructured: %w", err)
				}
				objs = append(objs, *obj)
			}
		}
	}
	return objs, nil
}
