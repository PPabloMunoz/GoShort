package handler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ppablomunoz/GoShort/internal/models"
	"github.com/ppablomunoz/GoShort/internal/utils"
)

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && (u.Scheme == "http" || u.Scheme == "https")
}

const (
	QUERY_INSERT_URL = `
		insert into urls (code, full_url) 
		values (?, ?) 
		returning code, full_url, is_active, created_at`
	QUERY_SELECT_URLS = `
		select u.code, u.full_url, u.is_active, u.created_at, count(c.id) as click_count
		from urls u
		left join clicks c on u.code = c.url_code
		group by u.code`
	QUERY_UPDATE_URL = `
		update urls
		set full_url = ?, is_active = ?
		where code = ?`
	QUERY_DELETE_URL = `
		delete from urls
		where code = ?`
)

func (h *Handler) NewURL(ctx *gin.Context) {
	var body models.NewURL
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil || body.FullURL == "" {
		log.Printf("[ERROR] Getting the body - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "The body is invalid"})
		return
	}

	if !isValidURL(body.FullURL) {
		log.Printf("[ERROR] Invalid URL format - %s\n", body.FullURL)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format. Only http and https are allowed."})
		return
	}

	// Create url
	var newUrl models.URL
	code, err := utils.GenerateShortCode()
	if err != nil {
		log.Printf("[ERROR] Creating the code - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error ocurred while creating the code for the url"})
		return
	}
	err = h.DB.QueryRow(QUERY_INSERT_URL, code, body.FullURL).Scan(
		&newUrl.Code,
		&newUrl.FullURL,
		&newUrl.IsActive,
		&newUrl.CreatedAt,
	)
	if err != nil {
		log.Printf("[ERROR] Inserting to DB - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not insert into db"})
		return
	}

	log.Printf("[INFO] New url created - '%s'", newUrl.Code)
	ctx.JSON(http.StatusCreated, newUrl)
}

func (h *Handler) GetURLs(ctx *gin.Context) {
	urls := make([]models.URL, 0)

	rows, err := h.DB.Query(QUERY_SELECT_URLS)
	if err != nil {
		log.Printf("[ERROR] Error getting data from DB - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not get data from db"})
		return
	}

	for rows.Next() {
		var url models.URL
		err = rows.Scan(&url.Code, &url.FullURL, &url.IsActive, &url.CreatedAt, &url.ClickCount)
		if err != nil {
			log.Printf("[ERROR] Error going through the rows - %v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not get data from db"})
			return
		}
		urls = append(urls, url)
	}

	ctx.JSON(http.StatusOK, urls)
}

func (h *Handler) UpdateURL(ctx *gin.Context) {
	code := ctx.Param("code")

	var body models.UpdateURL
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil || body.FullURL == "" {
		log.Printf("[ERROR] Getting the body - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "The body is invalid"})
		return
	}

	if !isValidURL(body.FullURL) {
		log.Printf("[ERROR] Invalid URL format - %s\n", body.FullURL)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format. Only http and https are allowed."})
		return
	}

	res, err := h.DB.Exec(QUERY_UPDATE_URL, body.FullURL, body.IsActive, code)
	if err != nil {
		log.Printf("[ERROR] Error updating DB - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while updating db"})
		return
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		log.Printf("[ERROR] No rows updated - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URL is not found"})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) DeleteURL(ctx *gin.Context) {
	code := ctx.Param("code")

	res, err := h.DB.Exec(QUERY_DELETE_URL, code)
	if err != nil {
		log.Printf("[ERROR] Error deleting in DB - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred while deleting in db"})
		return
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		log.Printf("[ERROR] No rows deleted - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URL is not found"})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
