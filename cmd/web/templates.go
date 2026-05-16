package main

import (
	"snippetbox.abdulmoiz.net/internal/models"
)

type templateData struct {
	Snippet *models.Snippet
	Snippets [] *models.Snippet
}
