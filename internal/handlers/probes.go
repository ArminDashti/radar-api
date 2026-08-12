package handlers

import (
	"net/http"

	"github.com/ArminDashti/radar-api/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) ListProbes(c *gin.Context) {
	rows, err := s.Pool.Query(requestContext(c),
		`SELECT id, code, name, flag_icon, created_at FROM probes ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list probes"})
		return
	}
	defer rows.Close()

	probes := make([]models.Probe, 0)
	for rows.Next() {
		var probe models.Probe
		if err := rows.Scan(&probe.ID, &probe.Code, &probe.Name, &probe.FlagIcon, &probe.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read probes"})
			return
		}
		probes = append(probes, probe)
	}
	c.JSON(http.StatusOK, probes)
}
