# WAN Metrics Implementation

This document describes the WAN metrics that have been implemented in the unifi library and unpoller.

## Overview

WAN metrics provide comprehensive monitoring of WAN interfaces including configuration, statistics, performance, and service provider information.

## API Endpoints

### Implemented Endpoints

1. **`/proxy/network/v2/api/site/{site}/wan/enriched-configuration`**
   - Complete WAN configuration with statistics
   - Provider capabilities (download/upload speeds)
   - Load balancing configuration
   - Failover priorities
   - Uptime and peak usage statistics

2. **`/proxy/network/v2/api/site/{site}/wan/load-balancing/status`**
   - Current load balancing mode (FAILOVER_ONLY, DISTRIBUTED)
   - Active WAN interfaces with priorities and weights

3. **`/proxy/network/v2/api/site/{site}/wan/{networkgroup}/isp-status`**
   - Detailed ISP connection status
   - Speedtest historical data
   - Uplink statistics (latency, throughput)
   - Connection warnings (downtimes, high latencies, packet loss)

4. **`/proxy/network/v2/api/site/{site}/wan-slas`**
   - WAN SLA monitoring data (future use)

## Prometheus Metrics

All metrics are prefixed with `unpoller_wan_`.

### Configuration Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `wan_failover_priority` | Gauge | WAN failover priority (lower number = higher priority) | wan_id, wan_name, wan_networkgroup, wan_type, wan_load_balance_type, site_name, source |
| `wan_load_balance_weight` | Gauge | WAN load balancing weight | " |
| `wan_provider_download_kbps` | Gauge | Configured ISP download speed in Kbps | " |
| `wan_provider_upload_kbps` | Gauge | Configured ISP upload speed in Kbps | " |
| `wan_smartq_enabled` | Gauge | SmartQueue QoS enabled (1) or disabled (0) | " |
| `wan_magic_enabled` | Gauge | Magic WAN enabled (1) or disabled (0) | " |
| `wan_vlan_enabled` | Gauge | VLAN enabled for WAN (1) or disabled (0) | " |

### Statistics Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `wan_uptime_percentage` | Gauge | WAN uptime percentage | wan_id, wan_name, wan_networkgroup, wan_type, wan_load_balance_type, site_name, source |
| `wan_peak_download_percent` | Gauge | Peak download usage as % of configured capacity | " |
| `wan_peak_upload_percent` | Gauge | Peak upload usage as % of configured capacity | " |
| `wan_max_rx_bytes_rate` | Gauge | Maximum receive bytes rate | " |
| `wan_max_tx_bytes_rate` | Gauge | Maximum transmit bytes rate | " |

### Service Provider Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `wan_service_provider_asn` | Gauge | Service provider autonomous system number | wan_id, wan_name, wan_networkgroup, isp_name, isp_city, site_name, source |

### Metadata Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `wan_creation_timestamp` | Gauge | WAN configuration creation timestamp | wan_id, wan_name, wan_networkgroup, wan_type, wan_load_balance_type, site_name, source |

## Example Queries

### Basic WAN Status

```promql
# Show all WAN interfaces with their failover priority
unpoller_wan_failover_priority

# WAN uptime percentage by interface
unpoller_wan_uptime_percentage{wan_name="Westelcom"}

# Current peak utilization
unpoller_wan_peak_download_percent
unpoller_wan_peak_upload_percent
```

### Load Balancing

```promql
# Show load balancing weights for all WAN interfaces
unpoller_wan_load_balance_weight

# Filter by load balance type
unpoller_wan_load_balance_weight{wan_load_balance_type="weighted"}
unpoller_wan_load_balance_weight{wan_load_balance_type="failover-only"}
```

### Performance Monitoring

```promql
# Peak throughput rates
rate(unpoller_wan_max_rx_bytes_rate[5m])
rate(unpoller_wan_max_tx_bytes_rate[5m])

# Utilization as percentage of configured capacity
unpoller_wan_peak_download_percent / 100
unpoller_wan_peak_upload_percent / 100
```

### Alerting Examples

```promql
# Alert if WAN uptime drops below 99%
unpoller_wan_uptime_percentage < 99

# Alert if utilization exceeds 80% of configured capacity
unpoller_wan_peak_download_percent > 80
unpoller_wan_peak_upload_percent > 80

# Alert on WAN failover (primary WAN down)
unpoller_wan_uptime_percentage{wan_failover_priority="1"} < 100
```

## Data Structures

### WANEnrichedConfiguration
Contains complete WAN configuration, statistics, and service provider details.

**Fields:**
- `Configuration`: WAN network configuration (type, load balancing, failover, etc.)
- `Details`: Service provider information (ISP name, city, ASN)
- `Statistics`: Uptime and peak usage statistics

### WANLoadBalancingStatus
Current load balancing status and interface configuration.

**Fields:**
- `Mode`: Load balancing mode (FAILOVER_ONLY, DISTRIBUTED)
- `WANInterfaces`: List of WAN interfaces with mode, priority, and weight

### WANISPStatusDetailed
Detailed ISP status including speedtest history and uplink monitoring.

**Fields:**
- `SpeedtestHistorical`: Historical speedtest results
- `UplinkStatus`: Real-time uplink statistics (latency, throughput)
- `ConnectionWarnings`: Downtimes, high latencies, packet loss

## Example Data

### WAN Enriched Configuration
```json
{
  "configuration": {
    "_id": "6711e729be659a441ec8e57c",
    "name": "Westelcom",
    "wan_networkgroup": "WAN",
    "wan_type": "dhcp",
    "wan_failover_priority": 1,
    "wan_load_balance_type": "weighted",
    "wan_load_balance_weight": 99,
    "wan_provider_capabilities": {
      "download_kilobits_per_second": 1000000,
      "upload_kilobits_per_second": 1000000
    }
  },
  "details": {
    "service_provider": {
      "asn": 11722,
      "city": "Potsdam",
      "name": "Westelcom Internet"
    }
  },
  "statistics": {
    "uptime_percentage": 100.0,
    "peak_usage": {
      "download_percentage": 11.1,
      "upload_percentage": 1.3
    }
  }
}
```

### Load Balancing Status
```json
{
  "mode": "FAILOVER_ONLY",
  "wan_interfaces": [
    {
      "name": "Westelcom",
      "wan_networkgroup": "WAN",
      "mode": "DISTRIBUTED",
      "priority": 1,
      "weight": 99
    },
    {
      "name": "Internet 2",
      "wan_networkgroup": "WAN2",
      "mode": "FAILOVER_ONLY",
      "priority": 2,
      "weight": 1
    }
  ]
}
```

## Future Enhancements

Additional WAN-related endpoints that could be added:

1. **WAN SLA Metrics** - Currently the `/wan-slas` endpoint returns empty, but future controller versions may provide:
   - Latency measurements
   - Packet loss percentages
   - Jitter metrics
   - Per-target monitoring results

2. **Real-time Uplink Status** - From `/stat/health`:
   - Per-target availability and latency
   - DNS and ICMP monitoring results
   - Connection warnings and alerts

3. **Speedtest Results** - Historical speedtest data with:
   - Download/upload Mbps
   - Latency measurements
   - Timestamps for tracking trends

## Implementation Notes

- WAN metrics are fetched from the `/proxy/network/v2/api/site/{site}/wan/enriched-configuration` endpoint
- The implementation handles multi-WAN setups with failover and load balancing
- Metrics are exported per-WAN interface with appropriate labels for filtering and aggregation
- Service provider information (ISP name, ASN) is included for network troubleshooting
- Peak usage percentages are calculated relative to configured provider capabilities
