Terraform Provider for Online.net
=================================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/src-d/terraform-provider-online`.

```sh
$ mkdir -p $GOPATH/src/github.com/src-d; cd $GOPATH/src/github.com/src-d
$ git clone git@github.com:src-d/terraform-provider-online
```

Enter the provider directory and build the provider:

```sh
$ cd $GOPATH/src/github.com/src-d/terraform-provider-online
$ make build
```

To install it in your home directory to test the provider:

```sh
$ cd $GOPATH/src/github.com/src-d/terraform-provider-online
$ make local-install
```

Installing the provider
-----------------------

### Installation from binaries (recommended)

The recommended way to install *terraform-provider-online* is use the binary
distributions from the [Releases](https://github.com/src-d/terraform-provider-online/releases) page. The packages are available for Linux and macOS.

Download and uncompress the latest release for your OS. This example uses the linux binary.

```sh
$ wget https://github.com/src-d/terraform-provider-online/releases/download/v0.1.0/terraform-provider-online_v0.1.0_linux_amd64.tar.gz
$ tar -xvf terraform-provider-online*.tar.gz
```

Now copy the binary to the Terraform's plugins folder (if this is your first plugin maybe it isn't present):

```sh
$ mkdir -p ~/.terraform.d/plugins/
$ mv terraform-provider-online*/terraform-provider-online ~/.terraform.d/plugins/
```

Using the provider
------------------

*To do: add documentation*

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-online
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

Acceptance tests need an authentication token and an Online.net server which will be used to run the tests against.
The server will be modified, using a production system is not adviced. All used servers must support RPNv2.
It will look for:

| Environment Variable | Description                                                                  | Example     | 
|----------------------|------------------------------------------------------------------------------|-------------|
| `ONLINE_SERVER_ID`   | ID of a dedicated server (only the numeric part)                             | `46952`     |
| `ONLINE_SERVER_ID_2` | ID of a 2nd dedicated server (only the numeric part)                         | `46953`     |
| `ONLINE_FAILVOVER_IP` | An available failover IP                                                     | `81.23.14.1`|
| `ONLINE_TOKEN`       | Online.net auth token received from https://console.online.net/en/api/access |             |

```sh
$ make testacc
```

You can also run a dedicated acceptance test, here an example:

```sh
$ make test-compile TEST=./provider
$ ONLINE_TOKEN=<token> ONLINE_SERVER_ID=<id> TF_ACC=1 ./provider.test -test.run TestDataOperatingSystemAcceptance -test.v
```

## License

Mozilla Public License Version 2.0, see [LICENSE](/LICENSE)

