package validators_test

import (
	"testing"

	testenum "frisboo-bank/openapi-generator-service/pkg/validation/testdata/enums/test_enum"
	"frisboo-bank/openapi-generator-service/pkg/validation/validators"

	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

func TestValidEnum(t *testing.T) {
	type TestModel struct {
		TestEnum testenum.TestEnum
	}
	model := TestModel{
		TestEnum: testenum.TestEnums.PASSED,
	}

	err := ozzo.Validate(&model.TestEnum, validators.ValidEnum())

	assert.NoError(t, err)
}

func TestInvalidEnum(t *testing.T) {
	type TestModel struct {
		TestEnum testenum.TestEnum
	}
	model := TestModel{
		TestEnum: testenum.TestEnums.UNKNOWN,
	}

	err := ozzo.Validate(&model.TestEnum, validators.ValidEnum())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid enum value: unknown")
}

func TestNotEnum(t *testing.T) {
	type TestModel struct {
		TestEnum string
	}
	model := TestModel{
		TestEnum: "value",
	}

	err := ozzo.Validate(&model.TestEnum, validators.ValidEnum())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value *string does not implement IsValid()")
}

func TestNotAllowedValueEnum(t *testing.T) {
	type TestModel struct {
		TestEnum testenum.TestEnum
	}
	model := TestModel{
		TestEnum: testenum.TestEnums.BOOKED,
	}

	err := ozzo.Validate(&model.TestEnum, validators.ValidEnum(
		testenum.TestEnums.PASSED,
	))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enum value booked is not allowed")

	err = ozzo.Validate(&model.TestEnum, validators.ValidEnum(
		testenum.TestEnums.BOOKED,
	))

	assert.NoError(t, err)
}
