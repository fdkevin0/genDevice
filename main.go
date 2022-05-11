package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"strings"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tidwall/sjson"
)

var (
	generatedImei    = client.GenIMEI()
	generatedImeiMD5 = fmt.Sprintf("%x", md5.Sum([]byte(generatedImei)))
	valuesMap        = map[string]any{
		"display":    "adb shell getprop ro.build.id",
		"product":    "adb shell getprop ro.build.product",
		"device":     "adb shell getprop ro.vendor.product.device",
		"board":      "adb shell getprop ro.product.board",
		"brand":      "adb shell getprop ro.product.brand",
		"model":      "adb shell getprop ro.product.model",
		"bootloader": "unknown",
		// "finger_print":        "adb shell getprop ro.vendor.build.fingerprint",
		"boot_id": strings.ToUpper(uuid.New().String()),
		// "procVersion":         "adb shell getprop cat /proc/version",
		"base_band":           "adb shell getprop ro.baseband",
		"version.incremental": "adb shell getprop ro.product.build.version.incremental",
		"version.release":     "adb shell getprop ro.product.build.version.release",
		"version.codename":    "adb shell getprop ro.build.version.codename",
		"sim_info":            "T-Mobile",
		"os_type":             "android",
		"macAddress":          GenerateMac().String(),
		"wifi_bssid":          GenerateMac().String(),
		"wifi_ssid":           "<unknown ssid>",
		"imsi_md5":            generatedImeiMD5,
		"imei":                generatedImei,
		"apn":                 "wifi",
	}
)

func main() {
	client.GenRandomDevice()
	result := string(client.SystemDeviceInfo.ToJson())
	for path, value := range valuesMap {
		if strings.HasPrefix(fmt.Sprint(value), "adb") {
			cmd := exec.Command("sh", "-c", fmt.Sprint(value))
			value = strings.TrimRight(string(lo.Must(cmd.Output())), "\r\n")
		}
		result, _ = sjson.Set(result, path, value)
	}
	fmt.Println(result)
}

func GenerateMac() net.HardwareAddr {
	buf := make([]byte, 6)
	var mac net.HardwareAddr

	_, err := rand.Read(buf)
	if err != nil {
	}

	// Set the local bit
	buf[0] |= 2

	mac = append(mac, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac
}
