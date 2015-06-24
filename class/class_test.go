package class

import (
	"testing"
)

func TestClassFactory(t *testing.T) {
	NewClassAssassin()
	NewClassDisabler()
	NewClassHealer()
	NewClassMage()
	NewClassTank()
}
