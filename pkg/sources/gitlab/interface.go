package gitlab

import (
	"errors"
	"io"
	"sync"
)

func Collect(src ...Source) Results {

	ch := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(len(src))
	for i := range src {
		go querySrc(&wg, ch, src[i])
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make(Results, 0)
	for res := range ch {
		results = append(results, res)
	}

	return results
}

func querySrc(wg *sync.WaitGroup, out chan<- Result, src Source) {
	defer wg.Done()

	if src.Tenant == "" || src.Token == "" || src.BaseURL == "" {
		out <- Result{
			Tenant: src.Tenant,
			Error:  "Missing or wrong configuration. Please check the config file (default: config.json)",
		}
		return
	}

	res := src.Query()
	if res.Error != nil {
		err := res.Error.Error()
		if errors.Is(res.Error, io.EOF) {
			err = "Target file is empty"
		}
		out <- Result{
			Tenant:   src.Tenant,
			CommitID: string(res.CommitID[:min(8, len(res.CommitID))]),
			Error:    err,
		}
		return
	}
	machines := res.Zones.ToFlatStructMachines()
	for z, m := range machines {
		pr := Result{
			Zone:       z,
			Tenant:     src.Tenant,
			CommitID:   string(res.CommitID[:min(8, len(res.CommitID))]),
			Machines:   m,
			Aggregates: aggregate(m),
		}
		out <- pr
	}
}

func aggregate(machines []FlatStructMachine) (aggr AggregatedResult) {

	for _, m := range machines {
		aggr.CpuCount += m.CpuCount
		aggr.MemorySizeGB += m.MemorySizeGB
	}

	return
}
