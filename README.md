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

References for writing custom providers

- <https://www.hashicorp.com/blog/terraform-custom-providers.html>
  * From the official Terraform blog, but it's out of date (written on September 26, 2014) and a lot of it doesn't work with current API
  * The horrendous code block formatting and coloring suggests no one at Terraform has read or tended to that page in quite some time
- <http://container-solutions.com/write-terraform-provider-part-1/>
  * A more modern (December 1, 2015) full example; not as much explanation or depth as the previous, but it works better
  * Full code at <https://github.com/ContainerSolutions/terraform-provider-template>
- <https://godoc.org/github.com/hashicorp/terraform/helper/schema>
  * godoc for the main provider-writing helper module


### License

Copyright © 2016 Christopher Brown. [MIT Licensed](https://chbrown.github.io/licenses/MIT/#2016).
