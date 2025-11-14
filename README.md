# fsearch

A command-line utility to search for files and folders efficiently.

---

## 🧭 Basic Usage

```bash
fsearch <search-term> [flags] <path>
```

**Arguments:**

- `<search-term>` — The filename or pattern to search for (required, first argument)
- `<path>` — The directory path to search in, relative or absolute (required, last argument)
- `[flags]` — Optional flags that modify search behavior (placed between the search term and path)

---

## ⚙️ API Reference

| Flag                | Type     | Description                                                                         | Default   | Example                                                      |
| ------------------- | -------- | ----------------------------------------------------------------------------------- | --------- | ------------------------------------------------------------ |
| `--partial`         | boolean  | Match files whose names contain the search term                                     | false     | `fsearch doc --partial ./src`                                |
| `--ignore-case`     | boolean  | Perform a case-insensitive search                                                   | false     | `fsearch README --ignore-case ./`                            |
| `--open`            | boolean  | Open the first matched file in the system’s default program                         | false     | `fsearch config.txt --open ./`                               |
| `--lines`           | integer  | Number of lines to show in preview if type is `file` and number is greater than `0` | 10        | `fsearch data.csv --lines=20 ./`                   |
| `--limit`           | integer  | Maximum number of matches to return                                                 | unlimited | `fsearch test --partial --limit=5 ./`                        |
| `--depth`           | integer  | Maximum folder depth to search                                                      | unlimited | `fsearch index --partial --depth=3 ./src`                    |
| `--ext`             | string[] | List of file extensions to include (comma-separated) without the `.`                | all       | `fsearch config --ext=txt,md,log ./`                         |
| `--exclude-ext`     | string[] | List of file extensions to exclude (comma-separated) without the `.`                | none      | `fsearch backup --exclude-ext=tmp,bak ./`                    |
| `--exclude-dir`     | string[] | List of directories to exclude (comma-separated)                                    | none      | `fsearch index --partial --exclude-dir=node_modules,.git ./` |
| `--min-size`        | string   | Minimum file size                                                                   | none      | `fsearch report --min-size=1 ./documents`                    |
| `--max-size`        | string   | Maximum file size                                                                   | none      | `fsearch config --max-size=10 ./`                            |
| `--size-type`       | string   | The type format used in size comparisons                                            | KB        | `fsearch config --max-size=10 --size-type=MB ./`             |
| `--modified-before` | string   | Include files modified before date (`YYYY-MM-DD`)                                   | none      | `fsearch log --modified-before=2024-01-01 ./logs`            |
| `--modified-after`  | string   | Include files modified after date (`YYYY-MM-DD`)                                    | none      | `fsearch report --modified-after=2024-06-01 ./`              |
| `--hidden`          | boolean  | Include hidden files and folders in search                                          | false     | `fsearch config --hidden ./`                                 |
| `--count`           | boolean  | Display only the count of matches (no file details)                                 | false     | `fsearch logs --count --partial ./`                          |
| `--regex`           | boolean  | Treat the search term as a regular expression pattern                               | false     | `fsearch "^[A-Z].*\\.js$" --regex ./src`                     |
| `--debug`           | boolean  | Show all passed flag values and environment info without performing a search        | false     | `fsearch doc --debug ./`                                     |
| `--type`            | string   | Type of item to search for — either `file` or `folder`                              | file      | `fsearch config --type=folder ./`                            |

---

## 📝 Notes

- The **search term** must always be the first argument.
- The **path** must always be the last argument.
- All **flags** should appear between the search term and the path.
- Dates for `--modified-before` and `--modified-after` must use the `YYYY-MM-DD` format.
- File size units are case-insensitive (e.g., KB, MB, GB).
- Multiple extensions or directories can be separated by commas.
- Flags can be combined for advanced filtering.
- When `--open` is used and multiple matches are found, only the first file is opened.
- Relative paths are resolved from the current working directory.
- Use `./` to search in the current directory.

---

## 💡 Examples

```bash
# Search for files containing "report" anywhere in their name
fsearch report --partial ./documents

# Search for .txt and .md files, ignoring case
fsearch notes --ignore-case --ext=txt,md ./docs

# Show only the number of matches
fsearch logs --count ./logs

# Preview first 20 lines of each match
fsearch data --lines=20 ./exports

# Search only for folders with "config" in the name
fsearch config --type=folder ./src
```

---

## 🚫 Limitations

- Does **not** search within file contents.
  Use tools like **ripgrep** or **grep** for content-based searching.
