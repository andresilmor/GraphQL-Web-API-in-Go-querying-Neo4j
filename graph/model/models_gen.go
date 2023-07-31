// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type BaseIdentification interface {
	IsBaseIdentification()
	GetUUID() *string
	GetLabel() *string
}

type GetPanoramicImagesResponse interface {
	IsGetPanoramicImagesResponse()
}

type MemberLoginResponse interface {
	IsMemberLoginResponse()
}

type BoundingBox struct {
	X      *int `json:"x,omitempty"`
	Y      *int `json:"y,omitempty"`
	Width  *int `json:"width,omitempty"`
	Height *int `json:"height,omitempty"`
}

type Error struct {
	Message     string  `json:"message"`
	Description *string `json:"description,omitempty"`
}

func (Error) IsMemberLoginResponse() {}

func (Error) IsGetPanoramicImagesResponse() {}

type HotspotPoint struct {
	BoundingBox *BoundingBox      `json:"boundingBox,omitempty"`
	Data        *HotspotPointData `json:"data,omitempty"`
}

type HotspotPointData struct {
	Alias   *string   `json:"alias,omitempty"`
	Content []*string `json:"content,omitempty"`
}

type HotspotPointTransform struct {
	Position *Position `json:"position,omitempty"`
	Scale    *Scale    `json:"scale,omitempty"`
}

type Institution struct {
	UUID  *string `json:"uuid,omitempty"`
	Label *string `json:"label,omitempty"`
	Name  *string `json:"name,omitempty"`
}

func (Institution) IsBaseIdentification()  {}
func (this Institution) GetUUID() *string  { return this.UUID }
func (this Institution) GetLabel() *string { return this.Label }

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Medication struct {
	UUID  *string `json:"uuid,omitempty"`
	Label *string `json:"label,omitempty"`
	Name  *string `json:"name,omitempty"`
}

func (Medication) IsBaseIdentification()  {}
func (this Medication) GetUUID() *string  { return this.UUID }
func (this Medication) GetLabel() *string { return this.Label }

type MedicationToTake struct {
	IntOfTime   *int        `json:"intOfTime,omitempty"`
	IsAvailable *bool       `json:"isAvailable,omitempty"`
	QuantityPer *int        `json:"quantityPer,omitempty"`
	TimeMeasure *string     `json:"timeMeasure,omitempty"`
	Medication  *Medication `json:"medication,omitempty"`
}

type Member struct {
	UUID     *string     `json:"uuid,omitempty"`
	Label    *string     `json:"label,omitempty"`
	Name     *string     `json:"name,omitempty"`
	Password *string     `json:"password,omitempty"`
	Token    *string     `json:"token,omitempty"`
	Username *string     `json:"username,omitempty"`
	Email    *string     `json:"email,omitempty"`
	MemberOf []*MemberOf `json:"memberOf,omitempty"`
}

func (Member) IsBaseIdentification()  {}
func (this Member) GetUUID() *string  { return this.UUID }
func (this Member) GetLabel() *string { return this.Label }

func (Member) IsMemberLoginResponse() {}

type MemberOf struct {
	Role        *string      `json:"role,omitempty"`
	Institution *Institution `json:"institution,omitempty"`
}

type Pacient struct {
	UUID  *string `json:"uuid,omitempty"`
	Label *string `json:"label,omitempty"`
	Name  *string `json:"name,omitempty"`
}

func (Pacient) IsBaseIdentification()  {}
func (this Pacient) GetUUID() *string  { return this.UUID }
func (this Pacient) GetLabel() *string { return this.Label }

type PanoramicSession struct {
	UUID        *string               `json:"uuid,omitempty"`
	Label       *string               `json:"label,omitempty"`
	ImageUID    *string               `json:"imageUID,omitempty"`
	PartOf      []*string             `json:"partOf,omitempty"`
	DirectedFor []*string             `json:"directedFor,omitempty"`
	Meta        *PanoramicSessionMeta `json:"meta,omitempty"`
	Mapping     []*HotspotPoint       `json:"mapping,omitempty"`
	ImageWidth  *int                  `json:"imageWidth,omitempty"`
	ImageHeight *int                  `json:"imageHeight,omitempty"`
}

func (PanoramicSession) IsBaseIdentification()  {}
func (this PanoramicSession) GetUUID() *string  { return this.UUID }
func (this PanoramicSession) GetLabel() *string { return this.Label }

func (PanoramicSession) IsGetPanoramicImagesResponse() {}

type PanoramicSessionMeta struct {
	CreatedBy *string `json:"createdBy,omitempty"`
	CreatedAt *string `json:"createdAt,omitempty"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
	IsActive  *bool   `json:"isActive,omitempty"`
}

type Position struct {
	X *float64 `json:"x,omitempty"`
	Y *float64 `json:"y,omitempty"`
	Z *float64 `json:"z,omitempty"`
}

type Scale struct {
	Width  *float64 `json:"width,omitempty"`
	Height *float64 `json:"height,omitempty"`
}
