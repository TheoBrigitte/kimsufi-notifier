package catalog

type Plan struct {
	AddonFamilies  []PlanAddonFamily   `json:"addonFamilies"`
	Blobs          PlanBlobs           `json:"blobs,omitempty"`
	Configurations []PlanConfiguration `json:"configurations"`
	Family         string              `json:"family"`
	InvoiceName    string              `json:"invoiceName"`
	PlanCode       string              `json:"planCode"`
	PricingType    string              `json:"pricingType"`
	Pricings       []PlanPricing       `json:"pricings"`
	Product        string              `json:"product"`
}

type PlanAddonFamily struct {
	Addons    []string `json:"addons"`
	Default   string   `json:"default"`
	Exclusive bool     `json:"exclusive"`
	Mandatory bool     `json:"mandatory"`
	Name      string   `json:"name"`
}

type PlanPricing struct {
	Capacities              []string                           `json:"capacities"`
	Commitement             int                                `json:"commitment"`
	Description             string                             `json:"description"`
	EngagementConfiguration PlanPricingEngagementConfiguration `json:"engagementConfiguration,omitempty"`
	Interval                int                                `json:"interval"`
	IntervalUnit            string                             `json:"intervalUnit"`
	Mode                    string                             `json:"mode"`
	MustBeCompleted         bool                               `json:"mustBeCompleted"`
	Phase                   int                                `json:"phase"`
	Price                   int                                `json:"price"`
	Quantity                PlanPricingMinMax                  `json:"quantity"`
	Repeat                  PlanPricingMinMax                  `json:"repeat"`
	Strategy                string                             `json:"strategy"`
	Tax                     int                                `json:"tax"`
	Type                    string                             `json:"type"`
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
	Commercial PlanBlobsCommercial `json:"commercial,omitempty"`
}

type PlanBlobsCommercial struct {
	Range string `json:"range,omitempty"`
}
