package unifi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WANEnrichedConfiguration represents the complete WAN configuration with statistics.
type WANEnrichedConfiguration struct {
	Configuration WANConfiguration `json:"configuration"`
	Details       WANDetails       `json:"details"`
	Statistics    WANStatistics    `json:"statistics"`
}

// WANConfiguration represents the WAN network configuration.
type WANConfiguration struct {
	ID                        string                  `json:"_id"`
	AttrHiddenID              string                  `json:"attr_hidden_id,omitempty"`
	AttrNoDelete              FlexBool                `json:"attr_no_delete"`
	Name                      string                  `json:"name"`
	Purpose                   string                  `json:"purpose"` // "wan"
	WANNetworkgroup           string                  `json:"wan_networkgroup"`
	WANType                   string                  `json:"wan_type"` // "dhcp", "static", "pppoe"
	WANTypeV6                 string                  `json:"wan_type_v6"`
	WANFailoverPriority       FlexInt                 `json:"wan_failover_priority"`
	WANLoadBalanceType        string                  `json:"wan_load_balance_type"` // "weighted", "failover-only"
	WANLoadBalanceWeight      FlexInt                 `json:"wan_load_balance_weight"`
	WANDNSPreference          string                  `json:"wan_dns_preference"` // "auto", "manual"
	WANIPv6DNSPreference      string                  `json:"wan_ipv6_dns_preference"`
	WANSmartqEnabled          FlexBool                `json:"wan_smartq_enabled"`
	WANMagicEnabled           FlexBool                `json:"wan_magic_enabled"`
	WANProviderCapabilities   WANProviderCapabilities `json:"wan_provider_capabilities"`
	WANVlanEnabled            FlexBool                `json:"wan_vlan_enabled"`
	WANDHCPOptions            []string                `json:"wan_dhcp_options"`
	WANIPAliases              []string                `json:"wan_ip_aliases"`
	IPv6WANDelegationType     string                  `json:"ipv6_wan_delegation_type"`
	ReportWANEvent            FlexBool                `json:"report_wan_event"`
	SettingPreference         string                  `json:"setting_preference"`
	WANDHCPCos                FlexInt                 `json:"wan_dhcp_cos,omitempty"`
	WANDNS1                   string                  `json:"wan_dns1,omitempty"`
	WANDNS2                   string                  `json:"wan_dns2,omitempty"`
	WANIPv6DNS1               string                  `json:"wan_ipv6_dns1,omitempty"`
	WANIPv6DNS2               string                  `json:"wan_ipv6_dns2,omitempty"`
	IGMPProxyFor              string                  `json:"igmp_proxy_for"`
	IGMPProxyUpstream         FlexBool                `json:"igmp_proxy_upstream"`
	IGMPProxyDownstreamNetIDs []string                `json:"igmp_proxy_downstream_networkconf_ids"`
	MacOverrideEnabled        FlexBool                `json:"mac_override_enabled"`
	LanTunnelWANNetwork       FlexBool                `json:"lanTunnelWanNetwork"`
	WAN5GNetwork              FlexBool                `json:"wan5gNetwork"`
	SingleNetworkLan          string                  `json:"single_network_lan"`
	WANDsliteRemoteHostAuto   FlexBool                `json:"wan_dslite_remote_host_auto"`
	WANEgressQOSEnabled       FlexBool                `json:"wan_egress_qos_enabled,omitempty"`
	WANStaticIP               string                  `json:"wan_ip,omitempty"`
	WANNetmask                string                  `json:"wan_netmask,omitempty"`
	WANGateway                string                  `json:"wan_gateway,omitempty"`
}

// Note: WANProviderCapabilities is defined in usg.go

// WANDetails represents WAN connection details and service provider information.
type WANDetails struct {
	CreationTimestamp FlexInt            `json:"creation_timestamp"`
	ServiceProvider   WANServiceProvider `json:"service_provider"`
}

// WANServiceProvider represents the detected ISP information.
type WANServiceProvider struct {
	ASN  FlexInt `json:"asn"`
	City string  `json:"city"`
	Name string  `json:"name"`
}

// WANStatistics represents WAN usage and uptime statistics.
type WANStatistics struct {
	PeakUsage        WANPeakUsage `json:"peak_usage"`
	UptimePercentage float64      `json:"uptime_percentage"`
}

// WANPeakUsage represents peak bandwidth usage statistics.
type WANPeakUsage struct {
	DownloadPercentage float64 `json:"download_percentage"`
	UploadPercentage   float64 `json:"upload_percentage"`
	MaxRxBytesR        FlexInt `json:"max_rx_bytes-r"`
	MaxTxBytesR        FlexInt `json:"max_tx_bytes-r"`
}

// WANLoadBalancingStatus represents the current load balancing status.
type WANLoadBalancingStatus struct {
	Mode          string                  `json:"mode"` // "FAILOVER_ONLY", "DISTRIBUTED"
	WANInterfaces []WANLoadBalancingIface `json:"wan_interfaces"`
}

// WANLoadBalancingIface represents a WAN interface in load balancing configuration.
type WANLoadBalancingIface struct {
	Name            string  `json:"name"`
	WANNetworkgroup string  `json:"wan_networkgroup"`
	Mode            string  `json:"mode"`
	Priority        FlexInt `json:"priority"`
	Weight          FlexInt `json:"weight"`
}

// WANISPStatus represents the ISP status for WAN interfaces.
type WANISPStatus struct {
	WANInterfaces []WANInterfaceStatus `json:"wan_interfaces"`
}

// WANInterfaceStatus represents the status of a single WAN interface.
type WANInterfaceStatus struct {
	Name            string `json:"name"`
	State           string `json:"state"` // "ACTIVE", "BACKUP", "DISCONNECTED"
	WANNetworkgroup string `json:"wan_networkgroup"`
}

// WANISPStatusDetailed represents detailed ISP status from the isp-status endpoint.
type WANISPStatusDetailed struct {
	ConnectionWarnings  WANConnectionWarnings `json:"connection_warnings"`
	InternetAlerts      WANInternetAlerts     `json:"internet_alerts"`
	LatencyMax          FlexInt               `json:"latency_max"`
	PingServer          string                `json:"ping_server"`
	SpeedtestHistorical []WANSpeedtest        `json:"speedtest_historical"`
	UplinkStatus        WANUplinkStatus       `json:"uplink_status"`
}

// WANConnectionWarnings represents connection quality warnings.
type WANConnectionWarnings struct {
	Downtimes     []json.RawMessage `json:"downtimes"`
	HighLatencies []json.RawMessage `json:"high_latencies"`
	PacketLoss    []json.RawMessage `json:"packet_loss"`
}

// WANInternetAlerts represents internet connectivity alerts.
type WANInternetAlerts struct {
	Data              []json.RawMessage `json:"data"`
	PageNumber        FlexInt           `json:"page_number"`
	TotalElementCount FlexInt           `json:"total_element_count"`
	TotalPageCount    FlexInt           `json:"total_page_count"`
}

// WANSpeedtest represents a historical speedtest result.
type WANSpeedtest struct {
	ID              string  `json:"id"`
	InterfaceName   string  `json:"interface_name"`
	WANNetworkgroup string  `json:"wan_networkgroup"`
	DownloadMbps    FlexInt `json:"download_mbps"`
	UploadMbps      FlexInt `json:"upload_mbps"`
	LatencyMs       FlexInt `json:"latency_ms"`
	Time            FlexInt `json:"time"` // Unix timestamp in milliseconds
}

// WANUplinkStatus represents uplink monitoring status and statistics.
type WANUplinkStatus struct {
	LatencyThreshold FlexInt              `json:"latency_threshold"`
	ReceivedBytes    FlexInt              `json:"received_bytes"`
	Statistics       []WANUplinkStatistic `json:"statistics"`
}

// WANUplinkStatistic represents a point-in-time uplink statistic.
type WANUplinkStatistic struct {
	Downtime                FlexBool `json:"downtime"`
	Latency                 FlexInt  `json:"latency"`
	LatencyMax              FlexInt  `json:"latency_max"`
	ReceivedBytesRateAvg    FlexInt  `json:"received_bytes_rate_avg"`
	TransmittedBytesRateAvg FlexInt  `json:"transmitted_bytes_rate_avg"`
	Timestamp               FlexInt  `json:"timestamp"` // Unix timestamp in milliseconds
}

// WANSLA represents WAN SLA monitoring configuration and metrics (currently empty in API response).
type WANSLA struct {
	// Future: Will contain SLA monitoring data like latency, packet loss, jitter
	ID   string          `json:"_id,omitempty"`
	Name string          `json:"name,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

// GetWANEnrichedConfiguration returns enriched WAN configuration for all WAN interfaces.
// The API returns a top-level array [{...}, {...}], not {"data": [...]}.
// The path must be formatted with site.Name; a literal %s in the path causes "invalid URL escape" when building the request.
func (u *Unifi) GetWANEnrichedConfiguration(sites []*Site) ([]*WANEnrichedConfiguration, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	data := []*WANEnrichedConfiguration{}

	for _, site := range sites {
		path := fmt.Sprintf(APIWANEnrichedConfigPath, site.Name)
		if strings.Contains(path, "%s") {
			return nil, fmt.Errorf("WAN enriched-config path still contains %%s (site name may be empty): %q", path)
		}
		u.DebugLog("Fetching WAN enriched configuration for site %s", site.Name)

		body, err := u.GetJSON(path)
		if err != nil {
			return nil, err
		}

		var raw []*WANEnrichedConfiguration
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, err
		}
		for _, wan := range raw {
			if wan != nil {
				data = append(data, wan)
			}
		}
	}

	return data, nil
}

// GetWANLoadBalancingStatus returns the current load balancing status for WAN interfaces.
func (u *Unifi) GetWANLoadBalancingStatus(sites []*Site) (*WANLoadBalancingStatus, error) {
	if len(sites) == 0 {
		return &WANLoadBalancingStatus{}, nil
	}

	site := sites[0]
	path := fmt.Sprintf(APIWANLoadBalancingStatusPath, site.Name)
	if strings.Contains(path, "%s") {
		return &WANLoadBalancingStatus{}, fmt.Errorf("WAN load-balancing path still contains %%s (site name empty): %q", path)
	}

	u.DebugLog("Fetching WAN load balancing status for site %s", site.Name)

	body, err := u.GetJSON(path)
	if err != nil {
		return nil, err
	}

	var response WANLoadBalancingStatus
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetWANISPStatus returns the ISP status for WAN interfaces.
func (u *Unifi) GetWANISPStatus(sites []*Site, wanNetworkgroup string) (*WANISPStatusDetailed, error) {
	if len(sites) == 0 {
		return &WANISPStatusDetailed{}, nil
	}

	site := sites[0]
	path := fmt.Sprintf(APIWANISPStatusPath, site.Name, wanNetworkgroup)
	if strings.Contains(path, "%s") {
		return &WANISPStatusDetailed{}, fmt.Errorf("WAN ISP status path still contains %%s: %q", path)
	}

	u.DebugLog("Fetching WAN ISP status for site %s, WAN group %s", site.Name, wanNetworkgroup)

	body, err := u.GetJSON(path)
	if err != nil {
		return nil, err
	}

	var response WANISPStatusDetailed
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetWANSLAs returns WAN SLA monitoring data.
func (u *Unifi) GetWANSLAs(sites []*Site) ([]*WANSLA, error) {
	data := []*WANSLA{}

	for _, site := range sites {
		path := fmt.Sprintf(APIWANSLAsPath, site.Name)
		if strings.Contains(path, "%s") {
			return data, fmt.Errorf("WAN SLAs path still contains %%s (site name empty): %q", path)
		}

		u.DebugLog("Fetching WAN SLAs for site %s", site.Name)

		body, err := u.GetJSON(path)
		if err != nil {
			return nil, err
		}

		var response []json.RawMessage
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}
		if len(response) == 0 {
			u.DebugLog("No WAN SLAs found for site %s", site.Name)
			continue
		}

		for _, raw := range response {
			var sla WANSLA
			if err := json.Unmarshal(raw, &sla); err != nil {
				u.DebugLog("Error parsing WAN SLA: %v", err)
				continue
			}
			data = append(data, &sla)
		}
	}

	return data, nil
}
