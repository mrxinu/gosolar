package gosolar

import (
	"encoding/json"
	"time"
)

// CommonNode represents a basic SolarWinds node with commonly used fields
type CommonNode struct {
	NodeID       int       `json:"nodeid"`
	Caption      string    `json:"caption"`
	IPAddress    string    `json:"ipaddress"`
	Status       int       `json:"status"`
	StatusLED    string    `json:"statusled"`
	Vendor       string    `json:"vendor"`
	MachineType  string    `json:"machinetype"`
	Location     string    `json:"location"`
	Contact      string    `json:"contact"`
	Description  string    `json:"description"`
	LastBoot     time.Time `json:"lastboot"`
	ResponseTime float64   `json:"responsetime"`
}

// Interface represents a SolarWinds network interface
type Interface struct {
	InterfaceID    int     `json:"interfaceid"`
	NodeID         int     `json:"nodeid"`
	Name           string  `json:"name"`
	Caption        string  `json:"caption"`
	Status         int     `json:"status"`
	AdminStatus    int     `json:"adminstatus"`
	OperStatus     int     `json:"operstatus"`
	Speed          int64   `json:"speed"`
	MTU            int     `json:"mtu"`
	Type           int     `json:"type"`
	TypeName       string  `json:"typename"`
	InUtilization  float64 `json:"inutilization"`
	OutUtilization float64 `json:"oututilization"`
}

// CustomProperty represents a SolarWinds custom property
type CustomProperty struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
}

// Alert represents a SolarWinds alert
type Alert struct {
	AlertID     int       `json:"alertid"`
	AlertName   string    `json:"alertname"`
	Message     string    `json:"message"`
	Severity    int       `json:"severity"`
	State       int       `json:"state"`
	NodeID      int       `json:"nodeid"`
	ObjectName  string    `json:"objectname"`
	TriggerTime time.Time `json:"triggertime"`
	AckBy       string    `json:"ackby"`
	AckTime     time.Time `json:"acktime"`
}

// Volume represents a SolarWinds volume/disk
type Volume struct {
	VolumeID    int     `json:"volumeid"`
	NodeID      int     `json:"nodeid"`
	Caption     string  `json:"caption"`
	Size        float64 `json:"size"`
	Used        float64 `json:"used"`
	Available   float64 `json:"available"`
	PercentUsed float64 `json:"percentused"`
	Type        string  `json:"type"`
	Status      int     `json:"status"`
	FileSystem  string  `json:"filesystem"`
}

// Application represents a SolarWinds application monitor
type Application struct {
	ApplicationID int     `json:"applicationid"`
	NodeID        int     `json:"nodeid"`
	Name          string  `json:"name"`
	Status        int     `json:"status"`
	Availability  float64 `json:"availability"`
	ResponseTime  float64 `json:"responsetime"`
}

// QueryResult is a generic wrapper for query results with metadata
type QueryResult[T any] struct {
	Results []T        `json:"results"`
	Count   int        `json:"count"`
	Meta    *QueryMeta `json:"meta,omitempty"`
}

// QueryMeta contains metadata about query execution
type QueryMeta struct {
	ExecutionTime time.Duration `json:"execution_time"`
	RowCount      int           `json:"row_count"`
	Query         string        `json:"query,omitempty"`
}

// UnmarshalQueryResult is a helper function to unmarshal query results into strongly typed structures
func UnmarshalQueryResult[T any](data []byte) (*QueryResult[T], error) {
	var results []T
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, WrapError(err, ErrorTypeInternal, "unmarshal", "failed to unmarshal query result")
	}

	return &QueryResult[T]{
		Results: results,
		Count:   len(results),
	}, nil
}

// NodeStatus constants for common node status values
const (
	NodeStatusUnknown     = 0
	NodeStatusUp          = 1
	NodeStatusDown        = 2
	NodeStatusWarning     = 3
	NodeStatusCritical    = 14
	NodeStatusUnreachable = 12
)

// InterfaceStatus constants for interface status values
const (
	InterfaceStatusUp             = 1
	InterfaceStatusDown           = 2
	InterfaceStatusTesting        = 3
	InterfaceStatusUnknown        = 4
	InterfaceStatusDormant        = 5
	InterfaceStatusNotPresent     = 6
	InterfaceStatusLowerLayerDown = 7
)

// AlertSeverity constants
const (
	AlertSeverityInformational = 0
	AlertSeverityNotice        = 1
	AlertSeverityWarning       = 2
	AlertSeverityMinor         = 3
	AlertSeverityMajor         = 4
	AlertSeverityCritical      = 5
)
