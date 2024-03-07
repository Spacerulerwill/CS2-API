package cs2

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// Containers
type Container struct {
	FormattedName string                                   `json:"formatted_name"`
	ImageURL      string                                   `json:"image_url"`
	Items         *orderedmap.OrderedMap[string, []string] `json:"items"`
	RequiresKey   bool                                     `json:"requires_key"`
}
type Case Container
type StickerCapsule Container
type PatchPack Container
type PinCapsule Container
type Collection Container
type SouvenirPackage struct {
	FormattedName string `json:"formatted_name"`
	ImageURL      string `json:"image_url"`
	Collection    string `json:"collection"`
}

// Items
type Item struct {
	FormattedName     string   `json:"formatted_name"`
	Description       string   `json:"description"`
	FlavorText        string   `json:"flavor_text"`
	Quality           string   `json:"quality"`
	InspectURLs       []string `json:"inspect_urls"`
	ImageURLs         []string `json:"image_urls"`
	StattrakAvailable bool     `json:"stattrak_available"`
	SouvenirAvailable bool     `json:"souvenir_available"`
	ContainersFoundIn []string `json:"containers_found_in"`
}
type SkinVariation struct {
	FormattedName   string   `json:"formatted_name"`
	ConditionImages []string `json:"condition_images"`
	InspectUrls     []string `json:"inspect_urls"`
}
type Skin struct {
	Item
	WeaponType          string                   `json:"weapon_type"`
	MinFloat            string                   `json:"min_float"`
	MaxFloat            string                   `json:"max_float"`
	WorstConditionIndex int                      `json:"worst_condition_index"`
	BestConditionIndex  int                      `json:"best_condition_index"`
	Variations          map[string]SkinVariation `json:"variations"`
}
type Sticker Item
type Patch Item
type Pin Item
type Agent Item
type MusicKit struct {
	Item
	Artist    string            `json:"artist"`
	AudioURLs map[string]string `json:"audio_urls"`
}
type GraffitiColorVariation struct {
	ImageUrl   string `json:"image_url"`
	InspectUrl string `json:"inspect_url"`
}
type Graffiti struct {
	Item
	ColorVarations map[string]GraffitiColorVariation `json:"color_variations"`
}
