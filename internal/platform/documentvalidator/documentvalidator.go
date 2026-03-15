package documentvalidator

type Document string

type DocumentValidator interface {
	IsValid(Document, int) (bool, error)
	Sanitize(Document) Document
}

type documentValidatorImpl struct{}

func (ref documentValidatorImpl) IsValid(document Document, lenght int) (bool, error) {
	if lenght == 0 {

		lenght = len(document)
	}
	switch lenght {
	case 11:
	case 14:
	default:
		return false, ErrInvalidDocumentLenght
	}
	return false, nil
}

func (ref documentValidatorImpl) Sanitize(document Document) Document {
	return ""
}
