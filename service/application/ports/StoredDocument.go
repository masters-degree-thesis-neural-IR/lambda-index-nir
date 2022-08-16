package ports

import "lambda-index-nir/service/application/domain"

type Store interface {
	StoreDocument(document domain.NormalizedDocument) error
}
