## terraform provider for network information

Right now it just provides your IP address by shelling out to `dig` and using OpenDNS's resolvers to echo back your IP by requesting the A record for the special name `myip.opendns.com`.

This is exposed as the value `wan_ip_address` on a (data) resource, `network_info`.


## Example

    data "network_info" "main" { }

    resource "aws_security_group" "cluster_sg" {
      ...

      ingress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["${data.network_info.main.wan_ip_address}/32"]
      }

      description = "Group to allow all traffic from current location"
   }


## References

References for writing custom providers

- <https://www.hashicorp.com/blog/terraform-custom-providers.html>
  * From the official Terraform blog, but it's out of date (written on September 26, 2014) and a lot of it doesn't work with current API
  * The horrendous code block formatting and coloring suggests no one at Terraform has read or tended to that page in quite some time
- <http://container-solutions.com/write-terraform-provider-part-1/>
  * A more modern (December 1, 2015) full example; not as much explanation or depth as the previous, but it works much better
  * Full code at <https://github.com/ContainerSolutions/terraform-provider-template>
- <https://godoc.org/github.com/hashicorp/terraform/helper/schema>
  * godoc for the main provider-writing helper module


## License

Copyright © 2016 Christopher Brown. [MIT Licensed](https://chbrown.github.io/licenses/MIT/#2016).
