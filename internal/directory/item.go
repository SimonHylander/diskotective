package directory

type Item struct {
	Name string
	Path string
	Type ItemType
	Size int64
}

type ItemType string

const (
	ItemTypeDirectory ItemType = "Directory"
	ItemTypeFile      ItemType = "File"
)

/*func (i Item) Name() string       { return i.Name }
func (i Item) Description() string { return i.Description }
func (i Item) FilterValue() string { return i.Name }
func (i Item) Size() int64         { return i.Size }*/

type BySize []Item

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size < a[j].Size }
