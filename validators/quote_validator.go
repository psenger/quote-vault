package validators

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type QuoteValidator struct{}

func NewQuoteValidator() *QuoteValidator {
	return &QuoteValidator{}
}

func (v *QuoteValidator) ValidateQuote(text, author, category string) error {
	if err := v.validateText(text); err != nil {
		return err
	}
	if err := v.validateAuthor(author); err != nil {
		return err
	}
	if err := v.validateCategory(category); err != nil {
		return err
	}
	return nil
}

func (v *QuoteValidator) validateText(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("quote text cannot be empty")
	}
	if utf8.RuneCountInString(text) < 10 {
		return errors.New("quote text must be at least 10 characters long")
	}
	if utf8.RuneCountInString(text) > 1000 {
		return errors.New("quote text cannot exceed 1000 characters")
	}
	return nil
}

func (v *QuoteValidator) validateAuthor(author string) error {
	author = strings.TrimSpace(author)
	if author == "" {
		return errors.New("author cannot be empty")
	}
	if utf8.RuneCountInString(author) > 100 {
		return errors.New("author name cannot exceed 100 characters")
	}
	return nil
}

func (v *QuoteValidator) validateCategory(category string) error {
	category = strings.TrimSpace(category)
	if category == "" {
		return errors.New("category cannot be empty")
	}
	if utf8.RuneCountInString(category) > 50 {
		return errors.New("category cannot exceed 50 characters")
	}
	// Allow only alphanumeric characters, spaces, and hyphens
	for _, r := range category {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
			 (r >= '0' && r <= '9') || r == ' ' || r == '-') {
			return errors.New("category can only contain letters, numbers, spaces, and hyphens")
		}
	}
	return nil
}