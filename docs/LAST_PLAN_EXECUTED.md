# PLAN — `css.DefaultLight()`: an app's default color scheme, declared in the stylesheet

> Master: `/home/cesar/.claude/plans/si-la-ui-se-drifting-treasure.md` (track A) · 2026-10-08

## Problem

`resetRules()` emits `:root { color-scheme: light dark }`, so every page follows the OS. An app
that wants light by default (mjosefa-cms) had no typed way to say so in CSS (THEMING.md "Tema
único" points at raw `:root { color-scheme: light }`), so it forced
`SetDocumentAttr("data-theme", "light")` from its WASM entry point. Its pre-login page is static
HTML+CSS with no WASM, so the login followed the OS (dark) while the app was light.

## API gate

1. **Prior art.** CSS itself: `color-scheme` on `:root`. MUI `CssVarsProvider`
   `defaultMode: "system" | "light" | "dark"`. next-themes `defaultTheme="system"`. Bootstrap 5.3
   `data-bs-theme` on `<html>`. All of them resolve it at runtime or through a document attribute.
   Here a pre-login page carries no JS and `html.Document` does not know the theme, so the
   declaration belongs in the stylesheet every page already loads.
2. **Names.** `css.Theme(css.DefaultLight())` reads "theme, light by default". It names the intent
   (a default a theme toggle may still override), takes no parameter so no invalid value exists,
   and `DefaultDark` is not exported because nobody uses it.
3. **Ledger.** Concepts +1 · files to make an app light: 1 (`config/css.go`), before 1
   (`web/client.go`, which never reached static pages) · call-site lines 0 · ways to do it 0 net
   (the app's `SetDocumentAttr` dies in the same wave).
4. **Where.** `webtyp/css` already owns the `color-scheme` rule and `Theme()`.
5. **Deletes.** mjosefa-cms' `SetDocumentAttr("data-theme","light")` and its comment; the raw-CSS
   advice in THEMING.md "Tema único".

## Steps

1. Red test (`css_test.go`): `Theme(DefaultLight())` contains
   `:root:not([data-theme]) { color-scheme: light; }`; `Theme()` does not.
2. `Override` gains `scheme string`; `DefaultLight()` sets it; `Theme()` appends the rule after the
   `:root` tail. `:not([data-theme])` keeps it out of the way of `[data-theme="dark"]` whatever the
   layer.
3. Docs: SPECS.md §5, THEMING.md "Tema único".
4. `gotest`; `gopush`.
