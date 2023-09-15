package APIHandlers

import (
	"encoding/json"
	"io"
	"library/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthorRequestBody struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Nationality string `json:"nationality"`
	Birthday    string `json:"birthday"`
}

type bookBody struct {
	Name            string            `json:"name"`
	Author          AuthorRequestBody `json:"author"`
	Category        string            `json:"category"`
	Volume          int               `json:"volume"`
	PublishedAt     string            `json:"published_at"`
	Summary         string            `json:"summary"`
	Publisher       string            `json:"publisher"`
	TableOfContents []string          `json:"table_of_contents"`
}

func (s *Server) HandleCreateAndGetAllBook(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		s.Logger.Error("can not find Authorization")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("there is no token for Authorization"))
		return
	}

	username, err := s.Authenticate.GetUsernameByToken(token)
	if err != nil {
		s.Logger.Error("can not find user by this Authorization")
		w.WriteHeader(http.StatusUnauthorized)
		s.Logger.Error("can not find user by this token")
		return
	}

	if r.Method == http.MethodPost {
		s.HandleCreateBook(w, r, *username)
	} else if r.Method == http.MethodGet {
		s.HandleGetAllBooks(w, r)
	} else {
		s.Logger.Error("http method is not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("http method is not allowed"))
		return
	}
}

func (s *Server) HandleCreateBook(w http.ResponseWriter, r *http.Request, username string) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.Logger.WithError(err).Warn("can not read the request data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var BD bookBody
	err = json.Unmarshal(body, &BD)
	if err != nil {
		s.Logger.WithError(err).Warn("can not unmarshal the body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	PublishedAt, _ := time.Parse("2006-01-02", BD.PublishedAt)
	date, _ := time.Parse("2006-01-02", BD.Author.Birthday)
	book := &db.Book{
		Name:              BD.Name,
		Category:          BD.Category,
		Volume:            BD.Volume,
		PublishedAt:       PublishedAt,
		Summary:           BD.Summary,
		PublisherName:     BD.Publisher,
		TableOfContents:   BD.TableOfContents,
		Owner:             username,
		AuthorFirstName:   BD.Author.FirstName,
		AuthorLastName:    BD.Author.LastName,
		AuthorNationality: BD.Author.Nationality,
		AuthorBirthday:    date,
	}

	err = s.DB.AddNewBook(book)
	if err != nil {
		s.Logger.WithError(err).Warn("can not add the new book: ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message": "book has been added successfully",
	}

	resBody, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

func (s *Server) HandleGetAllBooks(w http.ResponseWriter, r *http.Request) {

	var books *[]db.Book
	books, err := s.DB.GetAllBooks()
	if err != nil {
		s.Logger.WithError(err).Warn("can not get all books: ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var responseBooks []bookBody


	for _, BD := range *books {
		author := AuthorRequestBody{
			FirstName:   BD.AuthorFirstName,
			LastName:    BD.AuthorLastName,
			Nationality: BD.AuthorNationality,
			Birthday:    BD.AuthorBirthday.String(),
		}
		response := bookBody{
			Name:            BD.Name,
			Category:        BD.Category,
			Volume:          BD.Volume,
			PublishedAt:     BD.PublishedAt.String(),
			Summary:         BD.Summary,
			Publisher:       BD.PublisherName,
			TableOfContents: BD.TableOfContents,
			Author: author,
		}
		responseBooks = append(responseBooks, response)
	}
	resBody, err := json.Marshal(responseBooks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

func (s *Server) HandleUpdateAndDEleteAndGetBook(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		s.Logger.Error("can not find Authorization")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("there is no token for Authorization"))
		return
	}

	username, err := s.Authenticate.GetUsernameByToken(token)
	if err != nil {
		s.Logger.Error("can not find user by this Authorization")
		w.WriteHeader(http.StatusUnauthorized)
		s.Logger.Error("can not find user by this token")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	bookID := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(bookID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		s.HandleDeleteBookByid(w, r, id, *username)
	} else if r.Method == http.MethodGet {
		s.HandleGetBookByid(w, r, id)
	} else if r.Method == http.MethodPut {
		s.HandleUpdateBookByid(w, r, id, *username)
	} else {
		s.Logger.Error("http method is not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("http method is not allowed"))
		return
	}
}

func (s *Server) HandleGetBookByid(w http.ResponseWriter, r *http.Request, id int) {

	book, err := s.DB.GetBookByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("error : ")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("book by this id is not found"))
		return
	}

	author := AuthorRequestBody{
		FirstName:   book.AuthorFirstName,
		LastName:    book.AuthorLastName,
		Nationality: book.AuthorNationality,
		Birthday:    book.AuthorBirthday.String(),
	}
	response := bookBody{
		Name:            book.Name,
		Category:        book.Category,
		Volume:          book.Volume,
		PublishedAt:     book.PublishedAt.String(),
		Summary:         book.Summary,
		Publisher:       book.PublisherName,
		TableOfContents: book.TableOfContents,
		Author: author,
	}
	resBody, err := json.Marshal(response)
	if err != nil {
		s.Logger.WithError(err).Warn("can not unmarshal the body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)

}

func (s *Server) HandleDeleteBookByid(w http.ResponseWriter, r *http.Request, id int, username string) {
	book, err := s.DB.GetBookByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("error : ")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("book by this id is not found"))
		return
	}
	if book.Owner != username {
		s.Logger.Warn("user is not book owner to delete")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("you are not this book onwer to delete it"))
		return
	}
	err = s.DB.DeleteBookByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("error : ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can not delete because of some internal error"))
		return
	}
	response := map[string]interface{}{
		"message": "book has been deleted successfully",
	}
	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)

}

func (s *Server) HandleUpdateBookByid(w http.ResponseWriter, r *http.Request, id int, username string) {
	book, err := s.DB.GetBookByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("error : ")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("book by this id is not found"))
		return
	}
	if book.Owner != username {
		s.Logger.Warn("user is not book owner to delete")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("you are not this book onwer to update it"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.Logger.WithError(err).Warn("can not read the request data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var BD bookBody
	err = json.Unmarshal(body, &BD)
	if err != nil {
		s.Logger.WithError(err).Warn("can not unmarshal the body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if BD.Name != "" {
		book.Name = BD.Name
	}

	if BD.Author.FirstName != "" {
		book.AuthorFirstName = BD.Author.FirstName
	}

	if BD.Author.LastName != "" {
		book.AuthorLastName = BD.Author.LastName
	}

	if BD.Author.Birthday != "" {
		date, _ := time.Parse("2006-01-02", BD.Author.Birthday)
		book.AuthorBirthday = date
	}

	if BD.Author.Nationality != "" {
		book.AuthorNationality = BD.Author.Nationality
	}

	if BD.Category != "" {
		book.Category = BD.Category
	}

	if BD.Volume != 0 {
		book.Volume = BD.Volume
	}

	if BD.PublishedAt != "" {
		date, _ := time.Parse("2006-01-02", BD.PublishedAt)
		book.PublishedAt = date
	}

	if BD.Summary != "" {
		book.Summary = BD.Summary
	}

	if len(BD.TableOfContents) > 0 {
		book.TableOfContents = BD.TableOfContents
	}

	if BD.Publisher != "" {
		book.PublisherName = BD.Publisher
	}

	err = s.DB.UpdateBook(book)
	if err != nil {
		s.Logger.WithError(err).Error("error : ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can not update because of some internal error"))
		return
	}
	response := map[string]interface{}{
		"message": "book has been updated successfully",
	}
	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
