package repositories

import (
	"lambda-index-nir/service/application/domain"
)

type DocumentRepository interface {
	Save(document domain.Document) error
}
