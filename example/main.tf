terraform {
  required_providers {
    tfsdkvalidationdemo = {
      source  = "example.com/chrismarget/tfsdkvalidationdemo"
    }
  }
}

provider "tfsdkvalidationdemo" {}

resource "tfsdkvalidationdemo_bogus" "x" {
  count = 3
  validated = count.index
}
