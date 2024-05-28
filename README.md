# cv

```
❯ cv -h
cv - Translate data

Usage:
  cv INPUT_FORMAT OUTPUT_FORMAT [flags]

Read data from standard input and convert it to the specified format, then write it to standard output.

Valid values for INPUT_FORMAT and OUTPUT_FORMAT are:

  json, j       json
  yaml, yml, y  yaml
  toml, t       toml
  csv, c        csv
  ltsv, l       ltsv

For csv format, the delimiter can be specified with the -d option. If it is not specified, it will be , (comma).
  -d string
        delimiter (default ",")
```
