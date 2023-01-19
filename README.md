# Aurae AE

[![Test and Format](https://github.com/aurae-runtime/ae/actions/workflows/001-go-ubuntu-latest-make-test-format-lint.yaml/badge.svg?branch=main)](https://github.com/aurae-runtime/ae/actions/workflows/001-go-ubuntu-latest-make-test-format-lint.yaml)

Unix inspired command line client for Aurae written in Go!

Contributions and newcomers to the project are welcome. Please see [Getting Involved](https://github.com/aurae-runtime/community#getting-involved) to join the Discord and more.

[Signing the CLA](https://cla.aurae.io/) is required for contributions.

## What is ae? Why does it exist?

In order to understand `ae` we should first understand `aer`.

### What is aer?
We intend to auto generate a command line tool named `aer` built on the Rust client that has an identical scope of a single auraed node.

This tool is for "power-users" exists as a way of quickly developing and debugging the APIs as we change them. For example an auraed developer might make a change to an API and need a quick way to test the API locally against a single daemon.

### What is ae?

We intend to maintain a command line tool named `ae` build on the Go client that has a broader scope than `aer`.

The scope for `ae` will be for a cluster of Aurae nodes and will likely have identical functionality to `aer` in many places - but not all.

The `ae` command line tool will be for working with a group of Aurae nodes, and will likely include more pragmatic functionality that is meaningful to an enterprise operator.

### Who is ae for?

The `ae` tool should feel familiar to any cloud operator.

-   Tasks such as "rolling out nginx to prod" should be possible with `ae`.
-   Tasks such as "counting the number of threads of a process" should be possible with `ae`.
-   Tasks such as "changing the current database connection count for a service" should be possible with `ae`.
-   Tasks such as "swinging traffic from one service to another" should be possible with `ae`.

These are more "practical" or "outcome driven" tasks, and will likely require additional functionality to the lightweight `aer` tool.

## Using ae
### Commands

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
ae stop container # Note this has an alias: 'ae oci kill'
```

Stop a resource immediately.

#### Logs

```bash
ae logs <options>
```

#### Discovery

```bash
ae discover <cidr <cidr> | ip <ip>>
```

Scans a cluster of Aurae nodes and returns information about them, including the version running.

#### Health

```bash
ae check <cidr <cidr> | ip <ip>> <service,...>
```

Scans a cluster of Aurae nodes and returns the serving status of the given list of services.


## Guiding Principles

### Less is more.

We do not want `ae` to become a junk drawer. In situations where we are considering bringing in new functionality to the project, we prefer to keep it out.

For example imagine we were considering bringing in a `--filter`flag or a `--query` flag which would allow for us to filter the output returned by `ae`. In this situation we could very well identify advanced query patterns and libraries that would make it possible to adopt this feature.

However the feature comes at a maintenance cost, and frankly already exists in other tools. Instead of adopting a filter syntax and the associated burden that comes along with it we instead focus on outputting to formats (such as JSON) which can easily be queried by existing tools (such as [jq](https://stedolan.github.io/jq/)).

### Just because you can, doesn't mean you should.

Also known as **"The Jurassic Park"** rule.

In situations where functionality is made available or "is possible" because of current technical advancements or capabilities, we prefer to focus on the goals of ae instead of supporting optional features "just because we can".

For example imagine discovering a well written, and well regarded library that takes a Go struct and turns it into many formats which can be printed to the screen such as XML, JSON, YAML, TOML, BSON, and so on....

The project currently has established a need for JSON only, and plans to use jq for filtering and querying the data. In this case bringing the additional output types to the project "just because we can" would be a violation of this rule.

### No assumptions

Also known as "no conveniences" or "no magic" policies.

Assumptions and conveniences can delight a user, until the assumption is wrong in which case they can be catastrophically frustrating. We prefer to stay away from assumptions whenever possible. We prefer to stay away from "magic" and "convenience" style features.

A famous example of this is bringing a logger to a project which assumes that the logger will be able to take ownership of the -v flag for verbose. There are many situations (docker -v for volume, or grep -v for reverse) where this assumption is flawed and will conflict with the goals of the project and panic at runtime. We do not want ae to turn into the logger where we assume something incorrectly that ends up causing problems for others.

This will be a delicate principle to consider, as we also will need to hold firm opinions on things. A balance between a strong opinion and false assumptions is often times hard to get right.

For example the ae project will need to read from a configuration file in order to communicate with a remote server. In one case we could assume that the configuration file is always in a familiar location for convenience. In the other case we do not want to force a user to pass a very long flag every time they run the tool. What do we do?

We hold the opinion that the the file should be in a few well documented locations, and create logic to try the locations if no flag is passed. If a user passes a flag we should remove all assumptions from the program and only consider the input. We prefer clarity over magic.

### Remote Servers instead of Local Systems

This is an interesting principle to consider. As the functionality of the ae tool grows we will inevitably need to execute logic somewhere.

We prefer to keep the ae tool as "boring" and "unintelligent" as possible. We do not want the command line tool to do a lot of "work" or "processing" locally. In situations where we need to perform logic against a set of Aurae nodes, we need to find a way to take the logic out of ae for compatibility reasons.

We do not want to get into a situation where we have multiple clients attempting to perform slightly different tasks locally against a server. Especially when these clients are subject to take long periods of time for their work to complete.

In the event these types of situations arise, it is a sign that we likely need to deploy a scheduler or server mechanisms somewhere to manage the work on our behalf.



## Code Quality

### Linting

We are using the (golangci-lint)[https://golangci-lint.run/] tool to lint the code. You will need to install it for your system, and can find instructions at (this link)[https://golangci-lint.run/usage/install/]. This tool runs on every pull request, and it must pass before merging is allowed. You can run it locally with `make lint`

### Formatting

We are using the (gofmt)[https://pkg.go.dev/cmd/gofmt] tool to lint the code. This tool runs on every pull request, and it must pass before merging is allowed. You can run it locally with `make format`
