/*
Copyright 2023 The Radius Authors.

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

// Package paramresolver resolves {{context.*}} template expressions in recipe
// parameters against the Radius recipe context. It handles simple lookups,
// single-level ternary expressions, and property resolution from the application
// layer. This package is shared by both Terraform and Bicep drivers.
package paramresolver

import (
	"fmt"
	"strings"

	"github.com/radius-project/radius/pkg/recipes/recipecontext"
)

// ResolveParameters resolves all {{context.*}} template expressions in the
// given parameter map against the recipe context. Returns a new map with
// resolved values — the input map is never mutated.
//
// Resolution rules:
//   - If a parameter value is a string containing ONLY a single expression
//     (e.g., "{{context.resource.name}}"), it resolves to the original typed
//     value from the context (preserving int, bool, etc.).
//   - If a parameter value is a string containing an expression embedded in
//     surrounding text (e.g., "prefix-{{context.resource.name}}-suffix"), it
//     resolves to a concatenated string.
//   - Unrecognized expression paths are left as literal strings (unchanged).
//   - Nested maps and slices are traversed recursively.
//   - Non-string values pass through unchanged.
func ResolveParameters(params map[string]any, ctx *recipecontext.Context) map[string]any {
	if params == nil || ctx == nil {
		return params
	}

	lookup := buildContextLookup(ctx)
	result := make(map[string]any, len(params))
	for k, v := range params {
		result[k] = resolveValue(v, lookup)
	}
	return result
}

// ShallowMergeParameters merges base and override parameter maps using shallow
// merge semantics. Override keys take precedence for overlapping top-level keys,
// replacing the entire value (including nested objects). Returns a new map —
// neither input is mutated.
func ShallowMergeParameters(base map[string]any, override map[string]any) map[string]any {
	if base == nil && override == nil {
		return nil
	}

	result := make(map[string]any, len(base)+len(override))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		result[k] = v
	}
	return result
}

// contextEntry holds both the string representation and the original typed value
// for a context path. This enables typed resolution when an entire parameter
// value is a single expression.
type contextEntry struct {
	str   string
	typed any
}

// resolveValue recursively resolves {{...}} expressions in a value.
// Strings are resolved directly; maps and slices are walked recursively.
// Input values are never mutated — maps and slices are copied.
func resolveValue(val any, lookup map[string]contextEntry) any {
	switch v := val.(type) {
	case string:
		return resolveStringValue(v, lookup)
	case map[string]any:
		result := make(map[string]any, len(v))
		for k, inner := range v {
			result[k] = resolveValue(inner, lookup)
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, inner := range v {
			result[i] = resolveValue(inner, lookup)
		}
		return result
	default:
		return val
	}
}

// resolveStringValue handles string parameter values. If the entire value is a
// single expression, it returns the typed value. If embedded in text, it returns
// a concatenated string.
func resolveStringValue(input string, lookup map[string]contextEntry) any {
	// Fast path: no expressions at all
	if !strings.Contains(input, "{{") {
		return input
	}

	// Check if the entire string is a single expression (for typed resolution)
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "{{") && strings.HasSuffix(trimmed, "}}") {
		// Verify it's truly a single expression (no text before/after, no multiple expressions)
		inner := trimmed[2 : len(trimmed)-2]
		if !strings.Contains(inner, "{{") {
			// It might still contain }} inside quotes for ternary, so use proper parsing
			end := findExpressionEnd(trimmed, 2)
			if end == len(trimmed)-2 {
				// Single expression — resolve with typed value
				expr := strings.TrimSpace(inner)
				resolved := evaluateExpressionTyped(expr, lookup)
				if resolved != nil {
					return resolved
				}
				// Unresolved — return original string
				return input
			}
		}
	}

	// Multi-expression or embedded expression — resolve as string concatenation
	return resolveExpressionsString(input, lookup)
}

// resolveExpressionsString scans the input string for {{...}} expressions and
// replaces each with its string value. Unresolved expressions are left as-is.
func resolveExpressionsString(input string, lookup map[string]contextEntry) string {
	var result strings.Builder
	i := 0
	for i < len(input) {
		start := strings.Index(input[i:], "{{")
		if start == -1 {
			result.WriteString(input[i:])
			break
		}
		result.WriteString(input[i : i+start])

		exprStart := i + start + 2
		end := findExpressionEnd(input, exprStart)
		if end == -1 {
			// No matching }} — write the rest as-is
			result.WriteString(input[i+start:])
			break
		}

		expr := strings.TrimSpace(input[exprStart:end])
		original := input[i+start : end+2]

		resolved := evaluateExpressionString(expr, lookup)
		if resolved != "" {
			result.WriteString(resolved)
		} else {
			result.WriteString(original)
		}

		i = end + 2
	}
	return result.String()
}

// evaluateExpressionTyped evaluates an expression and returns the typed value,
// or nil if unresolved.
func evaluateExpressionTyped(expr string, lookup map[string]contextEntry) any {
	if isTernaryExpression(expr) {
		result := evaluateTernary(expr, lookup)
		if result == "" {
			return nil
		}
		return result
	}

	if entry, ok := lookup[expr]; ok {
		return entry.typed
	}
	return nil
}

// evaluateExpressionString evaluates an expression and returns the string value,
// or empty string if unresolved.
func evaluateExpressionString(expr string, lookup map[string]contextEntry) string {
	if isTernaryExpression(expr) {
		return evaluateTernary(expr, lookup)
	}

	if entry, ok := lookup[expr]; ok {
		return entry.str
	}
	return ""
}

// findExpressionEnd finds the position of the closing }} starting from pos,
// skipping over quoted strings so that "}" inside literals doesn't match.
func findExpressionEnd(input string, pos int) int {
	inQuote := false
	escaped := false
	for i := pos; i < len(input)-1; i++ {
		ch := input[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inQuote {
			escaped = true
			continue
		}
		if ch == '"' {
			inQuote = !inQuote
			continue
		}
		if !inQuote && ch == '}' && input[i+1] == '}' {
			return i
		}
	}
	return -1
}

// isTernaryExpression checks whether the expression contains a ternary
// operator (== and ?) at the top level (outside quotes).
func isTernaryExpression(expr string) bool {
	hasEq := false
	hasQuestion := false
	inQuote := false
	escaped := false
	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inQuote {
			escaped = true
			continue
		}
		if ch == '"' {
			inQuote = !inQuote
			continue
		}
		if inQuote {
			continue
		}
		if ch == '=' && i+1 < len(expr) && expr[i+1] == '=' {
			hasEq = true
			i++
		}
		if ch == '?' {
			hasQuestion = true
		}
	}
	return hasEq && hasQuestion
}

// evaluateTernary evaluates a single-level ternary expression of the form:
//
//	expr == "value" ? "trueResult" : "falseResult"
//
// Per V1 spec, only single-level ternary is supported. If the else branch
// contains another ternary (nested/chained), the entire expression is left
// unresolved (returns "").
func evaluateTernary(expr string, lookup map[string]contextEntry) string {
	eqIdx := findTopLevel(expr, "==")
	if eqIdx == -1 {
		return ""
	}

	lhs := strings.TrimSpace(expr[:eqIdx])
	rest := strings.TrimSpace(expr[eqIdx+2:])

	qIdx := findTopLevel(rest, "?")
	if qIdx == -1 {
		return ""
	}

	rhs := strings.TrimSpace(rest[:qIdx])
	afterQ := strings.TrimSpace(rest[qIdx+1:])

	colonIdx := findMatchingColon(afterQ)
	if colonIdx == -1 {
		return ""
	}

	trueVal := strings.TrimSpace(afterQ[:colonIdx])
	elseExpr := strings.TrimSpace(afterQ[colonIdx+1:])

	// V1: nested/chained ternaries are out of scope — leave unresolved
	if isTernaryExpression(elseExpr) {
		return ""
	}

	// Resolve lhs from lookup
	entry, ok := lookup[lhs]
	if !ok {
		return ""
	}

	// rhs must be a quoted string literal
	rhsVal := unquote(rhs)
	if rhsVal == "" && rhs != `""` {
		return ""
	}

	if entry.str == rhsVal {
		return unquote(trueVal)
	}

	return unquote(elseExpr)
}

// findTopLevel finds the first occurrence of needle in s that is outside
// quoted strings. Returns the index or -1.
func findTopLevel(s, needle string) int {
	inQuote := false
	escaped := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inQuote {
			escaped = true
			continue
		}
		if ch == '"' {
			inQuote = !inQuote
			continue
		}
		if !inQuote && i+len(needle) <= len(s) && s[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

// findMatchingColon finds the top-level : that matches the first ? in a
// ternary expression. It tracks nested ? : depth and respects quotes.
func findMatchingColon(s string) int {
	depth := 0
	inQuote := false
	escaped := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inQuote {
			escaped = true
			continue
		}
		if ch == '"' {
			inQuote = !inQuote
			continue
		}
		if inQuote {
			continue
		}
		if ch == '?' {
			depth++
		}
		if ch == ':' {
			if depth == 0 {
				return i
			}
			depth--
		}
	}
	return -1
}

// unquote removes surrounding double quotes from a string literal.
func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// buildContextLookup creates a flat map of expression paths to their values
// from the recipe context.
func buildContextLookup(ctx *recipecontext.Context) map[string]contextEntry {
	vals := map[string]contextEntry{
		"context.resource.name": {str: ctx.Resource.Name, typed: ctx.Resource.Name},
		"context.resource.id":   {str: ctx.Resource.ID, typed: ctx.Resource.ID},
		"context.resource.type": {str: ctx.Resource.Type, typed: ctx.Resource.Type},

		"context.application.name": {str: ctx.Application.Name, typed: ctx.Application.Name},
		"context.application.id":   {str: ctx.Application.ID, typed: ctx.Application.ID},

		"context.environment.name": {str: ctx.Environment.Name, typed: ctx.Environment.Name},
		"context.environment.id":   {str: ctx.Environment.ID, typed: ctx.Environment.ID},
	}

	// Flatten resource properties so expressions like
	// context.resource.properties.size resolve at deploy time.
	if ctx.Resource.Properties != nil {
		flattenProperties("context.resource.properties", ctx.Resource.Properties, vals)
	}

	if ctx.Runtime.Kubernetes != nil {
		vals["context.runtime.kubernetes.namespace"] = contextEntry{
			str:   ctx.Runtime.Kubernetes.Namespace,
			typed: ctx.Runtime.Kubernetes.Namespace,
		}
		vals["context.runtime.kubernetes.environmentNamespace"] = contextEntry{
			str:   ctx.Runtime.Kubernetes.EnvironmentNamespace,
			typed: ctx.Runtime.Kubernetes.EnvironmentNamespace,
		}
	}

	if ctx.Azure != nil {
		vals["context.azure.resourceGroup.name"] = contextEntry{
			str:   ctx.Azure.ResourceGroup.Name,
			typed: ctx.Azure.ResourceGroup.Name,
		}
		vals["context.azure.resourceGroup.id"] = contextEntry{
			str:   ctx.Azure.ResourceGroup.ID,
			typed: ctx.Azure.ResourceGroup.ID,
		}
		vals["context.azure.subscription.subscriptionId"] = contextEntry{
			str:   ctx.Azure.Subscription.SubscriptionID,
			typed: ctx.Azure.Subscription.SubscriptionID,
		}
	}

	if ctx.AWS != nil {
		vals["context.aws.region"] = contextEntry{
			str:   ctx.AWS.Region,
			typed: ctx.AWS.Region,
		}
		vals["context.aws.account"] = contextEntry{
			str:   ctx.AWS.Account,
			typed: ctx.AWS.Account,
		}
	}

	return vals
}

// flattenProperties recursively adds map entries as dot-notation paths into the
// lookup table. Nested map[string]any values are traversed; leaf values are
// stored with their original type preserved.
func flattenProperties(prefix string, props map[string]any, vals map[string]contextEntry) {
	for k, v := range props {
		path := prefix + "." + k
		switch typed := v.(type) {
		case map[string]any:
			flattenProperties(path, typed, vals)
		case string:
			vals[path] = contextEntry{str: typed, typed: typed}
		case nil:
			vals[path] = contextEntry{str: "", typed: nil}
		default:
			vals[path] = contextEntry{str: fmt.Sprintf("%v", typed), typed: typed}
		}
	}
}
