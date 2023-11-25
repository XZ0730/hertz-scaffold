package utils

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/XZ0730/hertz-scaffold/config"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func UploadToQiNiu(file multipart.File, fileheader *multipart.FileHeader, uid string) (int, string) {
	var AccessKey = config.QSS.AccessKey
	var SerectKey = config.QSS.SerectKey
	var Bucket = config.QSS.Bucket
	var ImgUrl = config.QSS.QiniuServer
	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SerectKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	var filebox string
	fileheader.Filename = fmt.Sprint(uid, "/", fileheader.Filename)
	key := filebox + fileheader.Filename
	err := formUploader.Put(context.Background(), &ret, upToken, key, file, fileheader.Size, &putExtra)
	if err != nil {
		code := 40003
		return code, err.Error()
	}

	url := ImgUrl + ret.Key
	return 200, url
}
