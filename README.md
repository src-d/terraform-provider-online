Terraform Provider for Online.net
=================================

- Website: https://www.terraform.io
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

  [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/src-d/terraform-provider-online`

```sh
$ mkdir -p $GOPATH/src/github.com/src-d; cd $GOPATH/src/github.com/src-d
$ git clone git@github.com:src-d/terraform-provider-online
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/src-d/terraform-provider-online
$ make build
```

To install it in your home directory to test the provifer

```sh
$ cd $GOPATH/src/github.com/src-d/terraform-provider-online
$ make local-install
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

```sh
$ make testacc
```

## License

Mozilla Public License Version 2.0, see [LICENSE](/LICENSE)

