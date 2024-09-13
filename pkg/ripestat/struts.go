package ripestat

type AsOverview struct {
	Messages       []interface{} `json:"messages"`
	SeeAlso        []interface{} `json:"see_also"`
	Version        string        `json:"version"`
	DataCallName   string        `json:"data_call_name"`
	DataCallStatus string        `json:"data_call_status"`
	Cached         bool          `json:"cached"`
	Data           struct {
		Type     string `json:"type"`
		Resource string `json:"resource"`
		Block    struct {
			Resource string `json:"resource"`
			Desc     string `json:"desc"`
			Name     string `json:"name"`
		} `json:"block"`
		Holder         string `json:"holder"`
		Announced      bool   `json:"announced"`
		QueryStarttime string `json:"query_starttime"`
		QueryEndtime   string `json:"query_endtime"`
	} `json:"data"`
	QueryID      string `json:"query_id"`
	ProcessTime  int    `json:"process_time"`
	ServerID     string `json:"server_id"`
	BuildVersion string `json:"build_version"`
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	Time         string `json:"time"`
}

type RIR struct {
	Messages       []interface{} `json:"messages"`
	SeeAlso        []interface{} `json:"see_also"`
	Version        string        `json:"version"`
	DataCallName   string        `json:"data_call_name"`
	DataCallStatus string        `json:"data_call_status"`
	Cached         bool          `json:"cached"`
	Data           struct {
		Resource       string `json:"resource"`
		Latest         string `json:"latest"`
		QueryStarttime string `json:"query_starttime"`
		QueryEndtime   string `json:"query_endtime"`
		Lod            int    `json:"lod"`
		Rirs           []struct {
			Rir          string `json:"rir"`
			FirstTime    string `json:"first_time"`
			LastTime     string `json:"last_time"`
			Registration string `json:"registration"`
			Status       string `json:"status"`
			Resource     string `json:"resource"`
			Country      string `json:"country"`
		} `json:"rirs"`
	} `json:"data"`
	QueryID      string `json:"query_id"`
	ProcessTime  int    `json:"process_time"`
	ServerID     string `json:"server_id"`
	BuildVersion string `json:"build_version"`
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	Time         string `json:"time"`
}

type PrefixRoutingConsistency struct {
	Messages       []interface{} `json:"messages"`
	SeeAlso        []interface{} `json:"see_also"`
	Version        string        `json:"version"`
	DataCallName   string        `json:"data_call_name"`
	DataCallStatus string        `json:"data_call_status"`
	Cached         bool          `json:"cached"`
	Data           struct {
		Resource string `json:"resource"`
		Routes   []struct {
			InBgp      bool     `json:"in_bgp"`
			InWhois    bool     `json:"in_whois"`
			Prefix     string   `json:"prefix"`
			Origin     int      `json:"origin"`
			IrrSources []string `json:"irr_sources"`
			AsnName    string   `json:"asn_name"`
		} `json:"routes"`
		Parameters struct {
			Resource          string `json:"resource"`
			DataOverloadLimit string `json:"data_overload_limit"`
		} `json:"parameters"`
		QueryStarttime string `json:"query_starttime"`
		QueryEndtime   string `json:"query_endtime"`
	} `json:"data"`
	QueryID      string `json:"query_id"`
	ProcessTime  int    `json:"process_time"`
	ServerID     string `json:"server_id"`
	BuildVersion string `json:"build_version"`
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	Time         string `json:"time"`
}

type MaxmindGeoLite struct {
	Messages       []interface{} `json:"messages"`
	SeeAlso        []interface{} `json:"see_also"`
	Version        string        `json:"version"`
	DataCallName   string        `json:"data_call_name"`
	DataCallStatus string        `json:"data_call_status"`
	Cached         bool          `json:"cached"`
	Data           struct {
		LocatedResources []struct {
			Resource  string `json:"resource"`
			Locations []struct {
				Country           string   `json:"country"`
				City              string   `json:"city"`
				Resources         []string `json:"resources"`
				Latitude          float64  `json:"latitude"`
				Longitude         float64  `json:"longitude"`
				CoveredPercentage float64  `json:"covered_percentage"`
			} `json:"locations"`
			UnknownPercentage int `json:"unknown_percentage"`
		} `json:"located_resources"`
		UnknownPercentage struct {
			V4 float64 `json:"v4"`
		} `json:"unknown_percentage"`
		ResultTime   string `json:"result_time"`
		LatestTime   string `json:"latest_time"`
		EarliestTime string `json:"earliest_time"`
		Parameters   struct {
			Resource   string      `json:"resource"`
			Resolution string      `json:"resolution"`
			Cache      interface{} `json:"cache"`
		} `json:"parameters"`
	} `json:"data"`
	QueryID      string `json:"query_id"`
	ProcessTime  int    `json:"process_time"`
	ServerID     string `json:"server_id"`
	BuildVersion string `json:"build_version"`
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	Time         string `json:"time"`
}
