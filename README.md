# Go package that mimics systemd v253 EnvironmentFile option

The parser is borrowed from `systemd` `v253` as is. Despite the original parser slightly oversimplify and allows to do weird things, see tests.

## Motivation

Common approach is to use environment variables to configure [`golang`](https://go.dev/) programs.
And [`systemd`](https://systemd.io/) is the most widespread system and service manager.
It is convenient to use literally the same file as environment holder at debugging time and
right as [`EnvironmentFile`](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#EnvironmentFile=) in `systemd` `.service`-file.

## Synopses

> Similar to `Environment=`, but reads the environment variables from
a text file. The text file should contain newline-separated variable assignments. Empty lines, lines
without an `=` separator, or lines starting with `;` or
`#` will be ignored, which may be used for commenting. The file must be UTF-8
encoded. Valid characters are
[unicode scalar values](https://www.unicode.org/glossary/#unicode_scalar_value) other than
[noncharacters](https://www.unicode.org/glossary/#noncharacter), `U+0000` `NUL`, and
`U+FEFF` [byte order mark](https://www.unicode.org/glossary/#byte_order_mark).
Control codes other than `NUL` are allowed.
>
> In the file, an unquoted value after the `=` is parsed with the same backslash-escape
rules as
[unquoted text](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_02_01)
in a POSIX shell, but unlike in a shell, interior whitespace is preserved and quotes after the
first non-whitespace character are preserved. Leading and trailing whitespace (space, tab, carriage return) is
discarded, but interior whitespace within the line is preserved verbatim. A line ending with a backslash will be
continued to the following one, with the newline itself discarded. A backslash
`\` followed by any character other than newline will preserve the following character, so that
`\\` will become the value `\`.
>
> In the file, a `'`-quoted value after the `=` can span multiple lines
and contain any character verbatim other than single quote, like
[single-quoted text](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_02_02)
in a POSIX shell. No backslash-escape sequences are recognized. Leading and trailing whitespace
outside of the single quotes is discarded.
>
> In the file, a `"`-quoted value after the `=` can span multiple lines,
and the same escape sequences are recognized as in
[double-quoted text](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_02_03)
of a POSIX shell. Backslash (`\`) followed by any of `"` `\` `` ` `` `$` will
preserve that character. A backslash followed by newline is a line continuation, and the newline itself is
discarded. A backslash followed by any other character is ignored; both the backslash and the following
character are preserved verbatim. Leading and trailing whitespace outside of the double quotes is
discarded.
>
> The argument passed should be an absolute filename or wildcard expression, optionally prefixed with
`-`, which indicates that if the file does not exist, it will not be read and no error or
warning message is logged. This option may be specified more than once in which case all specified files are
read. If the empty string is assigned to this option, the list of file to read is reset, all prior assignments
have no effect.
>
> The files listed with this directive will be read shortly before the process is executed (more
specifically, after all processes from a previous unit state terminated. This means you can generate these
files in one unit state, and read it with this option in the next. The files are read from the file
system of the service manager, before any file system changes like bind mounts take place).
>
> Settings from these files override settings made with `Environment=`. If the same
variable is set twice from these files, the files will be read in the order they are specified and the later
setting will override the earlier setting.

[systemd documentation](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#EnvironmentFile=)

## Import

```sh
go get github.com/michurin/systemd-env-file@latest
```

```go
import "github.com/michurin/systemd-env-file/sdenv"
```

## Binary

### Install

Install in standard go way

```sh
go install github.com/michurin/systemd-env-file/cmd/xenv@latest
```

The binary will be installed in the directory named by the `GOBIN` environment variable,
which defaults to `$GOPATH/bin` or `$HOME/go/bin` if the `GOPATH` environment variable is not set.

Build manually and install to custom place

```sh
go build ./cmd/...
install xenv /opt/bin # use your favorite options
```

### Usage

```sh
echo 'TEST = "OK"' >xenv.env
./xenv sh -c 'echo $TEST'
OK
./xenv env | grep TEST
TEST=OK
```

## Links

- [`parse_env_file_internal`](https://github.com/systemd/systemd/blob/v253/src/basic/env-file.c#L22) — `systemd` implementation
- Useful constants: [[1](https://github.com/systemd/systemd/blob/v253/src/basic/string-util.h#L13)], [[2](https://github.com/systemd/systemd/blob/v253/src/basic/escape.h#L15)]