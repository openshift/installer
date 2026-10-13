/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package cel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	asometrics "github.com/Azure/azure-service-operator/v2/internal/metrics"
	"github.com/Azure/azure-service-operator/v2/internal/set"
)

// Note: if we ever make a breaking change here, we should allow users to control which "version" of CEL parsing they
// Note: get with an annotation or property on the operatorSpec.

const (
	SelfIdent   = "self"
	SecretIdent = "secret"
)

var evaluator ExpressionEvaluator

// Evaluator returns the default expression evaluator
func Evaluator() ExpressionEvaluator {
	return evaluator
}

// TODO: We could put this instead in genruntime if we thought it made more sense to keep the "static value for use by webhooks"
// TODO: there.
// RegisterEvaluator registers an evaluator as the default expression evaluator.
func RegisterEvaluator(e ExpressionEvaluator) {
	evaluator = e
}

// ExpressionEvaluator defines an interface for evaluating expressions based on self (the resource) and an optional
// secret parameter (containing the secrets)
type ExpressionEvaluator interface {
	CompileAndRun(expression string, self any, secret map[string]string) (*ExpressionResult, error)
	Check(expression string, self any) (*cel.Type, error)
	FindSecretUsage(expression string, self any) (set.Set[string], error)
	Start()
	Stop()
}

type ExpressionEvaluatorOption func(e *expressionEvaluator) (*expressionEvaluator, error)

// Cache configures the program cache for the expression evaluator.
// The program cache is used to avoid re-compiling expressions or creating
// cel.Envs more than necessary.
func Cache(cache ProgramCacher) ExpressionEvaluatorOption {
	return func(e *expressionEvaluator) (*expressionEvaluator, error) {
		e.programCache = cache
		return e, nil
	}
}

// Metrics configures CEL prometheus metrics
func Metrics(metrics asometrics.CEL) ExpressionEvaluatorOption {
	return func(e *expressionEvaluator) (*expressionEvaluator, error) {
		e.metrics = metrics
		return e, nil
	}
}

// Log configures CEL logging
func Log(log logr.Logger) ExpressionEvaluatorOption {
	return func(e *expressionEvaluator) (*expressionEvaluator, error) {
		e.log = log
		return e, nil
	}
}

var _ ExpressionEvaluator = &expressionEvaluator{}

type expressionEvaluator struct {
	programCache ProgramCacher
	metrics      asometrics.CEL
	log          logr.Logger
}

func NewExpressionEvaluator(
	opts ...ExpressionEvaluatorOption,
) (ExpressionEvaluator, error) {
	result := &expressionEvaluator{}

	for _, opt := range opts {
		var err error
		result, err = opt(result)
		if err != nil {
			return nil, err
		}
	}

	if result.metrics == nil {
		// No metrics by default
		result.metrics = asometrics.NewCELNoOp()
	}

	if result.programCache == nil {
		// Return an error if metrics wasn't configured
		envCache := NewEnvCache(result.metrics, result.log, NewEnv)
		result.programCache = NewProgramCache(envCache, result.metrics, result.log, Compile)
	}

	return result, nil
}

// Start starts the expressionEvaluator.
func (e *expressionEvaluator) Start() {
	go e.programCache.Start()
}

// Stop stops the expressionEvaluator.
func (e *expressionEvaluator) Stop() {
	e.programCache.Stop()
}

// NewEnv returns a new cel.Env accepting two parameters:
//
// self: The resource itself. The resulting Env contains knowledge of all types referenced
// by self, meaning CEL expressions involving self must all pass the type-checker.
//
// secret: Optional (may be nil/empty), a map[string]string containing secrets retrieved from Azure.
// The keys of the secret are the same as the keys in the operatorSpec.secrets structure.
func NewEnv(resource reflect.Type) (*cel.Env, error) {
	types := findTypesRecursive(resource)
	selfPath := getTypeImportPath(resource)

	// coerce between type slices
	typesList := coerceList(types)
	typesList = append(typesList, ext.ParseStructField(ParseStructTag))

	// The typesList of supported env functions was taken from Kubernetes here:
	// https://kubernetes.io/docs/reference/using-api/cel/#language-overview.

	// Note: Adding to the below list of supported functions/features should be done with caution,
	// and any updates should be propagated to docs/hugo/content/guide/expressions.md as well.
	return cel.NewEnv(
		ext.Strings(ext.StringsVersion(3)), // At the time of writing, 3 was the latest version
		cel.DefaultUTCTimeZone(true),
		// Kubernetes specifies HomogenousAggregateLiterals, but it could cause issues for
		// construction of dynamic JSON payloads. This isn't a common operation, but then again
		// neither is having Heterogeneous Aggregate (list|map) literals, as we currently restrict output
		// types to string or map[string]string.
		// cel.HomogeneousAggregateLiterals(),
		cel.EagerlyValidateDeclarations(true),
		cel.OptionalTypes(),
		cel.CrossTypeNumericComparisons(true),

		ext.NativeTypes(typesList...),
		newJSONProvider(), // This must come after ext.NativeTypes

		// TODO: We could consider adding support for the Kubernetes List Library
		// TODO: https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-list-library
		// TODO: and the Kubernetes Regex Library
		// TODO: https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-regex-library
		// TODO: if there is user need. Adding them should be non-breaking.
		cel.Variable("self", cel.ObjectType(selfPath)),
		cel.Variable("secret", cel.MapType(cel.StringType, cel.StringType)),
	)
	// cel.ClearMacros can be used to disable macros.
	// See https://github.com/google/cel-spec/blob/master/doc/langdef.md#macros
	// In our case many of the macros such as has, all, etc may be useful so we don't
	// disable them. We don't support any custom macros though, only the builtins.
}

type ExpressionResult struct {
	Value  string
	Values map[string]string
}

// CompileAndRun compiles the specified expression and returns the result of the compilation.
// expression is a CEL expression that must return either a string or map[string]string.
// self is the resource being reconciled.
// secret is the set of secrets associated with the resource (which may be empty).
func (e *expressionEvaluator) CompileAndRun(expression string, self any, secret map[string]string) (*ExpressionResult, error) {
	if self == nil {
		return nil, errors.New("self cannot be nil")
	}

	// Cache lookup also compiles and checks the program for errors
	program, err := e.programCache.Get(reflect.TypeOf(self), expression)
	if err != nil {
		return nil, err
	}

	input := map[string]any{
		SelfIdent:   self,
		SecretIdent: secret,
	}
	// TODO: We may want to use ContextEval here, alongside prgm, err := env.cel.Program(ast, cel.InterruptCheckFrequency(10))
	out, _, err := program.Program.Eval(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to eval CEL expression: %q", expression)
	}

	// We could use out.Type(), which is a CEL type, rather than out.Value(), which is a Go value, but there's not
	// much practical reason to do that since we've already evaluated the expression.
	outStr, strOK := out.Value().(string)
	outMap, mapErr := valToMap(out)

	if strOK {
		return &ExpressionResult{
			Value: outStr,
		}, nil
	} else if mapErr == nil {
		return &ExpressionResult{
			Values: outMap,
		}, nil
	}

	// If we get here, the result type was something unexpected
	return nil, makeUnexpectedResultError(program.AST, AllowedOutputTypes()...)
}

// Check compiles the expression and ensures that it type-checks and passes output requirements.
// Returns the cel OutputType of the expression.
func (e *expressionEvaluator) Check(expression string, self any) (*cel.Type, error) {
	if self == nil {
		return nil, errors.New("self cannot be nil")
	}

	// Cache lookup also compiles, checks the program for errors, and ensures that it outputs an allowed type.
	program, err := e.programCache.Get(reflect.TypeOf(self), expression)
	if err != nil {
		return nil, err
	}

	return program.AST.OutputType(), nil
}

// FindSecretUsage finds CEL expressions using `secret.xyz` and returns the set of secret values (xyz in the example).
func (e *expressionEvaluator) FindSecretUsage(expression string, self any) (set.Set[string], error) {
	// Cache lookup also compiles and checks the program for errors
	program, err := e.programCache.Get(reflect.TypeOf(self), expression)
	if err != nil {
		return nil, err
	}

	result := findExprSelectedFields(SecretIdent, program.AST.NativeRep().Expr())
	return result, nil
}

func parseStructTag(field reflect.StructField, tag string) string {
	tag, found := field.Tag.Lookup(tag)
	if found {
		splits := strings.Split(tag, ",")
		if len(splits) > 0 {
			// We make the assumption that the leftmost entry in the tag is the name.
			// This seems to be true for most tags that have the concept of a name/key, such as:
			// https://pkg.go.dev/encoding/xml#Marshal
			// https://pkg.go.dev/encoding/json#Marshal
			// https://pkg.go.dev/go.mongodb.org/mongo-driver/bson#hdr-Structs
			// https://pkg.go.dev/gopkg.in/yaml.v2#Marshal
			name := splits[0]
			return name
		}
	}

	return field.Name
}

var fieldEscaper = strings.NewReplacer(
	"__", "__underscores__",
	".", "__dot__",
	"-", "__dash__",
	"/", "__slash__",
)

// These keywords come from https://github.com/google/cel-spec/blob/master/doc/langdef.md#syntax
// Note that we ONLY escape keywords that are part of CEL syntax, not keywords reserved for embedding.
var reservedKeywords = set.Make(
	"true",
	"false",
	"null",
	"in",
)

func escapeFieldName(field string) string {
	if reservedKeywords.Contains(field) {
		return fmt.Sprintf("__%s__", field)
	}

	return fieldEscaper.Replace(field)
}

func ParseStructTag(field reflect.StructField) string {
	tag := parseStructTag(field, "json")

	// We follow the same escaping sequences as Kubernetes
	// https://kubernetes.io/docs/reference/using-api/cel/#escaping
	return escapeFieldName(tag)
}

// findExprSelectedFields finds all select expressions that start with the specified ident and returns
// the set of selected values.
// For example, an expression that includes <ident>.foo and <ident>.bar would return the set {foo, bar}.
func findExprSelectedFields(ident string, expr ast.Expr) set.Set[string] {
	result := set.Set[string]{}
	findExprSelectedFieldsImpl(ident, expr, result)
	return result
}

func findExprSelectedFieldsImpl(ident string, expr ast.Expr, result set.Set[string]) {
	if expr == nil {
		// Shouldn't happen but let's just be safe
		return
	}

	switch expr.Kind() {
	case ast.CallKind:
		t := expr.AsCall()
		findExprSelectedFieldsImpl(ident, t.Target(), result)
		for _, item := range t.Args() {
			findExprSelectedFieldsImpl(ident, item, result)
		}
	case ast.UnspecifiedExprKind:
		return
	case ast.ComprehensionKind:
		t := expr.AsComprehension()
		findExprSelectedFieldsImpl(ident, t.IterRange(), result)
		findExprSelectedFieldsImpl(ident, t.AccuInit(), result)
		findExprSelectedFieldsImpl(ident, t.LoopStep(), result)
		findExprSelectedFieldsImpl(ident, t.LoopCondition(), result)
	case ast.IdentKind:
		return
	case ast.ListKind:
		t := expr.AsList()
		for _, item := range t.Elements() {
			findExprSelectedFieldsImpl(ident, item, result)
		}
	case ast.LiteralKind:
		return
	case ast.MapKind:
		t := expr.AsMap()
		for _, item := range t.Entries() {
			ent := item.AsMapEntry()
			findExprSelectedFieldsImpl(ident, ent.Value(), result)
			findExprSelectedFieldsImpl(ident, ent.Key(), result)
		}
	case ast.SelectKind:
		t := expr.AsSelect()
		op := t.Operand()
		if op.Kind() != ast.IdentKind {
			findExprSelectedFieldsImpl(ident, op, result)
		} else {
			selectIdent := op.AsIdent()
			if selectIdent == ident {
				result.Add(t.FieldName())
			}
		}
	case ast.StructKind:
		t := expr.AsStruct()
		for _, item := range t.Fields() {
			ent := item.AsStructField()
			findExprSelectedFieldsImpl(ident, ent.Value(), result)
		}
	}
}

func CheckOutputTypeAllowed(ast *cel.Ast, allowed ...*cel.Type) error {
	matched := false
	for _, t := range allowed {
		if ast.OutputType().IsExactType(t) {
			matched = true
		}
	}

	if !matched {
		return makeUnexpectedResultError(ast, allowed...)
	}

	return nil
}

func makeUnexpectedResultError(ast *cel.Ast, allowed ...*cel.Type) error {
	expectedTypes := lo.Map(
		allowed,
		func(item *cel.Type, _ int) string {
			return item.String()
		})
	expectedTypesStr := strings.Join(expectedTypes, ",")
	return errors.Errorf("expression %q must return one of [%s], but was %s", ast.Source().Content(), expectedTypesStr, ast.OutputType().String())
}

func valToMap(val ref.Val) (map[string]string, error) {
	// Convert the CEL value to a native Go map
	nativeVal, err := val.ConvertToNative(reflect.TypeOf(map[string]string{}))
	if err != nil {
		return nil, errors.Wrap(err, "error converting CEL value to native")
	}

	// Type assert the native value to a map[string]string
	result, ok := nativeVal.(map[string]string)
	if !ok {
		return nil, errors.Errorf("expectedStr map[string]string but got %T", nativeVal)
	}

	return result, nil
}

func simplePkgAlias(pkgPath string) string {
	paths := strings.Split(pkgPath, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[len(paths)-1]
}

// getResourceTypename should return the package name +
func getTypeImportPath(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return simplePkgAlias(t.PkgPath()) + "." + t.Name()
}

func findTypesRecursive(t reflect.Type) []reflect.Type { // Returns any here because that's what cel.Native expects
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var result []reflect.Type

	switch t.Kind() {
	case reflect.Struct:
		// Special case to not include v1.JSON type in the exported reflect data. This is required because in Go
		// this is just a string (a single .Raw field), but we want to transform it to a JSON map in CEL so that
		// users can use it the same way they specified it in the CRD. See json_type_provider.go.
		if t.PkgPath() == "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1" && t.Name() == "JSON" {
			break
		}

		result = append(result, t)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}

			fieldType := field.Type
			// Recurse into nested structs
			result = append(result, findTypesRecursive(fieldType)...)
		}
	case reflect.Array, reflect.Slice:
		result = append(result, findTypesRecursive(t.Elem())...)
	case reflect.Map:
		result = append(result, findTypesRecursive(t.Elem())...)
	default:
		// noop on all other types
	}

	return result
}
