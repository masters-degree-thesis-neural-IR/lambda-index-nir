package repositories

import "lambda-index-nir/service/application/domain"

type DocumentMetricsRepository interface {
	FindByDocumentIDs(documentIDs map[string]int8) ([]domain.NormalizedDocument, error)
	Save(domain.NormalizedDocument) error
}
