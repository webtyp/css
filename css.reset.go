//go:build !wasm

package css

// resetRules are the base reset rules RenderCSS() emits in the `tokens`
// layer — the lowest of the four the widget style DSL declares. Unlayered,
// `svg { display: block }` would beat `@layer widgets { .part { display:
// none } }` regardless of specificity, so any component hiding an icon by
// state could not.
func resetRules() []item {
	return []item{
		rule(selector("*, *::before, *::after"),
			boxSizing(str("border-box")),
		),
		// text-size-adjust stops mobile browsers reflowing text at their own
		// scale. tap-highlight-color is the iOS-only grey wash painted over
		// anything tappable for ~300ms after a touch: Chrome Android uses a
		// different colour and duration, so the same button flashes differently
		// on each. Components own their press state; the UA's is noise.
		rule(selector("html"),
			rawRule("  -webkit-text-size-adjust: 100%;\n  text-size-adjust: 100%;\n  -webkit-tap-highlight-color: transparent;"),
		),
		rule(selector("body"),
			margin(zero),
			fontFamily(FontSans),
			fontSize(TextBase),
			lineHeight(LeadingNormal),
			color(ColorOnSurface),
			background(ColorBackground),
		),
		// Amber, not ColorPrimary: every "this is the active element" signal in
		// the system is the Accent family — the current nav item, a selected
		// row — and a focus ring is the same statement for a control, so it
		// speaks the same colour instead of the primary/gradient hue that reads
		// as chrome here. A negative outline-offset draws the ring INSIDE the
		// border box: a positive offset put it outside, where any clipping
		// ancestor (a search bar with overflow: hidden) sheared it off and a
		// tight neighbour collided with it. Inset, it is always fully drawn and
		// never moves the layout.
		rule(selector(":focus-visible"),
			outline(str("2px solid "+ColorAccent.Var())),
			outlineOffset(px(-2)),
		),
		// User-agent margins and list chrome are geometry the style DSL cannot
		// express and cannot see. A <ul> used as a nav rail carries 40px of
		// padding-inline-start it never asked for, an <h1> carries 0.67em of
		// block margin: both silently widen or heighten the component around
		// them. Zero them here, in the lowest layer, so a part's box is exactly
		// what its rule declares.
		rule(selector("h1, h2, h3, h4, h5, h6, p, figure, blockquote, dl, dd"),
			margin(zero),
		),
		// Heading sizes are the DSL's business, not the user agent's. An <h1>
		// carries font-size: 2em by default, which compounds against whatever
		// the part around it declares and makes a title three times the size a
		// component asked for. FontSize()/FontWeight() say what a heading looks
		// like; its level stays a semantic choice.
		// An <a> arrives underlined and painted the user agent's link blue, which
		// fights every surface a component puts it on — a nav item ends up blue
		// text on a blue button. Components say what a link looks like; the
		// browser's default is chrome, like the heading sizes below.
		rule(selector("a"),
			rawRule("  color: inherit;\n  text-decoration: none;"),
		),
		rule(selector("h1, h2, h3, h4, h5, h6"),
			rawRule("  font-size: inherit;\n  font-weight: inherit;"),
		),
		// Form controls are the other UA-font holdout: an <input> arrives with
		// the user agent's family and size (Arial 13px on most platforms), which
		// shrinks the value inside a field the design measured for the app font.
		// Inheriting makes the control speak the same type as the text around it;
		// a part can still override size with FontSize().
		rule(selector("input, textarea, select"),
			rawRule("  font: inherit;"),
		),
		// A <button> is the widest cross-browser gap left. iOS Safari renders it
		// with appearance: push-button — its own corner radius, a vertical
		// gradient, a UA border and centred text — and that chrome outranks the
		// background and radius a part declares. Chrome Android paints a flat
		// grey box instead, so the same button is two different shapes. Zeroing
		// appearance, background, border and radius makes the element an empty
		// box on both, which is what a styled part expects to start from.
		// padding: 0 kills the UA's own button padding (1px block, 6px inline)
		// that a bare <button> inside a flex row still carries — a part that
		// declares Pad() rebuilds it anyway, but an unstyled one must not
		// arrive with a lopsided box the design never asked for.
		// The [type=…] selectors cover <input type="submit">, which is a button
		// everywhere except in the selector `button`.
		rule(selector("button, [type=\"button\"], [type=\"reset\"], [type=\"submit\"]"),
			rawRule("  -webkit-appearance: none;\n  appearance: none;\n  background-color: transparent;\n  background-image: none;\n  color: inherit;\n  font: inherit;\n  border: 0;\n  border-radius: 0;\n  padding: 0;"),
		),
		// iOS Safari paints an inset shadow and forces its own corner radius on
		// text fields; no border-radius a part declares removes it. :where()
		// keeps checkbox and radio out — appearance: none on those erases the
		// control entirely instead of flattening it, and their styling belongs
		// to webtyp/form.
		//
		// border: 0 for the same reason the button rule above carries it, and it
		// was the one half of that pair this rule was missing: appearance: none
		// does NOT drop the UA's own `border: 2px inset`, so a field arrived with
		// a heavy dark box no skin had asked for. A part painting a flat fill
		// (fieldset's input is As(Page), no border at all) could not remove it
		// either — there was nothing in the cascade to override, only chrome to
		// undo. Parts that DO want an edge declare it themselves and win here
		// anyway: this rule is @layer tokens, a skin is @layer widgets.
		rule(selector("input:where(:not([type=\"checkbox\"]):not([type=\"radio\"])), textarea"),
			rawRule("  -webkit-appearance: none;\n  appearance: none;\n  border: 0;\n  border-radius: 0;"),
		),
		// Firefox and Edge let a <select> inherit text-transform from an
		// ancestor; every other engine does not, so an uppercased container
		// silently uppercases the dropdown on two browsers only. The native
		// arrow stays: it is the only affordance the control has, and
		// webtyp/form owns replacing it.
		// Select controls: modern styled appearance with custom chevron arrow,
		// consistent with webtyp/components (ColorSurface, ColorOutline, RadiusMd, TextSm).
		rule(selector("select"),
			rawRule("  -webkit-appearance: none;\n  appearance: none;\n  font: inherit;\n  font-size: "+TextSm.Var()+";\n  color: "+ColorOnSurface.Var()+";\n  background-color: "+ColorSurface.Var()+";\n  border: 1px solid "+ColorOutline.Var()+";\n  border-radius: "+RadiusMd.Var()+";\n  padding: "+Space2.Var()+" "+Space8.Var()+" "+Space2.Var()+" "+Space3.Var()+";\n  min-height: 2.25rem;\n  line-height: "+LeadingNormal.Var()+";\n  cursor: pointer;\n  background-image: url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20' fill='%236E6E73'%3E%3Cpath fill-rule='evenodd' d='M5.22 8.22a.75.75 0 0 1 1.06 0L10 11.94l3.72-3.72a.75.75 0 1 1 1.06 1.06l-4.25 4.25a.75.75 0 0 1-1.06 0L5.22 9.28a.75.75 0 0 1 0-1.06Z' clip-rule='evenodd'/%3E%3C/svg%3E\");\n  background-repeat: no-repeat;\n  background-position: right 0.6rem center;\n  background-size: 1.25rem 1.25rem;\n  text-transform: none;\n  transition: border-color 0.15s ease, box-shadow 0.15s ease;"),
		),
		rule(selector("select:hover:not(:disabled)"),
			rawRule("  border-color: "+ColorMuted.Var()+";"),
		),
		rule(selector("select:disabled"),
			rawRule("  opacity: 0.5;\n  cursor: not-allowed;\n  background-color: "+ColorSurfaceSunken.Var()+";"),
		),
		rule(selector("select option"),
			rawRule("  background-color: "+ColorSurface.Var()+";\n  color: "+ColorOnSurface.Var()+";"),
		),
		// Checkbox controls: styled with primary color, smooth corners and clean checkmark.
		rule(selector("input[type=\"checkbox\"]"),
			rawRule("  -webkit-appearance: none;\n  appearance: none;\n  margin: 0;\n  width: 1.125rem;\n  height: 1.125rem;\n  flex-shrink: 0;\n  vertical-align: middle;\n  border: 1px solid "+ColorOutline.Var()+";\n  border-radius: "+RadiusSm.Var()+";\n  background-color: "+ColorSurface.Var()+";\n  cursor: pointer;\n  display: inline-grid;\n  place-content: center;\n  transition: background-color 0.15s ease, border-color 0.15s ease;"),
		),
		rule(selector("input[type=\"checkbox\"]:hover:not(:disabled)"),
			rawRule("  border-color: "+ColorMuted.Var()+";"),
		),
		rule(selector("input[type=\"checkbox\"]:checked"),
			rawRule("  background-color: "+ColorPrimary.Var()+";\n  border-color: "+ColorPrimary.Var()+";\n  background-image: url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='%23ffffff'%3E%3Cpath d='M13.78 4.22a.75.75 0 0 1 0 1.06l-7.25 7.25a.75.75 0 0 1-1.06 0L2.22 9.28a.751.751 0 0 1 .018-1.042.751.751 0 0 1 1.042-.018L6 10.94l6.72-6.72a.75.75 0 0 1 1.06 0Z'/%3E%3C/svg%3E\");\n  background-repeat: no-repeat;\n  background-position: center;\n  background-size: 75% 75%;"),
		),
		rule(selector("input[type=\"checkbox\"]:disabled"),
			rawRule("  opacity: 0.5;\n  cursor: not-allowed;\n  background-color: "+ColorSurfaceSunken.Var()+";"),
		),
		// Radio controls: styled with primary color, circular shape and centered indicator.
		rule(selector("input[type=\"radio\"]"),
			rawRule("  -webkit-appearance: none;\n  appearance: none;\n  margin: 0;\n  width: 1.125rem;\n  height: 1.125rem;\n  flex-shrink: 0;\n  vertical-align: middle;\n  border: 1px solid "+ColorOutline.Var()+";\n  border-radius: "+RadiusFull.Var()+";\n  background-color: "+ColorSurface.Var()+";\n  cursor: pointer;\n  display: inline-grid;\n  place-content: center;\n  transition: background-color 0.15s ease, border-color 0.15s ease;"),
		),
		rule(selector("input[type=\"radio\"]:hover:not(:disabled)"),
			rawRule("  border-color: "+ColorMuted.Var()+";"),
		),
		rule(selector("input[type=\"radio\"]:checked"),
			rawRule("  border-color: "+ColorPrimary.Var()+";\n  background-color: "+ColorPrimary.Var()+";\n  background-image: radial-gradient(circle, "+ColorOnPrimary.Var()+" 35%, transparent 40%);"),
		),
		rule(selector("input[type=\"radio\"]:disabled"),
			rawRule("  opacity: 0.5;\n  cursor: not-allowed;\n  background-color: "+ColorSurfaceSunken.Var()+";"),
		),
		// Switch / Toggle controls: checkbox with role="switch" or class="switch".
		rule(selector("input[type=\"checkbox\"][role=\"switch\"], input[type=\"checkbox\"].switch"),
			rawRule("  -webkit-appearance: none;\n  appearance: none;\n  margin: 0;\n  width: 2.375rem;\n  height: 1.375rem;\n  flex-shrink: 0;\n  vertical-align: middle;\n  border: 1px solid "+ColorOutline.Var()+";\n  border-radius: "+RadiusFull.Var()+";\n  background-color: "+ColorSurfaceSunken.Var()+";\n  cursor: pointer;\n  position: relative;\n  display: inline-block;\n  transition: background-color 0.2s ease, border-color 0.2s ease;"),
		),
		rule(selector("input[type=\"checkbox\"][role=\"switch\"]::before, input[type=\"checkbox\"].switch::before"),
			rawRule("  content: \"\";\n  position: absolute;\n  top: 2px;\n  left: 2px;\n  width: 1rem;\n  height: 1rem;\n  border-radius: "+RadiusFull.Var()+";\n  background-color: "+ColorOnPrimary.Var()+";\n  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.25);\n  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1), background-color 0.2s ease;"),
		),
		rule(selector("input[type=\"checkbox\"][role=\"switch\"]:hover:not(:disabled), input[type=\"checkbox\"].switch:hover:not(:disabled)"),
			rawRule("  border-color: "+ColorMuted.Var()+";"),
		),
		rule(selector("input[type=\"checkbox\"][role=\"switch\"]:checked, input[type=\"checkbox\"].switch:checked"),
			rawRule("  background-color: "+ColorPrimary.Var()+";\n  border-color: "+ColorPrimary.Var()+";\n  background-image: none;"),
		),
		rule(selector("input[type=\"checkbox\"][role=\"switch\"]:checked::before, input[type=\"checkbox\"].switch:checked::before"),
			rawRule("  transform: translateX(1rem);"),
		),
		rule(selector("input[type=\"checkbox\"][role=\"switch\"]:disabled, input[type=\"checkbox\"].switch:disabled"),
			rawRule("  opacity: 0.5;\n  cursor: not-allowed;"),
		),


		// Firefox ships placeholders at opacity 0.54, so the same muted colour
		// reads lighter there than on Chrome or Safari.
		rule(selector("::placeholder"),
			rawRule("  opacity: 1;"),
		),
		// Every engine renders a monospace default about 3px smaller than the
		// surrounding text — a historical quirk of the `monospace` keyword.
		// Re-stating the family resets that scaling, and 1em pins the size back
		// to the text around it.
		rule(selector("code, kbd, samp, pre"),
			rawRule("  font-family: monospace;\n  font-size: 1em;"),
		),
		rule(selector("ol, ul"),
			rawRule("  list-style: none;"),
			margin(zero),
			padding(zero),
		),
		// A <summary> draws a disclosure triangle the component did not ask
		// for and cannot style. Components supply their own affordance.
		rule(selector("summary"),
			rawRule("  list-style: none;"),
		),
		rule(selector("summary::-webkit-details-marker"),
			display(none),
		),
		rule(selector("img, svg, video"),
			display(block),
			maxWidth(pct(100)),
		),
		// Author styles outrank the user agent stylesheet whatever their layer,
		// so the display: block above silently defeats the UA's own
		// `[hidden] { display: none }` and an <img hidden> keeps rendering.
		// Restating the rule here is what makes the attribute work again.
		rule(selector("[hidden]"),
			display(none),
		),
		rule(selector(":root"),
			rawRule("color-scheme: light dark;"),
		),
		rule(selector("[data-theme=\"light\"]"),
			rawRule("color-scheme: light;"),
		),
		rule(selector("[data-theme=\"dark\"]"),
			rawRule("color-scheme: dark;"),
		),
	}
}
