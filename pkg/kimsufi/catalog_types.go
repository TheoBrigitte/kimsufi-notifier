package kimsufi

type Catalog struct {
	CatalogID    int          `json:"catalogID"`
	Locale       Locale       `json:"locale"`
	Plans        []Plan       `json:"plans"`
	Products     []Product    `json:"products"`
	Addons       []Addon      `json:"addons"`
	PlanFamilies []PlanFamily `json:"planFamilies"`
}

type Locale struct {
	CurrencyCode string `json:"currencyCode"`
	Subsidiary   string `json:"subsidiary"`
	TaxRate      int    `json:"taxRate"`
}

//
// Plan structure
//

type Plan struct {
	PlanCode      string            `json:"planCode"`
	InvoiceName   string            `json:"invoiceName"`
	AddonFamilies []PlanAddonFamily `json:"addonFamilies"`
	Product       string            `json:"product"`
	PricingType   string            `json:"pricingType"`
	//ConsumptionConfiguration string          `json:"consumptionConfiguration"`
	Pricings      []Pricing           `json:"pricings"`
	Configuration []PlanConfiguration `json:"configuration"`
	Family        string              `json:"family"`
	Blobs         PlanBlobs           `json:"blobs"`
}

type PlanAddonFamily struct {
	Name      string   `json:"name"`
	Exclusive bool     `json:"exclusive"`
	Mandatory bool     `json:"mandatory"`
	Addons    []string `json:"addons"`
	Default   string   `json:"default"`
}

type Pricing struct {
	Phase           int               `json:"phase"`
	Capacities      []string          `json:"capacities"`
	Commitement     int               `json:"commitement"`
	Description     string            `json:"description"`
	Interval        int               `json:"interval"`
	IntervalUnit    string            `json:"intervalUnit"`
	Quantity        PlanPricingMinMax `json:"quantity"`
	Repeat          PlanPricingMinMax `json:"repeat"`
	Price           int               `json:"price"`
	Tax             int               `json:"tax"`
	Mode            string            `json:"mode"`
	Strategy        string            `json:"strategy"`
	MustBeCompleted bool              `json:"mustBeCompleted"`
	Type            string            `json:"type"`
	//Promotions []string `json:"promotions"`
	EngagementConfiguration PlanPricingEngagementConfiguration `json:"engagementConfiguration"`
}

type PlanPricingMinMax struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type PlanPricingEngagementConfiguration struct {
	DefaultEndAction string `json:"defaultEndAction"`
	Duration         string `json:"duration"`
	Type             string `json:"type"`
}

type PlanConfiguration struct {
	Name        string   `json:"name"`
	IsCustom    bool     `json:"isCustom"`
	IsMandatory bool     `json:"isMandatory"`
	Values      []string `json:"values"`
}

type PlanBlobs struct {
	Commercial PlanBlobsCommercial `json:"commercial"`
}

type PlanBlobsCommercial struct {
	Range string `json:"range"`
}

//
// Product structure
//

type Product struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Blobs       ProductBlobs `json:"blobs"`
}

type ProductBlobs struct {
	Technical ProductBlobsTechnical `json:"technical"`
}

type ProductBlobsTechnical struct {
	Storage ProductBlobsTechnicalStorage `json:"storage"`
}

type ProductBlobsTechnicalStorage struct {
	Raid        string                                  `json:"raid"`
	Disks       []ProductBlobsTechnicalStorageDisk      `json:"disks"`
	HotSwap     bool                                    `json:"hotSwap"`
	RaidDetails ProductBlobsTechnicalStorageRaidDetails `json:"raidDetails"`
}

type ProductBlobsTechnicalStorageDisk struct {
	Specs      string `json:"specs"`
	Usage      string `json:"usage"`
	Number     int    `json:"number"`
	Capacity   int    `json:"capacity"`
	Interface  string `json:"interface"`
	Technology string `json:"technology"`
}

type ProductBlobsTechnicalStorageRaidDetails struct {
	Type string `json:"type"`
}

//
// Addon structure
//

type Addon struct {
	PlanCode    string `json:"planCode"`
	InvoiceName string `json:"invoiceName"`
	// AddonFamilies []AddonAddonFamily `json:"addonFamilies"`
	Product     string `json:"product"`
	PricingType string `json:"pricingType"`
	// ConsumptionConfiguration string          `json:"consumptionConfiguration"`
	Pricings []Pricing `json:"pricings"`
	//Configurations []AddonConfiguration `json:"configurations"`
	//Family string `json:"family"`
	//Blobs AddonBlobs `json:"blobs"`
}

//
// PlanFamily structure
//

type PlanFamily struct {
	Name string `json:"name"`
}
