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
	FormattedName string `json:"formattedName"`
	ImageURL      string `json:"imageUrl"`
	Skins         map[string][]string
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
