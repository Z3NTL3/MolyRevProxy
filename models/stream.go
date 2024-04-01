/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// mongo collection 'stream'
type StreamData struct {
	ID                 primitive.ObjectID `json:"omitempty,_id"`
	VidmolyAlias       string             `json:"vidmoly_alias"`
	HLS_PlaylistRemote string             `json:"hls_remote"`
	PlaylistConf       string             `json:"playlist_conf"` // stringified JSON of the dynamically manipulated config
}
