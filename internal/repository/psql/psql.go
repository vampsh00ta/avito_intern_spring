package psql

import (
	"avito_intern/pkg/client"
)

type Repository interface {
}

type Pg struct {
	client client.Client
}
