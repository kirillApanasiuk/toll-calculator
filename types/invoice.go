package types

type Invoice struct {
	OBUID         int     `json:"obuId"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
}
