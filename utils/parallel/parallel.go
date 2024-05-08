package parallel

import (
	"context"
	"sync"

	"github.com/nartvt/go-core/utils/slice"
	"github.com/nartvt/go-core/utils/typez"
)

type RoutineFnc func(ctx context.Context, p any) (any, error)
type parallel struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	jobsChannel chan typez.Tuple3[string, RoutineFnc, any]
	errChannel  chan error
	result      chan typez.Tuple2[string, any]
	keys        []string
}

type Parallel interface {
	AddFunc(key string, fn RoutineFnc, param any) Parallel
	Wait() (map[string]any, error)
}

func NewParallel(ctx context.Context) Parallel {
	p := &parallel{
		ctx:         ctx,
		wg:          &sync.WaitGroup{},
		jobsChannel: make(chan typez.Tuple3[string, RoutineFnc, any], 1),
		errChannel:  make(chan error, 1),
		result:      make(chan typez.Tuple2[string, any], 1),
		keys:        []string{},
	}
	go p.worker()
	return p
}

func (p *parallel) AddFunc(key string, fn RoutineFnc, param any) Parallel {
	if !slice.Contains(p.keys, key) {
		p.wg.Add(1)
		p.jobsChannel <- typez.T3(key, fn, param)
		p.keys = append(p.keys, key)
	}
	return p
}

func (p *parallel) worker() {
	for job := range p.jobsChannel {
		key, function, param := job.Unpack()
		go func(copy RoutineFnc) {
			defer p.wg.Done()
			res, errR := copy(p.ctx, param)
			if errR != nil {
				p.errChannel <- errR
				return
			}
			p.result <- typez.T2(key, res)
		}(function)
	}
}

func (p *parallel) Wait() (map[string]any, error) {
	waitChannel := make(chan bool)
	go func() {
		p.wg.Wait()
		close(waitChannel)
		p.close()
	}()
	select {
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	case err := <-p.errChannel:
		return nil, err
	default:
		return fromChan(p.result), nil
	}
}

// Close closes resources
func (p *parallel) close() {
	close(p.jobsChannel)
	close(p.result)
	close(p.errChannel)
}
func fromChan(s chan typez.Tuple2[string, any]) map[string]any {
	result := make(map[string]any)
	for item := range s {
		k, v := item.Unpack()
		result[k] = v
	}
	return result
}
