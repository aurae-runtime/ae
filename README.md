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

**What is `aer`?**

The intention is to build an automatically generated CLI tool called `aer`, which is based on the Rust client that has the identical scope of a single _[auraed](https://github.com/aurae-runtime/aurae/tree/main/auraed)_ node.

The tool will be aimed at "POWER-USERS" and exists as a rapid way to develop and debug against APIs that change frequently. For example, an [auraed](https://github.com/aurae-runtime/aurae/tree/main/auraed) developer can make a change to an existing API and test it locally against a single daemon.


**What is `ae`?**

We want to create and maintain a CLI tool, `ae` that is based on the Go client and has a broader view than `aer`.

The purpose will be to use `ae` for clusters of [Aurae](https://github.com/aurae-runtime/aurae) nodes and will probably work similar to `aer`, but not completely!

The `ae` CLI tool will work for a group of nodes and will probably contain more pragmatic functions that are more important for enterprise operators.

**Who is `ae` for?**

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

> This section is reserved for future installation as well as an example of integration.

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
  ae check <cidr <cidr> | ip <ip>> <service,...>
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
  
  Here the OCI CLI interface is implemented with the respective subcommands.
    
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

<!-- CONTRIBUTE -->
<h2 id="about-the-project">Contribute</h2>
