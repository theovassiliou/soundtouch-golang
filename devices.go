package soundtouch

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type speakers map[string]*Speaker
type speakerMap map[string]bool

var visibleSpeakers = make(speakers)
var filteredSpeakers speakerMap

type NetworkConfig struct {
	InterfaceName      string
	NoOfSystems        int
	SpeakerToListenFor []string
	UpdateHandlers     []UpdateHandlerConfig
}

func GetDevices(conf NetworkConfig) (speakers chan *Speaker) {
	return getDevices(conf, true)
}

func SearchDevices(conf NetworkConfig) (speakers chan *Speaker) {
	return getDevices(conf, false)
}

// GetDevices starts listening on the indicated interface for the speakers to listen for.
// passes to speakers the series of speakers that are handled for further processing
func getDevices(conf NetworkConfig, closeChannel bool) (speakers chan *Speaker) {
	iff, err := net.InterfaceByName(conf.InterfaceName)

	if err != nil {
		log.Fatalf("Error with interface. %s", err)
	}

	for _, value := range conf.SpeakerToListenFor {
		filteredSpeakers[value] = true
		log.Debugf("Reacting only speakers %v\n", value)

	}
	speakers = make(chan *Speaker)

	// var wg sync.WaitGroup
	log.Debugf("Scanning for Soundtouch systems.")
	go func() {
		for ok := true; ok; ok = (len(visibleSpeakers) < conf.NoOfSystems) {
			speakerCh := Lookup(iff)
			messageCh := make(chan *Update)

			for speaker := range speakerCh {
				// LookUp found a speaker
				speakerInfo, _ := speaker.Info()
				speaker.DeviceInfo = speakerInfo
				spkLogger := log.WithFields(log.Fields{
					"Speaker": speaker.Name(),
					"ID":      speaker.DeviceID(),
				})

				if contains(visibleSpeakers, speaker.DeviceID()) {
					spkLogger.Debugf("Already included. Ignoring.")
					continue
				}

				// check whether we might have to ignore the speaker
				if len(filteredSpeakers) > 0 && !(filteredSpeakers)[speakerInfo.Name] {
					// spkLogger.Traceln("Seen but ignoring messages from: ", speakerInfo.Name)
					continue
				}

				visibleSpeakers[speaker.DeviceID()] = speaker

				// register handles
				for _, uh := range conf.UpdateHandlers {
					speaker.AddUpdateHandler(uh)
				}

				go func(s *Speaker, msgChan chan *Update) {
					// defer wg.Done()
					webSocketCh, _ := s.Listen()
					s.webSocketCh = webSocketCh
					s.Handle(webSocketCh)
				}(speaker, messageCh)

				speakers <- speaker

			}
			if len(visibleSpeakers) < conf.NoOfSystems {
				time.Sleep(10 * time.Second)
			} else {
				ok = false
			}
		}

		if closeChannel {
			close(speakers)
		}
		log.Debugf("Found all Soundtouch systems. Normal Operation.")
	}()
	// wg.Wait()
	return
}

func contains(list speakers, deviceID string) bool {
	for _, ms := range list {
		if ms.DeviceInfo.DeviceID == deviceID {
			return true
		}
	}
	return false
}
