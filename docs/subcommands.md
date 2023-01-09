# Subcommands

The purpose of this document is to describe the handling of subcommands within `ae`.

## Create a new subcommand

All subcommands reside in a dedicated package within `./cmd/`. Subcommands of subcommands are stored within their parents package.

Example: `ae oci` is a 1st level subcommand and therefore stored in `./cmd/oci/`. `ae oci create` is a 2nd level subcommand of the `oci` subcommand and therefore stored in `./cmd/oci/create`. This structure is also reflected in its package name `òci_create`.

```
./cmd
├── cmd.go
├── oci
│   ├── create
│   │   └── oci_create.go
│   ├── delete
│   │   └── oci_delete.go
│   ├── kill
│   │   └── oci_kill.go
│   ├── oci.go
│   ├── start
│   │   └── oci_start.go
│   └── state
│       └── oci_state.go
├── root
│   └── root.go
├── runtime
│   ├── allocate
│   │   └── allocate.go
│   ├── free
│   │   └── free.go
│   ├── start
│   │   └── start.go
│   └── stop
│       └── stop.go
└── version
    └── version.go
```

To create a new command, proceed with the following steps:

Create a new directory for the new command:

```
$ mkdir ./cmd/new_command
```

Create a new go package within this directory:

```
$ echo "package new_command" > ./cmd/new_command/new_command.go
```

```
./cmd
├── cmd.go
├── new_command
│   └── new_command.go
├── oci
├── root
└── version
```

## Structure of a subcommand

Every command follows the interface defined in `./cmd/cmd.go`.

**Option struct**

The option struct is supposed to contain the command's input and holds any relevant information necessary to process the command.

```go
type option struct {}
```

**Complete**

`Complete()`  is the method where the option extracts the data from the arguments and sets its different attributes. For example, assuming your subcommand requires a TCP port as a parameter, you should check for its availability within `Complete()` and make it available within `option`.

`Complete()` is triggered before `Validate()` and `Execute()`.

```go
func (o *option) Complete(_ []string) error {}
```

**Validate**

`Validate()` ensures that the commands attribute values are coherent (when it makes sense). For example, assuming your subcommand requires a TCP port as a parameter, you should check for the validity of the parameter within `Validate()`.

Validate is triggered after `Complete()` and before `Execute()`

```go
func (o *option) Validate() error {}
```

**Execute**

Entry point to the subcommand action.

```go
func (o *option) Execute() error {}
```

**SetWriter**

Defines an `io.Writer` to output text. `io.Writer` is stored in `option` struct.

```go
func (o *option) SetWriter(writer io.Writer) {}
```

**NewCMD**

`NewCMD()` constructs a new `cobra.Command`. Furthermore, this is where command-specific flags should be defined, or further subcommands can be registered to the command.

```go
func NewCMD() *cobra.Command {}
```

## Registering the new command to the root command

A new subcommand can be added to the root command by adding it to the `init()` function of root command in `./cmd/root/root.go` via `rootCmd.AddCommand(new_command.NewCMD())`.
