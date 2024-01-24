package gitlab

import "sync"

type PrettyResult struct {
	Zone       string
	Tenant     string
	CommitID   string
	Machines   []FlatStructMachine
	Error      string
	Aggregates AggregatedResult
}

type PrettyResults []PrettyResult

type AggregatedResult struct {
	CpuCount     int
	MemorySizeGB int
}

func (c *Collector) Collect(src ...Source) PrettyResults {

	ch := make(chan PrettyResult)

	querySrc := func(wg *sync.WaitGroup, i int) {
		defer wg.Done()

		res := c.Query(src[i])
		if res.Error != nil {
			ch <- PrettyResult{
				Tenant:   src[i].Tenant,
				CommitID: string(res.CommitID[:min(8, len(res.CommitID))]),
				Error:    res.Error.Error(),
			}
			return
		}
		machines := res.Zones.ToFlatStructMachines()
		for z, m := range machines {
			pr := PrettyResult{
				Zone:       z,
				Tenant:     src[i].Tenant,
				CommitID:   string(res.CommitID[:min(8, len(res.CommitID))]),
				Machines:   m,
				Aggregates: aggregate(m),
			}
			ch <- pr
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(src))
	for i := range src {
		go querySrc(&wg, i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make(PrettyResults, 0)
	for res := range ch {
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
