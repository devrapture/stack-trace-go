data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlasloader",
  ]
}

env "local" {
  src = data.external_schema.gorm.url

  url = "postgresql://wordsearch:wordsearch@localhost:5433/wordsearch?sslmode=disable"
  dev = "postgresql://wordsearch:wordsearch@localhost:5433/wordsearch_dev?sslmode=disable"

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "production" {
  url = getenv("DATABASE_URL")
  dev = getenv("DEV_DATABASE_URL")
  src = data.external_schema.gorm.url
  migration {
    dir = "file://migrations"
    revisions_schema = "public"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}