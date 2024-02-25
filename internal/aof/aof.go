package aof

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/bobhonores/godis/internal/resp"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (a *Aof) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.file.Close()
}

func (a *Aof) Write(t resp.Token) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, err := a.file.Write(t.Marshal()); err != nil {
		return err
	}

	return nil
}

func (a *Aof) Read(fn func(t resp.Token)) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.file.Seek(0, io.SeekStart)

	reader := resp.NewReader(a.file)

	for {
		token, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		fn(token)
	}

	return nil
}
