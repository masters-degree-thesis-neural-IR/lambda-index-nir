package ports

import "lambda-index-nir/service/application/domain"

type DocumentEvent interface {
	Created(document domain.NormalizedDocument) error
}
