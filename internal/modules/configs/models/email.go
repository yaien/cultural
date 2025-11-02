package models

type Email struct {
	Subject string `bson:"subject" json:"subject"`
	Body    Node   `bson:"body" json:"body"`
}

var DefaultInvitationEmail = Email{
	Subject: "{{organization.name}} te ha invitado a unirte a su equipo de trabajo",
	Body: Node{
		Type: "email:body",
		Children: []Node{
			{
				Type: "email:section",
				Children: []Node{
					{
						Type:    "h1",
						Content: "{{organization.name}} te ha invitado a unirte a su equipo de trabajo",
					},
				},
			},
			{
				Type: "email:section",
				Children: []Node{
					{
						Type:    "p",
						Content: "Hola {{user.name}},",
					},
					{
						Type:    "p",
						Content: "Has sido invitado a unirte al equipo de trabajo de {{organization.name}} en nuestra plataforma. Para aceptar la invitación y crear tu cuenta, haz clic en el siguiente botón:",
					},
					{
						Type:    "email:button",
						Attrs:   map[string]any{"href": "{{invitation.url}}"},
						Content: "Aceptar invitación",
					},
					{
						Type:    "p",
						Content: "Si no solicitaste unirte a esta organización, puedes ignorar este correo electrónico.",
					},
					{
						Type:    "p",
						Content: "¡Esperamos verte pronto!",
					},
				},
			},
			{
				Type:    "email:section",
				Content: "El equipo de Cultural",
			},
		},
	},
}
