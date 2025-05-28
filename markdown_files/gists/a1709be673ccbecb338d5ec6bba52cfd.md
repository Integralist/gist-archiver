# Terraform: create backend definitions from a list of hostnames 

[View original Gist on GitHub](https://gist.github.com/Integralist/a1709be673ccbecb338d5ec6bba52cfd)

**Tags:** #terraform #fastly #work

## backend_template.vcl

```vcl
backend F_${backendname} {
    .between_bytes_timeout = 10s;
    .connect_timeout = 1s;
    .dynamic = true;
    .first_byte_timeout = 15s;
    .host = "${hostname}";
    .max_connections = 200;
    .port = "443";
    .share_key = "...";
    .ssl = true;
    .ssl_cert_hostname = "${hostname}";
    .ssl_check_cert = always;
    .ssl_sni_hostname = "${hostname}";
    .bypass_local_route_table = true;
    .probe = {
      .dummy = true;
    }
}
```

## main.tf

```hcl
terraform {
  required_providers {
    fastly = {
      source  = "fastly/fastly"
      version = ">= 1.0.0"
    }
  }
}

locals {
  backends = ["b1.example.com", "b2.example.com"]
}

resource "fastly_service_vcl" "myservice" {
  name = "myservice"

  domain {
    name = "test.hkakehas.tokyo"
  }

  snippet {
    name = "backends"
    type = "init"
    content = join("\n", [for backend in local.backends : templatefile("${path.module}/backend_template.vcl",
    { backendname : replace(backend, ".", "_"), hostname : backend, })])
  }

  force_destroy = true
}
```

## result.txt

```text
terraform apply
fastly_service_vcl.myservice: Refreshing state... [id=...]

Terraform used the selected providers to generate the following execution plan. Resource actions are
indicated with the following symbols:
  ~ update in-place
  ( snip )
      + snippet {
          + content  = <<-EOT
                backend F_b1_example_com {
                    .between_bytes_timeout = 10s;
                    .connect_timeout = 1s;
                    .dynamic = true;
                    .first_byte_timeout = 15s;
                    .host = "b1.example.com";
                    .max_connections = 200;
                    .port = "443";
                    .share_key = "...";
                    .ssl = true;
                    .ssl_cert_hostname = "b1.example.com";
                    .ssl_check_cert = always;
                    .ssl_sni_hostname = "b1.example.com";
                    .bypass_local_route_table = true;
                    .probe = {
                      .dummy = true;
                    }
                }
                backend F_b2_example_com {
                    .between_bytes_timeout = 10s;
                    .connect_timeout = 1s;
                    .dynamic = true;
                    .first_byte_timeout = 15s;
                    .host = "b2.example.com";
                    .max_connections = 200;
                    .port = "443";
                    .share_key = "...";
                    .ssl = true;
                    .ssl_cert_hostname = "b2.example.com";
                    .ssl_check_cert = always;
                    .ssl_sni_hostname = "b2.example.com";
                    .bypass_local_route_table = true;
                    .probe = {
                      .dummy = true;
                    }
                }
            EOT
          + name     = "backends"
          + priority = 100
          + type     = "init"
        }
        # (1 unchanged block hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

fastly_service_vcl.myservice: Modifying... [id=...]
fastly_service_vcl.myservice: Still modifying... [id=..., 10s elapsed]
fastly_service_vcl.myservice: Still modifying... [id=..., 20s elapsed]
fastly_service_vcl.myservice: Still modifying... [id=..., 30s elapsed]
fastly_service_vcl.myservice: Modifications complete after 31s [id=...]

```

