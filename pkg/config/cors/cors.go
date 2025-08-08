package cors

type CORSConfig struct {
	Address string `json:",env=CORS_ADDRESS"`
}
