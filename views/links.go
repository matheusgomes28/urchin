package views

import "github.com/matheusgomes28/urchin/common"

// This function is used to define links for the navigation.
// Reduces the amount of copying navigation structures to all the views.
func GetUrchinLinks() []common.Link {
	return []common.Link{
		{Name: "Home", Href: "/"},
		{Name: "About", Href: "/about"},
		{Name: "Services", Href: "/services"},
		{Name: "Images", Href: "/images"},
		{Name: "Contact", Href: "/contact"},
	}
}
