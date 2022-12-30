# Aurae AE

Unix inspired command line client for Aurae written in Go! 

Contributions and newcomers to the project are welcome. Please see [Gettinv Involved](https://github.com/aurae-runtime/community#getting-involved) to join the Discord and more.

## Commands

```bash 
ae oci
ae oci create
ae oci delete
ae oci kill
ae oci start
ae oci status
```

Implements the [OCI Command Line](https://github.com/opencontainers/runtime-tools/blob/master/docs/command-line-interface.md) interface with associated subcommands.


```bash
ae allocate
ae allocate cell
ae allocate pod
```

Reserve resources, and manage any prerequisites but do not start. Will fail if resources are unavailable. 

```bash
ae free
ae free cell
ae free pod
```

Free resources, and destroy any prerequisites that have been started. Will fail if resources cannot be freed or do not exist.

```bash
ae start
ae start executable, exe
ae start container # Note this has an alias: 'ae oci start'
```

Run a resource immediately.

```bash
ae stop
ae stop executable, exe
ae stop container # Note this has an alias: 'ae oci stop'
```

Stop a resource immediately. 
