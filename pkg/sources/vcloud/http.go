package vcloud

import (
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type PrettyResult struct {
	Endpoint   string
	Tenant     string
	Machines   []VM
	Aggregates AggregatedResult
	Error      error
}

type AggregatedResult struct {
	CPUs      int
	Memory    int
	StorageMB int
}

type PrettyResults []PrettyResult

func (c *Collector) Collect() PrettyResults {

	results := make(PrettyResults, 0)

	ch := make(chan PrettyResult)

	go func() {
		defer close(ch)
		var wg sync.WaitGroup
		for _, ep := range c.Endpoints {
			wg.Add(len(c.Tenants))
			for _, t := range c.Tenants {
				go func(ep, t string) {
					defer wg.Done()
					res := c.Query(ep, t)
					pr := PrettyResult{
						Machines:   res.VirtualMachines,
						Error:      res.Error,
						Endpoint:   endpointName(ep),
						Tenant:     t,
						Aggregates: aggregate(res.VirtualMachines),
					}
					ch <- pr
				}(ep, t)
			}
		}
		wg.Wait()
	}()

	for res := range ch {
		results = append(results, res)
	}

	return results
}

func endpointName(ep string) string {
	u, err := url.Parse(ep)
	if err != nil {
		return ep
	}
	return strings.ToUpper(strings.Split(u.Host, ".")[0])
}

func aggregate(machines []VM) (aggr AggregatedResult) {

	for _, m := range machines {
		aggr.CPUs += m.CPUs
		aggr.Memory += m.Memory

		s, _ := strconv.Atoi(m.Storage)
		aggr.StorageMB += s
	}

	return
}
