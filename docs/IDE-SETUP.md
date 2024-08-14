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

![IDE Setup Plugins](/docs/ide-setup-plugins.png)

### Line Length (Hard Wrap)

Go to **File | Settings | Editor | Code Style | General** and set the value of `Hard wrap at`
to `100`:

![IDE Setup Line Length](/docs/ide-setup-hard-wrap.png)

### Imports

Go to **File | Settings | Editor | Code Style | Go | Imports** and set the value of `Sorting type`
to `goimports`:

![IDE Setup Imports](/docs/ide-setup-goimports.png)

### File Watchers

Go to **File | Settings | Tools | File Watchers** and add watchers for `golines` and `gofumpt`; see
this [guide](https://github.com/mvdan/gofumpt?tab=readme-ov-file#goland) for more details. Add a new scope to exclude proto and generated files:

```text
file[campsite-booking-go]:*/&&!file[campsite-booking-go]:campgroundspb/v1//*&&!file[campsite-booking-go]:*/mock_*.go
```

![IDE Setup File Watchers Scope](/docs/ide-setup-file-watchers-scope.png)

![IDE Setup File Watchers](/docs/ide-setup-file-watchers.png)
