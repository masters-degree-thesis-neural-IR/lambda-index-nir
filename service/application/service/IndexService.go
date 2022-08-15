package service

import (
	"lambda-index-nir/service/application/domain"
	"lambda-index-nir/service/application/exception"
	"lambda-index-nir/service/application/repositories"
	"lambda-index-nir/service/application/stopwords"
	"lambda-index-nir/service/application/usecases"
	"strings"
)

type IndexService struct {
	IndexRepository repositories.IndexRepository
}

func NewIndexService(indexRepository repositories.IndexRepository) usecases.CreateIndexUc {
	var c usecases.CreateIndexUc = IndexService{
		IndexRepository: indexRepository,
	}
	return c
}

func (i IndexService) CreateIndex(id string, title string, body string) error {

	tokens := Tokenizer(body, true)
	normalizedTokens, err := RemoveStopWords(tokens, "en")

	if err != nil {
		return err
	}

	document := domain.Document{
		Id:     id,
		Length: len(normalizedTokens),
		Tf:     TermFrequency(normalizedTokens),
	}

	for _, term := range normalizedTokens {

		index, err := i.IndexRepository.FindByTerm(term)

		if err != nil {
			return err
		}

		if index != nil {
			documentList := index.Documents
			if NotContains(document, documentList) {
				index.Term = term
				index.Documents = append(documentList, document)
				i.IndexRepository.Update(*index)
			}
		} else {
			index := domain.Index{
				Term:      term,
				Documents: []domain.Document{document},
			}
			i.IndexRepository.Save(index)
		}
	}

	return nil
}

func NotContains(document domain.Document, documents []domain.Document) bool {

	for _, doc := range documents {
		if doc.Id == document.Id {
			return false
		}
	}

	return true
}

func Tokenizer(document string, normalize bool) []string {
	fields := strings.Fields(document)

	if normalize {

		localSlice := make([]string, len(fields))
		for i, token := range fields {
			localSlice[i] = strings.ToLower(token)
		}

		return localSlice
	}

	return fields

}

func StopWordLang(lang string) (map[string]bool, error) {

	if lang == "en" {
		return stopwords.English, nil
	}

	return nil, *exception.ThrowValidationError("Not found language from stop word")
}

func RemoveStopWords(tokens []string, lang string) ([]string, error) {

	stopWordLang, err := StopWordLang(lang)

	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return make([]string, 0), nil
	}

	var localSlice = make([]string, 0)

	for _, token := range tokens {
		if !stopWordLang[token] {
			localSlice = append(localSlice, token)
		}
	}

	return localSlice, nil

}

func TermFrequency(tokens []string) map[string]int {

	localMap := make(map[string]int)

	for _, token := range tokens {

		if localMap[token] == 0 {
			localMap[token] = 1
		} else {
			localMap[token] = localMap[token] + 1
		}
	}

	return localMap

}
