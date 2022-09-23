package pojo

import "errors"

type IpAddressStruct struct {
	IP         string `json:"ip"`
	Port       uint16 `json:"port"`
	IntranetIP string `json:"intranet_ip"`
	User       string `json:"user"`
	Password   string `json:"password"`
	DeployDir  string `json:"deploy_dir"`
}

type ProjectConfig struct {
	ProjectName string `json:"project_name"`
	ProjectApp  string `json:"project_app"`

	CI           IpAddressStruct   `json:"ci"`
	IPAddressArr []IpAddressStruct `json:"ip_address_arr"`

	DomainSchema    string `json:"domain_schema"`
	DomainAdmin     string `json:"domain_admin"`
	DomainClient    string `json:"domain_client"`
	DomainH5        string `json:"domain_h5"`
	DomainAdminSsl  bool   `json:"domain_admin_ssl"`
	DomainClientSsl bool   `json:"domain_client_ssl"`
	DomainH5Ssl     bool   `json:"domain_h5_ssl"`

	ServerAdminPort  int `json:"server_admin_port"`
	ServerClientPort int `json:"server_client_port"`

	DockerImageAdmin         string `json:"docker_image_admin"`
	DockerImageBackendClient string `json:"docker_image_backend_client"`
	DockerImageH5Build bool `json:"docker_image_h5_build"`
	DockerRegistry string `json:"docker_registry"`

	GitBackendRep     string `json:"git_backend_rep"`
	GitBackendBranch  string `json:"git_backend_branch"`
	GitFrontendRep    string `json:"git_frontend_rep"`
	GitFrontendBranch string `json:"git_frontend_branch"`
	GitH5Rep          string `json:"git_h5_rep"`
	GitH5Branch       string `json:"git_h5_branch"`

	MysqlDatabase string `json:"mysql_database"`
	MysqlIp       string `json:"mysql_ip"`
	MysqlPort     int    `json:"mysql_port"`
	MysqlUser     string `json:"mysql_user"`
	MysqlPassword string `json:"mysql_password"`

	RedisIp       string `json:"redis_ip"`
	RedisPort     int    `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
	RedisDatabase int    `json:"redis_database"`

	SmsCloudType             string `json:"sms_cloudType"`
	SmsAliyunAccessKeyID     string `json:"sms_aliyun_accessKeyId"`
	SmsAliyunAccessKeySecret string `json:"sms_aliyun_accessKeySecret"`
	SmsAliyunSignName        string `json:"sms_aliyun_signName"`
	SmsAliyunTemplateID      string `json:"sms_aliyun_templateId"`
	SmsTencentSecretID       string `json:"sms_tencent_secretId"`
	SmsTencentSecretKey      string `json:"sms_tencent_secretKey"`
	SmsTencentSdkAppID       string `json:"sms_tencent_sdkAppID"`
	SmsTencentSignName       string `json:"sms_tencent_signName"`
	SmsTencentTemplateID     string `json:"sms_tencent_templateId"`

	StorageStoreType          string `json:"storage_storeType"`
	StorageDownloadUrl        string `json:"storage_downloadUrl"`
	StorageOssEndpoint        string `json:"storage_oss_endpoint"`
	StorageOssAccessKeyID     string `json:"storage_oss_accessKeyId"`
	StorageOssAccessKeySecret string `json:"storage_oss_accessKeySecret"`
	StorageOssBucketName      string `json:"storage_oss_bucketName"`
	StorageCosRegion          string `json:"storage_cos_region"`
	StorageCosSecretID        string `json:"storage_cos_secretId"`
	StorageCosSecretKey       string `json:"storage_cos_secretKey"`
	StorageCosBucketName      string `json:"storage_cos_bucketName"`

	WxAppID        string `json:"wx_appId"`
	WxSecret       string `json:"wx_secret"`
	WxMpTxt        string `json:"wx_mp_txt"`
	WxMpQrcode     bool   `json:"wx_mp_qrcode"`
	WxPayEnable    bool   `json:"wx_pay_enable"`
	WxPayMchID     string `json:"wx_pay_mchId"`
	WxPayMchKey    string `json:"wx_pay_mchKey"`
	WxPayKeyString string `json:"wx_pay_keyString"`

	AlipayEnable             bool   `json:"alipay_enable"`
	AlipayAppID              string `json:"alipay_appId"`
	AlipayAlipayPublicKey    string `json:"alipay_alipayPublicKey"`
	AlipayMerchantPrivateKey string `json:"alipay_merchantPrivateKey"`

	SandpayQuickPayEnable     bool   `json:"sandpay_quickPayEnable"`
	SandpayPrivateKeyPassword string `json:"sandpay_privateKeyPassword"`
	SandpayMerchMid           string `json:"sandpay_merchMid"`
	SandpayCert               bool   `json:"sandpay_cert"`

	FaceidEnable    bool   `json:"faceid_enable"`
	FaceidSecretId  string `json:"faceid_secretId"`
	FaceidSecretKey string `json:"faceid_secretKey"`

	Web3JProd       bool   `json:"web3j_prod"`
	Web3JOwner      string `json:"web3j_owner"`
	Web3JPrivateKey string `json:"web3j_privateKey"`

	CollectDonateFirstInterval  uint16 `json:"collect_donate_first_interval"`
	CollectDonateSecondInterval uint16 `json:"collect_donate_second_interval"`

	OpenRealPassValidate  bool `json:"open_real_pass_validate"`
	RealPassValidateLimit int  `json:"real_pass_validate_limit"`

	OwnGetFeature bool `json:"own_get_feature"`

	RpFrontendZip bool `json:"rp_frontend_zip"`
	RpH5Zip       bool `json:"rp_h5_zip"`

	FrontendLogo              bool `json:"frontend_logo"`
	FrontendUploadImageRemark struct {
		BannerImage            string `json:"BannerImage"`
		PublisherLogo          string `json:"PublisherLogo"`
		UiGoodsComingSoonImage string `json:"UiGoodsComingSoonImage"`
		UiAdImage              string `json:"UiAdImage"`
		IpCoverImage           string `json:"IpCoverImage"`
		IpGoodsCoverImage      string `json:"IpGoodsCoverImage"`
		IpGoodsShowImage       string `json:"IpGoodsShowImage"`
		IpGoodsOriginImage     string `json:"IpGoodsOriginImage"`
		IpGoodsSquareImage     string `json:"IpGoodsSquareImage"`
		IpGoodsContentImage    string `json:"IpGoodsContentImage"`
	} `json:"frontend_uploadImageRemark"`
	FrontendPermission struct {
		ConfigUi                     bool `json:"ConfigUi"`
		ConfigUiGoodsComingSoonImage bool `json:"ConfigUi_GoodsComingSoonImage"`
		ConfigUiAd                   bool `json:"ConfigUi_Ad"`
		ConfigUiVerification         bool `json:"ConfigUi_Verification"`
		ConfigUiRealPassValidate     bool `json:"ConfigUi_Real_Pass_Validate"`
		ConfigUiOwnGetFeature        bool `json:"ConfigUi_Own_Get_Feature"`
		ConfigUiCollectDonate        bool `json:"ConfigUi_Collect_Donate"`
		ConfigUiApp                  bool `json:"ConfigUi_App"`
		ConfigUiChannel              bool `json:"ConfigUi_Channel"`
		IpCoverImage                 bool `json:"IpCoverImage"`
		IpBlindbox                   bool `json:"IpBlindbox"`
		IpGoodsCanComposite          bool `json:"IpGoods_CanComposite"`
		IpGoodsBuyFirst              bool `json:"IpGoods_BuyFirst"`
		IpGoodsMedia                 bool `json:"IpGoods_Media"`
		IpGoodsSku                   bool `json:"IpGoods_Sku"`
		Channel                      bool `json:"Channel"`
		Fin                          bool `json:"Fin"`
		Report                       bool `json:"Report"`
	} `json:"frontend_Permission"`
}

func FindIpAddressByIp(ipAddressArr []IpAddressStruct, ip string) (error, IpAddressStruct) {
	for _, eachItem := range ipAddressArr {
		if eachItem.IP == ip {
			return nil, eachItem
		}
	}
	return errors.New("ipAddress not existed"), IpAddressStruct{}
}
