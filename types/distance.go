package types

type Distance struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}
