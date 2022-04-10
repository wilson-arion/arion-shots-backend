package cloudinary

import (
	"github.com/cloudinary/cloudinary-go"
	"github.com/joho/godotenv"
	"log"
)

const (
	clCloudName = "CL_CLOUD_NAME"
	clApiKey    = "CL_API_KEY"
	clApiSecret = "CL_API_SECRET"
)

var (
	CL *cloudinary.Cloudinary

	cloudName string
	apiKey    string
	apiSecret string
)

// init establish the connection with the MySQL database.
func init() {
	var envs map[string]string
	envs, err := godotenv.Read()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cloudName = envs[clCloudName]
	apiKey = envs[clApiKey]
	apiSecret = envs[clApiSecret]
}

func Initialize() {
	cl, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

	CL = cl
	log.Println("cloudinary successfully configured")
}
