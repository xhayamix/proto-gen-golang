//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
//go:generate goimports -w --local "github.com/QualiArts/campus-server" mock_$GOPACKAGE/mock_$GOFILE
package entity

type Record interface {
	PK() string
	ToKeyValue() map[string]interface{}
	GetTypeMap() map[string]string
}

type Slice interface {
	EachRecord(func(Record) bool)
	Len() int
}
