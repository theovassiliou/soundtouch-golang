// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/theovassiliou/soundtouch-golang/server/restapi/operations"

	soundtouch "github.com/theovassiliou/soundtouch-golang"
)

//go:generate swagger generate server --target ../../server --name SoundtouchRESTfulJSONServer --spec ../swagger/swagger.yml
type speakerMap map[string]bool

type RestSpeaker struct {
	*soundtouch.Speaker
	SpeakerName string
	WebSocketCh chan *soundtouch.Update
}

func New(s *soundtouch.Speaker) *RestSpeaker {
	return &RestSpeaker{s, "", nil}
}

type Speakers map[string]*RestSpeaker

var visibleSpeakers = make(Speakers)

type config struct {
	Speakers            []string `short:"s" long:"speakers" description:"Speakers to listen for, all if not set"`
	Interface           string   `short:"i" long:"interface" description:"network interface to listen"`
	NoSoundtouchSystems int      `short:"n" long:"noSystems" description:"Number of Soundtouch systems to scan for."`
}

var soundtouchFlags = config{}

func configureFlags(api *operations.SoundtouchRESTfulJSONServerAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			ShortDescription: "Soundtouch Flags",
			LongDescription:  "",
			Options:          &soundtouchFlags,
		},
	}
}

func configureAPI(api *operations.SoundtouchRESTfulJSONServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf
	api.JSONConsumer = runtime.JSONConsumer()
	api.TxtProducer = runtime.TextProducer()

	iff, filteredSpeakers, _ := processConfig(soundtouchFlags)

	startScanning(iff, visibleSpeakers, filteredSpeakers)
	api.PlayPauseHandler = operations.PlayPauseHandlerFunc(func(params operations.PlayPauseParams) middleware.Responder {

		return middleware.NotImplemented("operation operations.PlayPause has not yet been implemented: " + soundtouchFlags.Speakers[0])
	})

	if api.PressKeyHandler == nil {
		api.PressKeyHandler = operations.PressKeyHandlerFunc(func(params operations.PressKeyParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PressKey has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func startScanning(iff *net.Interface, visibleSpeakers Speakers, filteredSpeakers speakerMap) {
	log.Printf("Scanning for Soundtouch systems.")
	log.Printf(soundtouchFlags.Interface)
	for ok := true; ok; ok = (len(visibleSpeakers) < soundtouchFlags.NoSoundtouchSystems) {

		speakerCh := soundtouch.Lookup(iff)
		messageCh := make(chan *soundtouch.Update)

		for speaker := range speakerCh {
			speakerInfo, _ := speaker.Info()
			speaker.DeviceInfo = speakerInfo

			if checkInMap(speaker.DeviceInfo.DeviceID, visibleSpeakers) {
				log.Printf("Already included. Ignoring.")
				continue
			}

			ms := New(speaker)

			// check wether we might have to ignore the speaker
			if len(filteredSpeakers) > 0 && !(filteredSpeakers)[speakerInfo.Name] {
				// spkLogger.Traceln("Seen but ignoring messages from: ", speakerInfo.Name)
				continue
			}

			visibleSpeakers[speaker.DeviceInfo.DeviceID] = ms
			log.Printf("Listening\n")
			log.Printf(" with IP: %v", speaker.IP)

			go func(s *soundtouch.Speaker, msgChan chan *soundtouch.Update) {
				webSocketCh, _ := s.Listen()
				magicSpeaker := New(s)
				magicSpeaker.WebSocketCh = webSocketCh
				magicSpeaker.SpeakerName = visibleSpeakers[ms.DeviceInfo.DeviceID].DeviceInfo.Name
				magicSpeaker.MessageLoop()
			}(speaker, messageCh)

		}
		time.Sleep(10 * time.Second)
	}
	log.Printf("Found all Soundtouch systems. Normal Operation.")
}

// Will create the interface, the speakerMap, and the scribble database
func processConfig(conf config) (*net.Interface, speakerMap, error) {
	filteredSpeakers := make(speakerMap)
	i, err := net.InterfaceByName(conf.Interface)

	if err != nil {
		log.Fatalf("Error with interface. %s", err)
	}

	log.Printf("Listening @ %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)

	for _, value := range conf.Speakers {
		filteredSpeakers[value] = true
		log.Printf("Reacting only speakers %v\n", value)

	}

	return i, filteredSpeakers, nil
}

func checkInMap(deviceID string, list Speakers) bool {
	for _, ms := range list {
		if ms.DeviceInfo.DeviceID == deviceID {
			return true
		}
	}
	return false
}

func (m *RestSpeaker) MessageLoop() {
	for message := range m.WebSocketCh {
		log.Printf("HAAALLLOO")
		m.HandleUpdate(*message, m.WebSocketCh)
	}
}

// handle message per speaker
func (m *RestSpeaker) HandleUpdate(msg soundtouch.Update, webSocketCh chan *soundtouch.Update) {
	typeName := reflect.TypeOf(msg.Value).Name()

	if !(msg.Is("NowPlaying") || msg.Is("Volume")) {
		if !msg.Is("ConnectionStateUpdated") {
			log.Printf("Ignoring %s\n", typeName)
		}
		return
	}

	if !HasContentItem(msg) {
	}

}

func ContentItem(u soundtouch.Update) soundtouch.ContentItem {
	if HasContentItem(u) {
		return u.Value.(soundtouch.NowPlaying).Content
	}
	return soundtouch.ContentItem{}

}

func HasContentItem(u soundtouch.Update) bool {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		return true
	}
	return false
}
