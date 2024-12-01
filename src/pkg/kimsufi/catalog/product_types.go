package catalog

type Product struct {
	Blobs       ProductBlobs `json:"blobs"`
	Description string       `json:"description"`
	Name        string       `json:"name"`
}

type ProductBlobs struct {
	Technical ProductBlobsTechnical `json:"technical"`
}

type ProductBlobsTechnical struct {
	Bandwidth ProductBlobsTechnicalBandwidth `json:"bandwidth,omitempty"`
	Memory    ProductBlobsTechnicalMemory    `json:"memory,omitempty"`
	Server    ProductBlobsTechnicalServer    `json:"server,omitempty"`
	Storage   ProductBlobsTechnicalStorage   `json:"storage,omitempty"`
}

type ProductBlobsTechnicalBandwidth struct {
	Burst      int     `json:"burst"`
	Guaranteed bool    `json:"guaranteed"`
	Level      float64 `json:"level"`
	Limit      int     `json:"limit"`
}

type ProductBlobsTechnicalMemory struct {
	ECC       bool   `json:"ecc"`
	Frequency int    `json:"frequency"`
	Interface string `json:"interface"`
	RamType   string `json:"ramType"`
	Size      int    `json:"size"`
}

type ProductBlobsTechnicalServer struct {
	CPU ProductBlobsTechnicalCPU `json:"cpu"`
}

type ProductBlobsTechnicalCPU struct {
	Brand     string  `json:"brand"`
	Cores     int     `json:"cores"`
	Frequency float64 `json:"frequency"`
	Model     string  `json:"model"`
	Number    int     `json:"number"`
	Threads   int     `json:"threads"`
	Type      string  `json:"type"`
}

type ProductBlobsTechnicalStorage struct {
	Disks       []ProductBlobsTechnicalStorageDisk      `json:"disks"`
	HotSwap     bool                                    `json:"hotSwap"`
	Raid        string                                  `json:"raid"`
	RaidDetails ProductBlobsTechnicalStorageRaidDetails `json:"raidDetails"`
}

type ProductBlobsTechnicalStorageDisk struct {
	Capacity   int    `json:"capacity"`
	Interface  string `json:"interface"`
	Number     int    `json:"number"`
	Specs      string `json:"specs"`
	Technology string `json:"technology"`
	Usage      string `json:"usage"`
}

type ProductBlobsTechnicalStorageRaidDetails struct {
	Type string `json:"type"`
}
