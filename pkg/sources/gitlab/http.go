package gitlab

import (
	"fmt"
)

type PrettyResult struct {
	From       string // name of the collector
	Machines   []FlatStructMachine
	Error      error
	Aggregates AggregatedResult
}

type AggregatedResult struct {
	CpuCount     int
	MemorySizeGB int
}

func (c *Collector) Collect() PrettyResult {
	res := c.Query()
	if res.Error != nil {
		return PrettyResult{Error: res.Error}
	}

	machines := res.Zones.ToFlatStructMachines()

	fromName := fmt.Sprintf("%s (%s)",
		c.Name,
		string(res.CommitID[:min(8, len(res.CommitID))]))

	return PrettyResult{
		From:       fromName,
		Machines:   machines,
		Aggregates: aggregate(machines),
		Error:      nil}
}

func aggregate(machines []FlatStructMachine) (aggr AggregatedResult) {

	for _, m := range machines {
		aggr.CpuCount += m.CpuCount
		aggr.MemorySizeGB += m.MemorySizeGB
	}

	return
}
