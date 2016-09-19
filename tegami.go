package tegami

import (
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/hkurokawa/go-connpass"
)

// GetEnvValues get values from environment values and return following values.
// TEGAMI_TIMEFRAME (int), TEGAMI_TW_CONSUMER_KEY (string), TEGAMI_TW_CONSUMER_SECRET (string), TEGAMI_TW_ACCESS_TOKEN (string), TEGAMI_TW_ACCESS_TOKEN_SECRET (string)
func GetEnvValues() (int, string, string, string, string) {
	timeframe := 1800
	argTimeframe := os.Getenv("TEGAMI_TIMEFRAME")
	if argTimeframe != "" {
		var err error
		timeframe, err = strconv.Atoi(argTimeframe)
		if err != nil {
			glog.Warningf("Failed to parse TEGAMI_TIMEFRAME. Use default timeframe value, 1800: %s", err)
		}
	}

	argTWConsumerKey := os.Getenv("TEGAMI_TW_CONSUMER_KEY")
	argTWConsumerSecret := os.Getenv("TEGAMI_TW_CONSUMER_SECRET")
	argTWAccessToken := os.Getenv("TEGAMI_TW_ACCESS_TOKEN")
	argTWAccessTokenSecret := os.Getenv("TEGAMI_TW_ACCESS_TOKEN_SECRET")

	return timeframe, argTWConsumerKey, argTWConsumerSecret, argTWAccessToken, argTWAccessTokenSecret
}

// CheckEvent search Connpass events and pass it to Notifier.
func CheckEvent(currentTime time.Time, timeframe int, client Notifier) (int, error) {
	query := &connpass.Query{}
	query.Order = connpass.UPDATE
	notifyCount := 0

	for {
		res, err := query.Search()

		if err != nil {
			glog.Error(err)
			return notifyCount, err
		}

		// filter and notify events
		continuePage := true
		for i := 0; i < len(res.Events); i++ {
			event := res.Events[i]
			glog.Infof("Start processing event ID: %d", event.Id)
			eventTime, err := time.Parse(time.RFC3339, event.Updated)
			if err != nil {
				glog.Errorf("An error occurred while parsing Updated field of the event: %s", err)
				return notifyCount, err
			}

			elapsed := currentTime.Sub(eventTime).Seconds()
			glog.Infof("event time: %s, sub: %f", eventTime, elapsed)
			if int(elapsed) <= timeframe {
				glog.Infof("Notify %s", event.Url)
				_, err := client.Notify(&event)

				if err != nil {
					glog.Errorf("Failed to notify event: %s", err)
					return notifyCount, err
				}
				notifyCount++
			} else {
				glog.Infof("Ignore %s due to its elapsed time", event.Url)
				continuePage = false
				break
			}
		}

		// check whether next page needs to be loaded.
		if !continuePage {
			break
		}
		offset := res.Start + res.Returned
		if offset > res.Available {
			break
		}
		query.Start = offset
	}
	return notifyCount, nil
}
