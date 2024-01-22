package vcloud

import (
	"fmt"
	"sync"
	"time"
)

type PrettyResult struct {
	From     string
	Machines []VM
	Error    error
}

func (c *Collector) Collect() PrettyResult {

	vms := make([]VM, 0)
	var err error

	ch := make(chan Result)

	go func() {
		defer close(ch)
		var wg sync.WaitGroup
		for _, ep := range c.Endpoints {
			wg.Add(len(c.Tenants))
			for _, t := range c.Tenants {
				go func(ep, t string) {
					defer wg.Done()
					res := c.Query(ep, t)
					ch <- res
				}(ep, t)
			}
		}
		wg.Wait()
	}()

	for res := range ch {
		if res.Error != nil {
			err = res.Error
			break
		}
		vms = append(vms, res.VirtualMachines...)
	}

	return PrettyResult{
		Machines: vms,
		From:     fmt.Sprintf("VMWare Cloud Directory (%s)", time.Now().Format(time.RFC3339)),
		Error:    err}
}
