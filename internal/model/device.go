package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Device struct {
	ID             bson.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	ModelName      string         `bson:"model_name" json:"model_name,omitempty"`
	ImageUrl       string         `bson:"imageUrl" json:"imageUrl,omitempty"`
	Specifications Specifications `bson:"specifications" json:"specifications,omitempty"`
}

type Specifications struct {
	Network      Network      `bson:"network" json:"network,omitempty"`
	Launch       Launch       `bson:"launch" json:"launch,omitempty"`
	Body         Body         `bson:"body" json:"body,omitempty"`
	Display      Display      `bson:"display" json:"display,omitempty"`
	Platform     Platform     `bson:"platform" json:"platform,omitempty"`
	Memory       Memory       `bson:"memory" json:"memory,omitempty"`
	MainCamera   MainCamera   `bson:"mainCamera" json:"mainCamera,omitempty"`
	SelfieCamera SelfieCamera `bson:"selfieCamera" json:"selfieCamera,omitempty"`
	Sound        Sound        `bson:"sound" json:"sound,omitempty"`
	Comms        Comms        `bson:"comms" json:"comms,omitempty"`
	Features     Features     `bson:"features" json:"features,omitempty"`
	Battery      Battery      `bson:"battery" json:"battery,omitempty"`
	Misc         Misc         `bson:"misc" json:"misc,omitempty"`
}

type Network struct {
	Technology string `bson:"technology" json:"technology,omitempty"`
	Bands2G    string `bson:"bands_2g" json:"bands_2g,omitempty"`
	Bands3G    string `bson:"bands_3g" json:"bands_3g,omitempty"`
	Bands4G    string `bson:"bands_4g" json:"bands_4g,omitempty"`
	Bands5G    string `bson:"bands_5g" json:"bands_5g,omitempty"`
	Speed      string `bson:"speed" json:"speed,omitempty"`
}

type Launch struct {
	Announced string `bson:"announced" json:"announced,omitempty"`
	Status    string `bson:"status" json:"status,omitempty"`
}

type Body struct {
	Dimensions string `bson:"dimensions" json:"dimensions,omitempty"`
	Weight     string `bson:"weight" json:"weight,omitempty"`
	Build      string `bson:"build" json:"build,omitempty"`
	SIM        string `bson:"sim" json:"sim,omitempty"`
	IPRating   string `bson:"ip_rating" json:"ip_rating,omitempty"`
}

type Display struct {
	Type       string `bson:"type" json:"type,omitempty"`
	Size       string `bson:"size" json:"size,omitempty"`
	Resolution string `bson:"resolution" json:"resolution,omitempty"`
}

type Platform struct {
	OS      string `bson:"os" json:"os,omitempty"`
	Chipset string `bson:"chipset" json:"chipset,omitempty"`
	CPU     string `bson:"cpu" json:"cpu,omitempty"`
	GPU     string `bson:"gpu" json:"gpu,omitempty"`
}

type Memory struct {
	CardSlot string `bson:"card_slot" json:"card_lot,omitempty"`
	Internal string `bson:"internal" json:"internal,omitempty"`
}

type MainCamera struct {
	Triple   string `bson:"triple" json:"triple,omitempty"`
	Features string `bson:"features" json:"features,omitempty"`
	Single   string `bson:"single" json:"single,omitempty"`
	Video    string `bson:"video" json:"video,omitempty"`
}

type SelfieCamera struct {
	Single string `bson:"single" json:"single,omitempty"`
	Video  string `bson:"video" json:"video,omitempty"`
}

type Sound struct {
	Loudspeaker string `bson:"loudspeaker" json:"loudspeaker,omitempty"`
	Jack35mm    string `bson:"jack_3.5mm" json:"jack_3.5mm,omitempty"`
}

type Comms struct {
	WLAN        string `bson:"wlan" json:"wlan,omitempty"`
	Bluetooth   string `bson:"bluetooth" json:"bluetooth,omitempty"`
	Positioning string `bson:"positioning" json:"positioning,omitempty"`
	NFC         string `bson:"nfc" json:"nfc,omitempty"`
	Radio       string `bson:"radio" json:"radio,omitempty"`
	USB         string `bson:"usb" json:"usb,omitempty"`
}

type Features struct {
	Sensors string `bson:"sensors" json:"sensors,omitempty"`
}

type Battery struct {
	Type     string `bson:"type" json:"type,omitempty"`
	Charging string `bson:"charging" json:"charging,omitempty"`
}

type Misc struct {
	Colors string `bson:"colors" json:"colors,omitempty"`
	Models string `bson:"models" json:"models,omitempty"`
	Price  string `bson:"price" json:"price,omitempty"`
}
