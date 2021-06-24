package remote

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform/terraform"
	tfversion "github.com/hashicorp/terraform/version"
	"github.com/mitchellh/copystructure"
)

type mockClient struct {
	Applies               *mockApplies
	ConfigurationVersions *mockConfigurationVersions
	CostEstimates         *mockCostEstimates
	Organizations         *mockOrganizations
	Plans                 *mockPlans
	PolicyChecks          *mockPolicyChecks
	Runs                  *mockRuns
	StateVersions         *mockStateVersions
	Variables             *mockVariables
	Workspaces            *mockWorkspaces
}

func newMockClient() *mockClient {
	c := &mockClient{}
	c.Applies = newMockApplies(c)
	c.ConfigurationVersions = newMockConfigurationVersions(c)
	c.CostEstimates = newMockCostEstimates(c)
	c.Organizations = newMockOrganizations(c)
	c.Plans = newMockPlans(c)
	c.PolicyChecks = newMockPolicyChecks(c)
	c.Runs = newMockRuns(c)
	c.StateVersions = newMockStateVersions(c)
	c.Variables = newMockVariables(c)
	c.Workspaces = newMockWorkspaces(c)
	return c
}

type mockApplies struct {
	client  *mockClient
	applies map[string]*tfe.Apply
	logs    map[string]string
}

func newMockApplies(client *mockClient) *mockApplies {
	return &mockApplies{
		client:  client,
		applies: make(map[string]*tfe.Apply),
		logs:    make(map[string]string),
	}
}

// create is a helper function to create a mock apply that uses the configured
// working directory to find the logfile.
func (m *mockApplies) create(cvID, workspaceID string) (*tfe.Apply, error) {
	c, ok := m.client.ConfigurationVersions.configVersions[cvID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	if c.Speculative {
		// Speculative means its plan-only so we don't create a Apply.
		return nil, nil
	}

	id := generateID("apply-")
	url := fmt.Sprintf("https://app.terraform.io/_archivist/%s", id)

	a := &tfe.Apply{
		ID:         id,
		LogReadURL: url,
		Status:     tfe.ApplyPending,
	}

	w, ok := m.client.Workspaces.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	if w.AutoApply {
		a.Status = tfe.ApplyRunning
	}

	m.logs[url] = filepath.Join(
		m.client.ConfigurationVersions.uploadPaths[cvID],
		w.WorkingDirectory,
		"apply.log",
	)
	m.applies[a.ID] = a

	return a, nil
}

func (m *mockApplies) Read(ctx context.Context, applyID string) (*tfe.Apply, error) {
	a, ok := m.applies[applyID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	// Together with the mockLogReader this allows testing queued runs.
	if a.Status == tfe.ApplyRunning {
		a.Status = tfe.ApplyFinished
	}
	return a, nil
}

func (m *mockApplies) Logs(ctx context.Context, applyID string) (io.Reader, error) {
	a, err := m.Read(ctx, applyID)
	if err != nil {
		return nil, err
	}

	logfile, ok := m.logs[a.LogReadURL]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		return bytes.NewBufferString("logfile does not exist"), nil
	}

	logs, err := ioutil.ReadFile(logfile)
	if err != nil {
		return nil, err
	}

	done := func() (bool, error) {
		a, err := m.Read(ctx, applyID)
		if err != nil {
			return false, err
		}
		if a.Status != tfe.ApplyFinished {
			return false, nil
		}
		return true, nil
	}

	return &mockLogReader{
		done: done,
		logs: bytes.NewBuffer(logs),
	}, nil
}

type mockConfigurationVersions struct {
	client         *mockClient
	configVersions map[string]*tfe.ConfigurationVersion
	uploadPaths    map[string]string
	uploadURLs     map[string]*tfe.ConfigurationVersion
}

func newMockConfigurationVersions(client *mockClient) *mockConfigurationVersions {
	return &mockConfigurationVersions{
		client:         client,
		configVersions: make(map[string]*tfe.ConfigurationVersion),
		uploadPaths:    make(map[string]string),
		uploadURLs:     make(map[string]*tfe.ConfigurationVersion),
	}
}

func (m *mockConfigurationVersions) List(ctx context.Context, workspaceID string, options tfe.ConfigurationVersionListOptions) (*tfe.ConfigurationVersionList, error) {
	cvl := &tfe.ConfigurationVersionList{}
	for _, cv := range m.configVersions {
		cvl.Items = append(cvl.Items, cv)
	}

	cvl.Pagination = &tfe.Pagination{
		CurrentPage:  1,
		NextPage:     1,
		PreviousPage: 1,
		TotalPages:   1,
		TotalCount:   len(cvl.Items),
	}

	return cvl, nil
}

func (m *mockConfigurationVersions) Create(ctx context.Context, workspaceID string, options tfe.ConfigurationVersionCreateOptions) (*tfe.ConfigurationVersion, error) {
	id := generateID("cv-")
	url := fmt.Sprintf("https://app.terraform.io/_archivist/%s", id)

	cv := &tfe.ConfigurationVersion{
		ID:        id,
		Status:    tfe.ConfigurationPending,
		UploadURL: url,
	}

	m.configVersions[cv.ID] = cv
	m.uploadURLs[url] = cv

	return cv, nil
}

func (m *mockConfigurationVersions) Read(ctx context.Context, cvID string) (*tfe.ConfigurationVersion, error) {
	cv, ok := m.configVersions[cvID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	return cv, nil
}

func (m *mockConfigurationVersions) Upload(ctx context.Context, url, path string) error {
	cv, ok := m.uploadURLs[url]
	if !ok {
		return errors.New("404 not found")
	}
	m.uploadPaths[cv.ID] = path
	cv.Status = tfe.ConfigurationUploaded
	return nil
}

type mockCostEstimates struct {
	client      *mockClient
	estimations map[string]*tfe.CostEstimate
	logs        map[string]string
}

func newMockCostEstimates(client *mockClient) *mockCostEstimates {
	return &mockCostEstimates{
		client:      client,
		estimations: make(map[string]*tfe.CostEstimate),
		logs:        make(map[string]string),
	}
}

// create is a helper function to create a mock cost estimation that uses the
// configured working directory to find the logfile.
func (m *mockCostEstimates) create(cvID, workspaceID string) (*tfe.CostEstimate, error) {
	id := generateID("ce-")

	ce := &tfe.CostEstimate{
		ID:                    id,
		MatchedResourcesCount: 1,
		ResourcesCount:        1,
		DeltaMonthlyCost:      "0.00",
		ProposedMonthlyCost:   "0.00",
		Status:                tfe.CostEstimateFinished,
	}

	w, ok := m.client.Workspaces.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	logfile := filepath.Join(
		m.client.ConfigurationVersions.uploadPaths[cvID],
		w.WorkingDirectory,
		"cost-estimate.log",
	)

	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		return nil, nil
	}

	m.logs[ce.ID] = logfile
	m.estimations[ce.ID] = ce

	return ce, nil
}

func (m *mockCostEstimates) Read(ctx context.Context, costEstimateID string) (*tfe.CostEstimate, error) {
	ce, ok := m.estimations[costEstimateID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	return ce, nil
}

func (m *mockCostEstimates) Logs(ctx context.Context, costEstimateID string) (io.Reader, error) {
	ce, ok := m.estimations[costEstimateID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	logfile, ok := m.logs[ce.ID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		return bytes.NewBufferString("logfile does not exist"), nil
	}

	logs, err := ioutil.ReadFile(logfile)
	if err != nil {
		return nil, err
	}

	ce.Status = tfe.CostEstimateFinished

	return bytes.NewBuffer(logs), nil
}

// mockInput is a mock implementation of terraform.UIInput.
type mockInput struct {
	answers map[string]string
}

func (m *mockInput) Input(ctx context.Context, opts *terraform.InputOpts) (string, error) {
	v, ok := m.answers[opts.Id]
	if !ok {
		return "", fmt.Errorf("unexpected input request in test: %s", opts.Id)
	}
	if v == "wait-for-external-update" {
		select {
		case <-ctx.Done():
		case <-time.After(time.Minute):
		}
	}
	delete(m.answers, opts.Id)
	return v, nil
}

type mockOrganizations struct {
	client        *mockClient
	organizations map[string]*tfe.Organization
}

func newMockOrganizations(client *mockClient) *mockOrganizations {
	return &mockOrganizations{
		client:        client,
		organizations: make(map[string]*tfe.Organization),
	}
}

func (m *mockOrganizations) List(ctx context.Context, options tfe.OrganizationListOptions) (*tfe.OrganizationList, error) {
	orgl := &tfe.OrganizationList{}
	for _, org := range m.organizations {
		orgl.Items = append(orgl.Items, org)
	}

	orgl.Pagination = &tfe.Pagination{
		CurrentPage:  1,
		NextPage:     1,
		PreviousPage: 1,
		TotalPages:   1,
		TotalCount:   len(orgl.Items),
	}

	return orgl, nil
}

// mockLogReader is a mock logreader that enables testing queued runs.
type mockLogReader struct {
	done func() (bool, error)
	logs *bytes.Buffer
}

func (m *mockLogReader) Read(l []byte) (int, error) {
	for {
		if written, err := m.read(l); err != io.ErrNoProgress {
			return written, err
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (m *mockLogReader) read(l []byte) (int, error) {
	done, err := m.done()
	if err != nil {
		return 0, err
	}
	if !done {
		return 0, io.ErrNoProgress
	}
	return m.logs.Read(l)
}

func (m *mockOrganizations) Create(ctx context.Context, options tfe.OrganizationCreateOptions) (*tfe.Organization, error) {
	org := &tfe.Organization{Name: *options.Name}
	m.organizations[org.Name] = org
	return org, nil
}

func (m *mockOrganizations) Read(ctx context.Context, name string) (*tfe.Organization, error) {
	org, ok := m.organizations[name]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	return org, nil
}

func (m *mockOrganizations) Update(ctx context.Context, name string, options tfe.OrganizationUpdateOptions) (*tfe.Organization, error) {
	org, ok := m.organizations[name]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	org.Name = *options.Name
	return org, nil

}

func (m *mockOrganizations) Delete(ctx context.Context, name string) error {
	delete(m.organizations, name)
	return nil
}

func (m *mockOrganizations) Capacity(ctx context.Context, name string) (*tfe.Capacity, error) {
	var pending, running int
	for _, r := range m.client.Runs.runs {
		if r.Status == tfe.RunPending {
			pending++
			continue
		}
		running++
	}
	return &tfe.Capacity{Pending: pending, Running: running}, nil
}

func (m *mockOrganizations) Entitlements(ctx context.Context, name string) (*tfe.Entitlements, error) {
	return &tfe.Entitlements{
		Operations:            true,
		PrivateModuleRegistry: true,
		Sentinel:              true,
		StateStorage:          true,
		Teams:                 true,
		VCSIntegrations:       true,
	}, nil
}

func (m *mockOrganizations) RunQueue(ctx context.Context, name string, options tfe.RunQueueOptions) (*tfe.RunQueue, error) {
	rq := &tfe.RunQueue{}

	for _, r := range m.client.Runs.runs {
		rq.Items = append(rq.Items, r)
	}

	rq.Pagination = &tfe.Pagination{
		CurrentPage:  1,
		NextPage:     1,
		PreviousPage: 1,
		TotalPages:   1,
		TotalCount:   len(rq.Items),
	}

	return rq, nil
}

type mockPlans struct {
	client *mockClient
	logs   map[string]string
	plans  map[string]*tfe.Plan
}

func newMockPlans(client *mockClient) *mockPlans {
	return &mockPlans{
		client: client,
		logs:   make(map[string]string),
		plans:  make(map[string]*tfe.Plan),
	}
}

// create is a helper function to create a mock plan that uses the configured
// working directory to find the logfile.
func (m *mockPlans) create(cvID, workspaceID string) (*tfe.Plan, error) {
	id := generateID("plan-")
	url := fmt.Sprintf("https://app.terraform.io/_archivist/%s", id)

	p := &tfe.Plan{
		ID:         id,
		LogReadURL: url,
		Status:     tfe.PlanPending,
	}

	w, ok := m.client.Workspaces.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	m.logs[url] = filepath.Join(
		m.client.ConfigurationVersions.uploadPaths[cvID],
		w.WorkingDirectory,
		"plan.log",
	)
	m.plans[p.ID] = p

	return p, nil
}

func (m *mockPlans) Read(ctx context.Context, planID string) (*tfe.Plan, error) {
	p, ok := m.plans[planID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	// Together with the mockLogReader this allows testing queued runs.
	if p.Status == tfe.PlanRunning {
		p.Status = tfe.PlanFinished
	}
	return p, nil
}

func (m *mockPlans) Logs(ctx context.Context, planID string) (io.Reader, error) {
	p, err := m.Read(ctx, planID)
	if err != nil {
		return nil, err
	}

	logfile, ok := m.logs[p.LogReadURL]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		return bytes.NewBufferString("logfile does not exist"), nil
	}

	logs, err := ioutil.ReadFile(logfile)
	if err != nil {
		return nil, err
	}

	done := func() (bool, error) {
		p, err := m.Read(ctx, planID)
		if err != nil {
			return false, err
		}
		if p.Status != tfe.PlanFinished {
			return false, nil
		}
		return true, nil
	}

	return &mockLogReader{
		done: done,
		logs: bytes.NewBuffer(logs),
	}, nil
}

type mockPolicyChecks struct {
	client *mockClient
	checks map[string]*tfe.PolicyCheck
	logs   map[string]string
}

func newMockPolicyChecks(client *mockClient) *mockPolicyChecks {
	return &mockPolicyChecks{
		client: client,
		checks: make(map[string]*tfe.PolicyCheck),
		logs:   make(map[string]string),
	}
}

// create is a helper function to create a mock policy check that uses the
// configured working directory to find the logfile.
func (m *mockPolicyChecks) create(cvID, workspaceID string) (*tfe.PolicyCheck, error) {
	id := generateID("pc-")

	pc := &tfe.PolicyCheck{
		ID:          id,
		Actions:     &tfe.PolicyActions{},
		Permissions: &tfe.PolicyPermissions{},
		Scope:       tfe.PolicyScopeOrganization,
		Status:      tfe.PolicyPending,
	}

	w, ok := m.client.Workspaces.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	logfile := filepath.Join(
		m.client.ConfigurationVersions.uploadPaths[cvID],
		w.WorkingDirectory,
		"policy.log",
	)

	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		return nil, nil
	}

	m.logs[pc.ID] = logfile
	m.checks[pc.ID] = pc

	return pc, nil
}

func (m *mockPolicyChecks) List(ctx context.Context, runID string, options tfe.PolicyCheckListOptions) (*tfe.PolicyCheckList, error) {
	_, ok := m.client.Runs.runs[runID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	pcl := &tfe.PolicyCheckList{}
	for _, pc := range m.checks {
		pcl.Items = append(pcl.Items, pc)
	}

	pcl.Pagination = &tfe.Pagination{
		CurrentPage:  1,
		NextPage:     1,
		PreviousPage: 1,
		TotalPages:   1,
		TotalCount:   len(pcl.Items),
	}

	return pcl, nil
}

func (m *mockPolicyChecks) Read(ctx context.Context, policyCheckID string) (*tfe.PolicyCheck, error) {
	pc, ok := m.checks[policyCheckID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	logfile, ok := m.logs[pc.ID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		return nil, fmt.Errorf("logfile does not exist")
	}

	logs, err := ioutil.ReadFile(logfile)
	if err != nil {
		return nil, err
	}

	switch {
	case bytes.Contains(logs, []byte("Sentinel Result: true")):
		pc.Status = tfe.PolicyPasses
	case bytes.Contains(logs, []byte("Sentinel Result: false")):
		switch {
		case bytes.Contains(logs, []byte("hard-mandatory")):
			pc.Status = tfe.PolicyHardFailed
		case bytes.Contains(logs, []byte("soft-mandatory")):
			pc.Actions.IsOverridable = true
			pc.Permissions.CanOverride = true
			pc.Status = tfe.PolicySoftFailed
		}
	default:
		// As this is an unexpected state, we say the policy errored.
		pc.Status = tfe.PolicyErrored
	}

	return pc, nil
}

func (m *mockPolicyChecks) Override(ctx context.Context, policyCheckID string) (*tfe.PolicyCheck, error) {
	pc, ok := m.checks[policyCheckID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	pc.Status = tfe.PolicyOverridden
	return pc, nil
}

func (m *mockPolicyChecks) Logs(ctx context.Context, policyCheckID string) (io.Reader, error) {
	pc, ok := m.checks[policyCheckID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	logfile, ok := m.logs[pc.ID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		return bytes.NewBufferString("logfile does not exist"), nil
	}

	logs, err := ioutil.ReadFile(logfile)
	if err != nil {
		return nil, err
	}

	switch {
	case bytes.Contains(logs, []byte("Sentinel Result: true")):
		pc.Status = tfe.PolicyPasses
	case bytes.Contains(logs, []byte("Sentinel Result: false")):
		switch {
		case bytes.Contains(logs, []byte("hard-mandatory")):
			pc.Status = tfe.PolicyHardFailed
		case bytes.Contains(logs, []byte("soft-mandatory")):
			pc.Actions.IsOverridable = true
			pc.Permissions.CanOverride = true
			pc.Status = tfe.PolicySoftFailed
		}
	default:
		// As this is an unexpected state, we say the policy errored.
		pc.Status = tfe.PolicyErrored
	}

	return bytes.NewBuffer(logs), nil
}

type mockRuns struct {
	sync.Mutex

	client     *mockClient
	runs       map[string]*tfe.Run
	workspaces map[string][]*tfe.Run

	// If modifyNewRun is non-nil, the create method will call it just before
	// saving a new run in the runs map, so that a calling test can mimic
	// side-effects that a real server might apply in certain situations.
	modifyNewRun func(client *mockClient, options tfe.RunCreateOptions, run *tfe.Run)
}

func newMockRuns(client *mockClient) *mockRuns {
	return &mockRuns{
		client:     client,
		runs:       make(map[string]*tfe.Run),
		workspaces: make(map[string][]*tfe.Run),
	}
}

func (m *mockRuns) List(ctx context.Context, workspaceID string, options tfe.RunListOptions) (*tfe.RunList, error) {
	m.Lock()
	defer m.Unlock()

	w, ok := m.client.Workspaces.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	rl := &tfe.RunList{}
	for _, run := range m.workspaces[w.ID] {
		rc, err := copystructure.Copy(run)
		if err != nil {
			panic(err)
		}
		rl.Items = append(rl.Items, rc.(*tfe.Run))
	}

	rl.Pagination = &tfe.Pagination{
		CurrentPage:  1,
		NextPage:     1,
		PreviousPage: 1,
		TotalPages:   1,
		TotalCount:   len(rl.Items),
	}

	return rl, nil
}

func (m *mockRuns) Create(ctx context.Context, options tfe.RunCreateOptions) (*tfe.Run, error) {
	m.Lock()
	defer m.Unlock()

	a, err := m.client.Applies.create(options.ConfigurationVersion.ID, options.Workspace.ID)
	if err != nil {
		return nil, err
	}

	ce, err := m.client.CostEstimates.create(options.ConfigurationVersion.ID, options.Workspace.ID)
	if err != nil {
		return nil, err
	}

	p, err := m.client.Plans.create(options.ConfigurationVersion.ID, options.Workspace.ID)
	if err != nil {
		return nil, err
	}

	pc, err := m.client.PolicyChecks.create(options.ConfigurationVersion.ID, options.Workspace.ID)
	if err != nil {
		return nil, err
	}

	r := &tfe.Run{
		ID:           generateID("run-"),
		Actions:      &tfe.RunActions{IsCancelable: true},
		Apply:        a,
		CostEstimate: ce,
		HasChanges:   false,
		Permissions:  &tfe.RunPermissions{},
		Plan:         p,
		Status:       tfe.RunPending,
		TargetAddrs:  options.TargetAddrs,
	}

	if options.Message != nil {
		r.Message = *options.Message
	}

	if pc != nil {
		r.PolicyChecks = []*tfe.PolicyCheck{pc}
	}

	if options.IsDestroy != nil {
		r.IsDestroy = *options.IsDestroy
	}

	w, ok := m.client.Workspaces.workspaceIDs[options.Workspace.ID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	if w.CurrentRun == nil {
		w.CurrentRun = r
	}

	if m.modifyNewRun != nil {
		// caller-provided callback may modify the run in-place to mimic
		// side-effects that a real server might take in some situations.
		m.modifyNewRun(m.client, options, r)
	}

	m.runs[r.ID] = r
	m.workspaces[options.Workspace.ID] = append(m.workspaces[options.Workspace.ID], r)

	return r, nil
}

func (m *mockRuns) Read(ctx context.Context, runID string) (*tfe.Run, error) {
	m.Lock()
	defer m.Unlock()

	r, ok := m.runs[runID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	pending := false
	for _, r := range m.runs {
		if r.ID != runID && r.Status == tfe.RunPending {
			pending = true
			break
		}
	}

	if !pending && r.Status == tfe.RunPending {
		// Only update the status if there are no other pending runs.
		r.Status = tfe.RunPlanning
		r.Plan.Status = tfe.PlanRunning
	}

	logs, _ := ioutil.ReadFile(m.client.Plans.logs[r.Plan.LogReadURL])
	if r.Status == tfe.RunPlanning && r.Plan.Status == tfe.PlanFinished {
		if r.IsDestroy || bytes.Contains(logs, []byte("1 to add, 0 to change, 0 to destroy")) {
			r.Actions.IsCancelable = false
			r.Actions.IsConfirmable = true
			r.HasChanges = true
			r.Permissions.CanApply = true
		}

		if bytes.Contains(logs, []byte("null_resource.foo: 1 error")) {
			r.Actions.IsCancelable = false
			r.HasChanges = false
			r.Status = tfe.RunErrored
		}
	}

	// we must return a copy for the client
	rc, err := copystructure.Copy(r)
	if err != nil {
		panic(err)
	}

	return rc.(*tfe.Run), nil
}

func (m *mockRuns) Apply(ctx context.Context, runID string, options tfe.RunApplyOptions) error {
	m.Lock()
	defer m.Unlock()

	r, ok := m.runs[runID]
	if !ok {
		return tfe.ErrResourceNotFound
	}
	if r.Status != tfe.RunPending {
		// Only update the status if the run is not pending anymore.
		r.Status = tfe.RunApplying
		r.Actions.IsConfirmable = false
		r.Apply.Status = tfe.ApplyRunning
	}
	return nil
}

func (m *mockRuns) Cancel(ctx context.Context, runID string, options tfe.RunCancelOptions) error {
	panic("not implemented")
}

func (m *mockRuns) ForceCancel(ctx context.Context, runID string, options tfe.RunForceCancelOptions) error {
	panic("not implemented")
}

func (m *mockRuns) Discard(ctx context.Context, runID string, options tfe.RunDiscardOptions) error {
	m.Lock()
	defer m.Unlock()

	r, ok := m.runs[runID]
	if !ok {
		return tfe.ErrResourceNotFound
	}
	r.Status = tfe.RunDiscarded
	r.Actions.IsConfirmable = false
	return nil
}

type mockStateVersions struct {
	client        *mockClient
	states        map[string][]byte
	stateVersions map[string]*tfe.StateVersion
	workspaces    map[string][]string
}

func newMockStateVersions(client *mockClient) *mockStateVersions {
	return &mockStateVersions{
		client:        client,
		states:        make(map[string][]byte),
		stateVersions: make(map[string]*tfe.StateVersion),
		workspaces:    make(map[string][]string),
	}
}

func (m *mockStateVersions) List(ctx context.Context, options tfe.StateVersionListOptions) (*tfe.StateVersionList, error) {
	svl := &tfe.StateVersionList{}
	for _, sv := range m.stateVersions {
		svl.Items = append(svl.Items, sv)
	}

	svl.Pagination = &tfe.Pagination{
		CurrentPage:  1,
		NextPage:     1,
		PreviousPage: 1,
		TotalPages:   1,
		TotalCount:   len(svl.Items),
	}

	return svl, nil
}

func (m *mockStateVersions) Create(ctx context.Context, workspaceID string, options tfe.StateVersionCreateOptions) (*tfe.StateVersion, error) {
	id := generateID("sv-")
	runID := os.Getenv("TFE_RUN_ID")
	url := fmt.Sprintf("https://app.terraform.io/_archivist/%s", id)

	if runID != "" && (options.Run == nil || runID != options.Run.ID) {
		return nil, fmt.Errorf("option.Run.ID does not contain the ID exported by TFE_RUN_ID")
	}

	sv := &tfe.StateVersion{
		ID:          id,
		DownloadURL: url,
		Serial:      *options.Serial,
	}

	state, err := base64.StdEncoding.DecodeString(*options.State)
	if err != nil {
		return nil, err
	}

	m.states[sv.DownloadURL] = state
	m.stateVersions[sv.ID] = sv
	m.workspaces[workspaceID] = append(m.workspaces[workspaceID], sv.ID)

	return sv, nil
}

func (m *mockStateVersions) Read(ctx context.Context, svID string) (*tfe.StateVersion, error) {
	sv, ok := m.stateVersions[svID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	return sv, nil
}

func (m *mockStateVersions) Current(ctx context.Context, workspaceID string) (*tfe.StateVersion, error) {
	w, ok := m.client.Workspaces.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	svs, ok := m.workspaces[w.ID]
	if !ok || len(svs) == 0 {
		return nil, tfe.ErrResourceNotFound
	}

	sv, ok := m.stateVersions[svs[len(svs)-1]]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	return sv, nil
}

func (m *mockStateVersions) Download(ctx context.Context, url string) ([]byte, error) {
	state, ok := m.states[url]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	return state, nil
}

type mockVariables struct {
	client     *mockClient
	workspaces map[string]*tfe.VariableList
}

var _ tfe.Variables = (*mockVariables)(nil)

func newMockVariables(client *mockClient) *mockVariables {
	return &mockVariables{
		client:     client,
		workspaces: make(map[string]*tfe.VariableList),
	}
}

func (m *mockVariables) List(ctx context.Context, workspaceID string, options tfe.VariableListOptions) (*tfe.VariableList, error) {
	vl := m.workspaces[workspaceID]
	return vl, nil
}

func (m *mockVariables) Create(ctx context.Context, workspaceID string, options tfe.VariableCreateOptions) (*tfe.Variable, error) {
	v := &tfe.Variable{
		ID:       generateID("var-"),
		Key:      *options.Key,
		Category: *options.Category,
	}
	if options.Value != nil {
		v.Value = *options.Value
	}
	if options.HCL != nil {
		v.HCL = *options.HCL
	}
	if options.Sensitive != nil {
		v.Sensitive = *options.Sensitive
	}

	workspace := workspaceID

	if m.workspaces[workspace] == nil {
		m.workspaces[workspace] = &tfe.VariableList{}
	}

	vl := m.workspaces[workspace]
	vl.Items = append(vl.Items, v)

	return v, nil
}

func (m *mockVariables) Read(ctx context.Context, workspaceID string, variableID string) (*tfe.Variable, error) {
	panic("not implemented")
}

func (m *mockVariables) Update(ctx context.Context, workspaceID string, variableID string, options tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	panic("not implemented")
}

func (m *mockVariables) Delete(ctx context.Context, workspaceID string, variableID string) error {
	panic("not implemented")
}

type mockWorkspaces struct {
	client         *mockClient
	workspaceIDs   map[string]*tfe.Workspace
	workspaceNames map[string]*tfe.Workspace
}

func newMockWorkspaces(client *mockClient) *mockWorkspaces {
	return &mockWorkspaces{
		client:         client,
		workspaceIDs:   make(map[string]*tfe.Workspace),
		workspaceNames: make(map[string]*tfe.Workspace),
	}
}

func (m *mockWorkspaces) List(ctx context.Context, organization string, options tfe.WorkspaceListOptions) (*tfe.WorkspaceList, error) {
	dummyWorkspaces := 10
	wl := &tfe.WorkspaceList{}

	// Get the prefix from the search options.
	prefix := ""
	if options.Search != nil {
		prefix = *options.Search
	}

	// Get all the workspaces that match the prefix.
	var ws []*tfe.Workspace
	for _, w := range m.workspaceIDs {
		if strings.HasPrefix(w.Name, prefix) {
			ws = append(ws, w)
		}
	}

	// Return an empty result if we have no matches.
	if len(ws) == 0 {
		wl.Pagination = &tfe.Pagination{
			CurrentPage: 1,
		}
		return wl, nil
	}

	// Return dummy workspaces for the first page to test pagination.
	if options.PageNumber <= 1 {
		for i := 0; i < dummyWorkspaces; i++ {
			wl.Items = append(wl.Items, &tfe.Workspace{
				ID:   generateID("ws-"),
				Name: fmt.Sprintf("dummy-workspace-%d", i),
			})
		}

		wl.Pagination = &tfe.Pagination{
			CurrentPage: 1,
			NextPage:    2,
			TotalPages:  2,
			TotalCount:  len(wl.Items) + len(ws),
		}

		return wl, nil
	}

	// Return the actual workspaces that matched as the second page.
	wl.Items = ws
	wl.Pagination = &tfe.Pagination{
		CurrentPage:  2,
		PreviousPage: 1,
		TotalPages:   2,
		TotalCount:   len(wl.Items) + dummyWorkspaces,
	}

	return wl, nil
}

func (m *mockWorkspaces) Create(ctx context.Context, organization string, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	if strings.HasSuffix(*options.Name, "no-operations") {
		options.Operations = tfe.Bool(false)
	} else if options.Operations == nil {
		options.Operations = tfe.Bool(true)
	}
	w := &tfe.Workspace{
		ID:         generateID("ws-"),
		Name:       *options.Name,
		Operations: *options.Operations,
		Permissions: &tfe.WorkspacePermissions{
			CanQueueApply: true,
			CanQueueRun:   true,
		},
	}
	if options.AutoApply != nil {
		w.AutoApply = *options.AutoApply
	}
	if options.VCSRepo != nil {
		w.VCSRepo = &tfe.VCSRepo{}
	}
	if options.TerraformVersion != nil {
		w.TerraformVersion = *options.TerraformVersion
	} else {
		w.TerraformVersion = tfversion.String()
	}
	m.workspaceIDs[w.ID] = w
	m.workspaceNames[w.Name] = w
	return w, nil
}

func (m *mockWorkspaces) Read(ctx context.Context, organization, workspace string) (*tfe.Workspace, error) {
	// custom error for TestRemote_plan500 in backend_plan_test.go
	if workspace == "network-error" {
		return nil, errors.New("I'm a little teacup")
	}

	w, ok := m.workspaceNames[workspace]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	return w, nil
}

func (m *mockWorkspaces) ReadByID(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	w, ok := m.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	return w, nil
}

func (m *mockWorkspaces) Update(ctx context.Context, organization, workspace string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	w, ok := m.workspaceNames[workspace]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	if options.Operations != nil {
		w.Operations = *options.Operations
	}
	if options.Name != nil {
		w.Name = *options.Name
	}
	if options.TerraformVersion != nil {
		w.TerraformVersion = *options.TerraformVersion
	}
	if options.WorkingDirectory != nil {
		w.WorkingDirectory = *options.WorkingDirectory
	}

	delete(m.workspaceNames, workspace)
	m.workspaceNames[w.Name] = w

	return w, nil
}

func (m *mockWorkspaces) UpdateByID(ctx context.Context, workspaceID string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	w, ok := m.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}

	if options.Name != nil {
		w.Name = *options.Name
	}
	if options.TerraformVersion != nil {
		w.TerraformVersion = *options.TerraformVersion
	}
	if options.WorkingDirectory != nil {
		w.WorkingDirectory = *options.WorkingDirectory
	}

	delete(m.workspaceNames, w.Name)
	m.workspaceNames[w.Name] = w

	return w, nil
}

func (m *mockWorkspaces) Delete(ctx context.Context, organization, workspace string) error {
	if w, ok := m.workspaceNames[workspace]; ok {
		delete(m.workspaceIDs, w.ID)
	}
	delete(m.workspaceNames, workspace)
	return nil
}

func (m *mockWorkspaces) DeleteByID(ctx context.Context, workspaceID string) error {
	if w, ok := m.workspaceIDs[workspaceID]; ok {
		delete(m.workspaceIDs, w.Name)
	}
	delete(m.workspaceIDs, workspaceID)
	return nil
}

func (m *mockWorkspaces) RemoveVCSConnection(ctx context.Context, organization, workspace string) (*tfe.Workspace, error) {
	w, ok := m.workspaceNames[workspace]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	w.VCSRepo = nil
	return w, nil
}

func (m *mockWorkspaces) RemoveVCSConnectionByID(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	w, ok := m.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	w.VCSRepo = nil
	return w, nil
}

func (m *mockWorkspaces) Lock(ctx context.Context, workspaceID string, options tfe.WorkspaceLockOptions) (*tfe.Workspace, error) {
	w, ok := m.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	if w.Locked {
		return nil, tfe.ErrWorkspaceLocked
	}
	w.Locked = true
	return w, nil
}

func (m *mockWorkspaces) Unlock(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	w, ok := m.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	if !w.Locked {
		return nil, tfe.ErrWorkspaceNotLocked
	}
	w.Locked = false
	return w, nil
}

func (m *mockWorkspaces) ForceUnlock(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	w, ok := m.workspaceIDs[workspaceID]
	if !ok {
		return nil, tfe.ErrResourceNotFound
	}
	if !w.Locked {
		return nil, tfe.ErrWorkspaceNotLocked
	}
	w.Locked = false
	return w, nil
}

func (m *mockWorkspaces) AssignSSHKey(ctx context.Context, workspaceID string, options tfe.WorkspaceAssignSSHKeyOptions) (*tfe.Workspace, error) {
	panic("not implemented")
}

func (m *mockWorkspaces) UnassignSSHKey(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	panic("not implemented")
}

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateID(s string) string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return s + string(b)
}
