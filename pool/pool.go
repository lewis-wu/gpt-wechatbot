package pool

import "github.com/panjf2000/ants/v2"

var pool *ants.Pool

func init() {
	var err error
	pool, err = ants.NewPool(1000)
	if err != nil {
		panic("ants pool create failed ")
	}
}
func GetPool() *ants.Pool {
	return pool
}
