package engine

type Engine struct{}

//
// func (e *Engine) Run(ctx context.Context, opts RunOptions) <-chan model.Event {
// 	ch := make(chan model.Event, 1000)
//
// 	go func() {
// 		defer close(ch)
//
// 		ch <- model.Event{Type: model.EventEngineStart}
//
// 		for _, target := range opts.Targets {
// 			ch <- model.Event{Type: model.EventTargetStart, Data: target}
// 			e.processTarget(ctx, target, opts, ch)
// 			ch <- model.Event{Type: model.EventTargetDone, Data: target}
// 		}
//
// 		ch <- model.Event{Type: model.EventEngineStop}
// 	}()
//
// 	return ch
// }
