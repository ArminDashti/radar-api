package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ArminDashti/radar-api/internal/models"
	"github.com/ArminDashti/radar-api/internal/rollup"
	"github.com/gin-gonic/gin"
)

type aggregate struct {
	RowID     int64
	Bucket    time.Time
	LatencyMS *float64
	OK        bool
}

func (s *Server) EndpointGrid(c *gin.Context) {
	interval, protocol, ok := gridParameters(c)
	if !ok {
		return
	}
	probe := c.DefaultQuery("probe", "all")
	if probe != "all" {
		var exists bool
		if err := s.Pool.QueryRow(requestContext(c),
			`SELECT EXISTS(SELECT 1 FROM probes WHERE code = $1)`, probe).Scan(&exists); err != nil || !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown probe"})
			return
		}
	}

	rows, err := s.endpointGridRows(c, protocol, probe, interval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not build endpoint grid"})
		return
	}
	c.JSON(http.StatusOK, makeGridResponse(interval, protocol, rows))
}

func (s *Server) ProbeGrid(c *gin.Context) {
	interval, protocol, ok := gridParameters(c)
	if !ok {
		return
	}
	rows, err := s.probeGridRows(c, protocol, interval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not build probe grid"})
		return
	}
	c.JSON(http.StatusOK, makeGridResponse(interval, protocol, rows))
}

func gridParameters(c *gin.Context) (rollup.Interval, string, bool) {
	interval, err := rollup.ParseInterval(c.DefaultQuery("interval", "minutes"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return rollup.Interval{}, "", false
	}
	protocol := c.DefaultQuery("protocol", "http")
	if protocol != "http" && protocol != "icmp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "protocol must be http or icmp"})
		return rollup.Interval{}, "", false
	}
	return interval, protocol, true
}

func (s *Server) endpointGridRows(c *gin.Context, protocol, probe string, interval rollup.Interval) ([]models.GridRow, error) {
	metadata, err := s.Pool.Query(requestContext(c), `
		SELECT e.id, e.name, COALESCE(p.flag_icon, ''), COALESCE(p.code, '')
		FROM endpoints e
		LEFT JOIN probes p ON p.id = e.probe_id
		WHERE e.active
		  AND (($1 = 'http' AND e.http_enabled) OR ($1 = 'icmp' AND e.icmp_enabled))
		  AND ($2 = 'all' OR p.code = $2 OR EXISTS (
			SELECT 1 FROM samples sx
			JOIN agents ax ON ax.id = sx.agent_id
			JOIN probes px ON px.id = ax.probe_id
			WHERE sx.endpoint_id = e.id AND px.code = $2
		  ))
		ORDER BY e.id`, protocol, probe)
	if err != nil {
		return nil, err
	}
	defer metadata.Close()
	gridRows := make([]models.GridRow, 0)
	for metadata.Next() {
		var row models.GridRow
		if err := metadata.Scan(&row.ID, &row.Name, &row.FlagIcon, &row.ProbeCode); err != nil {
			return nil, err
		}
		gridRows = append(gridRows, row)
	}

	base := `
		SELECT s.endpoint_id AS row_id, date_trunc('minute', s.observed_at) AS bucket,
		       avg(s.latency_ms) FILTER (WHERE s.ok) AS latency_ms, bool_or(s.ok) AS ok
		FROM samples s
		JOIN agents a ON a.id = s.agent_id
		JOIN probes agent_probe ON agent_probe.id = a.probe_id
		WHERE s.protocol = $1 AND s.observed_at >= $2
		  AND ($3 = 'all' OR agent_probe.code = $3)
		GROUP BY s.endpoint_id, date_trunc('minute', s.observed_at)`
	aggregates, err := s.queryAggregates(c, interval, base, protocol, interval.Start(time.Now()), probe)
	if err != nil {
		return nil, err
	}
	fillCells(gridRows, interval, aggregates)
	return gridRows, nil
}

func (s *Server) probeGridRows(c *gin.Context, protocol string, interval rollup.Interval) ([]models.GridRow, error) {
	metadata, err := s.Pool.Query(requestContext(c),
		`SELECT id, name, flag_icon, code FROM probes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer metadata.Close()
	gridRows := make([]models.GridRow, 0)
	for metadata.Next() {
		var row models.GridRow
		if err := metadata.Scan(&row.ID, &row.Name, &row.FlagIcon, &row.ProbeCode); err != nil {
			return nil, err
		}
		gridRows = append(gridRows, row)
	}
	base := `
		SELECT a.probe_id AS row_id, date_trunc('minute', s.observed_at) AS bucket,
		       avg(s.latency_ms) FILTER (WHERE s.ok) AS latency_ms, bool_or(s.ok) AS ok
		FROM samples s
		JOIN agents a ON a.id = s.agent_id
		WHERE s.protocol = $1 AND s.observed_at >= $2
		GROUP BY a.probe_id, date_trunc('minute', s.observed_at)`
	aggregates, err := s.queryAggregates(c, interval, base, protocol, interval.Start(time.Now()))
	if err != nil {
		return nil, err
	}
	fillCells(gridRows, interval, aggregates)
	return gridRows, nil
}

func (s *Server) queryAggregates(c *gin.Context, interval rollup.Interval, base string, args ...any) ([]aggregate, error) {
	query := rollupQuery(base, interval.Unit)
	rows, err := s.Pool.Query(requestContext(c), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]aggregate, 0)
	for rows.Next() {
		var value aggregate
		if err := rows.Scan(&value.RowID, &value.Bucket, &value.LatencyMS, &value.OK); err != nil {
			return nil, err
		}
		value.Bucket = value.Bucket.UTC()
		result = append(result, value)
	}
	return result, rows.Err()
}

func rollupQuery(base, unit string) string {
	query := "WITH minute_data AS (" + base + ")"
	source := "minute_data"
	for _, next := range []string{"hour", "day", "month"} {
		if unit == "minute" {
			break
		}
		name := next + "_data"
		query += fmt.Sprintf(`, %s AS (
			SELECT row_id, date_trunc('%s', bucket) AS bucket,
			       avg(latency_ms) FILTER (WHERE ok) AS latency_ms, bool_or(ok) AS ok
			FROM %s GROUP BY row_id, date_trunc('%s', bucket)
		)`, name, next, source, next)
		source = name
		if unit == next {
			break
		}
	}
	return query + " SELECT row_id, bucket, latency_ms, ok FROM " + source + " ORDER BY row_id, bucket"
}

func fillCells(rows []models.GridRow, interval rollup.Interval, values []aggregate) {
	buckets := interval.Buckets(time.Now())
	index := make(map[int64]int, len(buckets))
	for position, bucket := range buckets {
		index[bucket.Unix()] = position
	}
	rowIndex := make(map[int64]int, len(rows))
	for position := range rows {
		rows[position].Cells = make([]*models.GridCell, len(buckets))
		rowIndex[rows[position].ID] = position
	}
	for _, value := range values {
		rowPosition, rowFound := rowIndex[value.RowID]
		cellPosition, bucketFound := index[value.Bucket.Unix()]
		if rowFound && bucketFound {
			rows[rowPosition].Cells[cellPosition] = &models.GridCell{LatencyMS: value.LatencyMS, OK: value.OK}
		}
	}
}

func makeGridResponse(interval rollup.Interval, protocol string, rows []models.GridRow) models.GridResponse {
	bucketTimes := interval.Buckets(time.Now())
	buckets := make([]string, len(bucketTimes))
	for index, bucket := range bucketTimes {
		buckets[index] = bucket.Format(time.RFC3339)
	}
	return models.GridResponse{Interval: interval.Name, Protocol: protocol, Buckets: buckets, Rows: rows}
}
