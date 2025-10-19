package gocq

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/LagrangeDev/LagrangeGo/client/auth"
	"github.com/Mrs4s/go-cqhttp/global"
	"github.com/Mrs4s/go-cqhttp/internal/base"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

var DefaultAppInfo = &auth.AppInfo{
	OS:       "Linux",
	Kernel:   "Linux",
	VendorOS: "linux",

	CurrentVersion:   "3.2.19-39038",
	BuildVersion:     39038,
	MiscBitmap:       32764,
	PTVersion:        "2.0.0",
	PTOSVersion:      19,
	PackageName:      "com.tencent.qq",
	WTLoginSDK:       "nt.wtlogin.0.0.1",
	PackageSign:      "V1_LNX_NQ_3.2.19_39038_GW_B",
	AppID:            1600001615,
	SubAppID:         537313942,
	AppIDQrcode:      13697054,
	AppClientVersion: 39038,

	MainSigmap:  169742560,
	SubSigmap:   0,
	NTLoginType: 1,
}

/*
{
    "Os": "Linux",
    "Kernel": "Linux",
    "VendorOs": "linux",
    "CurrentVersion": "3.2.19-39038",
    "MiscBitmap": 32764,
    "PtVersion": "2.0.0",
    "SsoVersion": 19,
    "PackageName": "com.tencent.qq",
    "WtLoginSdk": "nt.wtlogin.0.0.1",
    "AppId": 1600001615,
    "SubAppId": 537313942,
    "AppIdQrCode": 537313942,
    "AppClientVersion": 39038,
    "MainSigMap": 169742560,
    "SubSigMap": 0,
    "NTLoginType": 1
}
*/

type AppInfoResp struct {
	OS               string `json:"Os" validate:"required"`
	Kernel           string `json:"Kernel" validate:"required"`
	VendorOS         string `json:"VendorOs" validate:"required"`
	CurrentVersion   string `json:"CurrentVersion" validate:"required"`
	MiscBitmap       int    `json:"MiscBitmap" validate:"required"`
	PTVersion        string `json:"PtVersion" validate:"required"`
	SsoVersion       int    `json:"SsoVersion" validate:"required"`
	PackageName      string `json:"PackageName" validate:"required"`
	WtLoginSdk       string `json:"WtLoginSdk" validate:"required"`
	AppID            int    `json:"AppId" validate:"required"`
	SubAppID         int    `json:"SubAppId" validate:"required"`
	AppIDQrCode      int    `json:"AppIdQrCode" validate:"required"`
	AppClientVersion int    `json:"AppClientVersion" validate:"required"`
	MainSigMap       int    `json:"MainSigMap" validate:"required"`
	SubSigMap        int    `json:"SubSigMap"`
	NTLoginType      int    `json:"NTLoginType" validate:"required"`
}

func fetchAppInfo(url string) (*auth.AppInfo, error) {

	data, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	resp := &AppInfoResp{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(resp)
	if err != nil {
		return nil, err
	}

	return &auth.AppInfo{
		OS:               resp.OS,
		Kernel:           resp.Kernel,
		VendorOS:         resp.VendorOS,
		CurrentVersion:   resp.CurrentVersion,
		BuildVersion:     resp.AppClientVersion,
		MiscBitmap:       resp.MiscBitmap,
		PTVersion:        resp.PTVersion,
		PTOSVersion:      resp.SsoVersion,
		PackageName:      resp.PackageName,
		WTLoginSDK:       resp.WtLoginSdk,
		PackageSign:      resp.CurrentVersion,
		AppID:            resp.AppID,
		SubAppID:         resp.SubAppID,
		AppIDQrcode:      resp.AppIDQrCode,
		AppClientVersion: resp.AppClientVersion,
		MainSigmap:       resp.MainSigMap,
		SubSigmap:        resp.SubSigMap,
		NTLoginType:      resp.NTLoginType,
	}, nil
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func loadAppInfoFromLocalFile() *auth.AppInfo {

	// 加载本地版本信息, 一般是在上次登录时保存的
	versionFile := path.Join(global.VersionsPath, "7.json")
	if global.FileExists(versionFile) {
		b, err := os.ReadFile(versionFile)
		if err != nil {
			log.Warnf("从文件 %s 读取本地版本信息文件出错: %v", versionFile, err)
			return nil
		}
		info, err := auth.UnmarshalAppInfo(b)
		if err != nil {
			log.Warnf("从文件 %s 解析本地版本信息出错: %v", versionFile, err)
			return nil
		}
		log.Infof("从文件 %s 读取协议版本 %s.", versionFile, cli.Version().CurrentVersion)
		return info
	} else {
		return nil
	}
}

func LoadAppInfo() *auth.AppInfo {

	log.Info("开始获取协议版本信息……")

	info := loadAppInfoFromLocalFile()

	if info != nil {
		return info
	}

	if base.Account.AppInfoUrl != "" {
		info, err := fetchAppInfo(base.Account.AppInfoUrl)

		if err != nil {
			log.Warnf("从远程获取版本信息失败: %v, 将使用内置的默认版本信息。", err)
		} else {
			log.Infof("从远程获取版本信息成功，版本：%s %s", info.OS, info.CurrentVersion)
			return info
		}
	}

	log.Infof("使用默认内置版本：%s %s", DefaultAppInfo.OS, DefaultAppInfo.CurrentVersion)
	return DefaultAppInfo
}
