package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/devrapture/stack-trace/internal/models"
)

func main() {
	statements, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.UserEmail{},
		&models.AuthIdentity{},
		&models.PasswordCredential{},
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"failed to load GORM schema: %v\n",
			err,
		)

		os.Exit(1)
	}

	_, _ = io.WriteString(
		os.Stdout,
		statements,
	)
}
