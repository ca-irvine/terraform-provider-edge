terraform {
  required_providers {
    edge = {
      source  = "ca-irvine/edge"
      version = "0.1.2"
    }
  }
}

provider "edge" {
  api_key_id = var.api_key_id
  api_key    = var.api_key
  endpoint   = var.endpoint
}

variable "api_key_id" {
  type = string
}

variable "api_key" {
  type = string
}

variable "endpoint" {
  type = string
}

resource "edge_value" "demo_bool" {
  value_id        = "demo-bool-value"
  enabled         = true
  description     = "demo bool value"
  default_variant = "off"

  boolean_value {
    variant = "on"
    value   = true
  }

  boolean_value {
    variant = "off"
    value   = false
  }

  targeting {
    variant = "on"
    spec    = "cel"
    expr    = "env == 'dev'"
  }

  targeting {
    variant = "on"
    spec    = "cel"
    expr    = "userId == 'XXX'"
  }
}

resource "edge_value" "demo_string" {
  value_id        = "demo-string-value"
  enabled         = true
  description     = "demo string value"
  default_variant = "string01"

  string_value {
    variant = "string01"
    value   = "string01"
  }

  string_value {
    variant = "string02"
    value   = "string02"
  }

  targeting {
    variant = "string02"
    spec    = "cel"
    expr    = "env == 'dev'"
  }

  targeting {
    variant = "string02"
    spec    = "cel"
    expr    = "userId == 'XXX'"
  }
}

resource "edge_value" "demo_json" {
  value_id        = "demo-json-value"
  enabled         = true
  description     = "demo json value"
  default_variant = "json01"

  json_value {
    variant = "json01"
    value   = jsonencode({
      "items" : [
        {"viewable": true, "content": "content1"},
        {"viewable": true, "content": "content2"},
        {"viewable": false, "content": "content3"}
      ]
    })
    transform {
      spec = "cel"
      expr = "{\"items\":items.map(item, item.viewable ? item : item.deleteKey([\"content\"]))}"
    }
  }

  json_value {
    variant = "json01"
    value   = jsonencode({
      "items" : [
        {"viewable": true, "content": "content1"},
        {"viewable": false, "content": "content2"},
        {"viewable": false, "content": "content3"}
      ]
    })
    transform {
      spec = "cel"
      expr = "{\"items\":items.map(item, item.viewable ? item.selectKey([\"content\"]) : item)}"
    }
  }

  targeting {
    variant = "json01"
    spec    = "cel"
    expr    = "env == 'dev'"
  }

  targeting {
    variant = "json01"
    spec    = "cel"
    expr    = "userId == 'XXX'"
  }
}
