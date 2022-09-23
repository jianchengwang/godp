package db

type SysRole struct {
	BaseModel
	RoleName string `json:"roleName"`
}
