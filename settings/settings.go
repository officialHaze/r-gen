package settings

import (
	"fmt"
)

func Generate() error {
	conf, err := readSettingsConf()
	if err != nil {
		return fmt.Errorf("error generating settings: %v", err)
	}

	MySettings = &Settings{
		ENV_FILE_NAME:         conf.Env_File_Name,
		SERVER_PORT:           conf.Server_Port,
		UPLOAD_MAX_SIZE:       conf.Upload_Max_Size,
		UPLOAD_LIMIT:          conf.Upload_Limit,
		PDF_MAX_MERGE_ALLOWED: conf.Pdf_Max_Merge_Allowed,
		GEN_QUOTA:             conf.Gen_Quota,
		RATE_LIMIT: 		   conf.Rate_Limit,
		
	}

	return nil
}

type Settings struct {
	ENV_FILE_NAME         string
	SERVER_PORT           int
	UPLOAD_MAX_SIZE       string
	UPLOAD_LIMIT          int
	PDF_MAX_MERGE_ALLOWED int
	GEN_QUOTA             string
	RATE_LIMIT	          int		
}
