// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package crdmanagement

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	"golang.org/x/exp/maps"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	ctrlleader "sigs.k8s.io/controller-runtime/pkg/leaderelection"
	"sigs.k8s.io/yaml"

	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/internal/util/match"
)

// ServiceOperatorVersionLabelOld is the label the CRDs have on them containing the ASO version. This value must match the value
// injected by config/crd/labels.yaml
const (
	ServiceOperatorVersionLabelOld = "serviceoperator.azure.com/version"
	ServiceOperatorVersionLabel    = "app.kubernetes.io/version"
	ServiceOperatorAppLabel        = "app.kubernetes.io/name"
	ServiceOperatorAppValue        = "azure-service-operator"
)

const CRDLocation = "crds"

const certMgrInjectCAFromAnnotation = "cert-manager.io/inject-ca-from"

type LeaderElector struct {
	Elector       *leaderelection.LeaderElector
	LeaseAcquired *sync.WaitGroup
	LeaseReleased *sync.WaitGroup
}

// NewLeaderElector creates a new LeaderElector
func NewLeaderElector(
	k8sConfig *rest.Config,
	log logr.Logger,
	ctrlOptions ctrl.Options,
	mgr ctrl.Manager,
) (*LeaderElector, error) {
	resourceLock, err := ctrlleader.NewResourceLock(
		k8sConfig,
		mgr,
		ctrlleader.Options{
			LeaderElection:             ctrlOptions.LeaderElection,
			LeaderElectionResourceLock: ctrlOptions.LeaderElectionResourceLock,
			LeaderElectionID:           ctrlOptions.LeaderElectionID,
		})
	if err != nil {
		return nil, err
	}

	log = log.WithName("crdManagementLeaderElector")
	leaseAcquiredWait := &sync.WaitGroup{}
	leaseAcquiredWait.Add(1)
	leaseReleasedWait := &sync.WaitGroup{}
	leaseReleasedWait.Add(1)

	var leaderContext context.Context
	var leaderContextLock sync.Mutex // used to ensure reads/writes of leaderContext are safe

	leaderElector, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:          resourceLock,
		LeaseDuration: *ctrlOptions.LeaseDuration,
		RenewDeadline: *ctrlOptions.RenewDeadline,
		RetryPeriod:   *ctrlOptions.RetryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				log.V(Status).Info("Elected leader")
				leaseAcquiredWait.Done()

				leaderContextLock.Lock()
				leaderContext = ctx
				leaderContextLock.Unlock()
			},
			OnStoppedLeading: func() {
				leaseReleasedWait.Done()

				// Cache the channel from current leader context so it can't be changed while we're using it
				leaderContextLock.Lock()
				lc := leaderContext
				leaderContextLock.Unlock()

				exitCode := 1
				if lc != nil {
					select {
					case <-lc.Done():
						exitCode = 0 // done is closed
					default:
					}
				}

				if exitCode == 0 {
					log.V(Status).Info("Lost leader due to cooperative lease release")
				} else {
					log.V(Status).Info("Lost leader")
				}
				os.Exit(exitCode)
			},
		},
		ReleaseOnCancel: ctrlOptions.LeaderElectionReleaseOnCancel,
		Name:            ctrlOptions.LeaderElectionID,
	})
	if err != nil {
		return nil, err
	}

	return &LeaderElector{
		Elector:       leaderElector,
		LeaseAcquired: leaseAcquiredWait,
		LeaseReleased: leaseReleasedWait,
	}, nil
}

type Manager struct {
	logger         logr.Logger
	kubeClient     kubeclient.Client
	leaderElection *LeaderElector

	crds []apiextensions.CustomResourceDefinition
}

// NewManager creates a new CRD manager.
// The leaderElection argument is optional, but strongly recommended.
func NewManager(logger logr.Logger, kubeClient kubeclient.Client, leaderElection *LeaderElector) *Manager {
	return &Manager{
		logger:         logger,
		kubeClient:     kubeClient,
		leaderElection: leaderElection,
	}
}

// ListCRDs lists ASO CRDs.
// This accepts a list rather than returning one to allow re-using the same list object (they're large and having multiple)
// copies of the collection results in huge memory usage.
// NOTE: This function also clears Versions[i].Schema.OpenAPIV3Schema for every CRD, as these are large and we don't use them for any of our comparisons
// DO NOT use this function to inspect the schema.
func (m *Manager) ListCRDs(ctx context.Context, list *apiextensions.CustomResourceDefinitionList) error {
	// Clear the existing list, if there is one.
	list.Items = nil
	list.Continue = ""
	list.ResourceVersion = ""

	selector := labels.NewSelector()
	requirement, err := labels.NewRequirement(ServiceOperatorAppLabel, selection.Equals, []string{ServiceOperatorAppValue})
	if err != nil {
		return err
	}
	selector = selector.Add(*requirement)

	match := client.MatchingLabelsSelector{
		Selector: selector,
	}

	err = m.kubeClient.List(ctx, list, match)
	if err != nil {
		return eris.Wrapf(err, "failed to list CRDs")
	}

	// Don't iterate by value here as we need to modify the CRDs in-place to clear the schema fields
	for i := range list.Items {
		crd := &list.Items[i]
		m.logger.V(Verbose).Info("Found an existing CRD", "CRD", crd.Name)
		// manually clear the spec.versions[].schema.openAPIV3Schema fields because they're large and we don't use them
		for j := range crd.Spec.Versions {
			crd.Spec.Versions[j].Schema.OpenAPIV3Schema = nil
		}
	}

	return nil
}

func (m *Manager) LoadOperatorCRDs(
	path string,
	podNamespace string,
	shouldLoad func(filename string) (bool, error),
) ([]apiextensions.CustomResourceDefinition, error) {
	if len(m.crds) > 0 {
		// Nothing to do as they're already loaded. Pod has to restart for them to change
		return m.crds, nil
	}

	crds, err := m.loadCRDs(path, shouldLoad)
	if err != nil {
		return nil, err
	}
	crds = m.fixCRDNamespaceRefs(crds, podNamespace)

	m.crds = crds
	return crds, nil
}

// FindMatchingCRDs finds the CRDs in "goal" that are in "existing" AND compare as equal according to the comparators with
// the corresponding CRD in "goal"
func (m *Manager) FindMatchingCRDs(
	existing []apiextensions.CustomResourceDefinition,
	goal []apiextensions.CustomResourceDefinition,
	comparators ...func(a apiextensions.CustomResourceDefinition, b apiextensions.CustomResourceDefinition) bool,
) map[string]apiextensions.CustomResourceDefinition {
	matching := make(map[string]apiextensions.CustomResourceDefinition)

	// Build a map so lookup is faster
	existingCRDs := make(map[string]apiextensions.CustomResourceDefinition, len(existing))
	for _, crd := range existing {
		existingCRDs[crd.Name] = crd
	}

	// Every goal CRD should exist and match an existing one
	for _, goalCRD := range goal {

		// Note that if the CRD is not found, we will get back a default initialized CRD.
		// We run the comparators on that as they may match, especially if the comparator is something like
		// "specs are not equal"
		existingCRD := existingCRDs[goalCRD.Name]

		// Deepcopy to ensure that modifications below don't persist
		existingCRD = *existingCRD.DeepCopy()
		goalCRD = *goalCRD.DeepCopy()

		equal := true
		for _, c := range comparators {
			if !c(existingCRD, goalCRD) { //nolint: gosimple
				equal = false
				break
			}
		}

		if equal {
			matching[goalCRD.Name] = goalCRD
		}
	}

	return matching
}

// FindNonMatchingCRDs finds the CRDs in "goal" that are not in "existing" OR are in "existing" but mismatch with the "goal"
// based on the comparator functions.
func (m *Manager) FindNonMatchingCRDs(
	existing []apiextensions.CustomResourceDefinition,
	goal []apiextensions.CustomResourceDefinition,
	comparators ...func(a apiextensions.CustomResourceDefinition, b apiextensions.CustomResourceDefinition) bool,
) map[string]apiextensions.CustomResourceDefinition {
	// Just invert the comparators and call FindMatchingCRDs
	invertedComparators := make([]func(a apiextensions.CustomResourceDefinition, b apiextensions.CustomResourceDefinition) bool, 0, len(comparators))
	for _, c := range comparators {
		c := c
		invertedComparators = append(
			invertedComparators,
			func(a apiextensions.CustomResourceDefinition, b apiextensions.CustomResourceDefinition) bool {
				return !c(a, b)
			})
	}

	return m.FindMatchingCRDs(existing, goal, invertedComparators...)
}

// DetermineCRDsToInstallOrUpgrade examines the set of goal CRDs and installed CRDs to determine the set which should
// be installed or upgraded.
func (m *Manager) DetermineCRDsToInstallOrUpgrade(
	goalCRDs []apiextensions.CustomResourceDefinition,
	existingCRDs []apiextensions.CustomResourceDefinition,
	patterns string,
) ([]*CRDInstallationInstruction, error) {
	m.logger.V(Info).Info("Goal CRDs", "count", len(goalCRDs))
	m.logger.V(Info).Info("Existing CRDs", "count", len(existingCRDs))

	// Filter the goal CRDs to only those goal CRDs that match an already installed CRD
	resultMap := make(map[string]*CRDInstallationInstruction, len(goalCRDs))
	for _, crd := range goalCRDs {
		resultMap[crd.Name] = &CRDInstallationInstruction{
			CRD: crd,
			// Assumption to begin with is that the CRD is excluded. This will get updated later if it's matched.
			FilterResult: Excluded,
			FilterReason: fmt.Sprintf("%q was not matched by CRD pattern and did not already exist in cluster", makeMatchString(crd)),
			DiffResult:   NoDifference,
		}
	}

	m.filterCRDsByExisting(existingCRDs, resultMap)
	err := m.filterCRDsByPatterns(patterns, resultMap)
	if err != nil {
		return nil, err
	}

	var filteredGoalCRDs []apiextensions.CustomResourceDefinition
	for _, result := range resultMap {
		if result.FilterResult == Excluded {
			continue
		}

		filteredGoalCRDs = append(filteredGoalCRDs, result.CRD)
	}

	goalCRDsWithDifferentVersion := m.FindNonMatchingCRDs(existingCRDs, filteredGoalCRDs, VersionEqual)
	for name := range goalCRDsWithDifferentVersion {
		result, ok := resultMap[name]
		if !ok {
			return nil, eris.Errorf("Couldn't find goal CRD %q. This is unexpected!", name)
		}

		result.DiffResult = VersionDifferent
	}

	// Collapse result to a slice
	results := maps.Values(resultMap)
	return results, nil
}

func (m *Manager) applyCRDs(
	ctx context.Context,
	goalCRDs []apiextensions.CustomResourceDefinition,
	instructions []*CRDInstallationInstruction,
	options Options,
) error {
	instructionsToApply := m.filterInstallationInstructions(instructions, true)

	if len(instructionsToApply) == 0 {
		m.logger.V(Status).Info("Successfully reconciled CRDs because there were no CRDs to update.")
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if m.leaderElection != nil {
		m.logger.V(Status).Info("Acquiring leader lock...")
		go m.leaderElection.Elector.Run(ctx)
		m.leaderElection.LeaseAcquired.Wait() // Wait for lease to be acquired

		// If lease was acquired we always want to wait til it's released, but defers run in LIFO order
		// so we need to make sure that the ctx is cancelled first here
		defer func() {
			cancel()
			m.leaderElection.LeaseReleased.Wait()
		}()

		// Double-checked locking, we need to make sure once we have the lock there's still work to do, as it may
		// already have been done while we were waiting for the lock.
		m.logger.V(Status).Info("Double-checked locking - ensure there's still CRDs to apply...")
		err := m.ListCRDs(ctx, options.ExistingCRDs)
		if err != nil {
			return eris.Wrap(err, "failed to list current CRDs")
		}
		instructions, err = m.DetermineCRDsToInstallOrUpgrade(goalCRDs, options.ExistingCRDs.Items, options.CRDPatterns)
		if err != nil {
			return eris.Wrap(err, "failed to determine CRDs to apply")
		}
		instructionsToApply = m.filterInstallationInstructions(instructions, false)
		if len(instructionsToApply) == 0 {
			m.logger.V(Status).Info("Successfully reconciled CRDs because there were no CRDs to update.")
			return nil
		}
	}

	m.logger.V(Status).Info("Will apply CRDs", "count", len(instructionsToApply))

	i := 0
	for _, instruction := range instructionsToApply {
		i += 1
		toApply := &apiextensions.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: instruction.CRD.Name,
			},
		}
		m.logger.V(Verbose).Info(
			"Applying CRD",
			"progress", fmt.Sprintf("%d/%d", i, len(instructionsToApply)),
			"crd", instruction.CRD.Name)

		result, err := controllerutil.CreateOrUpdate(ctx, m.kubeClient, toApply, func() error {
			resourceVersion := toApply.ResourceVersion
			*toApply = instruction.CRD
			toApply.ResourceVersion = resourceVersion

			return nil
		})
		if err != nil {
			// Logging here before returning because we may cancel the lease (due to defer cancel()) above, and we want to ensure
			// we can see the error before the process is killed due to lease release.
			m.logger.Error(err, "Failed to apply CRD", "name", instruction.CRD.Name)
			return eris.Wrapf(err, "failed to apply CRD %s", instruction.CRD.Name)
		}

		m.logger.V(Debug).Info("Successfully applied CRD", "name", instruction.CRD.Name, "result", result)
	}

	// Cancel the context, and wait for the lease to complete
	if m.leaderElection != nil {
		m.logger.V(Info).Info("Giving up leadership lease")
		cancel()
		m.leaderElection.LeaseReleased.Wait()
	}

	// If we make it to here, we have successfully updated all the CRDs we needed to. We need to kill the pod and let it restart so
	// that the new shape CRDs can be reconciled.
	m.logger.V(Status).Info("Restarting operator pod after updating CRDs", "count", len(instructionsToApply))
	os.Exit(0)

	// Will never get here
	return nil
}

type Options struct {
	Path         string
	Namespace    string
	CRDPatterns  string
	ExistingCRDs *apiextensions.CustomResourceDefinitionList
}

func (m *Manager) Install(ctx context.Context, options Options) error {
	shouldLoad := m.BuildCRDFileFilter(options.CRDPatterns, options.ExistingCRDs.Items)
	goalCRDs, err := m.LoadOperatorCRDs(options.Path, options.Namespace, shouldLoad)
	if err != nil {
		return eris.Wrap(err, "failed to load CRDs from disk")
	}

	installationInstructions, err := m.DetermineCRDsToInstallOrUpgrade(goalCRDs, options.ExistingCRDs.Items, options.CRDPatterns)
	if err != nil {
		return eris.Wrap(err, "failed to determine CRDs to apply")
	}

	included := IncludedCRDs(installationInstructions)
	if len(included) == 0 {
		return eris.New("No existing CRDs in cluster and no --crd-pattern specified")
	}

	// Note that this step will restart the pod when it succeeds
	// if any CRDs were applied.
	err = m.applyCRDs(ctx, goalCRDs, installationInstructions, options)
	if err != nil {
		return eris.Wrap(err, "failed to apply CRDs")
	}

	return nil
}

// loadCRDs loads CRD YAML files from the given directory.
// shouldLoad is called with the filename (e.g. "apiextensions.k8s.io_v1_customresourcedefinition_virtualnetworks.network.azure.com.yaml")
// before reading the file. If it returns false, the file is skipped entirely, avoiding the cost of reading
// and unmarshalling.
func (m *Manager) loadCRDs(path string, shouldLoad func(filename string) (bool, error)) ([]apiextensions.CustomResourceDefinition, error) {
	// Expectation is that every file in this folder is a CRD
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to read directory %s", path)
	}

	results := make([]apiextensions.CustomResourceDefinition, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Ignore directories
		}

		load, loadErr := shouldLoad(entry.Name())
		if loadErr != nil {
			return nil, eris.Wrapf(loadErr, "failed to determine if %s should be loaded", entry.Name())
		}
		if !load {
			m.logger.V(Verbose).Info("Skipping CRD file because shouldLoad returned false", "file", entry.Name())
			continue
		}

		filePath := filepath.Join(path, entry.Name())
		var content []byte
		content, err = os.ReadFile(filePath)
		if err != nil {
			return nil, eris.Wrapf(err, "failed to read %s", filePath)
		}

		crd := apiextensions.CustomResourceDefinition{}
		err = yaml.Unmarshal(content, &crd)
		if err != nil {
			return nil, eris.Wrapf(err, "failed to unmarshal %s to CRD", filePath)
		}

		m.logger.V(Verbose).Info("Loaded CRD", "crdPath", filePath, "name", crd.Name)
		results = append(results, crd)
	}

	return results, nil
}

func (m *Manager) filterInstallationInstructions(instructions []*CRDInstallationInstruction, log bool) []*CRDInstallationInstruction {
	var instructionsToApply []*CRDInstallationInstruction

	for _, item := range instructions {
		apply, reason := item.ShouldApply()
		if apply {
			instructionsToApply = append(instructionsToApply, item)
			if log {
				m.logger.V(Verbose).Info(
					"Will update CRD",
					"crd", item.CRD.Name,
					"diffResult", item.DiffResult,
					"filterReason", item.FilterReason,
					"reason", reason)
			}
		} else {
			if log {
				m.logger.V(Verbose).Info(
					"Will NOT update CRD",
					"crd", item.CRD.Name,
					"reason", reason)
			}
		}
	}

	return instructionsToApply
}

func (m *Manager) fixCRDNamespaceRefs(crds []apiextensions.CustomResourceDefinition, namespace string) []apiextensions.CustomResourceDefinition {
	results := make([]apiextensions.CustomResourceDefinition, 0, len(crds))

	for _, crd := range crds {
		crd = fixCRDNamespace(crd, namespace)
		results = append(results, crd)
	}

	return results
}

// BuildCRDFileFilter returns a predicate that determines whether a CRD file should be loaded based on
// patterns and existing CRDs. This allows skipping files that can't possibly be needed, avoiding the cost
// of reading and unmarshalling them.
//
// A file should be loaded if:
//   - Its derived CRD name matches an existing CRD in the cluster (needed for upgrades), OR
//   - Its group (from the filename) could match one of the specified patterns
//
// If there are no patterns and no existing CRDs, no files are loaded.
func (m *Manager) BuildCRDFileFilter(patterns string, existingCRDs []apiextensions.CustomResourceDefinition) func(filename string) (bool, error) {
	// Build set of existing CRD names for quick lookup
	existingCRDNames := make(map[string]bool, len(existingCRDs))
	for _, crd := range existingCRDs {
		existingCRDNames[crd.Name] = true
	}

	// Convert patterns from "group/Kind" format to just the group portion, so we can match
	// against the group extracted from the filename. For example:
	//   "containerservice.azure.com/*;network.azure.com/VirtualNetwork"
	// becomes:
	//   "containerservice.azure.com;network.azure.com"
	// TODO: if we moved v2/tools/generator/names to a shared location we could use that logic to de-pluralize the
	// TODO: kind name in the crd file path, and then match on kind + group. Not doing that now as this file loading
	// TODO: stuff is just an optimization and dealing with flect/pluralization is a bit of a pain.
	groupPatterns := patternsToGroupOnly(patterns)

	// Nothing to load; no patterns and no existing CRDs - return early
	if groupPatterns == "" && len(existingCRDNames) == 0 {
		return func(filename string) (bool, error) {
			return false, nil
		}
	}

	// Build a matcher for the group patterns
	var groupMatcher match.StringMatcher
	if groupPatterns != "" {
		groupMatcher = match.NewStringMatcher(groupPatterns)
	}

	return func(filename string) (bool, error) {
		// Check if this file's CRD is already in the cluster - if it is we always load it for potential upgrade
		crdName, err := crdNameFromFilename(filename)
		if err != nil {
			return false, err
		}
		if existingCRDNames[crdName] {
			return true, nil
		}

		// If no patterns specified, only existing CRDs matter
		if groupMatcher == nil {
			return false, nil
		}

		// Check if this file's group could match any pattern
		group, err := groupFromFilename(filename)
		if err != nil {
			return false, err
		}
		return groupMatcher.Matches(group).Matched, nil
	}
}

func (m *Manager) filterCRDsByExisting(existingCRDs []apiextensions.CustomResourceDefinition, resultMap map[string]*CRDInstallationInstruction) {
	for _, crd := range existingCRDs {
		result, ok := resultMap[crd.Name]
		if !ok {
			m.logger.V(Status).Info("Found existing CRD for which no goal CRD exists. This is unexpected!", "existing", makeMatchString(crd))
			continue
		}

		result.FilterResult = MatchedExistingCRD
		result.FilterReason = fmt.Sprintf("A CRD named %q was already installed, considering that existing CRD for update", makeMatchString(crd))
	}
}

func (m *Manager) filterCRDsByPatterns(patterns string, resultMap map[string]*CRDInstallationInstruction) error {
	if patterns == "" {
		return nil
	}

	matcher := match.NewStringMatcher(patterns)

	for _, goal := range resultMap {
		matchString := makeMatchString(goal.CRD)
		matchResult := matcher.Matches(matchString)
		if matchResult.Matched {
			goal.FilterResult = MatchedPattern
			goal.FilterReason = fmt.Sprintf("CRD named %q matched pattern %q", makeMatchString(goal.CRD), matchResult.MatchingPattern)
		}
	}

	err := matcher.WasMatched()
	if err != nil {
		return err
	}

	return nil
}

// patternsToGroupOnly takes a semicolon-separated CRD pattern string (format "group/Kind")
// and strips the "/Kind" suffix from each pattern, returning just the group portions.
// For example: "containerservice.azure.com/*;network.azure.com/VirtualNetwork"
// becomes: "containerservice.azure.com;network.azure.com"
func patternsToGroupOnly(patterns string) string {
	if patterns == "" {
		return ""
	}

	parts := strings.Split(patterns, ";")
	for i, part := range parts {
		if slashIdx := strings.Index(part, "/"); slashIdx >= 0 {
			parts[i] = part[:slashIdx]
		}
	}

	return strings.Join(parts, ";")
}

// fixCRDNamespace fixes up namespace references in the CRD to match the provided namespace.
// This could in theory be done with a string replace across the JSON representation of the CRD, but that's risky given
// we don't know what else might have the "azureserviceoperator-system" string in it. Instead, we hardcode specific places
// we know need to be fixed up. This is more brittle in the face of namespace additions but has the advantage of guaranteeing
// that we can't break our own CRDs with a string replace gone awry.
func fixCRDNamespace(crd apiextensions.CustomResourceDefinition, namespace string) apiextensions.CustomResourceDefinition {
	result := crd.DeepCopy()

	// Set spec.conversion.webhook.clientConfig.service.namespace
	if result.Spec.Conversion != nil &&
		result.Spec.Conversion.Webhook != nil &&
		result.Spec.Conversion.Webhook.ClientConfig != nil &&
		result.Spec.Conversion.Webhook.ClientConfig.Service != nil {
		result.Spec.Conversion.Webhook.ClientConfig.Service.Namespace = namespace
	}

	// Set cert-manager.io/inject-ca-from
	if len(result.Annotations) > 0 {
		if injectCAFrom, ok := result.Annotations[certMgrInjectCAFromAnnotation]; ok {
			split := strings.Split(injectCAFrom, "/")
			if len(split) == 2 {
				result.Annotations[certMgrInjectCAFromAnnotation] = fmt.Sprintf("%s/%s", namespace, split[1])
			}
		}
	}

	return *result
}

func VersionEqual(a apiextensions.CustomResourceDefinition, b apiextensions.CustomResourceDefinition) bool {
	if a.Labels == nil && b.Labels == nil {
		return true
	}

	if a.Labels == nil || b.Labels == nil {
		return false
	}

	aVersion, aOk := a.Labels[ServiceOperatorVersionLabel]
	bVersion, bOk := b.Labels[ServiceOperatorVersionLabel]

	if !aOk && !bOk {
		return true
	}

	return aVersion == bVersion
}
