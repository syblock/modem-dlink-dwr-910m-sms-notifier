package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type sms = struct {
	Id          int    `json:"messageid"`
	ClassId     int    `json:"classid"`
	PhoneNumber string `json:"phoneNumber"`
	SingleCount string `json:"singleCount"`
	Content     string `json:"smsContent"`
	Date        string `json:"smsDate"`
	Type        string `json:"smstype"`
}

type getSmsListResponse struct {
	CurPage        int   `json:"curPage"`
	Data           []sms `json:"data"`
	EndRowNum      int   `json:"endRowNum"`
	RecordsPerPage int   `json:"recordsPerpage"`
	StartRowNum    int   `json:"startRowNum"`
	TotalPage      int   `json:"totalPage"`
	TotalRecords   int   `json:"totalRecords"`
}

func main() {

	godotenv.Load()

	// get storage file path
	storageDirPath := os.Getenv("MODEM_STORAGE_DIR")
	storageFilePath := storageDirPath + "/last-update.txt"
	if storageDirPath == "" {
		log.Fatalln("MODEM_STORAGE_DIR is not set")
	}

	// create storage dir
	err := os.MkdirAll(storageDirPath, 0755)
	if err != nil {
		log.Fatalln("failed to create directory:", err)
	}

	// create storage file if not exist
	file, err := os.OpenFile(storageFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("failed to open file:", err)
	}
	defer file.Close()

	// read storage file
	storageFileContent, err := os.ReadFile(storageFilePath)
	if err != nil {
		log.Fatalln("failed to read file:", err)
	}

	// get last update time from storage file
	lastUpdateTime, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(string(storageFileContent)))
	if err != nil {
		lastUpdateTime = time.Now()
	}

	// get sms of first page
	r, err := getSmsList(1)
	if err != nil {
		log.Fatalln("Error getting sms list:", err)
	}

	// send notification for new sms
	for _, sms := range r.Data {

		// check if sms is new
		smsTime, err := time.Parse("2006-01-02 15:04:05", sms.Date)
		if err != nil {
			log.Fatalln("Error parsing time:", err)
		}
		if lastUpdateTime.After(smsTime) {
			continue
		}

		// send notification
		err = beeep.Notify(sms.PhoneNumber, sms.Content, "default")
		if err != nil {
			log.Fatalln("Error sending notification: ", err)
		}
	}

	// save new update time
	_, err = file.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatalln("failed to write to file:", err)
	}

	// delete older sms
	if r.TotalPage > 10 {

		// get sms for last page
		r, err = getSmsList(r.TotalPage)
		if err != nil {
			log.Fatalln("Error getting sms list for deleting:", err)
		}

		// delete sms of last page
		err = delSmsList(r.Data, r.CurPage)
		if err != nil {
			log.Fatalln("Error deleting sms list:", err)
		}

		log.Println("Last page sms list deleted")
	}

}

func getSmsList(page int) (*getSmsListResponse, error) {

	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"pageIndex": fmt.Sprint(page),
		}).
		Get(os.Getenv("MODEM_API_BASE_URL") + "/PageList")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status())
	}

	var r getSmsListResponse

	err = json.Unmarshal([]byte(resp.String()), &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func delSmsList(smsList []sms, currentPage int) error {

	if len(smsList) == 0 {
		return nil
	}

	client := resty.New()

	type deleEntry struct {
		Id      int `json:"id"`
		CurPage int `json:"curpage"`
	}

	deleList := make([]deleEntry, 0, len(smsList))

	for _, s := range smsList {
		deleList = append(deleList, deleEntry{
			Id:      s.Id,
			CurPage: currentPage,
		})
	}

	deleListJson, err := json.Marshal(deleList)

	if err != nil {
		return err
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"deleList": string(deleListJson),
		}).
		Get(os.Getenv("MODEM_API_BASE_URL") + "/DeleteList")

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("API request failed with status: %s", resp.Status())
	}

	return nil
}
