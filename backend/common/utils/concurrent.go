package utils

func ConcurrentForEach[T any](data []T, handler func(data T), limit int) {
	total := len(data)
	c := make(chan struct{}, limit)
	sc := make(chan struct{}, total)
	defer close(c)
	defer close(sc)

	for _, v := range data {
		c <- struct{}{}
		go func(v T) {
			defer func() {
				<-c
				sc <- struct{}{}
			}()
			handler(v)
		}(v)
	}

	for i := total; i > 0; i-- {
		<-sc
	}
}
