package monitoring


const (
	MetricTypeMatrix = "matrix"
	MetricTypeVector = "vector"
)


type Metric struct {
	MetricName string `json:"metric_name,omitempty" description:"metric name, eg. scheduler_up_sum"`
	MetricData `json:"data,omitempty" description:"actual metric result"`
	Error      string `json:"error,omitempty"`
}

type MetricData struct {
	MetricType   string        `json:"resultType,omitempty" description:"result type, one of matrix, vector"`
	MetricValues []MetricValue `json:"result,omitempty" description:"metric data including labels, time series and values"`
}

// The first element is the timestamp, the second is the metric value.
// eg, [1585658599.195, 0.528]
type Point [2]float64

type MetricValue struct {
	Metadata map[string]string `json:"metric,omitempty" description:"time series labels"`
	// The type of Point is a float64 array with fixed length of 2.
	// So Point will always be initialized as [0, 0], rather than nil.
	// To allow empty Sample, we should declare Sample to type *Point
	Sample *Point  `json:"value,omitempty" description:"time series, values of vector type"`
	Series []Point `json:"values,omitempty" description:"time series, values of matrix type"`
}