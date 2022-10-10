package types

import "godp/pkg/api/page"

type CloudVendorType string

const (
	Cloud_Vendor_None    CloudVendorType = "None"
	Cloud_Vendor_Aliyun  CloudVendorType = "Aliyun"
	Cloud_Vendor_Tencent CloudVendorType = "Tencent"
)

type AssetsHostGroupBy string

const (
	AssetsHostGroupBy_ClientUser  AssetsHostGroupBy = "ClientUser"
	AssetsHostGroupBy_CloudVendor AssetsHostGroupBy = "CloudVendor"
)

type AssetsHostQueryParam struct {
	page.PageInfo
	Q        string `form:"q"`
	UserId   uint   `form:"userId"`
	Verified int    `form:"verified"`
}
