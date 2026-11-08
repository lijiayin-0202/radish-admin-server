package ride

// PassengerListQuery 后台乘客列表查询条件。
type PassengerListQuery struct {
	Page           int    `json:"page" form:"page"`
	PageSize       int    `json:"pageSize" form:"pageSize"`
	Keyword        string `json:"keyword" form:"keyword"`
	RealNameStatus string `json:"realNameStatus" form:"realNameStatus"`
	AccountStatus  string `json:"accountStatus" form:"accountStatus"`
}

// PassengerItem 乘客列表行，字段与 ride-service 内部接口保持一致。
type PassengerItem struct {
	PassengerID      int64  `json:"passengerId"`
	PassengerNo      string `json:"passengerNo"`
	Mobile           string `json:"mobile"`
	Nickname         string `json:"nickname"`
	AvatarURL        string `json:"avatarUrl"`
	ExternalProvider string `json:"externalProvider"`
	RealNameStatus   int    `json:"realNameStatus"`
	AccountStatus    int    `json:"accountStatus"`
	PointBalance     int    `json:"pointBalance"`
	LastLoginAt      string `json:"lastLoginAt"`
	CreatedAt        string `json:"createdAt"`
}

// PassengerListResult 乘客分页结果。
type PassengerListResult struct {
	Items    []PassengerItem `json:"items"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}

// PassengerDetailResult 乘客详情结果。
type PassengerDetailResult struct {
	Passenger PassengerItem `json:"passenger"`
}
