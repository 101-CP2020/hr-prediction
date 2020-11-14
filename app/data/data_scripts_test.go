package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_okvedTitle(t *testing.T) {
	got := okvedTitle("РАЗДЕЛ C ОБРАБАТЫВАЮЩИЕ ПРОИЗВОДСТВА")

	exp := "Обрабатывающие производства"

	assert.Equal(t, exp, got)
}

func Test_okpdtrFromGroup(t *testing.T) {
	got := okpdtrFromGroup("Бухгалтер-ревизор (203392)")

	exp := 20339

	assert.Equal(t, exp, got)
}
