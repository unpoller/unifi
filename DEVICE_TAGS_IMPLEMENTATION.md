# Device Tags Implementation Plan

## Current State

**Query Results:**
- No "test" tag found in the current controller
- Found 1 tag: "tag" with 2 devices:
  - `f4:e2:c6:f7:c3:49`
  - `28:70:4e:85:a9:e9`

**API Endpoint:**
- `GET /proxy/network/v2/api/site/{site}/device-tags`
- Returns: `[{ "_id": "...", "name": "tag", "member_device_macs": ["mac1", "mac2"] }]`

## Conceptual Implementation

### Step 1: Add DeviceTag Type

Create a new type to represent device tags:

```go
// DeviceTag represents a device tag from the UniFi API
type DeviceTag struct {
    ID              string   `json:"_id"`
    Name            string   `json:"name"`
    MemberDeviceMacs []string `json:"member_device_macs"`
}
```

### Step 2: Add Tags Field to Device Structs

Add a `Tags` field to each device type. Since all device types share common fields through embedding, we can add it to a common interface or each struct:

**Option A: Add to each device struct individually**
```go
type UAP struct {
    // ... existing fields ...
    Mac   string   `fake:"{macaddress}" json:"mac"`
    Tags  []string `json:"tags"`  // Populated from device-tags API
    // ... rest of fields ...
}

type USW struct {
    // ... existing fields ...
    Mac   string   `fake:"{macaddress}" json:"mac"`
    Tags  []string `json:"tags"`  // Populated from device-tags API
    // ... rest of fields ...
}

// Repeat for UDM, UXG, USG, PDU, UBB, UCI
```

**Option B: Create a common embedded struct (if possible)**
```go
type DeviceBase struct {
    Mac  string   `json:"mac"`
    Tags []string `json:"tags"`
}

// Then embed in each device type (requires refactoring)
```

### Step 3: Add GetDeviceTags Method

```go
// GetDeviceTags returns all device tags for a site
func (u *Unifi) GetDeviceTags(site *Site) ([]*DeviceTag, error) {
    if site == nil || site.Name == "" {
        return nil, ErrNoSiteProvided
    }

    u.DebugLog("Polling Controller for Device Tags, site %s", site.SiteName)

    path := fmt.Sprintf("/proxy/network/v2/api/site/%s/device-tags", site.Name)
    
    var tags []*DeviceTag
    if err := u.GetData(path, &tags); err != nil {
        return nil, fmt.Errorf("fetching device tags: %w", err)
    }

    return tags, nil
}
```

### Step 4: Enrich Devices with Tags

Modify `GetDevices` to optionally enrich devices with tags:

```go
// GetDevices returns a response full of devices' data from the UniFi Controller.
func (u *Unifi) GetDevices(sites []*Site) (*Devices, error) {
    devices := new(Devices)

    for _, site := range sites {
        var response struct {
            Data []json.RawMessage `json:"data"`
        }

        devicePath := fmt.Sprintf(APIDevicePath, site.Name)
        if err := u.GetData(devicePath, &response); err != nil {
            return nil, err
        }

        loopDevices := u.parseDevices(response.Data, site)
        devices.UAPs = append(devices.UAPs, loopDevices.UAPs...)
        devices.USGs = append(devices.USGs, loopDevices.USGs...)
        devices.USWs = append(devices.USWs, loopDevices.USWs...)
        devices.UDMs = append(devices.UDMs, loopDevices.UDMs...)
        devices.UXGs = append(devices.UXGs, loopDevices.UXGs...)
        devices.PDUs = append(devices.PDUs, loopDevices.PDUs...)
        devices.UBBs = append(devices.UBBs, loopDevices.UBBs...)
        devices.UCIs = append(devices.UCIs, loopDevices.UCIs...)

        // Enrich devices with tags
        if err := u.enrichDevicesWithTags(devices, site); err != nil {
            u.ErrorLog("Failed to enrich devices with tags for site %s: %v", site.SiteName, err)
            // Don't fail the whole request if tags fail
        }
    }

    return devices, nil
}

// enrichDevicesWithTags adds tag information to devices based on their MAC addresses
func (u *Unifi) enrichDevicesWithTags(devices *Devices, site *Site) error {
    tags, err := u.GetDeviceTags(site)
    if err != nil {
        return err
    }

    // Create a map of MAC -> tag names
    macToTags := make(map[string][]string)
    for _, tag := range tags {
        for _, mac := range tag.MemberDeviceMacs {
            macToTags[mac] = append(macToTags[mac], tag.Name)
        }
    }

    // Enrich all device types with their tags
    for _, device := range devices.UAPs {
        device.Tags = macToTags[device.Mac]
    }
    for _, device := range devices.USWs {
        device.Tags = macToTags[device.Mac]
    }
    for _, device := range devices.UDMs {
        device.Tags = macToTags[device.Mac]
    }
    for _, device := range devices.USGs {
        device.Tags = macToTags[device.Mac]
    }
    for _, device := range devices.UXGs {
        device.Tags = macToTags[device.Mac]
    }
    for _, device := range devices.PDUs {
        device.Tags = macToTags[device.Mac]
    }
    for _, device := range devices.UBBs {
        device.Tags = macToTags[device.Mac]
    }
    for _, device := range devices.UCIs {
        device.Tags = macToTags[device.Mac]
    }

    return nil
}
```

### Step 5: Add API Path Constant

Add the device tags API path to `types.go`:

```go
const (
    // ... existing paths ...
    APIDeviceTagsPath string = "/proxy/network/v2/api/site/%s/device-tags"
)
```

### Step 6: Testing

Create tests to verify:
1. `GetDeviceTags` returns correct tags
2. Devices are enriched with tags correctly
3. Devices without tags have empty `Tags` slice
4. Multiple tags per device work correctly

## Implementation Considerations

### Performance
- Tags are fetched once per site, not per device
- Tag enrichment happens after device parsing
- If tag fetching fails, devices are still returned (graceful degradation)

### Backward Compatibility
- `Tags` field defaults to `nil`/empty slice for devices without tags
- Existing code continues to work without changes
- Tags are optional - if the API endpoint doesn't exist or fails, devices still work

### MAC Address Matching
- Tags use MAC addresses to match devices
- MAC addresses should be normalized (lowercase, colons) for matching
- Consider case-insensitive matching

### Example Usage

```go
// Fetch devices (automatically enriched with tags)
devices, err := uni.GetDevices(sites)
if err != nil {
    log.Fatal(err)
}

// Access tags on devices
for _, ap := range devices.UAPs {
    fmt.Printf("AP %s has tags: %v\n", ap.Name, ap.Tags)
    if contains(ap.Tags, "test") {
        fmt.Printf("  -> This is a test device!\n")
    }
}

// Or fetch tags separately
tags, err := uni.GetDeviceTags(site)
if err != nil {
    log.Fatal(err)
}

for _, tag := range tags {
    if tag.Name == "test" {
        fmt.Printf("Tag 'test' has %d devices\n", len(tag.MemberDeviceMacs))
    }
}
```

## Files to Modify

1. **types.go**: Add `DeviceTag` type and `APIDeviceTagsPath` constant
2. **uap.go**: Add `Tags []string` field to `UAP` struct
3. **usw.go**: Add `Tags []string` field to `USW` struct
4. **udm.go**: Add `Tags []string` field to `UDM` struct
5. **usg.go**: Add `Tags []string` field to `USG` struct
6. **uxg.go**: Add `Tags []string` field to `UXG` struct
7. **pdu.go**: Add `Tags []string` field to `PDU` struct
8. **ubb.go**: Add `Tags []string` field to `UBB` struct
9. **uci.go**: Add `Tags []string` field to `UCI` struct
10. **devices.go**: Add `GetDeviceTags` method and `enrichDevicesWithTags` helper
11. **devices_test.go**: Add tests for tag enrichment

## Next Steps

1. Implement `DeviceTag` type
2. Add `GetDeviceTags` method
3. Add `Tags` field to device structs
4. Implement enrichment logic
5. Add tests
6. Update documentation
