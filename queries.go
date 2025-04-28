package godruid

import (
	"bytes"
	"encoding/json"
)

// Check http://druid.io/docs/0.6.154/Querying.html#query-operators for detail description.

// The Query interface stands for any kinds of druid query.
type Query interface {
	setup()
	onResponse(content []byte) error
}

// ---------------------------------
// GroupBy Query
// ---------------------------------

type QueryGroupBy struct {
	QueryType        string                 `json:"queryType"`
	DataSource       string                 `json:"dataSource"`
	Dimensions       []DimSpec              `json:"dimensions"`
	Granularity      Granlarity             `json:"granularity"`
	LimitSpec        *Limit                 `json:"limitSpec,omitempty"`
	Having           *Having                `json:"having,omitempty"`
	Filter           *Filter                `json:"filter,omitempty"`
	Aggregations     []Aggregation          `json:"aggregations"`
	PostAggregations []PostAggregation      `json:"postAggregations,omitempty"`
	Intervals        []string               `json:"intervals"`
	Context          map[string]interface{} `json:"context,omitempty"`

	QueryResult []GroupByItem `json:"-"`
}

type GroupByItem struct {
	Version   string                 `json:"version"`
	Timestamp string                 `json:"timestamp"`
	Event     map[string]interface{} `json:"event"`
}

func (q *QueryGroupBy) setup() { q.QueryType = "groupBy" }
func (q *QueryGroupBy) onResponse(content []byte) error {
	res := new([]GroupByItem)
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// QueryGroupByGeneric is the model for groupBy query but with a generic response type.
type QueryGroupByGeneric[T any] struct {
	*QueryGroupBy

	QueryResult []GroupByItemGeneric[T] `json:"-"`
}

// GroupByItemGeneric is the response to groupBy query but with a generic type.
type GroupByItemGeneric[T any] struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Event     T      `json:"event"`
}

// Results flattens the groupBy results as 1 large array.
func (q *QueryGroupByGeneric[T]) Results() []T {
	result := make([]T, 0, len(q.QueryResult))
	for _, blob := range q.QueryResult {
		result = append(result, blob.Event)
	}
	return result
}

// ---------------------------------
// Search Query
// ---------------------------------

type QuerySearch struct {
	QueryType        string                 `json:"queryType"`
	DataSource       string                 `json:"dataSource"`
	Granularity      Granlarity             `json:"granularity"`
	Filter           *Filter                `json:"filter,omitempty"`
	Intervals        []string               `json:"intervals"`
	SearchDimensions []string               `json:"searchDimensions,omitempty"`
	Query            *SearchQuery           `json:"query"`
	Sort             *SearchSort            `json:"sort"`
	Context          map[string]interface{} `json:"context,omitempty"`

	QueryResult []SearchItem `json:"-"`
}

type SearchItem struct {
	Timestamp string     `json:"timestamp"`
	Result    []DimValue `json:"result"`
}

type DimValue struct {
	Dimension string `json:"dimension"`
	Value     string `json:"value"`
}

func (q *QuerySearch) setup() { q.QueryType = "search" }
func (q *QuerySearch) onResponse(content []byte) error {
	res := new([]SearchItem)
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// ---------------------------------
// SegmentMetadata Query
// ---------------------------------

type QuerySegmentMetadata struct {
	QueryType  string                 `json:"queryType"`
	DataSource string                 `json:"dataSource"`
	Intervals  []string               `json:"intervals"`
	ToInclude  *ToInclude             `json:"toInclude,omitempty"`
	Merge      interface{}            `json:"merge,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`

	QueryResult []SegmentMetaData `json:"-"`
}

type SegmentMetaData struct {
	Id        string                `json:"id"`
	Intervals []string              `json:"intervals"`
	Columns   map[string]ColumnItem `json:"columns"`
}

type ColumnItem struct {
	Type        string      `json:"type"`
	Size        int         `json:"size"`
	Cardinality interface{} `json:"cardinality"`
}

func (q *QuerySegmentMetadata) setup() { q.QueryType = "segmentMetadata" }
func (q *QuerySegmentMetadata) onResponse(content []byte) error {
	res := new([]SegmentMetaData)
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// ---------------------------------
// TimeBoundary Query
// ---------------------------------

type QueryTimeBoundary struct {
	QueryType  string                 `json:"queryType"`
	DataSource string                 `json:"dataSource"`
	Bound      string                 `json:"bound,omitempty"`
	Filter     *Filter                `json:"filter,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`

	QueryResult []TimeBoundaryItem `json:"-"`
}

type TimeBoundaryItem struct {
	Timestamp string       `json:"timestamp"`
	Result    TimeBoundary `json:"result"`
}

type TimeBoundary struct {
	MinTime string `json:"minTime"`
	MaxTime string `json:"maxTime"`
}

func (q *QueryTimeBoundary) setup() { q.QueryType = "timeBoundary" }
func (q *QueryTimeBoundary) onResponse(content []byte) error {
	res := new([]TimeBoundaryItem)
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// ---------------------------------
// Timeseries Query
// ---------------------------------

type QueryTimeseries struct {
	QueryType        string                 `json:"queryType"`
	DataSource       string                 `json:"dataSource"`
	Granularity      Granlarity             `json:"granularity"`
	Filter           *Filter                `json:"filter,omitempty"`
	Aggregations     []Aggregation          `json:"aggregations"`
	PostAggregations []PostAggregation      `json:"postAggregations,omitempty"`
	Intervals        []string               `json:"intervals"`
	Context          map[string]interface{} `json:"context,omitempty"`

	QueryResult []Timeseries `json:"-"`
}

type Timeseries struct {
	Timestamp string                 `json:"timestamp"`
	Result    map[string]interface{} `json:"result"`
}

func (q *QueryTimeseries) setup() { q.QueryType = "timeseries" }
func (q *QueryTimeseries) onResponse(content []byte) error {
	res := new([]Timeseries)
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// QueryTimeseriesGeneric is the model for timeseries query but with a generic response type.
type QueryTimeseriesGeneric[T any] struct {
	*QueryTimeseries

	QueryResult []TimeseriesGeneric[T] `json:"-"`
}

// TimeseriesGeneric is the response to timeseries query but with a generic type.
type TimeseriesGeneric[T any] struct {
	Timestamp string `json:"timestamp"`
	Result    T      `json:"result"`
}

// Results flattens the timeseries results as 1 large array.
func (q *QueryTimeseriesGeneric[T]) Results() []T {
	result := make([]T, 0, len(q.QueryResult))
	for _, blob := range q.QueryResult {
		result = append(result, blob.Result)
	}
	return result
}

// ---------------------------------
// TopN Query
// ---------------------------------

type QueryTopN struct {
	QueryType        string                 `json:"queryType"`
	DataSource       string                 `json:"dataSource"`
	Granularity      Granlarity             `json:"granularity"`
	Dimension        DimSpec                `json:"dimension"`
	Threshold        int                    `json:"threshold"`
	Metric           *TopNMetric            `json:"metric"`
	Filter           *Filter                `json:"filter,omitempty"`
	Aggregations     []Aggregation          `json:"aggregations"`
	PostAggregations []PostAggregation      `json:"postAggregations,omitempty"`
	Intervals        []string               `json:"intervals"`
	Context          map[string]interface{} `json:"context,omitempty"`

	QueryResult []TopNItem `json:"-"`
}

type TopNItem struct {
	Timestamp string                   `json:"timestamp"`
	Result    []map[string]interface{} `json:"result"`
}

func (q *QueryTopN) setup() { q.QueryType = "topN" }
func (q *QueryTopN) onResponse(content []byte) error {
	res := new([]TopNItem)
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// QueryTopNGeneric is the model for topn query but with a generic response type.
type QueryTopNGeneric[T any] struct {
	*QueryTopN

	QueryResult []TopNItemGeneric[T] `json:"-"`
}

// TopNItemGeneric is the response to topn query but with a generic type.
type TopNItemGeneric[T any] struct {
	Timestamp string `json:"timestamp"`
	Result    []T    `json:"result"`
}

func (q *QueryTopNGeneric[T]) onResponse(content []byte) error {
	res := new([]TopNItemGeneric[T])
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// Results flattens the results as 1 large array.
func (q *QueryTopNGeneric[T]) Results() []T {
	count := 0
	for _, blob := range q.QueryResult {
		count += len(blob.Result)
	}
	result := make([]T, 0, count)

	for _, blob := range q.QueryResult {
		result = append(result, blob.Result...)
	}
	return result
}

// ---------------------------------
// Scan Query
// ---------------------------------

// QueryScan is the model for scan query
type QueryScan struct {
	QueryType    string                 `json:"queryType"`
	DataSource   string                 `json:"dataSource"`
	Columns      []string               `json:"columns"`
	Intervals    []string               `json:"intervals"`
	BatchSize    int64                  `json:"batchSize"`
	Limit        int64                  `json:"limit"`
	Offset       int64                  `json:"offset"`
	Order        string                 `json:"order,omitempty"`
	Filter       *Filter                `json:"filter,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	ResultFormat string                 `json:"resultFormat,omitempty"`
	Legacy       bool                   `json:"legacy,omitempty"`

	QueryResult []ScanBlob `json:"-"`
}

// ScanBlob is the response to scan query
type ScanBlob struct {
	Columns   []string                 `json:"columns"`
	SegmentID string                   `json:"segmentId"`
	Events    []map[string]interface{} `json:"events"`
}

func (q *QueryScan) setup() { q.QueryType = "scan" }
func (q *QueryScan) onResponse(content []byte) error {
	res := new([]ScanBlob)
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// QueryScanGeneric is the model for scan query but with a generic response type.
type QueryScanGeneric[T any] struct {
	*QueryScan

	QueryResult []ScanBlobGeneric[T] `json:"-"`
}

// ScanBlobGeneric is the response to scan query but with a generic type.
type ScanBlobGeneric[T any] struct {
	Columns   []string `json:"columns"`
	SegmentID string   `json:"segmentId"`
	Events    []T      `json:"events"`
}

func (q *QueryScanGeneric[T]) onResponse(content []byte) error {
	res := new([]ScanBlobGeneric[T])
	d := json.NewDecoder(bytes.NewReader(content))
	d.UseNumber()
	if err := d.Decode(res); err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// Results flattens the scan blobs results as 1 large array.
func (q *QueryScanGeneric[T]) Results() []T {
	count := 0
	for _, blob := range q.QueryResult {
		count += len(blob.Events)
	}
	result := make([]T, 0, count)

	for _, blob := range q.QueryResult {
		result = append(result, blob.Events...)
	}
	return result
}
