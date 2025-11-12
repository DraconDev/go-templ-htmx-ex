package templates

import (
	_ "github.com/DraconDev/go-templ-htmx-ex/templates/components"
	_ "github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
	_ "github.com/DraconDev/go-templ-htmx-ex/templates/pages"
)

// This file exists to define the 'templates' package at the module root level
// and ensure all generated templ components in subdirectories are included
// in the main 'templates' package via side-effect imports.