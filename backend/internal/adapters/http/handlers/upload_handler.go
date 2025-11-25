package handlers

import (
	"net/http"

	"backend/internal/adapters/storage"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storage storage.FileStorage
}

func NewUploadHandler(storage storage.FileStorage) *UploadHandler {
	return &UploadHandler{storage: storage}
}

type UploadResponse struct {
	URL string `json:"url"`
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload file to MinIO (avatar, course_cover, lesson_material)
// @Tags upload
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param type formData string true "File type (avatar, cover, lesson)"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// // @Router /v1/upload [post]
func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "file is required"})
		return
	}

	fileType := c.PostForm("type")
	var folder string

	// Определяем папку в зависимости от типа
	switch fileType {
	case "avatar":
		folder = "avatars"
	case "cover":
		folder = "covers"
	case "lesson":
		folder = "lessons"
	default:
		folder = "misc"
	}

	url, err := h.storage.UploadFile(c.Request.Context(), file, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "upload failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{URL: url})
}
