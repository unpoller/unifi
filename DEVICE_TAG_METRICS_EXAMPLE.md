# Device Tag Metrics Example

This document shows how device tags could be added as metrics to devices, using existing metric structures as examples.

## Current Device Tag API Response

```json
[
  {
    "_id": "2f8ca14d-0740-48b8-85cc-8ab1956e47cd",
    "member_device_macs": [
      "f4:e2:c6:f7:c3:49",
      "28:70:4e:85:a9:e9"
    ],
    "name": "tag"
  }
]
```

## unpoller_device_load_average_1 Metric Example

The `unpoller_device_load_average_1` metric shows the 1-minute load average for devices. This is the metric we'll use as an example.

### Current unpoller_device_load_average_1 Metric (without tags)

**Prometheus Format:**
```yaml
unpoller_device_load_average_1{site_name="default",name="UDM-Pro",mac="f4:e2:c6:f7:c3:49",model="UDMPRO",type="udm"} 0.15
```

**Grafana Query Example:**
```promql
unpoller_device_load_average_1{site_name=~"$Site", name=~"$Gateway"}
```

**InfluxDB Line Protocol:**
```
unpoller_device_load_average_1,site_name=default,name="UDM-Pro",mac=f4:e2:c6:f7:c3:49,model=UDMPRO,type=udm value=0.15
```

### Adding Tags to Device Structs

Add a `Tags` field to each device type (UAP, USW, UDM, etc.):

```go
type UAP struct {
    // ... existing fields ...
    Mac   string   `fake:"{macaddress}" json:"mac"`
    Name  string   `json:"name"`
    Model string   `json:"model"`
    Tags  []string `json:"tags"`  // Added: list of tag names this device belongs to
    // ... rest of fields ...
}
```

### Example Device with Tags

```go
device := &UAP{
    Mac:   "f4:e2:c6:f7:c3:49",
    Name:  "Office AP",
    Model: "U6-Pro",
    Tags:  []string{"tag", "office", "floor-2"},  // Tags assigned to this device
    // ... other fields ...
}
```

## unpoller_device_load_average_1 Metric with Device Tags

### Actual Metric Output (Prometheus Format)

**Current (without device tags):**
```yaml
unpoller_device_load_average_1{
  site_name="default",
  name="UDM-Pro",
  mac="f4:e2:c6:f7:c3:49",
  model="UDMPRO",
  type="udm"
} 0.15
```

**With Device Tags Added:**
```yaml
unpoller_device_load_average_1{
  site_name="default",
  name="UDM-Pro",
  mac="f4:e2:c6:f7:c3:49",
  model="UDMPRO",
  type="udm",
  tag="office",
  tag="core-gateway",
  tag="production"
} 0.15
```

### Available Metrics in Prometheus

When you query Prometheus, you'll see multiple time series for devices with tags. Each tag creates a separate series:

**Example: Device with 3 tags creates 3 series (one per tag):**
```yaml
# Series 1: tag="office"
unpoller_device_load_average_1{site_name="default",name="UDM-Pro",mac="f4:e2:c6:f7:c3:49",model="UDMPRO",type="udm",tag="office"} 0.15

# Series 2: tag="core-gateway"  
unpoller_device_load_average_1{site_name="default",name="UDM-Pro",mac="f4:e2:c6:f7:c3:49",model="UDMPRO",type="udm",tag="core-gateway"} 0.15

# Series 3: tag="production"
unpoller_device_load_average_1{site_name="default",name="UDM-Pro",mac="f4:e2:c6:f7:c3:49",model="UDMPRO",type="udm",tag="production"} 0.15
```

**Note:** Prometheus will show the same metric value (0.15) for each tag, allowing you to filter/group by any tag independently.

### Browsing Available Metrics

**In Prometheus UI (`/graph`):**
When you type `unpoller_device_load_average_1{` and press Enter, you'll see autocomplete suggestions including:
- `site_name`
- `name`
- `mac`
- `model`
- `type`
- `tag` ← **New: appears here**

**In Grafana Explore:**
1. Select metric: `unpoller_device_load_average_1`
2. Click "Select label" dropdown
3. You'll see all available labels including `tag`
4. Select `tag` and choose a value (e.g., `office`)
5. Query becomes: `unpoller_device_load_average_1{tag="office"}`

### Example: Multiple Devices with Tags

**Device 1: UDM-Pro with tags ["office", "core-gateway", "production"]**
```yaml
unpoller_device_load_average_1{site_name="default",name="UDM-Pro",mac="f4:e2:c6:f7:c3:49",tag="office"} 0.15
unpoller_device_load_average_1{site_name="default",name="UDM-Pro",mac="f4:e2:c6:f7:c3:49",tag="core-gateway"} 0.15
unpoller_device_load_average_1{site_name="default",name="UDM-Pro",mac="f4:e2:c6:f7:c3:49",tag="production"} 0.15
```

**Device 2: Switch with tags ["office", "core-switch"]**
```yaml
unpoller_device_load_average_1{site_name="default",name="Main Switch",mac="28:70:4e:85:a9:e9",tag="office"} 0.08
unpoller_device_load_average_1{site_name="default",name="Main Switch",mac="28:70:4e:85:a9:e9",tag="core-switch"} 0.08
```

**Query Results:**
- `unpoller_device_load_average_1{tag="office"}` → Returns both devices (UDM-Pro: 0.15, Switch: 0.08)
- `unpoller_device_load_average_1{tag="production"}` → Returns only UDM-Pro (0.15)
- `avg by (tag) (unpoller_device_load_average_1{tag="office"})` → Returns average: 0.115

### Query Examples

#### Basic Queries

**Current Query (without device tags):**
```promql
unpoller_device_load_average_1{site_name=~"$Site", name=~"$Gateway"}
```

**With Device Tags - Filter by single tag:**
```promql
unpoller_device_load_average_1{site_name=~"$Site", name=~"$Gateway", tag="office"}
```

**With Device Tags - Filter by multiple tags (AND):**
```promql
unpoller_device_load_average_1{site_name=~"$Site", name=~"$Gateway", tag="office", tag="production"}
```

**With Device Tags - Filter by tag pattern (regex):**
```promql
unpoller_device_load_average_1{site_name=~"$Site", tag=~"office.*"}
```

#### Aggregation Queries

**Average load grouped by tag:**
```promql
avg by (tag) (unpoller_device_load_average_1{site_name=~"$Site"})
```

**Average load grouped by tag and device name:**
```promql
avg by (tag, name) (unpoller_device_load_average_1{site_name=~"$Site"})
```

**Count devices per tag:**
```promql
count by (tag) (unpoller_device_load_average_1{site_name=~"$Site"})
```

**Max load per tag:**
```promql
max by (tag) (unpoller_device_load_average_1{site_name=~"$Site"})
```

#### Filtering Queries

**Find devices with high load and specific tag:**
```promql
unpoller_device_load_average_1{site_name=~"$Site", tag="production"} > 1.0
```

**Compare load between tags:**
```promql
unpoller_device_load_average_1{site_name=~"$Site", tag=~"office|production"}
```

**Devices with any tag:**
```promql
unpoller_device_load_average_1{site_name=~"$Site", tag=~".+"}
```

**Devices without tags:**
```promql
unpoller_device_load_average_1{site_name=~"$Site"} unless unpoller_device_load_average_1{site_name=~"$Site", tag=~".+"}
```

### Available Labels in Grafana

**Current Available Labels (without device tags):**
- `site_name`
- `name`
- `mac`
- `model`
- `type`

**With Device Tags Added - New Labels Available:**
- `site_name`
- `name`
- `mac`
- `model`
- `type`
- `tag` ← **New: Device tags appear here**

### Using Tags in Grafana

**In Grafana's Label Filters:**
1. Click "Add filter" in the Label filters section
2. Select `tag` from the "Select label" dropdown
3. Choose operator (`=`, `!=`, `=~`, `!~`)
4. Enter tag value (e.g., `office`)
5. Result: `unpoller_device_load_average_1{site_name=~"$Site", name=~"$Gateway", tag="office"}`

**Example Label Filter Combinations:**
- `tag = "office"` → Shows only devices tagged "office"
- `tag =~ "office|production"` → Shows devices with either tag
- `tag != "test"` → Excludes test devices
- `tag =~ ".*"` → Shows all tagged devices

### InfluxDB Line Protocol

**Current (without device tags):**
```
unpoller_device_load_average_1,site_name=default,name="UDM-Pro",mac=f4:e2:c6:f7:c3:49,model=UDMPRO,type=udm value=0.15
```

**With Device Tags Added:**
```
unpoller_device_load_average_1,site_name=default,name="UDM-Pro",mac=f4:e2:c6:f7:c3:49,model=UDMPRO,type=udm,tag=office,tag=core-gateway,tag=production value=0.15
```

### Example: Switch Device

**Current (without device tags):**
```yaml
unpoller_device_load_average_1{
  site_name="default",
  name="Main Switch",
  mac="28:70:4e:85:a9:e9",
  model="USW-24-POE",
  type="usw"
} 0.08
```

**With Device Tags Added:**
```yaml
unpoller_device_load_average_1{
  site_name="default",
  name="Main Switch",
  mac="28:70:4e:85:a9:e9",
  model="USW-24-POE",
  type="usw",
  tag="office",
  tag="core-switch"
} 0.08
```

## Real-World Usage Examples

### Grafana Dashboard Queries

**Current Query (without device tags):**
```promql
unpoller_device_load_average_1{site_name=~"$Site", name=~"$Gateway"}
```

**Filter by device tag:**
```promql
unpoller_device_load_average_1{site_name=~"$Site", name=~"$Gateway", tag="office"}
```

**Average load by tag:**
```promql
avg by (tag) (unpoller_device_load_average_1{site_name=~"$Site"})
```

**Filter devices by multiple tags:**
```promql
unpoller_device_load_average_1{site_name=~"$Site", tag="office", tag="production"}
```

**Compare load averages across tags:**
```promql
avg by (tag, name) (unpoller_device_load_average_1{site_name=~"$Site"})
```

**Find high load devices with specific tag:**
```promql
unpoller_device_load_average_1{site_name=~"$Site", tag="production"} > 1.0
```

**Using Grafana Label Filters:**
- Use the "Label filters" section
- Select `tag` from the "Select label" dropdown
- Choose `=` operator and enter `office` as the value
- Result: `unpoller_device_load_average_1{site_name=~"$Site", name=~"$Gateway", tag="office"}`

### InfluxQL Queries

```sql
-- Average load for devices with tag "office"
SELECT mean(value) FROM unpoller_device_load_average_1 WHERE tag = 'office'

-- Average load grouped by tag
SELECT mean(value) FROM unpoller_device_load_average_1 GROUP BY tag

-- Find high load devices with specific tag
SELECT mean(value) FROM unpoller_device_load_average_1 
WHERE tag = 'production' AND value > 1.0 
GROUP BY name, tag

-- Compare load averages across tags
SELECT mean(value) FROM unpoller_device_load_average_1 
GROUP BY tag, name
```

## Option 2: Separate Device Tags Metric

Create a separate metric specifically for device tags:

```go
type DeviceTagMetric struct {
    SiteName   string   `json:"site_name"`
    DeviceMac  string   `json:"device_mac"`
    DeviceName string   `json:"device_name"`
    DeviceType string   `json:"device_type"`  // "uap", "usw", "udm", etc.
    Tags       []string `json:"tags"`
    SourceName string   `json:"source_name"`
}
```

### Metrics Output Example

**InfluxDB Line Protocol:**
```
unifi_device_tags,site=default,mac=f4:e2:c6:f7:c3:49,name="Office AP",type=uap,tag=tag,tag=office,tag=floor-2 tag_count=3i
```

**Prometheus:**
```yaml
unifi_device_tag_count{site="default", mac="f4:e2:c6:f7:c3:49", name="Office AP", type="uap", tag="tag"} 1
unifi_device_tag_count{site="default", mac="f4:e2:c6:f7:c3:49", name="Office AP", type="uap", tag="office"} 1
unifi_device_tag_count{site="default", mac="f4:e2:c6:f7:c3:49", name="Office AP", type="uap", tag="floor-2"} 1
```

## Option 3: Tag-Based Aggregation Metrics

Create metrics that aggregate by tag:

```go
type TagAggregationMetric struct {
    SiteName   string   `json:"site_name"`
    TagName    string   `json:"tag_name"`
    DeviceCount int     `json:"device_count"`
    DeviceMacs  []string `json:"device_macs"`
    SourceName string   `json:"source_name"`
}
```

### Metrics Output Example

**InfluxDB:**
```
unifi_tag_aggregation,site=default,tag=tag device_count=2i
unifi_tag_aggregation,site=default,tag=office device_count=5i
unifi_tag_aggregation,site=default,tag=floor-2 device_count=3i
```

**Prometheus:**
```yaml
unifi_tag_device_count{site="default", tag="tag"} 2
unifi_tag_device_count{site="default", tag="office"} 5
unifi_tag_device_count{site="default", tag="floor-2"} 3
```

## Recommended Approach

**Option 1** (adding Tags field to device structs) is recommended because:

1. **Simple Integration**: Tags become part of the device data structure
2. **Query Flexibility**: Can filter/group devices by tags in monitoring systems
3. **Consistent Pattern**: Similar to how client tags work (`ClientHistory.Tags`)
4. **Backward Compatible**: Empty tags array for devices without tags

### Implementation Example

```go
// Add to UAP struct
type UAP struct {
    // ... existing fields ...
    Mac   string   `fake:"{macaddress}" json:"mac"`
    Tags  []string `json:"tags"`  // Populated from device-tags API
    // ... rest of fields ...
}

// When fetching devices, enrich with tags:
func (u *Unifi) GetDevices(sites []*Site) (*Devices, error) {
    devices, err := u.fetchDevices(sites)
    if err != nil {
        return nil, err
    }
    
    // Fetch device tags for each site
    for _, site := range sites {
        tags, err := u.GetDeviceTags(site)
        if err == nil {
            u.enrichDevicesWithTags(devices, tags)
        }
    }
    
    return devices, nil
}

func (u *Unifi) enrichDevicesWithTags(devices *Devices, tags []DeviceTag) {
    // Create a map of MAC -> tag names
    macToTags := make(map[string][]string)
    for _, tag := range tags {
        for _, mac := range tag.MemberDeviceMacs {
            macToTags[mac] = append(macToTags[mac], tag.Name)
        }
    }
    
    // Enrich all devices with their tags
    for _, device := range devices.UAPs {
        device.Tags = macToTags[device.Mac]
    }
    for _, device := range devices.USWs {
        device.Tags = macToTags[device.Mac]
    }
    // ... etc for other device types
}
```

### Usage in Monitoring Queries

With tags as device fields, you can:

**Grafana Dashboard:**
- Filter devices by tag: `tag = "office"`
- Group metrics by tag: `GROUP BY tag`
- Create tag-based alerts: `WHERE tag = "critical"`

**PromQL:**
```promql
# Sum bytes for all devices with tag "office"
sum(unifi_device_bytes{tag="office"})

# Count devices per tag
count by (tag) (unifi_device_bytes)
```

**InfluxQL:**
```sql
SELECT mean(bytes) FROM unifi_device WHERE tag = 'office'
SELECT count(*) FROM unifi_device GROUP BY tag
```
