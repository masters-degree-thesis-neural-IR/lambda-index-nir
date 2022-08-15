package usecases

type CreateDocumentUc interface {
	CreateDocument(title string, body string) error
}
