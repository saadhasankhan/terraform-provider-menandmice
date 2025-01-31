---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "menandmice_dhcp_scope Data Source - terraform-provider-menandmice"
subcategory: ""
description: |-

---

# menandmice_dhcp_scope (Data Source)



## Example Usage

```terraform
terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}

data menandmice_dhcp_scope scope1{
  dhcp_server= "micetro.example.net."
  cidr = "192.168.2.0/24"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **cidr** (String) The cidr of DHCPScope.

### Optional

- **dhcp_server** (String) The DHCP server of this scope.
- **id** (String) The ID of this resource.

### Read-Only

- **available** (Number) Number of available addresses in the address pool(s) of the scope.
- **description** (String) A description for the DHCP scope.
- **enabled** (Boolean) If this scope is enabled.
- **name** (String) The name of the DHCP scope you want to query.
- **ref** (String) Internal reference to this DHCP reservation.
- **superscope** (String) The name of the superscope for the DHCP scope. Only applicable for MS DHCP servers.
