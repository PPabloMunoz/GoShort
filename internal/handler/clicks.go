package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	QUERY_GET_FULL_URL = `
		select full_url
		from urls
		where code = ? and is_active = 1`
	QUERY_INSERT_CLICK = `
		insert into clicks (url_code)
		values (?)`
)

func (h *Handler) EnterURL(ctx *gin.Context) {
	code := ctx.Param("code")

	// Get Full URL
	var url string
	err := h.DB.QueryRow(QUERY_GET_FULL_URL, code).Scan(&url)
	if err != nil {
		log.Printf("[ERROR] Error getting full url - %v\n", err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/?error=not_found")
		return
	}

	// Add click data
	res, err := h.DB.Exec(QUERY_INSERT_CLICK, code)
	if err != nil {
		log.Printf("[ERROR] Error creating click for url_code '%s' in db - %v\n", code, err)
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		log.Printf("[ERROR] No click data created from url_code '%s' - %v\n", code, err)
	}

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
