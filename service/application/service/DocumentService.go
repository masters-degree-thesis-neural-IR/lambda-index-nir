package service

import (
	"github.com/google/uuid"
	"lambda-index-nir/service/application/domain"
	"lambda-index-nir/service/application/exception"
	"lambda-index-nir/service/application/ports"
	"lambda-index-nir/service/application/repositories"
	"lambda-index-nir/service/application/usecases"
)

type DocumentService struct {
	DocumentEvent      ports.DocumentEvent
	DocumentRepository repositories.DocumentRepository
}

func NewDocumentService(documentEvent ports.DocumentEvent, documentRepository repositories.DocumentRepository) usecases.CreateDocumentUc {
	var c usecases.CreateDocumentUc = &DocumentService{
		DocumentEvent:      documentEvent,
		DocumentRepository: documentRepository,
	}
	return c
}

func (s *DocumentService) CreateDocument(title string, body string) error {

	if title == "" {
		return *exception.ThrowValidationError("Invalid title from document")
	}

	if body == "" {
		return *exception.ThrowValidationError("Invalid body from document")
	}

	id, _ := uuid.NewRandom()

	document := domain.Document{
		Id:    id.String(),
		Title: title,
		Body:  body,
	}

	var err = s.DocumentRepository.Save(document)

	if err != nil {
		return err
	}

	return s.DocumentEvent.Created(document)

}
