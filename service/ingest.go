package service

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/g33kzone/go-stox-ingest/db"
	models "github.com/g33kzone/go-stox-ingest/models/entity"
	"github.com/getsentry/sentry-go"
	"gorm.io/datatypes"
)

const (
	tempDir     = "tmp"
	csvFileType = ".CSV"
)

type Ingest interface {
	InitiateDownload(date time.Time) error
	// DownloadBsEBhavCopy()
}

type BSE struct {
	DB             db.Datastore
	BSEBhavCopyUrl string
}

// InitiateDownload - Start download of BSE bhavcopy
func (b *BSE) InitiateDownload(date time.Time) error {
	// fmt.Println(bseBhavFileName(date))
	BhavFileName := "EQ170521_CSV.ZIP"
	// BhavFileName := bseBhavFileName(date)

	csvFiles, err := b.downloadBseBhavCopy(BhavFileName)
	if err != nil {
		return err
	}

	fmt.Println(csvFiles)

	if err = b.ReadCSVFile(csvFiles); err != nil {
		return err
	}

	return nil
}

// DownloadBSEBhavCopy - Download BSE Bhav Copy for daily trades
func (b BSE) downloadBseBhavCopy(BhavFileName string) ([]string, error) {
	var csvFileCollection []string

	bseURL := b.BSEBhavCopyUrl + BhavFileName
	fmt.Println(bseURL)

	response, err := http.Get(bseURL)

	if err != nil {
		fmt.Println(err)
		return csvFileCollection, err
	}

	defer response.Body.Close()

	out, err := os.Create(filepath.Join(tempDir, BhavFileName))
	if err != nil {
		fmt.Println(err)
		return csvFileCollection, err
	}
	defer out.Close()

	fmt.Println("status", response.Status)
	if response.StatusCode != 200 {
		return csvFileCollection, err
	}

	// Write the body to file
	_, err = io.Copy(out, response.Body)
	fmt.Printf("err: %s", err)

	archive, err := zip.OpenReader(filepath.Join(tempDir, BhavFileName))
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	defer archive.Close()

	for _, file := range archive.Reader.File {

		if filepath.Ext(file.Name) == csvFileType {
			// add file names to array
			csvFileCollection = append(csvFileCollection, filepath.Join(tempDir, file.Name))
		}

		reader, err := file.Open()
		if err != nil {
			fmt.Printf("err: %s", err)
		}

		defer reader.Close()
		destPath := filepath.Join(tempDir, file.Name)

		writer, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fmt.Printf("err: %s", err)
		}
		defer writer.Close()

		if _, err = io.Copy(writer, reader); err != nil {
			fmt.Printf("err: %s", err)
		}
	}
	return csvFileCollection, nil
}

func bseBhavFileName(date time.Time) string {
	if date.IsZero() {
		date = time.Now()
	}
	return fmt.Sprintf("EQ%s_CSV.ZIP", date.Format("020106"))
}

func (b *BSE) ReadCSVFile(csvFiles []string) error {

	for _, file := range csvFiles {
		var bse_bhav []models.BSE_BHAV

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		csvRecords, err := csv.NewReader(f).ReadAll()
		if err != nil {
			return err
		}
		// skip header for CSV file
		csvData := csvRecords[1:]

		for _, line := range csvData {

			bhavCopy := models.BSE_BHAV{
				SC_CODE:    stringToUint32(line[0]),
				SC_NAME:    line[1],
				SC_GROUP:   line[2],
				SC_TYPE:    line[3],
				OPEN:       stringToFloat64(line[4]),
				HIGH:       stringToFloat64(line[5]),
				LOW:        stringToFloat64(line[6]),
				CLOSE:      stringToFloat64(line[7]),
				LAST:       stringToFloat64(line[8]),
				PREVCLOSE:  stringToFloat64(line[9]),
				NO_TRADES:  stringToUint64(line[10]),
				NET_TURNOV: stringToUint64(line[11]),
				TDCLOINDI:  line[12],
				CreatedAt:  datatypes.Date(time.Now()),
			}
			bse_bhav = append(bse_bhav, bhavCopy)
		}
		// fmt.Println(bse_bhav)
		b.DB.UploadBSEBhavCopy(bse_bhav)
	}

	return nil
}

func stringToUint32(str string) uint32 {
	val, err := strconv.Atoi(str)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return uint32(val)
}

func stringToUint64(str string) uint64 {
	val, err := strconv.Atoi(str)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return uint64(val)
}

func stringToFloat64(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return val
}
