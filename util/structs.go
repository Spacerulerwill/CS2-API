package util

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type SkinVariation struct {
	FormattedName   string    `json:"formatted_name"`
	ConditionImages [5]string `json:"condition_images"`
	InspectUrls     [5]string `json:"inspect_urls"`
}

type Skin struct {
	FormattedName       string                   `json:"formatted_name"`
	Description         string                   `json:"description"`
	FlavorText          string                   `json:"flavor_text"`
	MinFloat            string                   `json:"min_float"`
	MaxFloat            string                   `json:"max_float"`
	WeaponType          string                   `json:"weapon_type"`
	Rarity              string                   `json:"rarity"`
	ConditionImages     [5]string                `json:"condition_images"`
	WorstConditionIndex int                      `json:"worst_condition_index"`
	BestConditionIndex  int                      `json:"best_condition_index"`
	InspectUrls         [5]string                `json:"inspect_urls"`
	Variations          map[string]SkinVariation `json:"variations"`
	StattrakAvailable   bool                     `json:"stattrak_available"`
	SouvenirAvailable   bool                     `json:"souvenir_available"`
}

type Case struct {
	FormattedName string                                   `json:"formatted_name"`
	ImageURL      string                                   `json:"image_url"`
	Skins         *orderedmap.OrderedMap[string, []string] `json:"skins"`
	RequiresKey   bool                                     `json:"requires_key"`
}

type SouvenirPackage struct {
	FormattedName string `json:"formatted_name"`
	ImageURL      string `json:"image_url"`
	Collection    string `json:"collection"`
}

type Collection struct {
	FormattedName string                                   `json:"formatted_name"`
	ImageUrl      string                                   `json:"image_url"`
	Skins         *orderedmap.OrderedMap[string, []string] `json:"skins"`
}

type Sticker struct {
	FormattedName string `json:"formatted_name"`
	ImageURL      string `json:"image_url"`
	InspectUrl    string `json:"inspect_url"`
	Rarity        string `json:"rarity"`
}

type StickerCapsule struct {
	FormattedName string                                   `json:"formatted_name"`
	ImageURL      string                                   `json:"image_url"`
	Stickers      *orderedmap.OrderedMap[string, []string] `json:"stickers"`
}

type GraffitiColorVariation struct {
	ImageUrl   string `json:"image_url"`
	InspectUrl string `json:"inspect_url"`
}

type Graffiti struct {
	FormattedName  string                            `json:"formatted_name"`
	Rarity         string                            `json:"rarity"`
	ImageURL       string                            `json:"image_url"`
	InspectUrl     string                            `json:"inspect_url"`
	ColorVarations map[string]GraffitiColorVariation `json:"color_variations"`
	GraffitiBox    string                            `json:"graffiti_box"`
}

type MusicKit struct {
	FormattedName     string            `json:"formatted_name"`
	Artist            string            `json:"artist"`
	Description       string            `json:"description"`
	Rarity            string            `json:"rarity"`
	ImageURL          string            `json:"image_url"`
	StattrakAvailable bool              `json:"stattrak_available"`
	BoxesFoundIn      []string          `json:"boxes_found_in"`
	AudioURLs         map[string]string `json:"audio_urls"`
}

type Agent struct {
	FormattedName string `json:"formatted_name"`
	Rarity        string `json:"rarity"`
	ImageUrl      string `json:"image_url"`
	InspectUrl    string `json:"inspect_url"`
	Description   string `json:"description"`
	FlavorText    string `json:"flavor_text"`
}

type Patch struct {
	FormattedName string `json:"formatted_name"`
	Rarity        string `json:"rarity"`
	ImageUrl      string `json:"image_url"`
	InspectUrl    string `json:"inspect_url"`
	FlavorText    string `json:"flavor_text"`
}

type PatchPack struct {
	FormattedName string                                   `json:"formatted_name"`
	ImageUrl      string                                   `json:"image_url"`
	Patches       *orderedmap.OrderedMap[string, []string] `json:"patches"`
}

type Pin struct {
	FormattedName string `json:"formatted_name"`
	Rarity        string `json:"rarity"`
	ImageUrl      string `json:"image_url"`
	InspectUrl    string `json:"inspect_url"`
	PinCapsule    string `json:"pin_capsule"`
}

type PinCapsule struct {
	FormattedName string                                   `json:"formatted_name"`
	ImageUrl      string                                   `json:"image_url"`
	Pins          *orderedmap.OrderedMap[string, []string] `json:"pins"`
}
