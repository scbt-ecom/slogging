package slogging

//const (
//	defaultOrderExpiration = 1 * time.Minute
//)

//var c *cache.Cache
//
//func newOrderCache() *cache.Cache {
//	ordCache := cache.New(defaultOrderExpiration, 2*defaultOrderExpiration)
//
//	c = ordCache
//	return c
//}

const withoutRequestsOrder = 0

func (l *SLogger) SetOrder(order int) {
	//order = -1

	// TODO:
	l.order = order + 1
	//l.cache.Set(traceID, order, defaultOrderExpiration)
}

func (l *SLogger) GetOrder() int {
	if l.order == withoutRequestsOrder {
		return l.order
	}

	ord := l.order
	l.order += 1
	return ord

	//v, ok := l.cache.Get(traceID)
	//if !ok {
	//	return -1
	//}
	//
	//l.cache.Set(traceID, v.(int)+1, defaultOrderExpiration)
	//return v.(int)
}
