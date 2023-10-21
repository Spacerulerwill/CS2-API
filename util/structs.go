package util

type SkinVariation struct {
	FormattedName   string    `json:"formattedName"`
	ConditionImages [5]string `json:"conditionImages"`
}

type Skin struct {
	FormattedName     string                   `json:"formattedName"`
	Description       string                   `json:"description"`
	FlavorText        string                   `json:"flavorText"`
	MinFloat          string                   `json:"minFloat"`
	MaxFloat          string                   `json:"maxFloat"`
	WeaponType        string                   `json:"weaponType"`
	Rarity            string                   `json:"rarity"`
	ConditionImages   [5]string                `json:"conditionImages"`
	Variations        map[string]SkinVariation `json:"variations"`
	StattrakAvailable bool                     `json:"stattrakAvailable"`
	SouvenirAvailable bool                     `json:"souvenirAvailable"`
}

type Case struct {
	FormattedName string              `json:"formattedName"`
	ImageURL      string              `json:"imageUrl"`
	Skins         map[string][]string `json:"skins"`
}

type SouvenirPackage struct {
	FormattedName string `json:"formattedName"`
	ImageURL      string `json:"imageUrl"`
	Collection    string `json:"collection"`
}

type Collection struct {
	FormattedName string              `json:"formattedName"`
	ImageUrl      string              `json:"imageUrl"`
	Skins         map[string][]string `json:"skins"`
}

type Sticker struct {
	FormattedName string `json:"formattedName"`
	ImageURL      string `json:"imageUrl"`
	Rarity        string `json:"rarity"`
}

type StickerCapsule struct {
	FormattedName string              `json:"formattedName"`
	ImageURL      string              `json:"imageUrl"`
	Stickers      map[string][]string `json:"stickers"`
}

type Graffiti struct {
	FormattedName  string            `json:"formattedName"`
	Rarity         string            `json:"rarity"`
	ImageURL       string            `json:"imageUrl"`
	ColorVarations map[string]string `json:"colorVariations"`
	GraffitiBox    string            `json:"graffitiBox"`
}

type MusicKit struct {
	FormattedName     string            `json:"formattedName"`
	Artist            string            `json:"artist"`
	Description       string            `json:"description"`
	Rarity            string            `json:"rarity"`
	ImageURL          string            `json:"imageUrl"`
	StattrakAvailable bool              `json:"stattrakAvailable"`
	BoxesFoundIn      []string          `json:"boxesFoundIn"`
	AudioURLs         map[string]string `json:"audioUrls"`
}

type Agent struct {
	FormattedName string `json:"formattedName"`
	Rarity        string `json:"rarity"`
	ImageUrl      string `json:"imageUrl"`
	Description   string `json:"description"`
	FlavorText    string `json:"flavorText"`
}

type Patch struct {
	FormattedName string `json:"formattedName"`
	Rarity        string `json:"rarity"`
	ImageUrl      string `json:"imageUrl"`
	FlavorText    string `json:"flavorText"`
}
