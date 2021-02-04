# terraform-provider-custom
Experimentation with custom Terraform Provider

Project structure, scripts, and GNUmakefile all cloned from: https://github.com/hashicorp/terraform-provider-random

## Building the provider (Linux)

Clone repository to: $GOPATH/src/github.com/logicalChimp/terraform-provider-custom

```bash
$ mkdir -p $GOPATH/src/github.com/logicalChimo; cd $GOPATH/src/github.com/logicalChimp
$ git clone git@github.com:logicalChimp/terraform-provider-custom
```

Enter the provider directory and build the provider

```bash
$ cd $GOPATH/src/github.com/logicalChimp/terraform-provider-custom
$ make build
```

## Using the provider (Linux)

Install the provider where Terraform can find it

```bash
$ mkdir -p ~/.terraform.d/plugins/linux_amx64
$ cp $GOPATH/bin/terraform-provider-custom ~/.terraform.d/plugins/linux_amx64
```

Use it in Terraform

```terraform
terraform {
    required_version = "0.12.20"
}

provider "custom" {
}

resource "custom_sequential_integer" "this" {
  keepers = {
    test = "value1"
  }
  min = 2
  max = 10
}

resource "custom_pinned_timestamp" "this" {
  triggers = {
    seq_id = "${custom_sequential_integer.this.value}"
  }
}

output "generated_number" {
  description = "Sequential number"
  value       = custom_sequential_integer.this.value
}

output "pinned_timestamp" {
  description = "Pinned timestamp"
  value       = custom_pinned_timestamp.this.timestamp
}
```