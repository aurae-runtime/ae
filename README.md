<h1 align="center">
  <code>ae</code>
  <h3 align="center">UNIX inspired CLI Client for Aurae</h3>
</h1>

<div align='center'>

<a href='https://github.com/aurae-runtime/ae/blob/main/go.mod'>
<img alt="go version" src="https://img.shields.io/github/go-mod/go-version/aurae-runtime/ae?color=grey&logo=go&style=for-the-badge">
  
</a>
  
<a href="https://github.com/aurae-runtime/ae/blob/main/LICENSE">
<img alt="license" src="https://img.shields.io/github/license/aurae-runtime/ae?color=grey&logo=apache&style=for-the-badge"/>
  
</a>

</div>

---

> The project welcomes new contributors and developers. Check out the **[Getting Involved](https://github.com/aurae-runtime/community#getting-involved)** section and join the Discord. It is mandatory to sign the **[CLA](https://cla.aurae.io/)** to contribute.


<!-- TABLE OF CONTENTS -->
<h2 id="table-of-contents"> 📑 Table of Contents</h2>

<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project"> ► About The Project</a></li>
    <li><a href="#quickstart"> ► Quickstart</a></li>
    <li><a href="#usage"> ► Usage</a></li>
    <li><a href="#philosophy"> ► Philosophy</a></li>
    <li><a href="#contribute"> ► Contribute</a></li>
  </ol>
</details>

&nbsp;

<!-- ABOUT THE PROJECT -->
<h2 id="about-the-project">About The Project</h2>

`ae` is a UNIX inspired CLI client for **[Aurae](https://github.com/aurae-runtime/aurae)**, written in Go. However, in order to understand what `ae` should and can do, we must first understand `aer`.

### **What is `aer`?**

The intention is to build an automatically generated CLI tool called `aer`, which is based on the Rust client that has the identical scope of a single _[auraed](https://github.com/aurae-runtime/aurae/tree/main/auraed)_ node.

The tool will be aimed at "POWER-USERS" and exists as a rapid way to develop and debug against APIs that change frequently. For example, an [auraed](https://github.com/aurae-runtime/aurae/tree/main/auraed) developer can make a change to an existing API and test it locally against a single daemon.


### **What is `ae`?**

We want to create and maintain a CLI tool, `ae` that is based on the Go client and has a broader view than `aer`.

The purpose will be to use `ae` for clusters of [Aurae](https://github.com/aurae-runtime/aurae) nodes and will probably work similar to `aer`, but not completely!

The `ae` CLI tool will work for a group of nodes and will probably contain more pragmatic functions that are more important for enterprise operators.

### **Who is `ae` for?**

The `ae` utility should be familiar to every cloud operator.

_Typical tasks_ such as:

* "Rolling out NGINX to production"
* "Counting the numbers of threads of a process"
* "Changing the current databse connection count for a service"
* Swinging traffic from one to another"

should be possible with `ae`.

There are more "practical" and "impact-oriented" tasks, and these probably need extra functionality, which they add to the lightweight `aer` tool.

<!-- QUICKSTART -->
<h2 id="about-the-project">Quickstart</h2>

> This section is reserved for future installation instructions as well as an example of integration.

<!-- USAGE -->
<h2 id="about-the-project">Usage</h2>

There are a number of commands for `ae`.

These are shown here in _alphabetical_ order.

<details>
  <summary><code>allocate</code></summary>
   
  &nbsp;
  
  Resources are reserved and prerequisites can be managed, but it **does not** start. It will not work if the resources are not available.
    
  ```
  ae allocate
  ae allocate cell
  ae allocate pod
  ```
  
</details>

<details>
  <summary><code>check</code></summary>
  
  &nbsp;
  
  Checks the nodes of the cluster and returns the current serving status with the given list of services.
    
  ```
  ae check <cidr <cidr> | ip <ip>> <service, ...>
  ```
  
</details>

<details>
  <summary><code>discover</code></summary>
  
  &nbsp;
  
  Scans the complete network or cluster of nodes and returns information about it, including the version.
    
  ```
  ae discover <cidr <cidr> | ip <ip>>
  ```
    
</details>

<details>
  <summary><code>free</code></summary>
  
  &nbsp;
  
  It frees the resources and destroys the prerequisites that were started. It will fail if the resources cannot be freed or do not exist.
    
  ```
  ae free
  ae free cell
  ae free pod
  ```
    
</details>

<details>
  <summary><code>logs</code></summary>
  
  &nbsp;
  
  This option will accept aruguments and return or save some kind of logs.
    
  ```
  ae logs <options>
  ```
    
</details>

<details>
  <summary><code>oci</code></summary>
  
  &nbsp;
  
  Here the [OCI CLI interface](https://github.com/opencontainers/runtime-tools/blob/master/docs/command-line-interface.md) is implemented with the respective subcommands.
    
  ```
  ae oci
  ae oci create
  ae oci delete
  ae oci kill
  ae oci start
  ae oci status
  ```
    
</details>

<details>
  <summary><code>start</code></summary>
  
  &nbsp;
  
  It will run the rescource directly.
    
  ```
  ae start
  ae start executable, exe
  ae start container # Note this has an alias: 'ae oci start'
  ```
    
</details>

<details>
  <summary><code>stop</code></summary>
  
  &nbsp;
  
  It will stop the rescource directly.
    
  ```
  ae stop
  ae stop executable, exe
  ae stop container # Note this has an alias: 'ae oci kill'
  ```
    
</details>

<!-- PHILOSOPHY -->
<h2 id="about-the-project">Philosophy</h2>

**This project has a few philosophical principles.**
    
### **Less is more**
    
We do not want `ae` to become a junk drawer. In situations where we are considering bringing in new functionality to the project, we prefer to keep it out.

For example imagine we were considering bringing in a `--filter` flag or a `--query` flag which would allow for us to filter the output returned by `ae`. In this situation we could very well identify advanced query patterns and libraries that would make it possible to adopt this feature.

However the feature comes at a maintenance cost, and frankly already exists in other tools. Instead of adopting a filter syntax and the associated burden that comes along with it we instead focus on outputting to formats (such as JSON) which can easily be queried by existing tools (such as [jq](https://stedolan.github.io/jq/)).
    
### **Just because you can, doesn't mean you should.**
    
Also known as **"The Jurassic Park"** rule.

In situations where functionality is made available or "is possible" because of current technical advancements or capabilities, we prefer to focus on the goals of ae instead of supporting optional features "just because we can".

For example imagine discovering a well written, and well regarded library that takes a Go struct and turns it into many formats which can be printed to the screen such as XML, JSON, YAML, TOML, BSON, and so on....

The project currently has established a need for JSON only, and plans to use jq for filtering and querying the data. In this case bringing the additional output types to the project "just because we can" would be a violation of this rule.
    
### **No assumptions…**

Also known as "no conveniences" or "no magic" policies.

Assumptions and conveniences can delight a user, until the assumption is wrong in which case they can be catastrophically frustrating. We prefer to stay away from assumptions whenever possible. We prefer to stay away from "magic" and "convenience" style features.

A famous example of this is bringing a logger to a project which assumes that the logger will be able to take ownership of the -v flag for verbose. There are many situations (docker -v for volume, or grep -v for reverse) where this assumption is flawed and will conflict with the goals of the project and panic at runtime. We do not want ae to turn into the logger where we assume something incorrectly that ends up causing problems for others.

This will be a delicate principle to consider, as we also will need to hold firm opinions on things. A balance between a strong opinion and false assumptions is often times hard to get right.

For example the ae project will need to read from a configuration file in order to communicate with a remote server. In one case we could assume that the configuration file is always in a familiar location for convenience. In the other case we do not want to force a user to pass a very long flag every time they run the tool. What do we do?

We hold the opinion that the the file should be in a few well documented locations, and create logic to try the locations if no flag is passed. If a user passes a flag we should remove all assumptions from the program and only consider the input. We prefer clarity over magic.
    
### **Remote Servers instead of Local Systems!**
    
This is an interesting principle to consider. As the functionality of the ae tool grows we will inevitably need to execute logic somewhere.

We prefer to keep the ae tool as "boring" and "unintelligent" as possible. We do not want the command line tool to do a lot of "work" or "processing" locally. In situations where we need to perform logic against a set of Aurae nodes, we need to find a way to take the logic out of ae for compatibility reasons.

We do not want to get into a situation where we have multiple clients attempting to perform slightly different tasks locally against a server. Especially when these clients are subject to take long periods of time for their work to complete.

In the event these types of situations arise, it is a sign that we likely need to deploy a scheduler or server mechanisms somewhere to manage the work on our behalf.    
    
<!-- CONTRIBUTE -->
<h2 id="about-the-project">Contribute</h2>
    
The **[Aurae](https://github.com/aurae-runtime/aurae)** project is always looking for new members and developers. Here in this repository you can improve `ae`, but be sure to check out the [organisation](https://github.com/aurae-runtime) and the other internal projects. [This](https://github.com/aurae-runtime/community) is always a good starting point.
    
**Code Quality**
    
**_Linting_**

We are using the [golangci-lint](https://golangci-lint.run/) tool to lint the code. You will need to install it for your system, and can find instructions at [this link](https://golangci-lint.run/usage/install/). This tool runs on every pull request, and it must pass before merging is allowed. You can run it locally with `make lint`.
    
**_Formatting_**
    
We are using the [gofmt](https://pkg.go.dev/cmd/gofmt) tool to lint the code. This tool runs on every pull request, and it must pass before merging is allowed. You can run it locally with `make format`.

