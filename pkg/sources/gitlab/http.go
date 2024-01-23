package gitlab

type PrettyResult struct {
	Tenant     string
	CommitID   string
	Machines   []FlatStructMachine
	Error      error
	Aggregates AggregatedResult
}

type PrettyResults []PrettyResult

type AggregatedResult struct {
	CpuCount     int
	MemorySizeGB int
}

func (c *Collector) Collect(src ...Source) PrettyResults {

	ch := make(chan PrettyResult)

	querySrc := func(i int) {
		res := c.Query(src[i])
		machines := res.Zones.ToFlatStructMachines()
		pr := PrettyResult{
			Tenant:     src[i].Tenant,
			CommitID:   string(res.CommitID[:min(8, len(res.CommitID))]),
			Machines:   machines,
			Aggregates: aggregate(machines),
			Error:      res.Error,
		}
		ch <- pr
	}

	for i := range src {
		go querySrc(i)
	}

	results := make(PrettyResults, 0)
	for i := 0; i < len(src); i++ {
		res := <-ch
		results = append(results, res)
	}

	return results
}

func aggregate(machines []FlatStructMachine) (aggr AggregatedResult) {

	for _, m := range machines {
		aggr.CpuCount += m.CpuCount
		aggr.MemorySizeGB += m.MemorySizeGB
	}

	return
}
