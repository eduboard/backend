package config

import "os"

type config struct {
	Port,
	MongoHost,
	MongoPort,
	MongoDB,
	MongoUser,
	MongoPass,
	StaticDir string
}

func GetConfig() config {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
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

	return config{
		Port:      port,
		MongoHost: mongoHost,
		MongoPort: mongoPort,
		MongoDB:   mongoDB,
		MongoUser: mongoUser,
		MongoPass: mongPass,
		StaticDir: staticDir,
	}
}
