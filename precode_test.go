package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	status := responseRecorder.Code
	require.Equal(t, http.StatusOK, status)

	body := responseRecorder.Body
	bodySlice := strings.Split(body.String(), ",")
	require.NotEmpty(t, body)

	url := strings.Split(req.URL.String(), "=")
	expectedCountStr := string(url[1][0])
	expectedCount, _ := strconv.Atoi(expectedCountStr)
	require.Len(t, bodySlice, expectedCount)

	city := string(url[2])
	expectedAnswer := strings.Join(cafeList[city][:expectedCount], ",")
	require.Equal(t, expectedAnswer, body.String())
}

func TestMainHandlerWhenWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscoww", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	status := responseRecorder.Code
	assert.Equal(t, http.StatusBadRequest, status)

	body := responseRecorder.Body.String()
	require.Equal(t, "wrong city value", body)
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	status := responseRecorder.Code
	require.Equal(t, http.StatusOK, status)

	body := responseRecorder.Body
	bodySlice := strings.Split(body.String(), ",")
	assert.Len(t, bodySlice, totalCount)

	url := strings.Split(req.URL.String(), "=")
	city := string(url[2])
	expectedAnswer := strings.Join(cafeList[city], ",")
	require.Equal(t, expectedAnswer, body.String())
}
