package service

import (
	"lambda-index-nir/service/application/logger"
	"lambda-index-nir/service/application/nlp"
	"lambda-index-nir/service/application/repositories"
	"lambda-index-nir/service/application/usecases"
)

type IndexService struct {
	Logger                logger.Logger
	IndexMemoryRepository repositories.IndexMemoryRepository
}

func NewIndexService(logger logger.Logger, indexMemoryRepository repositories.IndexMemoryRepository) usecases.CreateIndexUc {

	return IndexService{
		Logger:                logger,
		IndexMemoryRepository: indexMemoryRepository,
	}
}

func (i IndexService) CreateIndex(id string, title string, body string) error {

	tokens := nlp.Tokenizer(body, true)
	normalizedTokens, err := nlp.RemoveStopWords(tokens, "en")

	if err != nil {
		i.Logger.Error(err.Error())
		return err
	}

	for _, term := range normalizedTokens {

		var documentList, err = i.IndexMemoryRepository.FindByTerm(term)

		if err != nil {
			i.Logger.Error(err.Error())
			return err
		}

		if documentList != nil && len(documentList) > 0 {
			if nlp.NotContains(id, documentList) {
				documentList = append(documentList, id)
				i.IndexMemoryRepository.Update(term, documentList)
			}
		} else {
			i.IndexMemoryRepository.Save(term, []string{id})
		}
	}

	return nil
}
