package workbooksapis

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CategoryType string

const (
	CategoryTypePerformance CategoryType = "performance"
	CategoryTypeRetention   CategoryType = "retention"
	CategoryTypeTSG         CategoryType = "TSG"
	CategoryTypeWorkbook    CategoryType = "workbook"
)

func PossibleValuesForCategoryType() []string {
	return []string{
		string(CategoryTypePerformance),
		string(CategoryTypeRetention),
		string(CategoryTypeTSG),
		string(CategoryTypeWorkbook),
	}
}

func parseCategoryType(input string) (*CategoryType, error) {
	vals := map[string]CategoryType{
		"performance": CategoryTypePerformance,
		"retention":   CategoryTypeRetention,
		"tsg":         CategoryTypeTSG,
		"workbook":    CategoryTypeWorkbook,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CategoryType(input)
	return &out, nil
}

type WorkbookSharedTypeKind string

const (
	WorkbookSharedTypeKindShared WorkbookSharedTypeKind = "shared"
)

func PossibleValuesForWorkbookSharedTypeKind() []string {
	return []string{
		string(WorkbookSharedTypeKindShared),
	}
}

func parseWorkbookSharedTypeKind(input string) (*WorkbookSharedTypeKind, error) {
	vals := map[string]WorkbookSharedTypeKind{
		"shared": WorkbookSharedTypeKindShared,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkbookSharedTypeKind(input)
	return &out, nil
}

type WorkbookUpdateSharedTypeKind string

const (
	WorkbookUpdateSharedTypeKindShared WorkbookUpdateSharedTypeKind = "shared"
)

func PossibleValuesForWorkbookUpdateSharedTypeKind() []string {
	return []string{
		string(WorkbookUpdateSharedTypeKindShared),
	}
}

func parseWorkbookUpdateSharedTypeKind(input string) (*WorkbookUpdateSharedTypeKind, error) {
	vals := map[string]WorkbookUpdateSharedTypeKind{
		"shared": WorkbookUpdateSharedTypeKindShared,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkbookUpdateSharedTypeKind(input)
	return &out, nil
}
