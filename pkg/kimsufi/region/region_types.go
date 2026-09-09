package region

// AllowedRegions is a list of known regions and their countries.
// Icons are the unicode flags for each country.
var AllowedRegions = Regions{
	{
		DisplayName: "Europe",
		Region:      "europe",
		Endpoint:    "ovh-eu",
		Countries: []Country{
			{
				Code: "CZ",
				Icon: "🇨🇿",
			},
			{
				Code: "DE",
				Icon: "🇩🇪",
			},
			{
				Code: "ES",
				Icon: "🇪🇸",
			},
			{
				Code: "FI",
				Icon: "🇫🇮",
			},
			{
				Code: "FR",
				Icon: "🇫🇷",
			},
			{
				Code: "GB",
				Icon: "🇬🇧",
			},
			{
				Code: "IE",
				Icon: "🇮🇪",
			},
			{
				Code: "IT",
				Icon: "🇮🇹",
			},
			{
				Code: "LT",
				Icon: "🇱🇹",
			},
			{
				Code: "MA",
				Icon: "🇲🇦",
			},
			{
				Code: "NL",
				Icon: "🇳🇱",
			},
			{
				Code: "PL",
				Icon: "🇵🇱",
			},
			{
				Code: "PT",
				Icon: "🇵🇹",
			},
			{
				Code: "SN",
				Icon: "🇸🇳",
			},
			{
				Code: "TN",
				Icon: "🇹🇳",
			},
		},
	},
	{
		DisplayName: "Other",
		Region:      "canada",
		Endpoint:    "ovh-ca",
		Countries: []Country{
			{
				Code: "ASIA",
			},
			{
				Code: "AU",
				Icon: "🇦🇺",
			},
			{
				Code: "CA",
				Icon: "🇨🇦",
			},
			{
				Code: "IN",
				Icon: "🇮🇳",
			},
			{
				Code: "QC",
			},
			{
				Code: "SG",
				Icon: "🇸🇬",
			},
			{
				Code: "WE",
			},
			{
				Code: "WS",
			},
		},
	},
	{
		DisplayName: "US",
		Region:      "united-states",
		Endpoint:    "ovh-us",
		Countries: []Country{
			{
				Code: "US",
				Icon: "🇺🇸",
			},
		},
	},
}

type Regions []Region

type Region struct {
	DisplayName string
	Region      string
	Endpoint    string
	Countries   []Country
}

type Country struct {
	Code string
	Icon string
}
