package soundtouch

import (
	"net"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSpeaker_Name(t *testing.T) {
	type fields Speaker
	tests := []struct {
		name     string
		fields   fields
		wantName string
	}{
		{
			"Valid name",
			fields{
				IP:           nil,
				Port:         0,
				BaseHTTPURL:  url.URL{},
				WebSocketURL: url.URL{},
				DeviceInfo: Info{
					DeviceID: "xxee",
					Name:     "Speakers Name",
					Type:     "",
					Raw:      nil,
				},
				conn: &websocket.Conn{},
			},
			"Speakers Name",
		},
		{
			"Empty name",
			fields{
				IP:           nil,
				Port:         0,
				BaseHTTPURL:  url.URL{},
				WebSocketURL: url.URL{},
				DeviceInfo:   Info{},
				conn:         &websocket.Conn{},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Speaker{
				IP:           tt.fields.IP,
				Port:         tt.fields.Port,
				BaseHTTPURL:  tt.fields.BaseHTTPURL,
				WebSocketURL: tt.fields.WebSocketURL,
				DeviceInfo:   tt.fields.DeviceInfo,
				conn:         tt.fields.conn,
			}
			if gotName := s.Name(); gotName != tt.wantName {
				t.Errorf("Speaker.Name() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func TestSpeaker_DeviceID(t *testing.T) {
	type fields struct {
		IP           net.IP
		Port         int
		BaseHTTPURL  url.URL
		WebSocketURL url.URL
		DeviceInfo   Info
		conn         *websocket.Conn
	}
	tests := []struct {
		name     string
		fields   fields
		wantName string
	}{
		{
			"Valid DeviceId",
			fields{
				IP:           nil,
				Port:         0,
				BaseHTTPURL:  url.URL{},
				WebSocketURL: url.URL{},
				DeviceInfo: Info{
					DeviceID: "xxee",
					Name:     "",
					Type:     "",
					Raw:      nil,
				},
				conn: &websocket.Conn{},
			},
			"xxee",
		},
		{
			"Empty name",
			fields{
				IP:           nil,
				Port:         0,
				BaseHTTPURL:  url.URL{},
				WebSocketURL: url.URL{},
				DeviceInfo:   Info{},
				conn:         &websocket.Conn{},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Speaker{
				IP:           tt.fields.IP,
				Port:         tt.fields.Port,
				BaseHTTPURL:  tt.fields.BaseHTTPURL,
				WebSocketURL: tt.fields.WebSocketURL,
				DeviceInfo:   tt.fields.DeviceInfo,
				conn:         tt.fields.conn,
			}
			if gotName := s.DeviceID(); gotName != tt.wantName {
				t.Errorf("Speaker.DeviceID() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

// LookupStaticSpeakers returns a channel that will deliver a stream of *Speaker.
// The *Speaker objects are created from the given list of IPs.
// The function will filter out IPs that are not valid soundtouch speakers.
func TestLookupStaticSpeakers(t *testing.T) {
	tests := []struct {
		name       string
		speakerIPs []string
		wantCount  int
	}{
		{
			name:       "Empty IP list",
			speakerIPs: []string{},
			wantCount:  0,
		},
		{
			name:       "Single IP",
			speakerIPs: []string{"192.168.1.100"},
			wantCount:  1,
		},
		{
			name:       "Multiple IPs",
			speakerIPs: []string{"192.168.1.100", "192.168.1.101", "192.168.1.102"},
			wantCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			speakerCh := LookupStaticSpeakers(tt.speakerIPs)

			count := 0
			for speaker := range speakerCh {
				// check the content of the speaker
				if speaker.IP.String() != tt.speakerIPs[count] {
					t.Errorf("LookupStaticSpeakers() got speaker with IP %v, want %v", speaker.IP.String(), tt.speakerIPs[count-1])
				}
				if speaker.Port != 8090 {
					t.Errorf("LookupStaticSpeakers() got speaker with port %v, want 8090", speaker.Port)
				}

				if speaker == nil {
					t.Error("received nil speaker")
				}
				count++
			}

			if count != tt.wantCount {
				t.Errorf("LookupStaticSpeakers() got %v speakers, want %v", count, tt.wantCount)
			}
		})
	}
}
