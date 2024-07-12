# dir-compare-go
A simple small tool to compare the sizes of files in two directories.
Outputs a table using [lipgloss](https://github.com/charmbracelet/lipgloss) with old and new filesizes as well as a comparison of the total difference.
The filesizes are displayed in the best matching unit of Bytes, Kibibytes, Mibibytes, etc.

## Planned expansion
Eventually the tool is supposed to include subdirectories and their contents.

## Example
```
dir-compare /path/to/source /path/to/dest

┌───────────┬────────────┬───────────┬──────────────────┐
│ File      │ Orig. Size │ New Size  │ Saved Space      │
├───────────┼────────────┼───────────┼──────────────────┤
│ SomeFile1 │ 127.1 MiB  │ 18.0 MiB  │ 109.1 MiB        │
│ SomeFile2 │ 135.7 MiB  │ 149.2 MiB │ -13.5 MiB        │
│ Total     │ 262.8 MiB  │ 167.2 MiB │ 95.6 MiB (1.6 %) │
└───────────┴────────────┴───────────┴──────────────────┘
```
