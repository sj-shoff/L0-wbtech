package model

import "errors"

func (o *Order) Validate() error {
	if o.OrderUID == "" {
		return errors.New("order_uid is required")
	}
	if o.TrackNumber == "" {
		return errors.New("track_number is required")
	}
	if o.Entry == "" {
		return errors.New("entry is required")
	}
	if len(o.Items) == 0 {
		return errors.New("at least one item is required")
	}
	return nil
}
