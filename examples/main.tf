terraform {
  required_providers {
    twsaws = {
      version = "1.0.0"
      source  = "truewhitespace/twsaws"
    }
  }
}

provider "twsaws" {
  backend = "localstack"
  default_key_expiry = "1m"
  default_key_grace = "30s"
}

module "u" {
  source = "./some-user"

  user_name = "cred-test-user"
}

output "key" {
  value = module.u.active_key_id
}

output "secret" {
  value = module.u.active_key_secret
}
