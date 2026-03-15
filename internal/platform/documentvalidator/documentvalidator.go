package documentvalidator

type Document string

type DocumentValidator interface {
	IsValid(Document, int) (bool, error)
	Sanitize(Document) Document
}

type documentValidatorImpl struct {
	DocumentNumber Document
}

func NewDocumentValidator(document Document) DocumentValidator {
	return documentValidatorImpl{
		DocumentNumber: document,
	}
}

func (ref documentValidatorImpl) IsValid(document Document, lenght int) (bool, error) {
	if lenght == 0 {
		lenght = len(ref.DocumentNumber)
	}
	switch lenght {
	case 11:
		return true, nil
	case 14:
		return true, nil
	default:
		return false, ErrInvalidDocumentLenght
	}
}

func (ref documentValidatorImpl) Sanitize(document Document) Document {

	return ""
}
