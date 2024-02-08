package vcloud

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

func Collect(ctx context.Context, src ...Source) Results {

	ch := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(len(src))
	for i := range src {
		go querySrc(ctx, &wg, ch, src[i])
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make(Results, 0)

	for {
		select {
		case <-ctx.Done():
			return results
		case res, ok := <-ch:
			if !ok {
				return results
			}
			results = append(results, res)
		}
	}

}

func querySrc(ctx context.Context, wg *sync.WaitGroup, out chan<- Result, src Source) {
	defer wg.Done()

	res := src.Query()
	pres := toResult(res)

	select {
	case out <- pres:
	case <-ctx.Done():
		log.Println(ctx.Err())
	}
}

func toResult(res Response) Result {
	if res.Error != nil {
		return Result{
			Name:     endpointName(res.Endpoint),
			Tenant:   res.Tenant,
			Error:    res.Error.Error(),
			Endpoint: res.Endpoint,
		}
	}

	return Result{
		Machines:   res.VirtualMachines,
		Name:       endpointName(res.Endpoint),
		Endpoint:   res.Endpoint,
		Tenant:     res.Tenant,
		Aggregates: aggregate(res.VirtualMachines),
	}

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
