#!/bin/sh

set -ex

cd "$(dirname "$0")/.."
{
	printf '// generated with %s; do not edit\n\npackage aws\n\nvar createPermissions = []string{\n' "${0}"
	find data/data/aws -type f -exec sed -n 's|.*# AWS permission: \(.*\)|\t"\1",|p' {} \+ | sort | uniq
	printf '}\n\nvar destroyPermissions = []string{\n'
	find pkg/destroy/aws -type f -exec sed -n 's|.*// AWS permission: \(.*\)|\t"\1",|p' {} \+ | sort | uniq
	printf '}\n'
} >pkg/asset/installconfig/aws/permissions_generated.go
