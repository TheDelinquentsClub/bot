package utils

import (
	"bytes"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

type Embed struct {
	title       string
	description string
	author      discord.EmbedAuthor
	fields      []discord.EmbedField
	footer      *discord.EmbedFooter
	timestamp   discord.Timestamp
	color       discord.Color
	Components  []discord.Component
	files       []sendpart.File
	image       *discord.EmbedImage
}

// NewEmbed initializes a new discord.Embed object
func NewEmbed() Embed {
	return Embed{}
}

// General embed configuring
func (e Embed) SetTitle(Title string) Embed {
	e.title = Title
	return e
}
func (e Embed) SetAuthor(Name string, Icon string) Embed {
	e.author = discord.EmbedAuthor{Name: Name, Icon: Icon}
	return e
}
func (e Embed) SetImage(URL string) Embed {
	e.image = &discord.EmbedImage{URL: URL}
	return e
}
func (e Embed) SetDescription(Description string) Embed {
	e.description = Description
	return e
}
func (e Embed) AddField(Name string, Inline bool, Value string) Embed {
	field := discord.EmbedField{Name: Name, Inline: Inline, Value: Value}
	e.fields = append(e.fields, field)
	return e
}

// AddFields lets you pass Multiple discord.EmbedField objects through a variable
func (e Embed) AddFields(Fields []discord.EmbedField) Embed {
	e.fields = append(e.fields, Fields...)
	return e
}
func (e Embed) SetColor(HexColor discord.Color) Embed {
	e.color = HexColor
	return e
}
func (e Embed) SetFooter(Text string, Icon string) Embed {
	e.footer = &discord.EmbedFooter{Text: Text, Icon: Icon}
	return e
}

// AddFile attatches a file to the discord message
func (e Embed) AddFile(Name string, file []byte) Embed {
	reader := bytes.NewReader(file)
	files := sendpart.File{Name: Name, Reader: reader}
	e.files = append(e.files, files)
	return e
}

func (e Embed) RemoveComponents() Embed {
	e.Components = []discord.Component{}
	return e
}

// Button message components

func (e Embed) AddPrimaryButton(Label string, ID string) Embed {
	button := &discord.ButtonComponent{Label: Label, CustomID: ID, Style: discord.PrimaryButton}
	e.Components = append(e.Components, button)
	return e
}
func (e Embed) AddSecondaryButton(Label string, ID string) Embed {
	button := &discord.ButtonComponent{Label: Label, CustomID: ID, Style: discord.SecondaryButton}
	e.Components = append(e.Components, button)
	return e
}
func (e Embed) AddDangerButton(Label string, ID string) Embed {
	button := &discord.ButtonComponent{Label: Label, CustomID: ID, Style: discord.DangerButton}
	e.Components = append(e.Components, button)
	return e
}
func (e Embed) AddURLButton(Label string, URL string) Embed {
	button := &discord.ButtonComponent{Label: Label, Style: discord.LinkButton, URL: URL}
	e.Components = append(e.Components, button)
	return e
}

// Select Message Component

type SelectComponent struct {
	embed       Embed
	entries     []discord.SelectComponentOption
	id          string
	placeholder string
	disabled    bool
}

func (e Embed) AddSelectComponent(ID string, Placeholder string, Disabled bool) SelectComponent {
	return SelectComponent{id: ID, placeholder: Placeholder, disabled: Disabled, embed: e}
}

// AddOption Generates a discord.SelectComponentOption object.
//
// Label and Value MUST be unique
func (e SelectComponent) AddOption(Label string, Value string, Description string, Emoji *discord.ButtonEmoji, Default bool) SelectComponent {
	entry := discord.SelectComponentOption{
		Label:       Label,
		Value:       Value,
		Description: Description,
		Emoji:       Emoji,
		Default:     Default,
	}
	e.entries = append(e.entries, entry)
	return e
}

// MakeSelectComponent generates a discord.SelectComponent object and returns the builder back to the Embed struct
func (e SelectComponent) MakeSelectComponent() Embed {
	selectComponent := &discord.SelectComponent{
		CustomID:    e.id,
		Options:     e.entries,
		Placeholder: e.placeholder,
		Disabled:    e.disabled,
	}
	e.embed.Components = append(e.embed.Components, selectComponent)
	return e.embed
}

// MakeResponse generates an api.InteractionResponse object
func (e Embed) MakeResponse() api.InteractionResponse {
	var components = &[]discord.Component{
		&discord.ActionRowComponent{
			Components: e.Components,
		},
	}

	if e.Components == nil {
		components = &[]discord.Component{}
	}

	res := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Flags: api.EphemeralResponse,
			Embeds: &[]discord.Embed{
				{
					Title:       e.title,
					Description: e.description,
					Color:       e.color,
					Fields:      e.fields,
					Footer:      e.footer,
					Image:       e.image,
					Timestamp:   discord.NowTimestamp(),
				},
			},
			Components: components,
			Files:      e.files,
		},
	}
	return res
}
func (e Embed) UpdateResponse() api.InteractionResponse {
	var components = &[]discord.Component{
		&discord.ActionRowComponent{
			Components: e.Components,
		},
	}

	if e.Components == nil {
		components = nil
	}

	res := api.InteractionResponse{
		Type: api.UpdateMessage,
		Data: &api.InteractionResponseData{
			Flags: api.EphemeralResponse,
			Embeds: &[]discord.Embed{
				{
					Title:       e.title,
					Description: e.description,
					Color:       e.color,
					Fields:      e.fields,
					Footer:      e.footer,
					Image:       e.image,
					Timestamp:   discord.NowTimestamp(),
				},
			},
			Components: components,
			Files:      e.files,
		},
	}
	return res
}

// EditInteraction generates an api.EditInteractionResponseData object
func (e Embed) EditInteraction() api.EditInteractionResponseData {
	var components = &[]discord.Component{
		&discord.ActionRowComponent{
			Components: e.Components,
		},
	}

	if e.Components == nil {
		components = &[]discord.Component{}
	}

	edit := api.EditInteractionResponseData{
		Embeds: &[]discord.Embed{
			{
				Title:       e.title,
				Description: e.description,
				Color:       e.color,
				Fields:      e.fields,
				Footer:      e.footer,
				Image:       e.image,
				Timestamp:   discord.NowTimestamp(),
			},
		},
		Components: components,
		Files:      e.files,
	}
	return edit
}

// MakeMessage generates an api.SendMessageData object
func (e Embed) MakeMessage() api.SendMessageData {

	var components = []discord.Component{
		&discord.ActionRowComponent{
			Components: e.Components,
		},
	}

	if e.Components == nil {
		components = []discord.Component{}
	}

	res := api.SendMessageData{
		Embeds: []discord.Embed{
			{
				Title:       e.title,
				Description: e.description,
				Color:       e.color,
				Fields:      e.fields,
				Footer:      e.footer,
				Image:       e.image,
				Timestamp:   discord.NowTimestamp(),
			},
		},
		Components: components,
		Files:      e.files,
	}
	return res
}
