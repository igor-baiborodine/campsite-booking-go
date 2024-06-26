## IntelliJ IDEA/GoLand Configuration

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Plugins](#plugins)
- [Line Length (Hard Wrap)](#line-length-hard-wrap)
- [Imports](#imports)
- [File Watchers](#file-watchers)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

### Plugins

Go to **File | Settings | Plugins** and install **Go**(Intelli IDEA) and **File Watchers** plugins:

![IDE Setup Plugins](/readme/ide-setup-plugins.png)

### Line Length (Hard Wrap)

Go to **File | Settings | Editor | Code Style | General** and set the value of `Hard wrap at`
to `100`:

![IDE Setup Line Length](/readme/ide-setup-hard-wrap.png)

### Imports

Go to **File | Settings | Editor | Code Style | Go | Imports** and set the value of `Sorting type`
to `goimports`:

![IDE Setup Imports](/readme/ide-setup-goimports.png)

### File Watchers

Go to **File | Settings | Tools | File Watchers** and add watchers for `golines` and `gofumpt`; see
this [guide](https://github.com/mvdan/gofumpt?tab=readme-ov-file#goland) for more details.