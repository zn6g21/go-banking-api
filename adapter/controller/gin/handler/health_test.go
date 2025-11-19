package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/health", nil)
	gincontext, _ := gin.CreateTestContext(w)
	gincontext.Request = request

	Health(gincontext)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}
