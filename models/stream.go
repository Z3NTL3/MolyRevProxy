/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Playlist struct {
	ManConf interface{} `json:"man_conf"`      // JSON of the dynamically manipulated  manifest
	OrnConf interface{} `json:"original_conf"` // JSON of the manifest
}

// mongo collection 'streams'
type StreamData struct {
	ID                 primitive.ObjectID `json:"omitempty,_id"`
	VidmolyAlias       string             `json:"vidmoly_alias"`
	HLS_PlaylistRemote string             `json:"hls_remote"`
	/*
		This will enable to reverse proxy and manipulate target HLS playlist configurations
	*/

	Playlist Playlist `json:"playlist"`
}

const StreamCol string = "streams"
