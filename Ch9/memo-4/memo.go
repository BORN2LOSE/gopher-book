package memo

import (
	"sync"
)

type Func func(string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// Это первый запрос данного ключа
		// Эта горутина становится ответственной за
		// вычесление значения и опевещения о готовности.
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)
		close(e.ready) // Широковещательное оповещение о готовности
	} else {
		// Повторный запрос данного ключа
		memo.mu.Unlock()

		<-e.ready // Ожидание готовности
	}
	return e.res.value, e.res.err
}
