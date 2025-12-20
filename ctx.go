package neuron

type ctxkey string

const (
	requestStartKey   ctxkey = "neuron.request_start"
	requestIDKey      ctxkey = "neuron.request_id"
	correlationIDKey  ctxkey = "neuron.correlation_id"
	maxRetriesKey     ctxkey = "neuron.max_retries"
	retryConditionKey ctxkey = "neuron.retry_condition"
)
