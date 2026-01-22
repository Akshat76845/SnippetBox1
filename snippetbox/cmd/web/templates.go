package main
import "snippetbox.alexedwards.net/internal/models"

type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
}