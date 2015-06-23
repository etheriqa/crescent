package crescent

type InstanceInput chan Input
type InstanceOutput chan Output

type InstanceInputWriter interface {
	Write(ClientID, interface{})
}

type InstanceOutputWriter interface {
	Write(interface{})
	BindClientID(ClientID) InstanceOutputWriter
}

type BoundInstanceOutput struct {
	o  InstanceOutput
	id ClientID
}

// MakeInstanceInput returns a InstanceInput
func MakeInstanceInput(n int) InstanceInput {
	return make(chan Input, n)
}

// Write sends the Input
func (o InstanceInput) Write(id ClientID, d interface{}) {
	o <- Input{
		ClientID: id,
		Input:    d,
	}
}

// MakeInstanceOutput returns a InstanceOutput
func MakeInstanceOutput(n int) InstanceOutput {
	return make(chan Output, n)
}

// Write sends the Output
func (o InstanceOutput) Write(d interface{}) {
	o <- Output{
		Output: d,
	}
}

// BindClientID returns a InstanceOutputWriter
func (o InstanceOutput) BindClientID(id ClientID) InstanceOutputWriter {
	return BoundInstanceOutput{
		o:  o,
		id: id,
	}
}

// Write sends the Output
func (bo BoundInstanceOutput) Write(d interface{}) {
	bo.o <- Output{
		ClientID: bo.id,
		Output:   d,
	}
}

// BindClientID returns a InstanceOutputWriter
func (bo BoundInstanceOutput) BindClientID(id ClientID) InstanceOutputWriter {
	return BoundInstanceOutput{
		o:  bo.o,
		id: id,
	}
}
