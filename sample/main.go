package main

import (
	"flag"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/glog"
	"github.com/hkurokawa/go-connpass"
	"github.com/taskd19/tegami"
)

type dryTwitterClient struct {
	tegami.TwitterClient
}

func (client *dryTwitterClient) Notify(event *connpass.Event) (anaconda.Tweet, error) {
	glog.Info("Dry")
	return anaconda.Tweet{}, nil
}

func main() {
	flag.Set("stderrthreshold", "INFO")
	flag.Parse()

	timeframe, consumerKey, consumerSecret, accessToken, accessTokenSecret := tegami.GetEnvValues()
	client := tegami.NewTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	// dry run
	// timeframe := 1800
	// client := &dryTwitterClient{}

	currentTime := time.Now()
	glog.Infof("currentTime: %s", currentTime)
	glog.Infof("timeframe: %d", timeframe)
	result, err := tegami.CheckEvent(currentTime, timeframe, client)

	if err != nil {
		glog.Errorf("Tegami finished with error: %s", err)
		os.Exit(-1)
	}
	glog.Infof("Tegami notified %d events at %s", result, currentTime)
}
