# Semantic Versioning Explanation

[View original Gist on GitHub](https://gist.github.com/Integralist/c5f3bfa368cb1dc7ce51)

## Semantic Versioning Explanation.md

```
1.2.3-beta.1+meta
```

- `1`: Major
- `2`: Minor
- `3`: Patch
- `beta.1`: Pre-release
- `meta`: Metadata

## Major

The major number is incremented when the API to the package or application changes in backwards incompatible ways.

## Minor

The minor number is incremented when new features are added to the API without breaking backwards compatibility. If the major number is incremented the minor number returns to 0.

## Patch

The patch number is incremented when no new features are added but bug fixes are released. If the major or minor numbers are incremented this returns to 0.

## Pre-release

A pre-release is a . separated list of identifiers following a `-`. For example, `1.2.3-beta.1`. These are optional and are only needed for pre-release versions. In this case 1.2.3 would be a release version following a pre-release like `1.2.3-beta.1`.

## Metadata

The final section of information is build metadata. This is a . separated list of identifiers following a +. This is different from pre-release information and should be ignored when determining precedence.

