# f

`f` is a small utility for crafting shell commands in an editor prior to execution.
Useful for long commands and fixing typos.

## Usage
`f` can be called without arguments are can be passed arbitrary input for seeding the editor (e.g. `f !!` to fix a previously failed command).

```
Usage: f  
   or: f [options] [text] Open editor with text provided e.g. f !! to open with last command

Options:
  -l, --label   Label command and save to log.
  --history     Print command history (only named commands are saved).
  --dry-run     Print command without execution.
```

## Contributing

Only introduce added complexity when absolutely necessary.

PRs welcome.