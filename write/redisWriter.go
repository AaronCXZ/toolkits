package write

import "github.com/go-redis/redis"

type RedisWriter struct {
	cli     *redis.Client
	listKey string
}

func NewRedisWrite(key string, cli *redis.Client) *RedisWriter {
	return &RedisWriter{
		cli:     cli,
		listKey: key,
	}
}

func (w *RedisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.listKey, p).Result()
	return int(n), err
}

func (w *RedisWriter) Close() error {
	return w.cli.Close()
}
