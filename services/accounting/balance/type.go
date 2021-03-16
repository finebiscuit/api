package balance

//go:generate go run github.com/dmarkham/enumer -type=Type -output=type_gen.go -gqlgen -sql

type Type uint8

const (
	NoType Type = iota

	Cash
	CashChecking
	CashSaving
	CashPhysical
	CashDigital
	CashOther

	Investment
	InvestmentBrokerage
	InvestmentRetirement
	InvestmentHealth
	InvestmentEducation
	InvestmentOther

	RealEstate
	RealEstatePrimary
	RealEstateInvestment
	RealEstateVacation
	RealEstateCommercial
	RealEstateLand
	RealEstateOther

	Vehicle
	VehicleAutomobile
	VehicleBoat
	VehiclePlane
	VehicleOther

	Property
	PropertyFurniture
	PropertyJewellery
	PropertyCollectible
	PropertyElectronics
	PropertyOther

	Cryptocurrency

	Other
)

var parentTypes = []Type{NoType, Cash, Investment, RealEstate, Vehicle, Property, Cryptocurrency, Other}
var subtypeMapping map[Type][]Type

func init() {
	subtypeMapping = make(map[Type][]Type)
}

func (i Type) Parent() Type {
	return NoType
}

func (i Type) Subtypes() []Type {
	return []Type{}
}
