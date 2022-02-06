Terraform `Incapsula` Provider
=========================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the team at [Imperva](https://www.imperva.com/).

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.14.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-incapsula`

```sh
$ git clone git@github.com:imperva/terraform-provider-incapsula $GOPATH/src/github.com/terraform-providers/terraform-provider-incapsula
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/imperva/terraform-provider-incapsula
$ make build
```

Using the provider
----------------------
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it. Documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/incapsula/index.html).

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-incapsula
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

An automation script is provided for Mac darwin 64amd based developers that 
encapsulates initial setups along make described commands. 
Please note that OS_ARCH=darwin_amd64 is uncommented in GNUmakefile for default Mac users, if needed for Linux users comment back and uncomment OS_ARCH=linux_amd64

Brew is a pre-requisite for this script, as the main package manager to install 
the dependent libraries such as Golang, Terraform and Git.
More details about this script is provided as inner code comments and description.

Script location **/scripts/tf-provider-incap-orch.sh**.

Script installation command will clone this repository to /workspace folder
as a first step and pull from git in subsequent runs.

It's recommended to download the script to some directory in local machine and start
with installation command execution

```sh
./tf-provider-incap-orch.sh -i "youApiID" "youApiKey"
```

