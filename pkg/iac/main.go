package iac

type Manager struct {
	Data []Result
}

func NewManager() *Manager {
	return &Manager{Data: make([]Result, 0)}
}

type Collector interface {
	Collect() Result
}

type Result struct {
	From     string // name of the collector
	Machines []IACMachine
	Error    error
}

func (m *Manager) Collect(fromSources ...Collector) {

	out := make(chan Result)

	for i := range fromSources {
		go func(i int) { out <- fromSources[i].Collect() }(i)
	}

	for range fromSources {
		res := <-out
		m.Data = append(m.Data, res)
	}

}
