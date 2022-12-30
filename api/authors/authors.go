package authors

import (
	"Rest-with-postgres/internal/database"
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service struct {
	queries *database.Queries
}

func NewService(queries *database.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// RegisterHandlers - Defining handlers
func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/authors", s.Create)
	router.GET("/author/:id", s.Get)
	router.PUT("/author/:id", s.FullUpdate)
	router.PATCH("/author/:id", s.PartialUpdate)
	router.DELETE("/author/:id", s.Delete)
	router.GET("/authors", s.List)
}

type apiAuthor struct {
	ID   int64
	Name string `json:"name,omitempty" binding:"required"`
	Bio  string `json:"bio,omitempty" binding:"required"`
}

func fromDB(author database.Author) *apiAuthor {
	//Method converts a database author to an API author
	return &apiAuthor{
		ID:   author.ID,
		Name: author.Name,
		Bio:  author.Bio,
	}
}

func (s *Service) Create(c *gin.Context) {
	//Parse the request
	var request apiAuthor
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Create author
	params := database.CreateAuthorParams{
		Name: request.Name,
		Bio:  request.Bio,
	}
	author, err := s.queries.CreateAuthor(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Build a response
	response := fromDB(author)
	c.IndentedJSON(http.StatusCreated, response)
}

type pathParameters struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Service) Get(c *gin.Context) {
	//Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Get author
	author, err := s.queries.GetAuthor(context.Background(), pathParams.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Build a response
	response := fromDB(author)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) Delete(c *gin.Context) {
	//Parse
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Delete author
	if err := s.queries.DeleteAuthor(context.Background(), pathParams.ID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Response
	c.Status(http.StatusOK)
}

func (s *Service) List(c *gin.Context) {
	// List authors
	authors, err := s.queries.ListAuthors(context.Background())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	if len(authors) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Build response
	var response []*apiAuthor
	for _, author := range authors {
		response = append(response, fromDB(author))
	}
	c.IndentedJSON(http.StatusOK, authors)
}

func (s *Service) FullUpdate(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request apiAuthor
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update author
	params := database.UpdateAuthorParams{
		ID:   pathParams.ID,
		Name: request.Name,
		Bio:  request.Bio,
	}
	author, err := s.queries.UpdateAuthor(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(author)
	c.IndentedJSON(http.StatusOK, response)
}

type apiAuthorPartialUpdate struct {
	Name *string `json:"name,omitempty" binding:"omitempty,max=32"`
	Bio  *string `json:"bio,omitempty" binding:"omitempty"`
}

func (s *Service) PartialUpdate(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request apiAuthorPartialUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update author
	params := database.PartialUpdateAuthorParams{ID: pathParams.ID}
	if request.Name != nil {
		params.UpdateName = true
		params.Name = *request.Name
	}
	if request.Bio != nil {
		params.UpdateBio = true
		params.Bio = *request.Bio
	}
	author, err := s.queries.PartialUpdateAuthor(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(author)
	c.IndentedJSON(http.StatusOK, response)
}
