package importers

// Lister defines the interface common to all importers
type Lister interface {
	List(pkg string, recursive bool) ([]string, error)
}
