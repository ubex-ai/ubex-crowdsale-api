package application

import (
    "os"
    "fmt"
    "bytes"
    "io/ioutil"
    "path/filepath"

    "github.com/spf13/viper"
    "github.com/sirupsen/logrus"
    "github.com/gin-gonic/gin"
)

func InitConfig() {
    setDefaults()
    bindEnvVariables()
    InitDefaultConfigFiles()
    verifyConfig()
    setReleaseMode()
    printAppState()
}

func InitDefaultConfigFiles() {
    dir, err := os.Getwd()

    if err != nil {
        logrus.Fatal(err)
    }

    viper.SetConfigType("yaml")

    files := []string{
        dir + "/conf/global.yaml",
        dir + fmt.Sprintf("/conf/config.%s.yaml", viper.GetString(ENVNAME)),
        dir + "/conf/local.yaml",
    }

    ReadAndMerge(files)
}

func ReadAndMerge(files []string) {
    for _, file := range files {
        data, err := ioutil.ReadFile(file)

        if err != nil && filepath.Base(file) != "local.yaml" {
            logrus.Panicf(
                "Can not load config file %s. An error message: %s", file, err,
            )
            return
        }

        err = viper.MergeConfig(bytes.NewReader(data))
        if err != nil {
            logrus.Panicf(
                "An error '%s' has occurred while parsing config file %s",
                err,
                file,
            )
        }
    }
}

func setReleaseMode() {
    if Env() == MAINNET_ENV {
        gin.SetMode(gin.ReleaseMode)
    }
}

func printAppState() {
    logrus.Infof("Application Environment. %s = %s", ENVNAME, Env())
    logrus.Infof("Application runs with the following config:\n%+v\n", viper.AllSettings())
}

func verifyConfig() {
    if !viper.IsSet("project.name") {
        logrus.Warn("Project name is not configured")
    }
}

func bindEnvVariables() {
    viper.BindEnv(ENVNAME)
}

func setDefaults() {
    viper.SetDefault(ENVNAME, RINKEBY_ENV)
}
