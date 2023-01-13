package tracer

import "go.opentelemetry.io/otel/attribute"

type tracerAttribute struct {
	Key   string
	Value string
}

func SetAttribute(key string, value string) tracerAttribute {
	return tracerAttribute{
		Key:   key,
		Value: value,
	}
}

func (a tracerAttribute) getAttribute() attribute.KeyValue {
	return attribute.String(a.Key, a.Value)
}
