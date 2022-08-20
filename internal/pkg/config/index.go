package serviceconfig

import (
	mediagrouprawsendhelper "forwarding-bot/internal/pkg/helper/media-group-raw"
)

type ChannelConfigMap map[string]*ChannelConfig
type SenderMap map[int64]*mediagrouprawsendhelper.MediaGroupRawSendHelper

type ChannelConfig struct {
	ChannelID    int64   `json:"channel_id,omitempty"`
	IsDBChannel  bool    `json:"is_db_channel,omitempty"`
	ForwardToIDs []int64 `json:"forward_to_ids,omitempty"`
}
