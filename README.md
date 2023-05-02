# Process Settings

This is the golang implementation of the `process_settings` library originally written in [Ruby](https://github.com/Invoca/process_settings).
The overall functionality of the library is the same, but due to the nature of Golang, the implementation is slightly different.
Please read through the documentation below to understand how to use this library in your project.

## Installation

To install this library, run the following command:

```bash
go get github.com/Invoca/process_settings.go
```

## Usage

The `process_settings.ProcessSettings` object can be freely created and used at any time.
But typical usage is through the global instance that should be configured at the start of your process.

### Basic Configuration

Before using the global instance, you must first set it by assigning a `process_settings.ProcessSettings` object to it.

```go
import (
    "github.com/Invoca/process_settings.go"
)

func main() {
    process_settings.SetGlobalProcessSettings(
      process_settings.NewProcessSettingsFromFile("/etc/process_settings/combined_process_settings.yml")
    )
}
```

### Configuration with Static Context

When initializing a new `process_settings.ProcessSettings` object, you can provide a `map[string]instance{}` of static context.
This context will be used to select the settings that are specifically targeted for your process.
For example, if you're running multiple services, you might want to target settings to a specific service.

```go
import (
    "github.com/Invoca/process_settings.go"
)

func main() {
    ps, err := process_settings.NewProcessSettingsFromFile(
        "/etc/process_settings/combined_process_settings.yml",
        map[string]instance{}{
            "service_name": "frontend",
            "datacenter": "AWS-US-EAST-1",
        }
    )

    if err != nil {
        panic(err)
    }

    process_settings.SetGlobalProcessSettings(ps)
}
```

### Reading Settings

For the following section, consider the `combined_process_settings.yml` file:
```yaml
---
- filename: frontend.yml
  settings:
    frontend:
      log_level: info
- filename: frontend-microsite.yml
  target:
    domain: microsite.example.com
  settings:
    frontend:
      log_level: debug
- meta:
    version: 27
    END: true
```

To read a setting, application code should call the `process_settings.Get()` method.
For example:

```go
log_level := process_settings.Get("frontend", "log_level")
```

### Dynamic Settings

The `process_settings.ProcessSettings` object has a `Monitor` built in that loads settings changes dynamically whenever the file changes,
by using the [fsnotify](https://github.com./fsnotify/fsnotify) library which in turn uses the `INotify` module of the Linux kernel, or `FSEvents` on MacOS. There is no need to restart the process or send it a signal to tell it to reload changes.

To start the monitor goroutine, call the `StartMonitor()` method on the `process_settings.ProcessSettings` object.

```go
func main() {
    ps, err := process_settings.NewProcessSettingsFromFile(
        "/etc/process_settings/combined_process_settings.yml",
        map[string]instance{}{
            "service_name": "frontend",
            "datacenter": "AWS-US-EAST-1",
        }
    )

    if err != nil {
        panic(err)
    }

    process_settings.SetGlobalProcessSettings(ps)
    ps.StartMonitor()
}
```

#### Read Latest Settings Through `process_settings.Get()` and `process_settings.SafeGet()`

The simplest approach--as shown above--is to read the latest settings at any time through `process_settings.Get()`
and `process_settings.SafeGet()` (which delegates to `process_settings.instance`):

```
http_version := process_settings.Get('frontend', 'http_version')
```

#### Register a `WhenUpdated` Callback
Alternatively, if you need to execute initially and whenever the value is updated, register a callback with `process_settings.WhenUpdated()`:

```go
process_settings.WhenUpdated(func() {
    logger.level = process_settings.Get("frontend", "log_level")
})
```

By default, the `WhenUpdated` func is called initially when registered. We've found this to be convenient in most cases; it can be disabled by passing an optional second
argument `false`, in which case the block will be called 0 or more times in the future,
when any of the process settings for this process change.

`WhenUpdated` **_is not_** idempotent, so adding the same func multiple times will result
in multiple registered callbacks doing the same operation.

In case you need to cancel the callback later, `WhenUpdated` returns a handle (the index of the registered func) which can later be passed into `CancelWhenUpdated`.

Note that all callbacks run sequentially on the shared change monitoring thread, so please be considerate!

## Targeting
Each settings YAML file has an optional `target` key at the top level, next to `settings`.

If there is no `target` key, the target defaults to `true`, meaning all processes are targeted for these settings. (However, the settings may be overridden by other YAML files. See "Precedence" below.)

### Hash Key-Values Are AND'd
To `target` on context values, provide a hash of key-value pairs. All keys must match for the target to be met. For example, consider this target hash:
```
target:
  service_name: frontend
  datacenter: AWS-US-EAST-1
```
This will be applied in any process that has `service_name == "frontend"` AND is running in `datacenter == "AWS-US-EAST-1"`.

### Multiple Values Are OR'd
Values may be set to an array, in which case the key matches if _any_ of the values matches. For example, consider this target hash:
```
target:
  service_name: [frontend, auth]
  datacenter: AWS-US-EAST-1
```
This will be applied in any process that has (`service_name == "frontend"` OR `service_name == "auth"`) AND `datacenter == "AWS-US-EAST-1"`.

### Precedence
The settings YAML files are always combined in alphabetical order by file path. Later settings take precedence over the earlier ones.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.
