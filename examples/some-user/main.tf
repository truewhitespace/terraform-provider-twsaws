terraform {
  required_providers {
    twsaws = {
      version = "0.1"
      source  = "github.com/truewhitespace/twsaws"
    }
  }
}

variable "user_name" {
  type    = string
}

resource "twsaws_rotating_keys" "key" {
    user_name = var.user_name
}

# Returns all coffees
output "active_key_id" {
  value = twsaws_rotating_keys.key.active_key_id
}

# Only returns packer spiced latte
output "active_key_secret" {
  value = twsaws_rotating_keys.key.active_key_secret
}
