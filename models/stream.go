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

type M3U8 interface {
	*m3u8.MasterPlaylist | *m3u8.MediaPlaylist
}

type Manifests[T M3U8] struct {
	ManConf T `json:"man_conf"`      // JSON of the dynamically manipulated  manifest
	OrnConf T `json:"original_conf"` // JSON of the original manifest
}

type Playlist struct {
	Master Manifests[*m3u8.MasterPlaylist]
	Main   []Manifests[*m3u8.MediaPlaylist]
}

// mongo collection 'streams'
type StreamData []struct {
	ID                 primitive.ObjectID `json:"_id,omitempty"`
	VidmolyAlias       string             `json:"vidmoly_alias"`
	HLS_PlaylistRemote string             `json:"hls_remote"`
	/*
	   This will enable to reverse proxy and manipulate target HLS playlist configurations
	*/

	Playlist Playlist `json:"playlist"`
}

const StreamCol string = "streams"
