# cv

```
❯ cv -h
cv - Translate data

Usage:
  cv [flags]

Read data from standard input and convert it to the specified format, then write it to standard output.

The input format is specified with the -i option. If it is not specified, it will be automatically determined.
The output format is specified with the -o option. If it is not specified, it will be json.
Valid values for each format are json, yaml, toml, and csv.
For csv format, the delimiter can be specified with the -d option. If it is not specified, it will be , (comma).

  -d string
        delimiter (default ",")
  -i string
        source format (default "auto")
  -o string
        destination format (default "json")
```
