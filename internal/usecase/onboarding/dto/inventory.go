package dto

import inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"

const (
	machineSizeName = "machine"
)

type CreateInventory struct {
	Name      string
	Namespace string
}

func NewCreateInventory(name string, namespace string) CreateInventory {
	return CreateInventory{Name: name, Namespace: namespace}
}

type Inventory struct {
	UUID         string
	Namespace    string
	ProductSKU   string
	SerialNumber string
	Sizes        map[string]string
	NICs         []inventories.NICSpec
}

func NewInventory(
	UUID string,
	namespace string,
	productSKU string,
	serialNumber string,
	sizes map[string]string,
	NICs []inventories.NICSpec) Inventory {
	return Inventory{
		UUID:         UUID,
		Namespace:    namespace,
		ProductSKU:   productSKU,
		SerialNumber: serialNumber,
		Sizes:        sizes,
		NICs:         NICs}
}

func (i *Inventory) IsMachine() bool {
	_, ok := i.Sizes[inventories.GetSizeMatchLabel(machineSizeName)]
	return ok
}
