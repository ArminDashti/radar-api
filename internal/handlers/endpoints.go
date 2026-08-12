package handlers

import (
	"net/http"
	"strings"

	"github.com/ArminDashti/radar-api/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) ListEndpoints(c *gin.Context) {
	rows, err := s.Pool.Query(requestContext(c), `
		SELECT id, name, host, http_enabled, icmp_enabled, probe_id, active, created_at
		FROM endpoints ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list endpoints"})
		return
	}
	defer rows.Close()

	endpoints := make([]models.Endpoint, 0)
	for rows.Next() {
		var endpoint models.Endpoint
		if err := rows.Scan(&endpoint.ID, &endpoint.Name, &endpoint.Host, &endpoint.HTTPEnabled,
			&endpoint.ICMPEnabled, &endpoint.ProbeID, &endpoint.Active, &endpoint.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read endpoints"})
			return
		}
		endpoints = append(endpoints, endpoint)
	}
	c.JSON(http.StatusOK, endpoints)
}

func (s *Server) CreateEndpoint(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Host        string `json:"host" binding:"required"`
		HTTPEnabled bool   `json:"http_enabled"`
		ICMPEnabled bool   `json:"icmp_enabled"`
		ProbeID     *int64 `json:"probe_id"`
		Active      bool   `json:"active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and host are required"})
		return
	}
	input.Name, input.Host = strings.TrimSpace(input.Name), strings.TrimSpace(input.Host)
	if input.Name == "" || input.Host == "" || (!input.HTTPEnabled && !input.ICMPEnabled) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, host, and at least one protocol are required"})
		return
	}
	var endpoint models.Endpoint
	err := s.Pool.QueryRow(requestContext(c), `
		INSERT INTO endpoints (name, host, http_enabled, icmp_enabled, probe_id, active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, host, http_enabled, icmp_enabled, probe_id, active, created_at`,
		input.Name, input.Host, input.HTTPEnabled, input.ICMPEnabled, input.ProbeID, input.Active,
	).Scan(&endpoint.ID, &endpoint.Name, &endpoint.Host, &endpoint.HTTPEnabled, &endpoint.ICMPEnabled,
		&endpoint.ProbeID, &endpoint.Active, &endpoint.CreatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not create endpoint"})
		return
	}
	c.JSON(http.StatusCreated, endpoint)
}
