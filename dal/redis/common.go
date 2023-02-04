package redis

import (
	"github.com/redis/go-redis/v9"
)

type ScanCmder interface {
	Scan(dst interface{}) error
}

// MustScan enhances the Scan method of a ScanCmder with these features:
//   - it returns the error redis.Nil when the key does not exist. See https://github.com/go-redis/redis/issues/1668
//   - it supports embedded struct better. See https://github.com/go-redis/redis/issues/2005#issuecomment-1019667052
func MustScan(s ScanCmder, dest ...interface{}) error {
	switch cmd := s.(type) {
	case *redis.MapStringStringCmd:
		if len(cmd.Val()) == 0 {
			return redis.Nil
		}
	case *redis.SliceCmd:
		keyExists := false
		for _, v := range cmd.Val() {
			if v != nil {
				keyExists = true
				break
			}
		}
		if !keyExists {
			return redis.Nil
		}
	}

	for _, d := range dest {
		if err := s.Scan(d); err != nil {
			return err
		}
	}

	return nil
}
