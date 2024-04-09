package httpModels

var DeleteExpire = map[string]int{
	"year":  0,
	"month": -1,
	"day":   0,
}

type EmptyStruct []interface{}

type ID struct {
	ID uint64 `json:"id"`
}
