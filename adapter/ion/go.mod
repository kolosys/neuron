module github.com/kolosys/neuron/adapter/ion

go 1.24

require (
	github.com/kolosys/ion v0.2.0
	github.com/kolosys/neuron/adapter v0.0.0-00010101000000-000000000000
)

replace github.com/kolosys/neuron => ../../

replace github.com/kolosys/neuron/adapter => ../
