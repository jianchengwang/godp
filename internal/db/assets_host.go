package db

type AssetsHost struct {
	BaseModel
	IP           string `json:"ip"`
	Port         uint16 `json:"port"`
	IntranetIP   string `json:"intranetIp"`
	User         string `json:"user"`
	Password     string `json:"password"`
	WorkDir      string `json:"workDir"`
	UserId       uint   `json:"userId"`
	ClientUserId uint   `json:"clientUserId"`
	Group        string `json:"Group"`
}
