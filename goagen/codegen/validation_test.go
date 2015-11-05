package codegen_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/raphael/goa/design"
	"github.com/raphael/goa/goagen/codegen"
)

var _ = Describe("validation code generation", func() {
	BeforeEach(func() {
		codegen.TempCount = 0
	})

	Describe("ValidationChecker", func() {
		Context("given an attribute definition and validations", func() {
			var attType design.DataType
			var validations []design.ValidationDefinition

			att := new(design.AttributeDefinition)
			target := "val"
			context := "context"
			var code string // generated code

			JustBeforeEach(func() {
				att.Type = attType
				att.Validations = validations
				code = codegen.ValidationChecker(att, target, context)
			})

			Context("of enum", func() {
				BeforeEach(func() {
					attType = design.Integer
					enumVal := &design.EnumValidationDefinition{
						Values: []interface{}{1, 2, 3},
					}
					validations = []design.ValidationDefinition{enumVal}
				})

				It("produces the validation go code", func() {
					Ω(code).Should(Equal(enumValCode))
				})
			})

			Context("of pattern", func() {
				BeforeEach(func() {
					attType = design.String
					patternVal := &design.PatternValidationDefinition{
						Pattern: ".*",
					}
					validations = []design.ValidationDefinition{patternVal}
				})

				It("produces the validation go code", func() {
					Ω(code).Should(Equal(patternValCode))
				})
			})
		})
	})
})

const (
	enumValCode = `	if !(val == 1 || val == 2 || val == 3) {
		err = goa.InvalidEnumValueError(` + "`context`" + `, val, []interface{}{1, 2, 3}, err)
	}`

	patternValCode = `	if val != "" {
		if ok := goa.ValidatePattern(` + "`.*`" + `, val); !ok {
			err = goa.InvalidPatternError(` + "`context`" + `, val, ` + "`.*`" + `, err)
		}
	}`
)
