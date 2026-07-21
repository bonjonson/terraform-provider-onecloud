terraform {
  required_providers {
    onecloud = {
      source  = "onecloud"
      version = "1.0"
    }
  }
}

provider "onecloud" {
  api_token = var.onecloud_token
}
