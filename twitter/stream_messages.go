package twitter

type StatusDeletion struct {
	ID        int64  `json:"id"`
	IDStr     string `json:"id_str"`
	UserID    int64  `json:"user_id"`
	UserIDStr string `json:"user_id_str"`
}

type statusDeletionNotice struct {
	Delete struct {
		StatusDeletion *StatusDeletion `json:"status"`
	} `json:"delete"`
}

type LocationDeletion struct {
	UserID          int64  `json:"user_id"`
	UserIDStr       string `json:"user_id_str"`
	UpToStatusID    int64  `json:"up_to_status_id"`
	UpToStatusIDStr string `json:"up_to_status_id_str"`
}

type locationDeletionNotice struct {
	ScrubGeo *LocationDeletion `json:"scrub_geo"`
}

type StreamLimit struct {
	Track int64 `json:"track"`
}

type streamLimitNotice struct {
	Limit *StreamLimit `json:"limit"`
}

type StatusWithheld struct {
	ID                  int64    `json:"id"`
	UserID              int64    `json:"user_id"`
	WithheldInCountries []string `json:"withheld_in_countries"`
}

type statusWithheldNotice struct {
	StatusWithheld *StatusWithheld `json:"status_withheld"`
}

type UserWithheld struct {
	ID                  int64    `json:"id"`
	WithheldInCountries []string `json:"withheld_in_countries"`
}
type userWithheldNotice struct {
	UserWithheld *UserWithheld `json:"user_withheld"`
}

type StreamDisconnect struct {
	Code       int64  `json:"code"`
	StreamName string `json:"stream_name"`
	Reason     string `json:"reason"`
}

type streamDisconnectNotice struct {
	StreamDisconnect *StreamDisconnect `json:"disconnect"`
}

type StallWarning struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	PercentFull int    `json:"percent_full"`
}

type stallWarningNotice struct {
	StallWarning *StallWarning `json:"warning"`
}

type FriendsList struct {
	Friends []int64 `json:"friends"`
}

type directMessageNotice struct {
	DirectMessage *DirectMessage `json:"direct_message"`
}

type Event struct {
	Event     string `json:"event"`
	CreatedAt string `json:"created_at"`
	// TODO: add List or deprecate it
	TargetObject *Tweet `json:"target_object"`
}

type DirectMessage struct {
	CreatedAt           string `json:"created_at"`
	ID                  int64  `json:"id"`
	IDStr               string `json:"id_str"`
	RecipientID         int64  `json:"recipient_id"`
	RecipientScreenName string `json:"recipient_screen_name"`
	SenderID            int64  `json:"sender_id"`
	SenderScreenName    string `json:"sender_screen_name"`
	Text                string `json:"text"`
}

type Tweet struct {
	CreatedAt            string   `json:"created_at"`
	FavoriteCount        int      `json:"favorite_count"`
	Favorited            bool     `json:"favorited"`
	FilterLevel          string   `json:"filter_level"`
	ID                   int64    `json:"id"`
	IDStr                string   `json:"id_str"`
	InReplyToScreenName  string   `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64    `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr string   `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64    `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   string   `json:"in_reply_to_user_id_str"`
	Lang                 string   `json:"lang"`
	PossiblySensitive    bool     `json:"possibly_sensitive"`
	RetweetCount         int      `json:"retweet_count"`
	Retweeted            bool     `json:"retweeted"`
	RetweetedStatus      *Tweet   `json:"retweeted_status"`
	Source               string   `json:"source"`
	Text                 string   `json:"text"`
	Truncated            bool     `json:"truncated"`
	WithheldCopyright    bool     `json:"withheld_copyright"`
	WithheldInCountries  []string `json:"withheld_in_countries"`
	WithheldScope        string   `json:"withheld_scope"`
	QuotedStatusID       int64    `json:"quoted_status_id"`
	QuotedStatusIDStr    string   `json:"quoted_status_id_str"`
	QuotedStatus         *Tweet   `json:"quoted_status"`
}
