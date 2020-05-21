package dailyco

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const DAILYCO_API_ROOT = "https://api.daily.co/v1/"

const DAILYCO_ROOMS_FRAGMENT = "/rooms"
const DAILYCO_ROOM_FRAGMENT = "/rooms/%s"
const DAILYCO_MEETING_TOKENS_FRAGMENT = "/meeting-tokens"

type Client struct {
	AuthorizationToken string
}

type ErrorResponse struct {
	Error string `json:"error"`
	Info  string `json:"info"`
}

type createRoomRequest struct {
	Name       string          `json:"name"`
	Privacy    string          `json:"privacy"`
	Properties *RoomProperties `json:"properties,omitempty"`
}

type RoomProperties struct {
	Nbf                *int64  `json:"nbf,omitempty"`
	Exp                *int64  `json:"exp,omitempty"`
	MaxParticipants    *uint   `json:"max_participants,omitempty"`
	Autojoin           *bool   `json:"autojoin,omitempty"`
	EnableKnocking     *bool   `json:"enable_knocking,omitempty"`
	EnableScreenshare  *bool   `json:"enable_screenshare,omitempty"`
	EnableChat         *bool   `json:"enable_chat,omitempty"`
	StartVideoOff      *bool   `json:"start_video_off,omitempty"`
	StartAudioOff      *bool   `json:"start_audio_off,omitempty"`
	OwnerOnlyBroadcast *bool   `json:"owner_only_broadcast,omitempty"`
	EnableRecording    *string `json:"enable_recording,omitempty"`
	EjectAtRoomExp     *bool   `json:"eject_at_room_exp,omitempty"`
	EjectAfterElapsed  *uint   `json:"eject_after_elapsed,omitempty"`
	Lang               *string `json:"lang,omitempty"`
}

type RoomResponse struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	ApiCreated bool            `json:"api_created"`
	Privacy    string          `json:"privacy"`
	Url        string          `json:"url"`
	CreatedAt  string          `json:"created_at"`
	Config     *RoomProperties `json:"config"`
}

type createMeetingTokenRequest struct {
	Properties *MeetingTokenProperties `json:"properties,omitempty"`
}

type MeetingTokenProperties struct {
	Nbf                 *int64  `json:"nbf,omitempty"`
	Exp                 *int64  `json:"exp,omitempty"`
	RoomName            *string `json:"room_name,omitempty"`
	IsOwner             *bool   `json:"is_owner,omitempty"`
	UserName            *string `json:"user_name,omitempty"`
	UserId              *string `json:"user_id,omitempty"`
	EnableScreenshare   *bool   `json:"enable_screenshare,omitempty"`
	StartVideoOff       *bool   `json:"start_video_off,omitempty"`
	StartAudioOff       *bool   `json:"start_audio_off,omitempty"`
	EnableRecording     *string `json:"enable_recording,omitempty"`
	StartCloudRecording *bool   `json:"start_cloud_recording,omitempty"`
	CloseTabOnExit      *bool   `json:"close_tab_on_exit,omitempty"`
	EjectAtTokenExp     *bool   `json:"eject_at_token_exp,omitempty"`
	EjectAfterElapsed   *string `json:"eject_after_elapsed,omitempty"`
	Lang                *string `json:"lang,omitempty"`
}

type MeetingTokenResponse struct {
	Token string `json:"token"`
}

type deleteMeetingResponse struct {
	Deleted bool   `json:"deleted"`
	Name    string `json:"name"`
}

func (c *Client) call(fragment string, verb string, body interface{}) ([]byte, error) {
	// Prepare URL
	u, uErr := url.Parse(DAILYCO_API_ROOT)
	if uErr != nil {
		return nil, uErr
	}
	u.Path = path.Join(u.Path, fragment)

	// Prepare Body
	requestJson := []byte("")
	if body != nil {
		var jErr error
		requestJson, jErr = json.Marshal(body)
		if jErr != nil {
			return nil, jErr
		}
	}

	// Setup Request
	req, reqErr := http.NewRequest(verb, u.String(), bytes.NewBuffer(requestJson))
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.AuthorizationToken))
	req.Header.Set("Content-Type", "application/json")

	// Send Request
	client := &http.Client{}
	resp, cErr := client.Do(req)
	if cErr != nil {
		return nil, cErr
	}
	defer resp.Body.Close()

	// Read Response
	rawBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	var errorResponse ErrorResponse
	unErr := json.Unmarshal(rawBody, &errorResponse)
	if unErr != nil {
		return nil, unErr
	} else if errorResponse.Error != "" {
		return nil, fmt.Errorf("%s: %s", errorResponse.Error, errorResponse.Info)
	}

	return rawBody, nil
}

func (c *Client) CreateRoom(name string, privacy string, properties *RoomProperties) (*RoomResponse, error) {
	requestBody := &createRoomRequest{
		Name:       name,
		Privacy:    privacy,
		Properties: properties,
	}

	rawResponse, err := c.call(DAILYCO_ROOMS_FRAGMENT, "POST", requestBody)
	if err != nil {
		return nil, err
	}

	var roomResponse RoomResponse
	uErr := json.Unmarshal(rawResponse, &roomResponse)
	if uErr != nil {
		return nil, uErr
	}

	return &roomResponse, nil
}

func (c *Client) DeleteRoom(roomName string) (bool, error) {
	urlFragment := fmt.Sprintf(DAILYCO_ROOM_FRAGMENT, roomName)

	rawResponse, err := c.call(urlFragment, "DELETE", nil)
	if err != nil {
		return false, err
	}

	var response deleteMeetingResponse
	uErr := json.Unmarshal(rawResponse, &response)
	if uErr != nil {
		return false, uErr
	}

	return response.Deleted, nil
}

func (c *Client) CreateMeetingToken(properties *MeetingTokenProperties) (*MeetingTokenResponse, error) {
	requestBody := &createMeetingTokenRequest{
		Properties: properties,
	}

	rawResponse, err := c.call(DAILYCO_MEETING_TOKENS_FRAGMENT, "POST", requestBody)
	if err != nil {
		return nil, err
	}

	var tokenResponse MeetingTokenResponse
	uErr := json.Unmarshal(rawResponse, &tokenResponse)
	if uErr != nil {
		return nil, uErr
	}

	return &tokenResponse, nil
}
