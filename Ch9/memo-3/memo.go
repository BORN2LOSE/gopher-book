/* Пакет memo обеспечивает безопасное с точки зрения
* параллельности запоминание функции типа Func.
 */

package memo

import (
	"sync"
)

// Memo кеширует результаты вызовов Func.
type Memo struct {
	f     Func
	mu    sync.Mutex // Защита cache
	cache map[string]result
}

// Func является типом функции с запоминанием
type Func func(key string) (interface{}, error)
type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Get метод безопасен, с точки зрения параллельности
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()
	if !ok {
		res.value, res.err = memo.f(key)

		// Между этими двумя критическими разделами
		// несколько горутин могут вычеслить
		// f(key) и обновлять отображение.
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	return res.value, res.err
}
