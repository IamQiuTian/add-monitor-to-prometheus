package g

type ConsulData struct {
	ID string
	Name string
	Address string
	Port int
	Tags []string
}

type ConsulDataJson struct {
	Group    string `json:"group"`
	Hostname string `json:"hostname"`
	Address   string `json:"address"`
	Port	 int   `json:"port"`
}


