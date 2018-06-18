package config

import "os"

type config struct {
	Host,
	MongoHost,
	MongoPort,
	MongoDB,
	MongoUser,
	MongoPass,
	StaticDir,
	LogFile string
}

func GetConfig() config {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = ":8080"
	}

	mongoHost, ok := os.LookupEnv("MONGO_HOST")
	if !ok {
		mongoHost = "localhost"
	}

	mongoPort, ok := os.LookupEnv("MONGO_PORT")
	if !ok {
		mongoPort = "27017"
	}

	mongoDB, ok := os.LookupEnv("MONGO_DB")
	if !ok {
		mongoDB = "eduboard"
	}

	mongoUser, ok := os.LookupEnv("MONGO_USER")
	if !ok {
		mongoUser = ""
	}

	mongPass, ok := os.LookupEnv("MONGO_PASS")
	if !ok {
		mongPass = ""
	}

	staticDir, ok := os.LookupEnv("STATIC_DIR")
	if !ok {
		staticDir = "./static"
	}

	logFile, ok := os.LookupEnv("LOGFILE")
	if !ok {
		logFile = ""
	}

	return config{
		Host:      host,
		MongoHost: mongoHost,
		MongoPort: mongoPort,
		MongoDB:   mongoDB,
		MongoUser: mongoUser,
		MongoPass: mongPass,
		StaticDir: staticDir,
		LogFile:   logFile,
	}
}
