package unifi

import (
	"encoding/json"
	"fmt"
)

// GetFirewallPolicies returns firewall policies for all provided sites.
// Uses the v2 API endpoint: GET /proxy/network/v2/api/site/{site}/firewall-policies
func (u *Unifi) GetFirewallPolicies(sites []*Site) ([]*FirewallPolicy, error) {
	policies := make([]*FirewallPolicy, 0)

	for _, site := range sites {
		path := fmt.Sprintf(APIFirewallPoliciesPath, site.Name)

		body, err := u.GetJSON(path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch firewall policies for site %s: %w", site.SiteName, err)
		}

		var raw []*FirewallPolicy
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, fmt.Errorf("failed to parse firewall policies for site %s: %w", site.SiteName, err)
		}

		for _, policy := range raw {
			if policy == nil {
				continue
			}

			policy.SiteName = site.SiteName
			policy.SourceName = u.URL
			policies = append(policies, policy)
		}
	}

	return policies, nil
}

// FirewallPolicyEndpoint represents the source or destination of a firewall rule.
type FirewallPolicyEndpoint struct {
	MatchOppositePorts bool   `json:"match_opposite_ports"`
	MatchingTarget     string `json:"matching_target"`
	PortMatchingType   string `json:"port_matching_type"`
	ZoneID             string `json:"zone_id"`
}

// FirewallPolicySchedule represents when a firewall rule is active.
type FirewallPolicySchedule struct {
	Mode string `json:"mode"`
}

// FirewallPolicy represents a firewall policy rule from the UniFi controller.
type FirewallPolicy struct {
	ID                  string                 `json:"_id"`
	Action              string                 `json:"action"`               // ALLOW, BLOCK, REJECT
	ConnectionStateType string                 `json:"connection_state_type"` // ALL, CUSTOM
	ConnectionStates    []string               `json:"connection_states"`
	CreateAllowRespond  FlexBool               `json:"create_allow_respond"`
	Destination         FirewallPolicyEndpoint `json:"destination"`
	Enabled             FlexBool               `json:"enabled"`
	ICMPTypename        string                 `json:"icmp_typename"`
	ICMPv6Typename      string                 `json:"icmp_v6_typename"`
	Index               FlexInt                `json:"index"`
	IPVersion           string                 `json:"ip_version"` // BOTH, IPV4, IPV6
	Logging             FlexBool               `json:"logging"`
	MatchIPSec          FlexBool               `json:"match_ip_sec"`
	Name                string                 `json:"name"`
	Predefined          FlexBool               `json:"predefined"`
	Protocol            string                 `json:"protocol"` // all, tcp, udp, icmp, etc.
	Schedule            FirewallPolicySchedule `json:"schedule"`
	Source              FirewallPolicyEndpoint `json:"source"`

	SiteName   string `json:"-"`
	SourceName string `json:"-"`
}
