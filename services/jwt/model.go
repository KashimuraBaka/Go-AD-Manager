package jwt

type PayLoad struct {
	SUB string
	NBF int64
	IP  string
	ISS string
	EXP int64
	IAT int64
	JTI string
}
