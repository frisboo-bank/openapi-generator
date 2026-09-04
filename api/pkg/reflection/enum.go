package reflection

import "iter"

type EnumValue interface {
	comparable
	IsValid() bool
	String() string
}

type EnumContainer[T EnumValue] interface {
	All() iter.Seq[T]
}
