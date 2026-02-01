package unifi

import "fmt"

// Sysinfo holds controller system information and health from GET /proxy/network/api/s/{site}/stat/sysinfo.
// UniFi OS only. See https://github.com/unpoller/unpoller/issues/927
type Sysinfo struct {
	SiteName     string   `json:"-"`
	SourceName   string   `json:"-"`
	Timezone     string   `json:"timezone"`
	Autobackup   bool     `json:"autobackup"`
	Build        string   `json:"build"`
	Version      string   `json:"version"`
	PreviousVer  string   `json:"previous_version"`
	DataRetDays  int      `json:"data_retention_days"`
	DataRet5min  int      `json:"data_retention_time_in_hours_for_5minutes_scale"`
	DataRetHour  int      `json:"data_retention_time_in_hours_for_hourly_scale"`
	DataRetDay   int      `json:"data_retention_time_in_hours_for_daily_scale"`
	DataRetMonth int      `json:"data_retention_time_in_hours_for_monthly_scale"`
	UpdateAvail  bool     `json:"update_available"`
	UpdateDown   bool     `json:"update_downloaded"`
	Hostname     string   `json:"hostname"`
	Name         string   `json:"name"`
	IPAddrs      []string `json:"ip_addrs"`
	InformPort   int      `json:"inform_port"`
	HTTPSPort    int      `json:"https_port"`
	PortalPort   int      `json:"portal_http_port"`
	Uptime       int64    `json:"uptime"`
	AnonymousID  string   `json:"anonymous_controller_id"`
	HasWebRTC    bool     `json:"has_webrtc_support"`
	DeviceType   string   `json:"ubnt_device_type"`
	UDMVersion   string   `json:"udm_version"`
	Unsupported  int      `json:"unsupported_device_count"`
	IsCloud      bool     `json:"is_cloud_console"`
	ConsoleVer   string   `json:"console_display_version"`
}

// GetSysinfoSite returns controller system info for a single site.
func (u *Unifi) GetSysinfoSite(site *Site) (*Sysinfo, error) {
	path := fmt.Sprintf(APISysinfoPath, site.Name)

	var s Sysinfo

	if err := u.GetData(path, &s); err != nil {
		return nil, err
	}

	s.SiteName = site.SiteName
	s.SourceName = site.SourceName

	return &s, nil
}

// GetSysinfo returns controller system info for all sites.
func (u *Unifi) GetSysinfo(sites []*Site) ([]*Sysinfo, error) {
	data := make([]*Sysinfo, 0, len(sites))

	for _, site := range sites {
		s, err := u.GetSysinfoSite(site)
		if err != nil {
			return nil, fmt.Errorf("GetSysinfo(%s): %w", site.Name, err)
		}

		data = append(data, s)
	}

	return data, nil
}
