package soundtouch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/mdns"
)

const websocketPort int = 8080
const messageBufferSize int = 256

// Speaker defines a soundtouch speaker
type Speaker struct {
	IP           net.IP
	Port         int
	BaseHTTPURL  url.URL
	WebSocketURL url.URL
	DeviceInfo   Info
	conn         *websocket.Conn
	WebSocketCh  chan *Update
	Plugins      []Plugin
}

// LookupSpeakers listens via mdns for soundtouch speakers and returns Speaker channel
func LookupSpeakers(iface *net.Interface) <-chan *Speaker {
	speakerCh := make(chan *Speaker)
	entriesCh := make(chan *mdns.ServiceEntry, 7)
	defer close(entriesCh)
	go func() {
		defer close(speakerCh)
		for entry := range entriesCh {
			fSpeaker := NewMdnsSpeaker(entry)

			// filter non-soundtouch speakers
			if fSpeaker.DeviceInfo.String() != "" {
				speakerCh <- fSpeaker
			}
		}
	}()

	params := mdns.DefaultParams("_soundtouch._tcp")
	params.Entries = entriesCh
	if iface != nil {
		params.Interface = iface
	}
	mdns.Query(params)
	return speakerCh
}

// LookupSpeakers listens via mdns for soundtouch speakers and returns Speaker channel
func LookupStaticSpeakers(speakerIPs []string) <-chan *Speaker {
	speakerCh := make(chan *Speaker)
	go func() {
		defer close(speakerCh)
		for _, ip := range speakerIPs {
			fSpeaker := NewIPSpeaker(ip)

			// filter non-soundtouch speakers
			if fSpeaker.DeviceInfo.String() != "" {
				speakerCh <- fSpeaker
			}
		}
	}()

	return speakerCh
}

// NewMdnsSpeaker returns a new Speaker entity based on a mdns service entry
func NewMdnsSpeaker(entry *mdns.ServiceEntry) *Speaker {
	if entry == nil {
		return &Speaker{}
	}

	fSpeaker := &Speaker{
		entry.AddrV4,
		entry.Port,
		url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), entry.Port),
		},
		url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), websocketPort),
		},
		Info{},
		nil,
		nil,
		nil,
	}

	// Ask ones for Info() to fill the speaker
	_, err := fSpeaker.Info()

	if err != nil {
		log.Warnf("Error %v while retrieving info for speaker %v", err, fSpeaker.DeviceInfo.Name)
	}

	return fSpeaker
}

// NewIPSpeaker creates a new speaker for the given ipAdress
func NewIPSpeaker(ipAddress string) *Speaker {
	if ipAddress == "" {
		return &Speaker{}
	}

	AddrV4 := net.ParseIP(ipAddress)
	return &Speaker{
		AddrV4,
		8090,
		url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%v:%v", AddrV4.String(), 8090),
		},
		url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("%v:%v", AddrV4.String(), websocketPort),
		},
		Info{},
		nil,
		nil,
		nil,
	}
}

// Listen creates a listenes that distributes Update messages via channel
func (s *Speaker) Listen() (chan *Update, error) {
	spkLogger := log.WithFields(log.Fields{
		"Speaker": s.DeviceInfo.Name,
		"ID":      s.DeviceInfo.DeviceID,
	})
	spkLogger.Tracef("Dialing %v", s.WebSocketURL.String())
	conn, _, err := websocket.DefaultDialer.Dial(
		s.WebSocketURL.String(),
		http.Header{
			"Sec-WebSocket-Protocol": []string{"gabbo"},
		})
	if err != nil {
		return nil, err
	}

	s.conn = conn
	messageCh := make(chan *Update, messageBufferSize)
	go func() {
		for {
			mLogger := log.WithFields(log.Fields{
				"Speaker": s.DeviceInfo.Name,
				"ID":      s.DeviceInfo.DeviceID,
			})
			_, body, err := conn.ReadMessage()
			if err != nil {
				spkLogger.Warn(err)
				spkLogger.Warn("Trying to reconnect")
				spkLogger.Tracef("Re-Dialing %v", s.WebSocketURL.String())
				conn, _, err = websocket.DefaultDialer.Dial(
					s.WebSocketURL.String(),
					http.Header{
						"Sec-WebSocket-Protocol": []string{"gabbo"},
					})

				if err != nil {
					spkLogger.Warn("Re-Dialing failed")
					spkLogger.Fatal(err)
				}

				s.conn = conn
			}
			mLogger.Tracef("Raw Message: %v", string(body))

			update, err := NewUpdate(body)
			if err != nil {
				mLogger.Tracef("Message: unkown")
				mLogger.Tracef(err.Error())
			} else {
				mLogger.Tracef("Message: %v", update)
			}
			if update != nil {
				messageCh <- update
			}
		}
	}()
	return messageCh, nil

}

// Close closes the socket to the soundtouch speaker
func (s *Speaker) Close() error {
	log.Debugf("Closing socket")
	return s.conn.Close()
}

// GetData returns received raw data retrieved a GET for a given soundtouch action
func (s *Speaker) GetData(action string) ([]byte, error) {
	actionURL := s.BaseHTTPURL
	actionURL.Path = action

	mLogger := log.WithFields(log.Fields{
		"Speaker": s.DeviceInfo.Name,
		"ID":      s.DeviceInfo.DeviceID,
	})

	mLogger.Tracef("GET: %s\n", actionURL.String())

	resp, err := http.Get(actionURL.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

// SetData sets raw data via  POST for a given soundtouch action
func (s *Speaker) SetData(action string, input []byte) ([]byte, error) {
	actionURL := s.BaseHTTPURL
	actionURL.Path = action
	buffer := bytes.NewBuffer(input)

	mLogger := log.WithFields(log.Fields{
		"Speaker": s.DeviceInfo.Name,
		"ID":      s.DeviceInfo.DeviceID,
	})

	mLogger.Tracef("POST: %s, %v\n", actionURL.String(), buffer)

	resp, err := http.Post(actionURL.String(), "application/xml", buffer)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Name returns the speakers name as indicated in the info message, or "" if name unknwon
func (s *Speaker) Name() (name string) {
	return s.DeviceInfo.Name
}

// DeviceID returns the speakers DeviceID as indicated in the info message, or "" if name unknwon
func (s *Speaker) DeviceID() (name string) {
	return s.DeviceInfo.DeviceID
}
