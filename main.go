package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tidwall/sjson"
)

var (
	generatedImei    = GenerateImei()
	generatedImeiMD5 = fmt.Sprintf("%x", md5.Sum([]byte(generatedImei)))
	valuesMap        = map[string]any{
		"deviceInfoVersion":        2,
		"data.display":             "adb shell getprop ro.build.id",
		"data.product":             "adb shell getprop ro.build.product",
		"data.device":              "adb shell getprop ro.vendor.product.device",
		"data.board":               "adb shell getprop ro.product.board",
		"data.brand":               "adb shell getprop ro.product.brand",
		"data.model":               "adb shell getprop ro.product.model",
		"data.bootloader":          "unknown",
		"data.fingerprint":         "adb shell getprop ro.vendor.build.fingerprint",
		"data.bootId":              strings.ToUpper(uuid.New().String()),
		"data.procVersion":         "adb shell getprop cat /proc/version",
		"data.baseBand":            "adb shell getprop ro.baseband",
		"data.version.incremental": "adb shell getprop ro.product.build.version.incremental",
		"data.version.release":     "adb shell getprop ro.product.build.version.release",
		"data.version.codename":    "adb shell getprop ro.build.version.codename",
		"data.simInfo":             "T-Mobile",
		"data.osType":              "android",
		"data.macAddress":          GenerateMac().String(),
		"data.wifiBSSID":           GenerateMac().String(),
		"data.wifiSSID":            "<unknown ssid>",
		"data.imsiMd5":             generatedImeiMD5,
		"data.imei":                generatedImei,
		"data.apn":                 "wifi",
	}
)

func main() {
	result := "{}"
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

const IMEI_BASE_DIGITS_COUNT int = 14 // The number of digits without the last - the control one.

func GenerateImei() string {
	var sum int = 0
	var toAdd int = 0
	var digits [IMEI_BASE_DIGITS_COUNT + 1]int

	lolrndsrc := rand.NewSource(time.Now().UnixNano())
	lolrnd := rand.New(lolrndsrc)

	for i := 0; i < IMEI_BASE_DIGITS_COUNT; i++ {
		digits[i] = lolrnd.Intn(10)
		toAdd = digits[i]
		if (i+1)%2 == 0 {
			toAdd *= 2
			if toAdd >= 10 {
				toAdd = (toAdd % 10) + 1
			}
		}
		sum += toAdd
		fmt.Printf("%d", digits[i])
	}
	var ctrlDigit int = (sum * 9) % 10
	digits[IMEI_BASE_DIGITS_COUNT] = ctrlDigit
	return strconv.Itoa(ctrlDigit)
}
