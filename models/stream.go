/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package models

import (
	"github.com/grafov/m3u8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Playlist[T any] struct {
	ManConf T `json:"man_conf"`      // stringified JSON of the dynamically manipulated  manifest
	OrnConf T `json:"original_conf"` // stringified JSON of the manifest
}

// mongo collection 'streams'
type StreamData struct {
	ID                 primitive.ObjectID `json:"omitempty,_id"`
	VidmolyAlias       string             `json:"vidmoly_alias"`
	HLS_PlaylistRemote string             `json:"hls_remote"`
	/*
		This will enable to reverse proxy and manipulate target HLS playlist configurations
	*/
	Master Playlist[m3u8.MasterPlaylist] `json:"master_manifest"`
	Media  Playlist[m3u8.MediaPlaylist]  `json:"media_manifest"`
}

const StreamCol string = "streams"
