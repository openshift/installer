## AWS Policy Equivalence Package

This package checks for structural equivalence of two AWS policy documents. See [Godoc](https://pkg.go.dev/github.com/hashicorp/awspolicyequivalence) for more information on usage.

### Post v1.5 Validation vs. Equivalence

In versions 1.5 and earlier, this package has had a validation role. For example, `{}` is a valid JSON but an invalid AWS policy. But, AWS emits this empty JSON in some cases. Should this package determine `{}` is equivalent to itself or throw an error and say it's _not_ equivalent to itself? Since the purpose of this package is primarily _equivalence_ and not validation, we are removing some of the validation role.

In other words, for v1.5 and earlier, `{}` is not equivalent to itself and returns an error. Post v1.5, `{}` is equivalent to itself and _does not_ return an error. **_This may impact you if you have relied on this package for validation!_**

### CI

![Go Build/Test](https://github.com/hashicorp/awspolicyequivalence/actions/workflows/go.yml/badge.svg)
