package vcloud

import (
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type PrettyResult struct {
	Name       string
	Endpoint   string
	Tenant     string
	Machines   []VM
	Aggregates AggregatedResult
	Error      string
}

type AggregatedResult struct {
	CPUs      int
	Memory    int
	StorageMB int
}

type PrettyResults []PrettyResult

func (c *Collector) Collect(src ...Source) PrettyResults {

	ch := make(chan PrettyResult)

	var wg sync.WaitGroup

	querySrc := func(i int) {
		defer wg.Done()

		results := make(chan PrettyResult)

		for _, ep := range src[i].Endpoints {
			go func(ep, t string) {

				res := c.Query(ep, t)
				if res.Error != nil {
					results <- PrettyResult{
						Name:     endpointName(ep),
						Tenant:   t,
						Error:    res.Error.Error(),
						Endpoint: ep,
					}
					return
				}
				pr := PrettyResult{
					Machines:   res.VirtualMachines,
					Name:       endpointName(ep),
					Endpoint:   ep,
					Tenant:     t,
					Aggregates: aggregate(res.VirtualMachines),
				}
				results <- pr
			}(ep, src[i].Tenant)
		}

		for i := 0; i < len(src[i].Endpoints); i++ {
			ch <- <-results
		}
	}

	wg.Add(len(src))
	for i := range src {
		go querySrc(i)
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
