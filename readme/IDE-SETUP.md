# IntelliJ IDEA/GoLand Configuration

<!-- START doctoc -->
<!-- END doctoc -->

## Plugins

**File | Settings | Plugins**: install **Go**(Intelli IDEA) and **File Watchers** plugins:

![IDE Setup Plugins](/readme/ide-setup-plugins.png)

## Line Length (Hard Wrap)

**File | Settings | Editor | Code Style | General**: set the value of `Hard wrap at` to `100`:

![IDE Setup Line Length](/readme/ide-setup-hard-wrap.png)

## Imports

**File | Settings | Editor | Code Style | Go | Imports**: set the value of `Sorting type` to `goimports`:

![IDE Setup Imports](/readme/ide-setup-goimports.png)

## File Watchers

**File | Settings | Tools | File Watchers**: add watchers for `golines` and `gofumpt`; see this 
[guide](https://github.com/mvdan/gofumpt?tab=readme-ov-file#goland) for more details.