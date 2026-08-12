package models

import "time"

type Probe struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	FlagIcon  string    `json:"flag_icon"`
	CreatedAt time.Time `json:"created_at"`
}

type Endpoint struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Host        string    `json:"host"`
	HTTPEnabled bool      `json:"http_enabled"`
	ICMPEnabled bool      `json:"icmp_enabled"`
	ProbeID     *int64    `json:"probe_id"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

type Target struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Host      string   `json:"host"`
	Protocols []string `json:"protocols"`
}

type SampleInput struct {
	EndpointID int64     `json:"endpoint_id" binding:"required"`
	Protocol   string    `json:"protocol" binding:"required"`
	ObservedAt time.Time `json:"observed_at" binding:"required"`
	LatencyMS  *float64  `json:"latency_ms"`
	OK         bool      `json:"ok"`
}

type GridCell struct {
	LatencyMS *float64 `json:"latency_ms"`
	OK        bool     `json:"ok"`
}

type GridRow struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	FlagIcon  string      `json:"flag_icon"`
	ProbeCode string      `json:"probe_code"`
	Cells     []*GridCell `json:"cells"`
}

type GridResponse struct {
	Interval string    `json:"interval"`
	Protocol string    `json:"protocol"`
	Buckets  []string  `json:"buckets"`
	Rows     []GridRow `json:"rows"`
}
