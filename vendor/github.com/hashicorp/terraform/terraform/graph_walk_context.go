package terraform

import (
	"context"
	"sync"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/instances"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/provisioners"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform-plugin-sdk/tfdiags"
)

// ContextGraphWalker is the GraphWalker implementation used with the
// Context struct to walk and evaluate the graph.
type ContextGraphWalker struct {
	NullGraphWalker

	// Configurable values
	Context            *Context
	State              *states.SyncState   // Used for safe concurrent access to state
	RefreshState       *states.SyncState   // Used for safe concurrent access to state
	Changes            *plans.ChangesSync  // Used for safe concurrent writes to changes
	InstanceExpander   *instances.Expander // Tracks our gradual expansion of module and resource instances
	Operation          walkOperation
	StopContext        context.Context
	RootVariableValues InputValues

	// This is an output. Do not set this, nor read it while a graph walk
	// is in progress.
	NonFatalDiagnostics tfdiags.Diagnostics

	errorLock          sync.Mutex
	once               sync.Once
	contexts           map[string]*BuiltinEvalContext
	contextLock        sync.Mutex
	variableValues     map[string]map[string]cty.Value
	variableValuesLock sync.Mutex
	providerCache      map[string]providers.Interface
	providerSchemas    map[string]*ProviderSchema
	providerLock       sync.Mutex
	provisionerCache   map[string]provisioners.Interface
	provisionerSchemas map[string]*configschema.Block
	provisionerLock    sync.Mutex
}

func (w *ContextGraphWalker) EnterPath(path addrs.ModuleInstance) EvalContext {
	w.contextLock.Lock()
	defer w.contextLock.Unlock()

	// If we already have a context for this path cached, use that
	key := path.String()
	if ctx, ok := w.contexts[key]; ok {
		return ctx
	}

	ctx := w.EvalContext().WithPath(path)
	w.contexts[key] = ctx.(*BuiltinEvalContext)
	return ctx
}

func (w *ContextGraphWalker) EvalContext() EvalContext {
	w.once.Do(w.init)

	// Our evaluator shares some locks with the main context and the walker
	// so that we can safely run multiple evaluations at once across
	// different modules.
	evaluator := &Evaluator{
		Meta:               w.Context.meta,
		Config:             w.Context.config,
		Operation:          w.Operation,
		State:              w.State,
		Changes:            w.Changes,
		Schemas:            w.Context.schemas,
		VariableValues:     w.variableValues,
		VariableValuesLock: &w.variableValuesLock,
	}

	ctx := &BuiltinEvalContext{
		StopContext:           w.StopContext,
		Hooks:                 w.Context.hooks,
		InputValue:            w.Context.uiInput,
		InstanceExpanderValue: w.InstanceExpander,
		Components:            w.Context.components,
		Schemas:               w.Context.schemas,
		ProviderCache:         w.providerCache,
		ProviderInputConfig:   w.Context.providerInputConfig,
		ProviderLock:          &w.providerLock,
		ProvisionerCache:      w.provisionerCache,
		ProvisionerLock:       &w.provisionerLock,
		ChangesValue:          w.Changes,
		StateValue:            w.State,
		RefreshStateValue:     w.RefreshState,
		Evaluator:             evaluator,
		VariableValues:        w.variableValues,
		VariableValuesLock:    &w.variableValuesLock,
	}

	return ctx
}

func (w *ContextGraphWalker) init() {
	w.contexts = make(map[string]*BuiltinEvalContext)
	w.providerCache = make(map[string]providers.Interface)
	w.providerSchemas = make(map[string]*ProviderSchema)
	w.provisionerCache = make(map[string]provisioners.Interface)
	w.provisionerSchemas = make(map[string]*configschema.Block)
	w.variableValues = make(map[string]map[string]cty.Value)

	// Populate root module variable values. Other modules will be populated
	// during the graph walk.
	w.variableValues[""] = make(map[string]cty.Value)
	for k, iv := range w.RootVariableValues {
		w.variableValues[""][k] = iv.Value
	}
}

func (w *ContextGraphWalker) Execute(ctx EvalContext, n GraphNodeExecutable) tfdiags.Diagnostics {
	// Acquire a lock on the semaphore
	w.Context.parallelSem.Acquire()

	err := n.Execute(ctx, w.Operation)

	// Release the semaphore
	w.Context.parallelSem.Release()

	if err == nil {
		return nil
	}

	// Acquire the lock because anything is going to require a lock.
	w.errorLock.Lock()
	defer w.errorLock.Unlock()

	// If the error is non-fatal then we'll accumulate its diagnostics in our
	// non-fatal list, rather than returning it directly, so that the graph
	// walk can continue.
	if nferr, ok := err.(tfdiags.NonFatalError); ok {
		w.NonFatalDiagnostics = w.NonFatalDiagnostics.Append(nferr.Diagnostics)
		return nil
	}

	//  If we early exit, it isn't an error.
	if _, isEarlyExit := err.(EvalEarlyExitError); isEarlyExit {
		return nil
	}

	// Otherwise, we'll let our usual diagnostics machinery figure out how to
	// unpack this as one or more diagnostic messages and return that. If we
	// get down here then the returned diagnostics will contain at least one
	// error, causing the graph walk to halt.
	var diags tfdiags.Diagnostics
	diags = diags.Append(err)
	return diags

}
