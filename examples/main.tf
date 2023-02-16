terraform {
  required_providers {
    cairvine = {
      source  = "local/edu/cairvine"
      version = "0.0.1"
    }
  }
  required_version = ">= 1.2.5"
}


provider "cairvine" {
  api_key_id = "XXXX"
  api_key    = "XXXX"
  endpoint   = "http://localhost:8018"
}

resource "cairvine_edge_value" "demo_bool" {
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
    exp     = "env == 'dev'"
  }

  targeting {
    variant = "on"
    spec    = "cel"
    exp     = "userId == 'XXX'"
  }
}

resource "cairvine_edge_value" "demo_string" {
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
    exp     = "env == 'dev'"
  }

  targeting {
    variant = "string02"
    spec    = "cel"
    exp     = "userId == 'XXX'"
  }
}


resource "cairvine_edge_value" "demo_json" {
  value_id        = "demo-json-value"
  enabled         = true
  description     = "demo json value"
  default_variant = "json01"

  json_value {
    variant = "json01"
    value   = "{\"name\": \"json01\"}"
  }

  json_value {
    variant = "json02"
    value   = "{\"name\": \"json02\"}"
  }

  targeting {
    variant = "json02"
    spec    = "cel"
    exp     = "env == 'dev'"
  }

  targeting {
    variant = "json02"
    spec    = "cel"
    exp     = "userId == 'XXX'"
  }
}

