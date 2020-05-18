package dto

type AddForward struct {
	Network       string `json:"network"`
	ListenAddress string `json:"listen_address"`
	ListenPort    int    `json:"listen_port"`
	TargetAddress string `json:"target_address"`
	TargetPort    int    `json:"target_port"`
}

type ForwardDetail struct {
	Status int `json:"status"`
	Network       string `json:"network"`
	ListenAddress string `json:"listen_address"`
	ListenPort    int    `json:"listen_port"`
	TargetAddress string `json:"target_address"`
	TargetPort    int    `json:"target_port"`
}

type ForwardList struct {
	Total int             `json:"total"`
	List  []ForwardDetail `json:"list"`
}

type PortForwardFilters struct {
	PageNum        uint       `form:"page_num,default=1" binding:"gte=0"`
	PageSize       uint       `form:"page_size,default=20" binding:"gte=0,lte=300"`
	Status         *int       `form:"status"`
}