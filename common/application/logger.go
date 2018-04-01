package application

import (
    "os"
    "github.com/sirupsen/logrus"
)

func InitLogger() {

    logrus.SetOutput(os.Stdout)

    logrus.SetLevel(logrus.InfoLevel)

    //levels := []logrus.Level{
    //    logrus.WarnLevel,
    //    logrus.PanicLevel,
    //    logrus.FatalLevel,
    //    logrus.ErrorLevel,
    //}

    if env := Env(); env == RINKEBY_ENV {
        //levels = append(levels, logrus.DebugLevel)
        //levels = append(levels, logrus.InfoLevel)
        logrus.SetLevel(logrus.DebugLevel)
    }

    // Add sentry log hook to send logs asynchronous
    //if viper.IsSet("log.sentry.dsn") {
    //    if raven, err := raven2.New(viper.GetString("log.sentry.dsn")); err == nil {
    //        raven.SetRelease(Version)
    //        raven.SetEnvironment(Env())
    //
    //        if hook, err := logrus_sentry.NewWithClientSentryHook(raven, levels); err == nil {
    //            hook.Timeout = 1000 * time.Millisecond
    //            hook.StacktraceConfiguration.Enable = true
    //            hook.StacktraceConfiguration.Level = log.ErrorLevel
    //            log.AddHook(hook)
    //        }
    //    }
    //}
}
