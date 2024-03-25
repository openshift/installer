/*
Copyright 2023.

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

package common

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	kbatch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/go-logr/logr"
	cp "github.com/otiai10/copy"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift-kni/lifecycle-agent/api/v1alpha1"
)

// TODO: Need a better way to change this but will require relatively big refactoring
var OstreeDeployPathPrefix = ""

var RetryBackoffTwoMinutes = wait.Backoff{
	Steps:    120,
	Duration: time.Second,
	Factor:   1.0,
	Jitter:   0.1,
}

// GetConfigMap retrieves the configmap from cluster
func GetConfigMap(ctx context.Context, c client.Client, configMap v1alpha1.ConfigMapRef) (*corev1.ConfigMap, error) {

	cm := &corev1.ConfigMap{}
	if err := c.Get(ctx, types.NamespacedName{
		Name:      configMap.Name,
		Namespace: configMap.Namespace,
	}, cm); err != nil {
		return nil, fmt.Errorf("failed to get configMap %w", err)
	}

	return cm, nil
}

// GetConfigMaps retrieves a collection of configmaps from cluster
func GetConfigMaps(ctx context.Context, c client.Client, configMaps []v1alpha1.ConfigMapRef) ([]corev1.ConfigMap, error) {
	var cms []corev1.ConfigMap
	var cmSet = map[v1alpha1.ConfigMapRef]bool{}
	var uniqueCms []v1alpha1.ConfigMapRef

	// Remove duplicate configmaps
	for _, cm := range configMaps {
		if _, found := cmSet[cm]; !found {
			cmSet[cm] = true
			uniqueCms = append(uniqueCms, cm)
		}
	}

	for _, cm := range uniqueCms {
		existingCm, err := GetConfigMap(ctx, c, cm)
		if err != nil {
			return nil, err
		}
		cms = append(cms, *existingCm)
	}

	return cms, nil
}

// PathInsideChroot returns filepath removing host fs prefix
func PathInsideChroot(filename string) (string, error) {
	relPath, err := filepath.Rel(PathOutsideChroot("/"), filename)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path of %s inside of %s: %w", filename, PathOutsideChroot("/"), err)
	}
	return filepath.Join("/", relPath), nil
}

// PathOutsideChroot returns filepath with host fs
func PathOutsideChroot(filename string) string {
	if _, err := os.Stat(Host); err != nil {
		return filename
	}
	return filepath.Join(Host, filename)
}

func CopyOutsideChroot(src, dest string) error {
	if err := cp.Copy(PathOutsideChroot(src), PathOutsideChroot(dest)); err != nil {
		return fmt.Errorf("failed to copy outsidechroot src:%s dest:%s: %w", src, dest, err)
	}
	return nil
}

func GetStaterootPath(osname string) string {
	return fmt.Sprintf("%s/ostree/deploy/%s", OstreeDeployPathPrefix, osname)
}

// GetStaterootOptOpenshift returns the path to the `/opt/openshift` directory
// in a given stateroot. Note that since `/opt` in ostree systems is actually a
// symlink to `/var/opt`, and the `/var` directory of a stateroot is outside
// the stateroot deployment, we need to access it in this odd manner.
func GetStaterootOptOpenshift(staterootPath string) string {
	return filepath.Join(staterootPath, "var", OptOpenshift)
}

// FuncTimer check execution time
func FuncTimer(start time.Time, name string, r logr.Logger) {
	elapsed := time.Since(start)
	r.Info(fmt.Sprintf("%s took %s", name, elapsed))
}

func isConflictOrRetriable(err error) bool {
	return apierrors.IsConflict(err) || apierrors.IsInternalError(err) || apierrors.IsServiceUnavailable(err) || net.IsConnectionRefused(err)
}

func RetryOnConflictOrRetriable(backoff wait.Backoff, fn func() error) error {
	return retry.OnError(backoff, isConflictOrRetriable, fn) //nolint:wrapcheck
}

func IsRetriable(err error) bool {
	return apierrors.IsInternalError(err) || apierrors.IsServiceUnavailable(err) || net.IsConnectionRefused(err)
}

func RetryOnRetriable(backoff wait.Backoff, fn func() error) error {
	return retry.OnError(backoff, IsRetriable, fn) //nolint:wrapcheck
}

func GetDesiredStaterootName(ibu *v1alpha1.ImageBasedUpgrade) string {
	return GetStaterootName(ibu.Spec.SeedImageRef.Version)
}

func GetStaterootCertsDir(ibu *v1alpha1.ImageBasedUpgrade) string {
	return PathOutsideChroot(filepath.Join(GetStaterootOptOpenshift(GetStaterootPath(GetDesiredStaterootName(ibu))), KubeconfigCryptoDir))
}

func GetStaterootName(seedImageVersion string) string {
	return fmt.Sprintf("rhcos_%s", strings.ReplaceAll(seedImageVersion, "-", "_"))
}

func RemoveDuplicates[T comparable](list []T) []T {
	result := []T{}
	mp := make(map[T]bool)
	for _, item := range list {
		if _, present := mp[item]; present {
			continue
		}
		result = append(result, item)
		mp[item] = true
	}
	return result
}

// LogPodLogs a function to print out logs from pods generated by job. If any error it will simply log the error and exit silently
func LogPodLogs(job *kbatch.Job, log logr.Logger, clientset *kubernetes.Clientset) {
	var sinceSeconds int64 = 30 // this number is derived from medium requeue value we generally use
	podLogOpts := corev1.PodLogOptions{
		SinceSeconds: &sinceSeconds,
		Timestamps:   false,
	}

	// find pods with labels
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"job-name": job.Name}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	pods, err := clientset.CoreV1().Pods(LcaNamespace).List(context.Background(), listOptions)
	if err != nil {
		log.Info("Failed to get pod during pod log retrieval", "err", err.Error())
		return
	}

	// we expect exactly one job for now
	if len(pods.Items) == 1 {
		req := clientset.CoreV1().Pods(pods.Items[0].Namespace).GetLogs(pods.Items[0].Name, &podLogOpts)
		podLogs, err := req.Stream(context.Background())
		if err != nil {
			log.Info("Failed to get pod log", "err", err.Error())
			return
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, podLogs); err != nil {
			log.Info("Failed to copy pod log", "err", err.Error())
			return
		}
		if buf.Len() > 0 {
			log.Info(fmt.Sprintf("------ start pod `%s` log  -----", pods.Items[0].Name), "job name", job.Name)
			log.Info(buf.String())
			log.Info(fmt.Sprintf("------ end pod `%s` log  -----", pods.Items[0].Name), "job name", job.Name)
		} else {
			log.Info("No new pod logs available", "job name", job.Name, "pod name", pods.Items[0].Name)
		}
	}
}

// IsJobFinished job "finished" if it has a "Complete" or "Failed" condition marked as true.
func IsJobFinished(job *kbatch.Job) (bool, kbatch.JobConditionType) {
	for _, c := range job.Status.Conditions {
		if (c.Type == kbatch.JobComplete || c.Type == kbatch.JobFailed) && c.Status == corev1.ConditionTrue {
			return true, c.Type
		}
	}

	return false, ""
}

func GenerateDeleteOptions() *client.DeleteOptions {
	propagationPolicy := metav1.DeletePropagationBackground

	delOpt := client.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	}
	return &delOpt
}
