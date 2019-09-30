# File Exists Check

Tests for the existence of one or more files matching the specified globbing pattern:

`OK`: One or more files matched the globbing pattern.

`CRITICAL`: No files matched the globbing pattern.

When the `--negate` flag is used, the test logic is inverted:

`OK`: No files matched the globbing pattern.

`CRITICAL`: One or more files matched the globbing pattern.

## Options

* `--pattern (-p)`: Filepath or globbing pattern to check for one or more existing file
* `--negate (-n)`: If set, asserts filepath or globbing pattern should not match an existing file

## Globbing patterns

[File globbing patterns specify sets of filenames using wildcard characters.][1]

This plugin currently uses golang's `filepath.Glob` and therefore [does not support the ** double-star or globstar operator][2].

## Examples

Return OK if `/etc/resolv.conf` exists:

`check_file_exists --pattern /etc/resolv.conf`

Return a critical if files with `.dmp` extension exist in `/var/tmp`:

`check_file_exists --pattern '/var/tmp/*.dmp' --negate`

[1]: https://en.wikipedia.org/wiki/Glob_%28programming%29
[2]: https://www.client9.com/golang-globs-and-the-double-star-glob-operator/