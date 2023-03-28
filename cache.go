package cache

import "time"

type Cache interface {
	Set(key string, value []byte) (ok bool)
	Get(key string) (value []byte)
}

type Options struct {
	timeout time.Duration
}

type Option func(options *Options)

func defaultOptions() Options {
	return Options{timeout: time.Duration(30) * time.Minute}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		if timeout > 0 {
			options.timeout = timeout
		}
	}
}
