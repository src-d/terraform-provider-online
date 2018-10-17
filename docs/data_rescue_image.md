# Data: rescue_image

A Rescue Image is a minimal operating system capable of booting over the network to allow modification of files to fix an issue where the server cannot boot.

## Example Usage

Explicit name example

```HCL
data "online_rescue_image" "example" {
    name   = "rescue-image-name"
    server ="${online_server.example_server.id}"
}
```

Partial name match example

```HCL
data "online_rescue_image" "example" {
    name_filter = "rescue-image-partial-"
    server ="${online_server.example_server.id}"
}

## Argument Reference
* `name` - (Required) Exact name of the desired image, conflicts with `name_filter`
* `name_filter` - (Required) Partial name of the desired image to filter with, in case multiple get found the last one will be used, conflicts with `name`
* server - (Optional) Server for the desired image

## Attributes Reference
* `image` - The requested image