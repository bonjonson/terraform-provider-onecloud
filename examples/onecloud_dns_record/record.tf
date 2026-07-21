resource "onecloud_dns_record" "www" {
  domain_id = onecloud_dns_domain.example.id
  type      = "A"
  name      = "www"
  content   = "192.168.1.1"
  ttl       = "300"
}

resource "onecloud_dns_record" "cname" {
  domain_id     = onecloud_dns_domain.example.id
  type          = "CNAME"
  name          = "example.com."
  mnemonic_name = "imagination"
  ttl           = "5"
}

resource "onecloud_dns_record" "txt" {
  domain_id = onecloud_dns_domain.example.id
  type      = "TXT"
  name      = "text-item"
  text      = "abcdefghijklmnopqrstuvwxyz"
  ttl       = "60"
}

resource "onecloud_dns_record" "srv" {
  domain_id = onecloud_dns_domain.example.id
  type      = "SRV"
  name      = "service"
  proto     = "tcp"
  priority  = "1"
  weight    = "1"
  port      = "54444"
  target    = "target.example.com."
  service   = "example"
}
