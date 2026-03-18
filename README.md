# Code Generator

`hsc` is a small CLI for generating random grouped codes.

## Install

```bash
go install github.com/cuimingda/hsc-cli/cmd/hsc@latest
```

## Usage

```bash
hsc
```

The command generates a code with 4 groups separated by `-` by default.

Show the current version:

```bash
hsc --version
```

## Rules

- The code always has 4 groups.
- `--group-size` controls the size of each group and supports `4` or `5`.
- Each group always contains exactly 2 letters.
- The remaining characters in each group are digits.
- The first character of the first group is always a letter.
- Each letter can appear at most once in a generated code.
- Letter case is randomized for generated output.
- `--underscore` switches the group separator from `-` to `_`.

## Flags

```text
      --digits string    candidate digits for generated code (digits only, no duplicates, length 1-10) (default "23456789")
      --group-size int   characters per group (allowed values: 4 or 5) (default 4)
      --letters string   candidate letters for generated code (letters only, case-insensitive deduplication, at least 8 unique letters) (default "cuimngda")
      --underscore       use _ instead of - as the group separator
  -h, --help             help for hsc
      --version          print the current hsc version
```

## Examples

Generate a code with default settings:

```bash
hsc
```

Generate 5 characters per group:

```bash
hsc --group-size 5
```

Show the current version:

```bash
hsc --version
```

Use a custom letter pool:

```bash
hsc --letters AbCdEfGhIj
```

Use a custom digit pool:

```bash
hsc --digits 0123456789
```

Use `_` as the group separator:

```bash
hsc --underscore
```

Use custom letters, digits, and group size together:

```bash
hsc --group-size 5 --letters AbCdEfGhIj --digits 0123456789
```
