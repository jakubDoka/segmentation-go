package network

// Clients is map of clients used by server
var Clients = Registry{
	Conns: map[string]*ClientData{},
}

// Server is global state of host used by client
var Server = Host{}
