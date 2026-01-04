package model

import "strconv"

type EndpointKind string

const (
	EndpointNet EndpointKind = "net"
	EndpointURL EndpointKind = "url"
)

type Endpoint struct {
	Kind EndpointKind `json:"kind"`
	Net  *NetEndpoint `json:"net,omitempty"`
	URL  *URLEndpoint `json:"url,omitempty"`
}

type NetEndpoint struct {
	Addr string
	Port int
}

func (e Endpoint) IsZero() bool {
	return e.Kind == "" && e.Net == nil && e.URL == nil
}

func (e Endpoint) String() string {
	switch e.Kind {
	case EndpointNet:
		if e.Net == nil {
			return ""
		}
		return e.Net.String()
	case EndpointURL:
		if e.URL == nil {
			return ""
		}
		return e.URL.String()
	default:
		// fallback: if Kind not set but one pointer exists
		if e.Net != nil {
			return e.Net.String()
		}
		if e.URL != nil {
			return e.URL.String()
		}
		return ""
	}
}

func (e NetEndpoint) String() string {
	return e.Addr + ":" + strconv.Itoa(e.Port)
}

type URLEndpoint struct {
	Raw string

	Scheme string
	Host   string
	Port   int
	Path   string
}

func NewEndpointNet(addr string, port int) Endpoint {
	ne := &NetEndpoint{Addr: addr, Port: port}
	return Endpoint{Kind: EndpointNet, Net: ne}
}

func NewEndpointURL(raw string) Endpoint {
	ue := &URLEndpoint{Raw: raw}
	return Endpoint{Kind: EndpointURL, URL: ue}
}

func (e URLEndpoint) String() string {
	return e.Raw
}
