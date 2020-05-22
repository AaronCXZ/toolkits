package write

import (
	"context"

	"github.com/olivere/elastic"
)

type ESWriter struct {
	cli   *elastic.Client
	index string
}

func (w *ESWriter) Sync() error {
	return nil
}

func (w *ESWriter) Close() error {
	return nil
}

func (w *ESWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	_, err = w.cli.Index().Index(w.index).BodyString(string(p)).Do(context.Background())
	return n, err
}

func NewESWriter(addr, index string) (*ESWriter, error) {
	client, err := elastic.NewClient(elastic.SetURL(addr))
	if err != nil {
		return nil, err
	}
	return &ESWriter{
		cli:   client,
		index: index,
	}, nil
}
