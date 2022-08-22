package geodata

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"io"
	"net/http"
	"os"
	"github.com/Dreamacro/clash/common/convert"
)

var initFlag bool

func InitGeoSite() error {
	if _, err := os.Stat(C.Path.GeoSite()); os.IsNotExist(err) {
		log.Infoln("Can't find GeoSite.dat, start download")
		if err := downloadGeoSite(C.Path.GeoSite()); err != nil {
			return fmt.Errorf("can't download GeoSite.dat: %s", err.Error())
		}
		log.Infoln("Download GeoSite.dat finish")
	}
	if !initFlag {
		if err := Verify(C.GeositeName); err != nil {
			log.Warnln("GeoSite.dat invalid, remove and download: %s", err)
			if err := os.Remove(C.Path.GeoSite()); err != nil {
				return fmt.Errorf("can't remove invalid GeoSite.dat: %s", err.Error())
			}
			if err := downloadGeoSite(C.Path.GeoSite()); err != nil {
				return fmt.Errorf("can't download GeoSite.dat: %s", err.Error())
			}
		}
		initFlag = true
	}
	return nil
}

func downloadGeoSite(path string) (err error) {
	resp, err := getUrl("https://ghproxy.com/https://raw.githubusercontent.com/yuumimi/rules/release/clash/geosite.dat")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)

	return err
}

func getUrl(url string) (resp *http.Response, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	convert.SetUserAgent(req.Header)
	resp, err = http.DefaultClient.Do(req)
	return
}
