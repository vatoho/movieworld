package entity

type Film struct {
	ID            uint64
	Name          string
	Description   string
	Duration      uint16
	MinAge        uint8
	Country       string
	ProducerName  string
	DateOfRelease string
	SumMark       uint64
	NumOfMarks    uint64
	Rating        float64
}
