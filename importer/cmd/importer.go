package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Station-Manager/adif"
	"github.com/Station-Manager/database/sqlite/models"
	"github.com/Station-Manager/utils"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/goccy/go-json"
)

func importer(filename string) error {
	if workingDir == "" {
		return fmt.Errorf("working directory not set")
	}

	if err := checkConfig(workingDir); err != nil {
		return err
	}

	if err := checkDbDir(workingDir); err != nil {
		return err
	}

	fpath, err := checkAdiFile(workingDir, filename)
	if err != nil {
		return err
	}

	data, err := readFile(fpath)
	if err != nil {
		return err
	}

	adiObj, err := adif.Marshal(data)
	if err != nil {
		return err
	}

	if err = insertIntoDatabase(adiObj); err != nil {
		return err
	}

	return nil
}

func checkConfig(workingDir string) error {
	cfgFile := filepath.Join(workingDir, "config.json")

	exists, err := utils.PathExists(cfgFile)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("config file not found at %s", cfgFile)
	}
	return nil
}

func checkDbDir(workingDir string) error {
	dbDir := filepath.Join(workingDir, "db")
	exists, err := utils.PathExists(dbDir)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("database directory not found at %s", dbDir)
	}

	return nil
}

func checkAdiFile(workingDir, filename string) (string, error) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}

	fpath := filepath.Join(workingDir, filename)
	exists, err := utils.PathExists(fpath)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("ADI file not found at %s", fpath)
	}

	fmt.Printf("Importing file: %s\n", fpath)

	return fpath, nil
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var data []byte
	if data, err = io.ReadAll(file); err != nil {
		return nil, err
	}

	return data, nil
}

func insertIntoDatabase(adiObj adif.Adif) error {
	if err := db.Open(); err != nil {
		return fmt.Errorf("importer: failed to open database: %w", err)
	}

	sessionId, err := db.GenerateSession()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	tx, cancelCtx, err := db.BeginTxContext(ctx)
	if err != nil {
		return err
	}
	defer cancelCtx()
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	var model *models.Qso
	count := 0
	for _, record := range adiObj.Records {
		if model, err = adiRecordToQsoModel(&record, 1, sessionId); err != nil {
			return fmt.Errorf("importer: failed to convert ADI record to QSO model: %w", err)
		}

		if err = model.Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("importer: failed to insert QSO model into database: %w", err)
		}
		count++
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("importer: failed to commit transaction: %w", err)
	}

	fmt.Printf("Inserted %d records into database\n", count)

	return nil
}

func adiRecordToQsoModel(r *adif.Record, logbookId, sessionId int64) (*models.Qso, error) {

	c := *r

	// blank the fields NOT to be stored in the additionalData field
	c.Call = ""
	c.Band = ""
	c.Mode = ""
	c.QsoDate = ""
	c.TimeOn = ""
	c.TimeOff = ""
	c.RstSent = ""
	c.RstRcvd = ""
	c.Country = ""
	c.Freq = ""

	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	freqMHz, err := strconv.ParseFloat(r.Freq, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid frequency value: %w", err)
	}
	freqHz := int64(math.Round(freqMHz * 1_000_000))

	model := &models.Qso{
		Call:           r.Call,
		Band:           r.Band,
		Mode:           r.Mode,
		Freq:           freqHz,
		QsoDate:        r.QsoDate,
		TimeOn:         r.TimeOn,
		TimeOff:        r.TimeOff,
		RstSent:        r.RstSent,
		RstRcvd:        r.RstRcvd,
		Country:        r.Country,
		AdditionalData: data,
		LogbookID:      logbookId,
		SessionID:      sessionId,
	}
	return model, nil
}
