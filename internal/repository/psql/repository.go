package psql

import "avito_intern/pkg/client"

type Repository interface {
	Tx
	Banner
	BannerHistory
	User
}

type Pg struct {
	client client.Client
}

func New(client client.Client) Repository {
	return Pg{client}
}
