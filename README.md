## Terraform Provider Plugin for network information

This [Terraform](https://www.terraform.io/) [Provider Plugin](https://www.terraform.io/docs/extend/plugin-types.html#providers) resolves your public IP address via DNS or HTTP,
and makes it available for interpolation elsewhere in your Terraform configuration as a [Data Source](https://www.terraform.io/docs/configuration/data-sources.html) named `network_info`,
which exposes a single attribute, `wan_ip_address` (a string).


### Example

```HCL
# "local" can be anything you want
data "network_info" "local" { }

resource "aws_security_group" "cluster_sg" {
  // ...

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["${data.network_info.local.wan_ip_address}/32"]
  }

  description = "Group to allow all traffic from current location"
}
```

When unconfigured, as above, `network_info` will call the DNS resolver.
This uses [OpenDNS](https://www.opendns.com/)'s resolvers to echo back your IP by requesting the `A` record for the special name `myip.opendns.com`, exactly like:

    dig +short myip.opendns.com @resolver1.opendns.com

Alternatively, you can select the HTTP resolver, which sends a plain `GET` request via HTTP to [`https://checkip.amazonaws.com`](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/authorizing-access-to-an-instance.html#authorizing-access-prereqs):

```HCL
data "network_info" "local" {
  dns = false
}
```

Finally, you can explicitly specify a URL (which should return a quad-dotted IPv4 address somewhere in the response when sent a `GET` request):

```HCL
data "network_info" "local" {
  http = "http://whatismyip.akamai.com"
}
```

The DNS resolver is not user-configurable.


### Install

First, use `go` to fetch and build:

    go get github.com/chbrown/terraform-provider-network

Then edit your `~/.terraformrc` to contain:

    providers {
      network = "/Your/absolute/expanded/GOPATH/bin/terraform-provider-network"
    }

Where the value of `network` is an absolute path to the `terraform-provider-network` binary that `go get` just built.
Assuming you have the `GOPATH` environment variable set, this will be the value of `$GOPATH/bin/terraform-provider-network`.


### Debugging

If you get the cryptic error message:

    Error configuring: 1 error(s) occurred:

    * Incompatible API version with plugin. Plugin version: 1, Ours: 2

This means `terraform` has been updated to a version different from the libraries the plugin was built with.
The solution is to reinstall, by using `go get` with the `-u` flag:

    go get -u github.com/chbrown/terraform-provider-network


### References

- [Writing Custom Providers](https://www.terraform.io/docs/extend/writing-custom-providers.html) (2018-04-12)
  * Official documentation.
  * Features comprehensive example describing best practices and suggested code layout.
- [Writing Custom Terraform Providers](https://www.hashicorp.com/blog/writing-custom-terraform-providers) (2014-09-26)
  * Out of date (a lot of it doesn't work with the current API) and stylistically neglected.
  * > This guide exists for historical purposes ...
- [Write your own Terraform provider: Part 1](http://container-solutions.com/write-terraform-provider-part-1/) (2015-12-01)
  * Third-party perspective; not as much explanation or depth as the Hashicorp documentation.
  * As far as I can tell, there is no Part 2.
  * [Full source code](https://github.com/ContainerSolutions/terraform-provider-template)
- [GoDoc for `helper/schema` package](https://godoc.org/github.com/hashicorp/terraform/helper/schema)


### License

Copyright © 2016 Christopher Brown. [MIT Licensed](https://chbrown.github.io/licenses/MIT/#2016).
