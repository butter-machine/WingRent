package postgres

import "github.com/lib/pq"

var DuplicateKey = pq.ErrorCode("23505")
