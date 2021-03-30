package model

import (
	"time"
)

type Message struct {
	ID         uint       `faker:"-"`
	CreatedAt  *time.Time `faker:"-"`
	UpdatedAt  *time.Time `faker:"-"`
	DeletedAt  *time.Time `faker:"-" json:",omitempty"`
	SentAt     *time.Time `faker:"-" json:",omitempty"`
	ReceivedAt *time.Time `faker:"-" json:",omitempty"`

	Name      string `faker:"name"`
	Summary   string `faker:"sentence"`
	Paragraph string `faker:"paragraph"`

	IPV4 string `faker:"ipv4"`
	IPV6 string `faker:"ipv6"`
	MAC  string `faker:"mac_address"`

	Latitude  float32 `faker:"lat"`
	Longitude float32 `faker:"long"`

	FirstName        string `faker:"first_name"`
	LastName         string `faker:"last_name"`
	Email            string `faker:"email"`
	PhoneNumber      string `faker:"phone_number"`
	CreditCardNumber string `faker:"cc_number"`
	CreditCardType   string `faker:"cc_type"`

	Criticality int
	Status      bool
}
