package local

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/command/format"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/plans/planfile"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/states/statemgr"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/tfdiags"
)

func (b *Local) opPlan(
	stopCtx context.Context,
	cancelCtx context.Context,
	op *backend.Operation,
	runningOp *backend.RunningOperation) {

	log.Printf("[INFO] backend/local: starting Plan operation")

	var diags tfdiags.Diagnostics

	if op.PlanFile != nil {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Can't re-plan a saved plan",
			"The plan command was given a saved plan file as its input. This command generates "+
				"a new plan, and so it requires a configuration directory as its argument.",
		))
		b.ReportResult(runningOp, diags)
		return
	}

	// Local planning requires a config, unless we're planning to destroy.
	if !op.Destroy && !op.HasConfig() {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"No configuration files",
			"Plan requires configuration to be present. Planning without a configuration would "+
				"mark everything for destruction, which is normally not what is desired. If you "+
				"would like to destroy everything, run plan with the -destroy option. Otherwise, "+
				"create a Terraform configuration file (.tf file) and try again.",
		))
		b.ReportResult(runningOp, diags)
		return
	}

	// Setup our count hook that keeps track of resource changes
	countHook := new(CountHook)
	if b.ContextOpts == nil {
		b.ContextOpts = new(terraform.ContextOpts)
	}
	old := b.ContextOpts.Hooks
	defer func() { b.ContextOpts.Hooks = old }()
	b.ContextOpts.Hooks = append(b.ContextOpts.Hooks, countHook)

	// Get our context
	tfCtx, configSnap, opState, ctxDiags := b.context(op)
	diags = diags.Append(ctxDiags)
	if ctxDiags.HasErrors() {
		b.ReportResult(runningOp, diags)
		return
	}
	// the state was locked during succesfull context creation; unlock the state
	// when the operation completes
	defer func() {
		err := op.StateLocker.Unlock(nil)
		if err != nil {
			b.ShowDiagnostics(err)
			runningOp.Result = backend.OperationFailure
		}
	}()

	runningOp.State = tfCtx.State()

	// Perform the plan in a goroutine so we can be interrupted
	var plan *plans.Plan
	var planDiags tfdiags.Diagnostics
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		log.Printf("[INFO] backend/local: plan calling Plan")
		plan, planDiags = tfCtx.Plan()
	}()

	if b.opWait(doneCh, stopCtx, cancelCtx, tfCtx, opState) {
		// If we get in here then the operation was cancelled, which is always
		// considered to be a failure.
		log.Printf("[INFO] backend/local: plan operation was force-cancelled by interrupt")
		runningOp.Result = backend.OperationFailure
		return
	}
	log.Printf("[INFO] backend/local: plan operation completed")

	diags = diags.Append(planDiags)
	if planDiags.HasErrors() {
		b.ReportResult(runningOp, diags)
		return
	}

	// Record whether this plan includes any side-effects that could be applied.
	runningOp.PlanEmpty = plan.Changes.Empty()

	// Save the plan to disk
	if path := op.PlanOutPath; path != "" {
		if op.PlanOutBackend == nil {
			// This is always a bug in the operation caller; it's not valid
			// to set PlanOutPath without also setting PlanOutBackend.
			diags = diags.Append(fmt.Errorf(
				"PlanOutPath set without also setting PlanOutBackend (this is a bug in Terraform)"),
			)
			b.ReportResult(runningOp, diags)
			return
		}
		plan.Backend = *op.PlanOutBackend

		// We may have updated the state in the refresh step above, but we
		// will freeze that updated state in the plan file for now and
		// only write it if this plan is subsequently applied.
		plannedStateFile := statemgr.PlannedStateUpdate(opState, plan.State)

		log.Printf("[INFO] backend/local: writing plan output to: %s", path)
		err := planfile.Create(path, configSnap, plannedStateFile, plan)
		if err != nil {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Failed to write plan file",
				fmt.Sprintf("The plan file could not be written: %s.", err),
			))
			b.ReportResult(runningOp, diags)
			return
		}
	}

	// Perform some output tasks if we have a CLI to output to.
	if b.CLI != nil {
		schemas := tfCtx.Schemas()

		if runningOp.PlanEmpty {
			b.CLI.Output("\n" + b.Colorize().Color(strings.TrimSpace(planNoChanges)))
			// Even if there are no changes, there still could be some warnings
			b.ShowDiagnostics(diags)
			return
		}

		b.renderPlan(plan, plan.State, schemas)

		// If we've accumulated any warnings along the way then we'll show them
		// here just before we show the summary and next steps. If we encountered
		// errors then we would've returned early at some other point above.
		b.ShowDiagnostics(diags)

		// Give the user some next-steps, unless we're running in an automation
		// tool which is presumed to provide its own UI for further actions.
		if !b.RunningInAutomation {

			b.CLI.Output("\n------------------------------------------------------------------------")

			if path := op.PlanOutPath; path == "" {
				b.CLI.Output(fmt.Sprintf(
					"\n" + strings.TrimSpace(planHeaderNoOutput) + "\n",
				))
			} else {
				b.CLI.Output(fmt.Sprintf(
					"\n"+strings.TrimSpace(planHeaderYesOutput)+"\n",
					path, path,
				))
			}
		}
	}
}

func (b *Local) renderPlan(plan *plans.Plan, baseState *states.State, schemas *terraform.Schemas) {
	RenderPlan(plan, baseState, schemas, b.CLI, b.Colorize())
}

// RenderPlan renders the given plan to the given UI.
//
// This is exported only so that the "terraform show" command can re-use it.
// Ideally it would be somewhere outside of this backend code so that both
// can call into it, but we're leaving it here for now in order to avoid
// disruptive refactoring.
//
// If you find yourself wanting to call this function from a third callsite,
// please consider whether it's time to do the more disruptive refactoring
// so that something other than the local backend package is offering this
// functionality.
//
// The difference between baseState and priorState is that baseState is the
// result of implicitly running refresh (unless that was disabled) while
// priorState is a snapshot of the state as it was before we took any actions
// at all. priorState can optionally be nil if the caller has only a saved
// plan and not the prior state it was built from. In that case, changes to
// output values will not currently be rendered because their prior values
// are currently stored only in the prior state. (see the docstring for
// func planHasSideEffects for why this is and when that might change)
func RenderPlan(plan *plans.Plan, baseState *states.State, schemas *terraform.Schemas, ui cli.Ui, colorize *colorstring.Colorize) {
	counts := map[plans.Action]int{}
	var rChanges []*plans.ResourceInstanceChangeSrc
	for _, change := range plan.Changes.Resources {
		if change.Action == plans.Delete && change.Addr.Resource.Resource.Mode == addrs.DataResourceMode {
			// Avoid rendering data sources on deletion
			continue
		}

		rChanges = append(rChanges, change)
		counts[change.Action]++
	}

	headerBuf := &bytes.Buffer{}
	fmt.Fprintf(headerBuf, "\n%s\n", strings.TrimSpace(planHeaderIntro))
	if counts[plans.Create] > 0 {
		fmt.Fprintf(headerBuf, "%s create\n", format.DiffActionSymbol(plans.Create))
	}
	if counts[plans.Update] > 0 {
		fmt.Fprintf(headerBuf, "%s update in-place\n", format.DiffActionSymbol(plans.Update))
	}
	if counts[plans.Delete] > 0 {
		fmt.Fprintf(headerBuf, "%s destroy\n", format.DiffActionSymbol(plans.Delete))
	}
	if counts[plans.DeleteThenCreate] > 0 {
		fmt.Fprintf(headerBuf, "%s destroy and then create replacement\n", format.DiffActionSymbol(plans.DeleteThenCreate))
	}
	if counts[plans.CreateThenDelete] > 0 {
		fmt.Fprintf(headerBuf, "%s create replacement and then destroy\n", format.DiffActionSymbol(plans.CreateThenDelete))
	}
	if counts[plans.Read] > 0 {
		fmt.Fprintf(headerBuf, "%s read (data resources)\n", format.DiffActionSymbol(plans.Read))
	}

	ui.Output(colorize.Color(headerBuf.String()))

	ui.Output("Terraform will perform the following actions:\n")

	// Note: we're modifying the backing slice of this plan object in-place
	// here. The ordering of resource changes in a plan is not significant,
	// but we can only do this safely here because we can assume that nobody
	// is concurrently modifying our changes while we're trying to print it.
	sort.Slice(rChanges, func(i, j int) bool {
		iA := rChanges[i].Addr
		jA := rChanges[j].Addr
		if iA.String() == jA.String() {
			return rChanges[i].DeposedKey < rChanges[j].DeposedKey
		}
		return iA.Less(jA)
	})

	for _, rcs := range rChanges {
		if rcs.Action == plans.NoOp {
			continue
		}

		providerSchema := schemas.ProviderSchema(rcs.ProviderAddr.Provider)
		if providerSchema == nil {
			// Should never happen
			ui.Output(fmt.Sprintf("(schema missing for %s)\n", rcs.ProviderAddr))
			continue
		}
		rSchema, _ := providerSchema.SchemaForResourceAddr(rcs.Addr.Resource.Resource)
		if rSchema == nil {
			// Should never happen
			ui.Output(fmt.Sprintf("(schema missing for %s)\n", rcs.Addr))
			continue
		}

		// check if the change is due to a tainted resource
		tainted := false
		if !baseState.Empty() {
			if is := baseState.ResourceInstance(rcs.Addr); is != nil {
				if obj := is.GetGeneration(rcs.DeposedKey.Generation()); obj != nil {
					tainted = obj.Status == states.ObjectTainted
				}
			}
		}

		ui.Output(format.ResourceChange(
			rcs,
			tainted,
			rSchema,
			colorize,
		))
	}

	// stats is similar to counts above, but:
	// - it considers only resource changes
	// - it simplifies "replace" into both a create and a delete
	stats := map[plans.Action]int{}
	for _, change := range rChanges {
		switch change.Action {
		case plans.CreateThenDelete, plans.DeleteThenCreate:
			stats[plans.Create]++
			stats[plans.Delete]++
		default:
			stats[change.Action]++
		}
	}
	ui.Output(colorize.Color(fmt.Sprintf(
		"[reset][bold]Plan:[reset] "+
			"%d to add, %d to change, %d to destroy.",
		stats[plans.Create], stats[plans.Update], stats[plans.Delete],
	)))

	// If there is at least one planned change to the root module outputs
	// then we'll render a summary of those too.
	var changedRootModuleOutputs []*plans.OutputChangeSrc
	for _, output := range plan.Changes.Outputs {
		if !output.Addr.Module.IsRoot() {
			continue
		}
		if output.ChangeSrc.Action == plans.NoOp {
			continue
		}
		changedRootModuleOutputs = append(changedRootModuleOutputs, output)
	}
	if len(changedRootModuleOutputs) > 0 {
		ui.Output(colorize.Color("[reset]\n[bold]Changes to Outputs:[reset]" + format.OutputChanges(changedRootModuleOutputs, colorize)))
	}
}

const planHeaderIntro = `
An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
`

const planHeaderNoOutput = `
Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
`

const planHeaderYesOutput = `
This plan was saved to: %s

To perform exactly these actions, run the following command to apply:
    terraform apply %q
`

const planNoChanges = `
[reset][bold][green]No changes. Infrastructure is up-to-date.[reset][green]

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
`

const planRefreshing = `
[reset][bold]Refreshing Terraform state in-memory prior to plan...[reset]
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.
`
