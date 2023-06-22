## v0.20.0 (December 19, 2022)

Support [`arrayType`](https://github.com/aws-cloudformation/cloudformation-resource-schema#arraytype).

## v0.19.0 (November 21, 2022)

The `Sanitize` function formats the JSON document using [`json.Indent`](https://pkg.go.dev/encoding/json#Indent) before sanitizing.

## v0.18.0 (July 28, 2022)

`Property.Maximum` and `Property.Minimum` changed to type (pointer to) `json.Number`.

## v0.17.0 (May 20, 2022)

Add helper to `propertyTransform` to get values based on path.

## v0.16.0 (February 20, 2022)

`Santize` now only removes property `pattern` regexes not supported by Go.

## v0.15.0 (January 26, 2022)

Expand properties in arrays of objects.

## v0.14.0 (October 21, 2021)

Expand properties wrapped inside `oneOf`.

## v0.13.0 (October 12, 2021)

Allow arbitrary levels of nesting when expanding Definition and Property JSON Pointer references.
Support `propertyTransform` and `tagging` keywords.

## v0.12.0 (October 1, 2021)

Support relative `file://` URLs in JSON Schema document when loading from a file path.

## v0.11.0 (September 28, 2021)

Make JSON Pointer prefix and separator constants public.

## v0.10.0 (September 17, 2021)

Add `Sanitize` function.

## v0.9.0 (September 14, 2021)

Add default field to Property.

## v0.8.0 (August 26, 2021)

Correct `PropertyJsonPointer.EqualsPath`.
Add `PropertySubschema` to handle `allOf`, `anyO`f and `oneOf` keywords.

## v0.7.0 (August 24, 2021)

Add maximum and minimum fields to Property.

## v0.6.0 (August 12, 2021)

Resolve patternProperty refs during resource expansion.

## v0.5.0 (August 10, 2021)

Allow property examples to be array of any.
Remove all attempts to rewrite patterns - It is the responsibility of the caller to deal with any regex syntax mismatches.

## v0.4.0  (August 9, 2021)

Remove any negative lookahead from patterns while loading document.

## v0.3.0 (August 5, 2021)

Correct resource handlers JSON tag.

## v0.2.0 (July 21, 2021)

Add `Property.IsRequired()`.

## v0.1.0 (April 20, 2021)

Initial release.
