package service

import (
	"errors"
	"github.com/lyh-demo/go-webapp-demo/container"
	"github.com/lyh-demo/go-webapp-demo/model"
	"github.com/lyh-demo/go-webapp-demo/model/dto"
	"github.com/lyh-demo/go-webapp-demo/repository"
	"github.com/lyh-demo/go-webapp-demo/util"
)

// BookService is a service for managing books.
type BookService interface {
	FindByID(id string) (*model.Book, error)
	FindAllBooks() (*[]model.Book, error)
	FindAllBooksByPage(page string, size string) (*model.Page, error)
	FindBooksByTitle(title string, page string, size string) (*model.Page, error)
	CreateBook(dto *dto.BookDto) (*model.Book, map[string]string)
	UpdateBook(dto *dto.BookDto, id string) (*model.Book, map[string]string)
	DeleteBook(id string) (*model.Book, map[string]string)
}

type bookService struct {
	container container.Container
}

// NewBookService is constructor.
func NewBookService(container container.Container) BookService {
	return &bookService{container: container}
}

// FindByID returns one record matched book's id.
func (b *bookService) FindByID(id string) (*model.Book, error) {
	if !util.IsNumeric(id) {
		return nil, errors.New("failed to fetch data")
	}

	rep := b.container.GetRepository()
	book := model.Book{}
	var result *model.Book
	var err error
	if result, err = book.FindByID(rep, util.ConvertToUint(id)).Take(); err != nil {
		return nil, err
	}
	return result, nil
}

// FindAllBooks returns the list of all books.
func (b *bookService) FindAllBooks() (*[]model.Book, error) {
	rep := b.container.GetRepository()
	book := model.Book{}
	result, err := book.FindAll(rep)
	if err != nil {
		b.container.GetLogger().GetZapLogger().Errorf(err.Error())
		return nil, err
	}
	return result, nil
}

// FindAllBooksByPage returns the page object of all books.
func (b *bookService) FindAllBooksByPage(page string, size string) (*model.Page, error) {
	rep := b.container.GetRepository()
	book := model.Book{}
	result, err := book.FindAllByPage(rep, page, size)
	if err != nil {
		b.container.GetLogger().GetZapLogger().Errorf(err.Error())
		return nil, err
	}
	return result, nil
}

// FindBooksByTitle returns the page object of books matched given book title.
func (b *bookService) FindBooksByTitle(title string, page string, size string) (*model.Page, error) {
	rep := b.container.GetRepository()
	book := model.Book{}
	result, err := book.FindByTitle(rep, title, page, size)
	if err != nil {
		b.container.GetLogger().GetZapLogger().Errorf(err.Error())
		return nil, err
	}
	return result, nil
}

// CreateBook register the given book data.
func (b *bookService) CreateBook(dto *dto.BookDto) (*model.Book, map[string]string) {
	if e := dto.Validate(); e != nil {
		return nil, e
	}

	rep := b.container.GetRepository()
	var result *model.Book
	var err error

	if trErr := rep.Transaction(func(txRep repository.Repository) error {
		result, err = txCreateBook(txRep, dto)
		return err
	}); trErr != nil {
		b.container.GetLogger().GetZapLogger().Errorf(trErr.Error())
		return nil, map[string]string{"error": "Failed to the registration"}
	}
	return result, nil
}

func txCreateBook(txRep repository.Repository, dto *dto.BookDto) (*model.Book, error) {
	var result *model.Book
	var err error
	book := dto.Create()

	category := model.Category{}
	if book.Category, err = category.FindByID(txRep, dto.CategoryID).Take(); err != nil {
		return nil, err
	}

	format := model.Format{}
	if book.Format, err = format.FindByID(txRep, dto.FormatID).Take(); err != nil {
		return nil, err
	}

	if result, err = book.Create(txRep); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateBook updates the given book data.
func (b *bookService) UpdateBook(dto *dto.BookDto, id string) (*model.Book, map[string]string) {
	if e := dto.Validate(); e != nil {
		return nil, e
	}

	rep := b.container.GetRepository()
	var result *model.Book
	var err error

	if trErr := rep.Transaction(func(txRep repository.Repository) error {
		result, err = txUpdateBook(txRep, dto, id)
		return err
	}); trErr != nil {
		b.container.GetLogger().GetZapLogger().Errorf(trErr.Error())
		return nil, map[string]string{"error": "Failed to the update"}
	}
	return result, nil
}

func txUpdateBook(txRep repository.Repository, dto *dto.BookDto, id string) (*model.Book, error) {
	var book, result *model.Book
	var err error

	b := model.Book{}
	if book, err = b.FindByID(txRep, util.ConvertToUint(id)).Take(); err != nil {
		return nil, err
	}

	book.Title = dto.Title
	book.Isbn = dto.Isbn
	book.CategoryID = dto.CategoryID
	book.FormatID = dto.FormatID

	category := model.Category{}
	if book.Category, err = category.FindByID(txRep, dto.CategoryID).Take(); err != nil {
		return nil, err
	}

	format := model.Format{}
	if book.Format, err = format.FindByID(txRep, dto.FormatID).Take(); err != nil {
		return nil, err
	}

	if result, err = book.Update(txRep); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteBook deletes the given book data.
func (b *bookService) DeleteBook(id string) (*model.Book, map[string]string) {
	rep := b.container.GetRepository()
	var result *model.Book
	var err error

	if trErr := rep.Transaction(func(txRep repository.Repository) error {
		result, err = txDeleteBook(txRep, id)
		return err
	}); trErr != nil {
		b.container.GetLogger().GetZapLogger().Errorf(trErr.Error())
		return nil, map[string]string{"error": "Failed to the delete"}
	}
	return result, nil
}

func txDeleteBook(txRep repository.Repository, id string) (*model.Book, error) {
	var book, result *model.Book
	var err error

	b := model.Book{}
	if book, err = b.FindByID(txRep, util.ConvertToUint(id)).Take(); err != nil {
		return nil, err
	}

	if result, err = book.Delete(txRep); err != nil {
		return nil, err
	}

	return result, nil
}
