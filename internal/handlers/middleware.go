package handlers

import (
	"context"
	"net/http"
	"strings"

	radarauth "github.com/ArminDashti/radar-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	agentIDKey      = "agent_id"
	agentProbeIDKey = "agent_probe_id"
)

type Server struct {
	Pool      *pgxpool.Pool
	JWTSecret string
}

func (s *Server) WebAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		claims, err := radarauth.ParseToken(s.JWTSecret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}

func (s *Server) AgentAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing agent bearer token"})
			return
		}
		hash := radarauth.HashAgentToken(token)
		var storedHash string
		var agentID, probeID int64
		err := s.Pool.QueryRow(c.Request.Context(),
			`SELECT id, probe_id, token_hash FROM agents WHERE token_hash = $1`, hash,
		).Scan(&agentID, &probeID, &storedHash)
		if err != nil || !radarauth.EqualTokenHash(storedHash, hash) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid agent token"})
			return
		}
		c.Set(agentIDKey, agentID)
		c.Set(agentProbeIDKey, probeID)
		c.Next()
	}
}

func bearerToken(c *gin.Context) (string, bool) {
	scheme, token, found := strings.Cut(c.GetHeader("Authorization"), " ")
	return token, found && strings.EqualFold(scheme, "Bearer") && token != ""
}

func requestContext(c *gin.Context) context.Context {
	return c.Request.Context()
}
