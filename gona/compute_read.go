package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type BootProfile struct {
	BootProfileID int    `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	Builder       string `json:"builder"`
	Kernel        string `json:"kernel"`
	Boot          string `json:"boot"`
	Serial        string `json:"serial"`
	DiskRepresent string `json:"disk_represent"`
	// An int, confirmed live: the API sends 0 or a template id, not a string.
	ImageTemplate int     `json:"image_template"`
	LastUpdated   string  `json:"last_updated"`
	Extra         *string `json:"extra"`
	VNCDisplay    *string `json:"vncdisplay"`
	DiskRoot      *string `json:"disk_root"`
	Bootloader    *string `json:"bootloader"`
	Ramdisk       *string `json:"ramdisk"`
	Initrd        *string `json:"initrd"`
	Created       *string `json:"created"`
	PAE           int     `json:"pae"`
	ACPI          int     `json:"acpi"`
	APIC          int     `json:"apic"`
	XLocaltime    int     `json:"xlocaltime"`
	SDL           int     `json:"sdl"`
	VNC           int     `json:"vnc"`
	VNCConsole    int     `json:"vncconsole"`
	VNCUnused     int     `json:"vncunused"`
	Hide          int     `json:"hide"`
	KVM           int     `json:"kvm"`
}

type ServerDisk struct{}

type DedicatedIDName struct {
	ID   int
	Name string
}

type DedicatedOSProfile struct {
	OSID              int
	Name              string
	GroupName         string
	Tags              []string
	DiskLayouts       []DedicatedIDName
	Scripts           []DedicatedIDName
	DefaultDiskLayout int
	DefaultScripts    []int
	AllowSSHKeys      int
	SetRootPassword   int
	RescueImage       int
	Public            int
	Enabled           int
	Created           string
	LastUpdated       string
	ProfileID         int
	Arch              string
	Flavor            string
	LocationID        *int
}

type DedicatedRescueOS struct {
	OSID              int
	Name              string
	GroupName         string
	Tags              []string
	DiskLayouts       []DedicatedIDName
	Scripts           []DedicatedIDName
	DefaultDiskLayout int
	DefaultScripts    []int
	AllowSSHKeys      int
	SetRootPassword   int
	RescueImage       int
	Public            int
	Enabled           int
	Created           string
	LastUpdated       string
	ProfileID         int
	Arch              string
	Flavor            string
	LocationID        *int
}

type DedicatedDiskLayout struct {
	LayoutID int    `json:"id"`
	Name     string `json:"name"`
	Profile  string `json:"profile"`
	MinDisks int    `json:"min_disks"`
}

func (c *Client) GetBootProfiles() ([]BootProfile, error) {
	var profiles []BootProfile
	if err := c.get(context.Background(), "cloud/boot-profiles", &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (c *Client) GetServerDisks(mbpkgid int) ([]ServerDisk, error) {
	var disks []ServerDisk
	if err := c.get(context.Background(), fmt.Sprintf("cloud/disks/%d", mbpkgid), &disks); err != nil {
		return nil, err
	}
	return disks, nil
}

func (c *Client) GetDedicatedOSProfiles() ([]DedicatedOSProfile, error) {
	var profiles []DedicatedOSProfile
	if err := c.get(context.Background(), "dedicated/os", &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (c *Client) GetDedicatedRescueOS() ([]DedicatedRescueOS, error) {
	var profiles []DedicatedRescueOS
	if err := c.get(context.Background(), "dedicated/os/rescue-system-list", &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (c *Client) GetDedicatedDiskLayouts(osID int) ([]DedicatedDiskLayout, error) {
	var layouts []DedicatedDiskLayout
	if err := c.get(context.Background(), fmt.Sprintf("dedicated/disklayouts/%d", osID), &layouts); err != nil {
		return nil, err
	}
	return layouts, nil
}

func (p *DedicatedOSProfile) UnmarshalJSON(data []byte) error {
	type wireProfile struct {
		OSID              int             `json:"id"`
		Name              string          `json:"name"`
		GroupName         string          `json:"group_name"`
		Tags              []string        `json:"tags"`
		DiskLayouts       json.RawMessage `json:"disklayouts"`
		Scripts           json.RawMessage `json:"scripts"`
		DefaultDiskLayout *int            `json:"default_disklayout"`
		DefaultScripts    []int           `json:"default_scripts"`
		AllowSSHKeys      int             `json:"allow_ssh_keys"`
		SetRootPassword   int             `json:"set_root_password"`
		RescueImage       int             `json:"rescue_image"`
		Public            int             `json:"public"`
		Enabled           int             `json:"enabled"`
		Created           string          `json:"created"`
		LastUpdated       string          `json:"last_updated"`
		ProfileID         int             `json:"profile_id"`
		Arch              string          `json:"arch"`
		Flavor            string          `json:"flavor"`
		LocationID        *int            `json:"location_id"`
	}

	var wire wireProfile
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	diskLayouts, err := parseDedicatedIDNames(wire.DiskLayouts)
	if err != nil {
		return fmt.Errorf("parse dedicated os disklayouts: %w", err)
	}
	scripts, err := parseDedicatedIDNames(wire.Scripts)
	if err != nil {
		return fmt.Errorf("parse dedicated os scripts: %w", err)
	}

	*p = DedicatedOSProfile{
		OSID:              wire.OSID,
		Name:              wire.Name,
		GroupName:         wire.GroupName,
		Tags:              wire.Tags,
		DiskLayouts:       diskLayouts,
		Scripts:           scripts,
		DefaultDiskLayout: intOrZero(wire.DefaultDiskLayout),
		DefaultScripts:    wire.DefaultScripts,
		AllowSSHKeys:      wire.AllowSSHKeys,
		SetRootPassword:   wire.SetRootPassword,
		RescueImage:       wire.RescueImage,
		Public:            wire.Public,
		Enabled:           wire.Enabled,
		Created:           wire.Created,
		LastUpdated:       wire.LastUpdated,
		ProfileID:         wire.ProfileID,
		Arch:              wire.Arch,
		Flavor:            wire.Flavor,
		LocationID:        wire.LocationID,
	}
	return nil
}

func (p *DedicatedRescueOS) UnmarshalJSON(data []byte) error {
	type wireProfile struct {
		OSID              int             `json:"id"`
		Name              string          `json:"name"`
		GroupName         string          `json:"group_name"`
		Tags              []string        `json:"tags"`
		DiskLayouts       json.RawMessage `json:"disklayouts"`
		Scripts           json.RawMessage `json:"scripts"`
		DefaultDiskLayout *int            `json:"default_disklayout"`
		DefaultScripts    []int           `json:"default_scripts"`
		AllowSSHKeys      int             `json:"allow_ssh_keys"`
		SetRootPassword   int             `json:"set_root_password"`
		RescueImage       int             `json:"rescue_image"`
		Public            int             `json:"public"`
		Enabled           int             `json:"enabled"`
		Created           string          `json:"created"`
		LastUpdated       string          `json:"last_updated"`
		ProfileID         int             `json:"profile_id"`
		Arch              string          `json:"arch"`
		Flavor            string          `json:"flavor"`
		LocationID        *int            `json:"location_id"`
	}

	var wire wireProfile
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	diskLayouts, err := parseDedicatedIDNames(wire.DiskLayouts)
	if err != nil {
		return fmt.Errorf("parse dedicated rescue os disklayouts: %w", err)
	}
	scripts, err := parseDedicatedIDNames(wire.Scripts)
	if err != nil {
		return fmt.Errorf("parse dedicated rescue os scripts: %w", err)
	}

	*p = DedicatedRescueOS{
		OSID:              wire.OSID,
		Name:              wire.Name,
		GroupName:         wire.GroupName,
		Tags:              wire.Tags,
		DiskLayouts:       diskLayouts,
		Scripts:           scripts,
		DefaultDiskLayout: intOrZero(wire.DefaultDiskLayout),
		DefaultScripts:    wire.DefaultScripts,
		AllowSSHKeys:      wire.AllowSSHKeys,
		SetRootPassword:   wire.SetRootPassword,
		RescueImage:       wire.RescueImage,
		Public:            wire.Public,
		Enabled:           wire.Enabled,
		Created:           wire.Created,
		LastUpdated:       wire.LastUpdated,
		ProfileID:         wire.ProfileID,
		Arch:              wire.Arch,
		Flavor:            wire.Flavor,
		LocationID:        wire.LocationID,
	}
	return nil
}

func intOrZero(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func parseDedicatedIDNames(raw json.RawMessage) ([]DedicatedIDName, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return []DedicatedIDName{}, nil
	}

	var byID map[string]string
	if err := json.Unmarshal(raw, &byID); err == nil {
		keys := make([]string, 0, len(byID))
		for key := range byID {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			left, leftErr := strconv.Atoi(keys[i])
			right, rightErr := strconv.Atoi(keys[j])
			if leftErr == nil && rightErr == nil {
				return left < right
			}
			return keys[i] < keys[j]
		})

		items := make([]DedicatedIDName, 0, len(keys))
		for _, key := range keys {
			id, err := strconv.Atoi(key)
			if err != nil {
				return nil, fmt.Errorf("invalid id key %q: %w", key, err)
			}
			items = append(items, DedicatedIDName{ID: id, Name: byID[key]})
		}
		return items, nil
	}

	var list []json.RawMessage
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	items := make([]DedicatedIDName, 0, len(list))
	for _, item := range list {
		parsed, err := parseDedicatedIDName(item)
		if err != nil {
			return nil, err
		}
		items = append(items, parsed)
	}
	return items, nil
}

func parseDedicatedIDName(raw json.RawMessage) (DedicatedIDName, error) {
	var obj struct {
		ID       *int    `json:"id"`
		LayoutID *int    `json:"layout_id"`
		ScriptID *int    `json:"script_id"`
		Name     *string `json:"name"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil && (obj.ID != nil || obj.LayoutID != nil || obj.ScriptID != nil || obj.Name != nil) {
		id := 0
		switch {
		case obj.ID != nil:
			id = *obj.ID
		case obj.LayoutID != nil:
			id = *obj.LayoutID
		case obj.ScriptID != nil:
			id = *obj.ScriptID
		}
		name := ""
		if obj.Name != nil {
			name = *obj.Name
		}
		return DedicatedIDName{ID: id, Name: name}, nil
	}

	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return DedicatedIDName{Name: strings.TrimSpace(text)}, nil
	}

	var id int
	if err := json.Unmarshal(raw, &id); err == nil {
		return DedicatedIDName{ID: id}, nil
	}

	return DedicatedIDName{}, fmt.Errorf("unsupported id/name shape %s", string(raw))
}
