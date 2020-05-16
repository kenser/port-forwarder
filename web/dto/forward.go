package dto

type AddForward struct {
	Network       string
	ListenAddress string
	ListenPort    int
	TargetAddress string
	TargetPort    int
}
