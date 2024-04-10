package test

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func assertBoolAttribute(t *testing.T, attributeName string, node *goquery.Selection, expectExists bool) {
	_, attrExists := node.Attr(attributeName)
	assert.Equal(t, expectExists, attrExists)

}

func assertStringAttribute(t *testing.T, attributeName string, node *goquery.Selection, expectedContent string) {
	sizeAttr, sizteAttrExists := node.Attr(attributeName)
	assert.True(t, sizteAttrExists)
	assert.Equal(t, expectedContent, sizeAttr)
}
