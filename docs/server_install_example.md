# Install Online Server

Here a full example to setup a Online server with your ssh keys.
You can put a wrong OS information to get an error message
with the available OS for your server. Then you can find or
create the `partitioning_template_ref` field from the
[template page](https://console.online.net/en/template).

```HCL
provider "online" {
}

variable "server_id" {
  type = string
}

data "online_operating_system" "os" {
  server_id = var.server_id
  name = "centos"
  version = "CentOS 7.6"
}

data "online_ssh_keys" "keys" {
}

resource "online_server" "nicofonk" {
  server_id = var.server_id
  hostname = "online-server"
  os_id = "${data.online_operating_system.os.os_id}"
  user_login = "user1"
  user_password = "userpassword"
  root_password = "rootpassword"
  partitioning_template_ref = "13e84239-f208-4374-9ac2-8737a46211c6"
  ssh_keys = "${data.online_ssh_keys.keys.ssh_keys}"
}
```
