package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Device struct {
	ID             bson.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	ModelName      string         `bson:"model_name" json:"model_name,omitempty"`
	ImageUrl       string         `bson:"imageUrl" json:"imageUrl,omitempty"`
	Specifications Specifications `bson:"specifications" json:"specifications,omitempty"`
}

type Specifications struct {
	Network      Network      `bson:"Network" json:"network,omitempty"`
	Launch       Launch       `bson:"Launch" json:"launch,omitempty"`
	Body         Body         `bson:"Body" json:"body,omitempty"`
	Display      Display      `bson:"Display" json:"display,omitempty"`
	Platform     Platform     `bson:"Platform" json:"platform,omitempty"`
	Memory       Memory       `bson:"Memory" json:"memory,omitempty"`
	MainCamera   MainCamera   `bson:"Main Camera" json:"mainCamera,omitempty"`
	SelfieCamera SelfieCamera `bson:"Selfie Camera" json:"selfieCamera,omitempty"`
	Sound        Sound        `bson:"Sound" json:"sound,omitempty"`
	Comms        Comms        `bson:"Comms" json:"comms,omitempty"`
	Features     Features     `bson:"Features" json:"features,omitempty"`
	Battery      Battery      `bson:"Battery" json:"battery,omitempty"`
	Misc         Misc         `bson:"Misc" json:"misc,omitempty"`
}

type Network struct {
	Technology string `bson:"Technology" json:"technology,omitempty"`
	Bands2G    string `bson:"2G Bands" json:"bands_2g,omitempty"`
	Bands3G    string `bson:"3G Bands" json:"bands_3g,omitempty"`
	Bands4G    string `bson:"4G Bands" json:"bands_4g,omitempty"`
	Bands5G    string `bson:"5G Bands" json:"bands_5g,omitempty"`
	Speed      string `bson:"Speed" json:"speed,omitempty"`
}

type Launch struct {
	Announced string `bson:"Announced" json:"announced,omitempty"`
	Status    string `bson:"Status" json:"status,omitempty"`
}

type Body struct {
	Dimensions string `bson:"Dimensions" json:"dimensions,omitempty"`
	Weight     string `bson:"Weight" json:"weight,omitempty"`
	Build      string `bson:"Build" json:"build,omitempty"`
	SIM        string `bson:"SIM" json:"sim,omitempty"`
	IPRating   string `bson:"IP Rating" json:"ip_rating,omitempty"`
}

type Display struct {
	Type       string `bson:"Type" json:"type,omitempty"`
	Size       string `bson:"Size" json:"size,omitempty"`
	Resolution string `bson:"Resolution" json:"resolution,omitempty"`
}

type Platform struct {
	OS      string `bson:"OS" json:"os,omitempty"`
	Chipset string `bson:"Chipset" json:"chipset,omitempty"`
	CPU     string `bson:"CPU" json:"cpu,omitempty"`
	GPU     string `bson:"GPU" json:"gpu,omitempty"`
}

type Memory struct {
	CardSlot string `bson:"Card Slot" json:"card_lot,omitempty"`
	Internal string `bson:"Internal" json:"internal,omitempty"`
}

type MainCamera struct {
	Triple   string `bson:"Triple" json:"triple,omitempty"`
	Features string `bson:"Features" json:"features,omitempty"`
	Single   string `bson:"Single" json:"single,omitempty"`
	Video    string `bson:"Video" json:"video,omitempty"`
}

type SelfieCamera struct {
	Single string `bson:"Single" json:"single,omitempty"`
	Video  string `bson:"Video" json:"video,omitempty"`
}

type Sound struct {
	Loudspeaker string `bson:"Loudspeaker" json:"loudspeaker,omitempty"`
	Jack35mm    string `bson:"3.5mm Jack" json:"jack_3.5mm,omitempty"`
}

type Comms struct {
	WLAN        string `bson:"WLAN" json:"wlan,omitempty"`
	Bluetooth   string `bson:"Bluetooth" json:"bluetooth,omitempty"`
	Positioning string `bson:"Positioning" json:"positioning,omitempty"`
	NFC         string `bson:"NFC" json:"nfc,omitempty"`
	Radio       string `bson:"Radio" json:"radio,omitempty"`
	USB         string `bson:"USB" json:"usb,omitempty"`
}

type Features struct {
	Sensors string `bson:"Sensors" json:"sensors,omitempty"`
}

type Battery struct {
	Type     string `bson:"Type" json:"type,omitempty"`
	Charging string `bson:"Charging" json:"charging,omitempty"`
}

type Misc struct {
	Colors string `bson:"colors" json:"colors,omitempty"`
	Models string `bson:"models" json:"models,omitempty"`
	Price  string `bson:"price" json:"price,omitempty"`
}
