package unifi

// GetWANStatus returns WAN interface status (ACTIVE/BACKUP) for a site.
func (u *Unifi) GetWANStatus(_ *Site) (*WANStatus, error) {
	return nil, nil
}
