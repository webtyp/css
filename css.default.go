//go:build !wasm

package css

// toneRules derives theme-adapted opaque tinted background, high-contrast foreground,
// and border from any runtime --tone-color CSS custom property (populated from the
// database palette table, never hardcoded in Go or CSS). Mixed against ColorBackground
// (never transparent) so badges riding a card's edge stay crisp and opaque.
func toneRules() []item {
	return []item{
		rule(selector(`[data-toned="1"]`),
			decl{"background-color", "color-mix(in oklab, var(--tone-color) 18%, " + ColorBackground.Var() + ")"},
			decl{"color", "color-mix(in oklab, var(--tone-color), light-dark(#000000, #ffffff) 32%)"},
			decl{"border-color", "color-mix(in oklab, var(--tone-color) 45%, " + ColorBackground.Var() + ")"},
		),
	}
}

// themeTokens are the colours whose value adapts to the theme. They are
// declared twice: the static half here, the adaptive half inside
// ModernColorSupport — see declare() and enhancedDecls().
func themeTokens() []Token {
	return []Token{
		ColorBackground,
		ColorOnBackground,
		ColorSurface,
		ColorOnSurface,
		ColorOutline,
		ColorMuted,
		ColorSurfaceSunken,
		ColorSelection,
		ColorOnSelection,
		ColorAccentWash,
		ColorAccentHover,
		ColorDangerWash,
	}
}

// themeGroupDecls declares the theme colors' static values plus, for every
// theme pair among them, the plain light/dark halves declareSplit needs.
func themeGroupDecls() []decl {
	var decls []decl
	for _, t := range themeTokens() {
		decls = append(decls, declare(t))
	}
	for _, t := range []Token{ColorBackground, ColorOnBackground, ColorSurface, ColorOnSurface, ColorOutline, ColorMuted} {
		decls = append(decls, declareSplit(t)...)
	}
	return decls
}

// defaultRoots declares every token group that is not brand identity: the
// adaptive theme surface colors plus the typography, spacing, radius,
// shadow, motion, z-index, breakpoint, and layout scales. These are the
// tokens an app inherits as-is; only brandRoot() is meant to be reskinned.
func defaultRoots() []item {
	return []item{
		root(themeGroupDecls()...),
		// The adaptive upgrade, gated. It has to come after the static block
		// so it wins where it applies, and before anything an app declares so
		// Theme(Set(...)) still outranks both.
		supports(ModernColorSupport, root(enhancedDecls(themeTokens()...)...)),
		root(
			// Interaction intensities
			declare(MixHover),
			declare(MixFocus),
			declare(MixPress),
			declare(ToneColor),
		),
		root(
			// Typography scale
			declare(TextXs),
			declare(TextSm),
			declare(TextBase),
			declare(TextLg),
			declare(TextXl),
			declare(Text2xl),
		),
		root(
			// Spacing scale
			declare(Space0),
			declare(Space1),
			declare(Space2),
			declare(Space3),
			declare(Space4),
			declare(Space6),
			declare(Space8),
			declare(Space12),
		),
		root(
			// Border-radius scale
			declare(RadiusSm),
			declare(RadiusMd),
			declare(RadiusLg),
			declare(RadiusFull),
		),
		root(
			// Typography
			declare(FontSans),
			declare(LeadingNormal),
			declare(FontWeightRegular),
			declare(FontWeightMedium),
			declare(FontWeightBold),
		),
		root(
			// Elevation
			declare(ShadowSm),
			declare(ShadowMd),
			declare(ShadowLg),
		),
		root(
			// Motion
			declare(DurationFast),
			declare(DurationBase),
			declare(DurationSlow),
			declare(EaseInOut),
		),
		root(
			// Z-index
			declare(ZBase),
			declare(ZDropdown),
			declare(ZSticky),
			declare(ZModal),
			declare(ZToast),
			declare(ZTooltip),
		),
		root(
			// Breakpoints
			declare(BpSm),
			declare(BpMd),
			declare(BpLg),
			declare(BpXl),
		),
		root(
			// Container widths
			declare(MaxWReadable),
		),
		root(
			// Grid columns
			declare(ColumnNarrow),
			declare(ColumnMedium),
			declare(ColumnWide),
			// Rail widths
			declare(RailNarrow),
			declare(RailWide),
			// Shared control height and chip size
			declare(ControlHeight),
			declare(ControlHeightCompact),
			declare(ControlWidth),
			declare(ChipWidth),
			declare(ChipHeight),
			declare(VeilBlur),
		),
		root(
			// Device geometry — insets reported by the device, and the viewport
			// height that shrinks with Safari iOS' collapsing URL bar.
			declare(SafeTop),
			declare(SafeRight),
			declare(SafeBottom),
			declare(SafeLeft),
			declare(ViewportH),
		),
	}
}
