# 🧩 Simple Improvements: Template Organization (Implemented)

This starter now follows a clean, minimal, and extensible template structure so every future project stays organized without adding complexity.

## 1️⃣ Template Organization

### Goals

- Make it obvious where to edit:
  - Layout / shell
  - Navigation
  - Pages
  - Reusable components
- Keep it minimal (no over-engineering)
- Make it easy to customize per startup

### Current Structure

```bash
templates/
  layout.templ              # Base layout (SEO + shell + shared styles)
  home.templ                # Home page content
  login.templ               # Login page content
  profile.templ             # Profile page content
  admin_dashboard.templ     # Admin dashboard page content
  auth_callback.templ       # OAuth callback page content

  components/
    navigation.templ        # Shared navigation components (extensible)
```

### How Layout Is Used

[`templates/layout.templ`](templates/layout.templ:15) defines:
- `Layout(title string, description string, navigation templ.Component, content templ.Component)`
  - Handles:
    - `<html>` / `<head>` / `<body>`
    - SEO meta tags (title, description, OG/Twitter, JSON-LD)
    - Shared CSS (dark theme, glass morphism, glow effects)
  - Renders:
    - `@navigation` (passed from handlers)
    - `@content` (page-specific content)

Handlers use it like:

```go
// Home
component := templates.Layout(
    "Home",
    "Production-ready startup platform with Google OAuth, PostgreSQL database, and admin dashboard. Built with Go + HTMX + Templ.",
    navigation,
    templates.HomeContent(),
)

// Profile
component := templates.Layout(
    "Profile",
    "User profile page with authentication details and account management.",
    templates.NavigationLoggedIn(userInfo),
    templates.ProfileContent(userInfo.Name, userInfo.Email, userInfo.Picture),
)

// Login
component := templates.Layout(
    "Login",
    "Secure authentication page with Google OAuth integration for user access.",
    templates.NavigationLoggedOut(),
    templates.LoginContent(),
)

// Admin
component := templates.Layout(
    "Admin Dashboard",
    "Administrative dashboard with user statistics, analytics, and platform management tools.",
    templates.NavigationLoggedIn(userInfo),
    templates.AdminDashboardContent(userInfo, dashboardData),
)

// Auth callback
component := templates.Layout(
    "Authenticating",
    "Authentication processing page for OAuth callback and session establishment.",
    templates.NavigationLoggedOut(),
    templates.AuthCallbackContent(),
)
```

This matches Simple Improvement #1:
- All pages consistently use a single layout entrypoint
- SEO + layout concerns live in one place
- Navigation and content are passed in explicitly

### Navigation Components

[`templates/layout.templ`](templates/layout.templ:130) includes:

- `NavigationLoggedIn(user UserInfo)`
- `NavigationLoggedOut()`
- `UserAvatar(user UserInfo)`
- `UserInitials(user UserInfo)`

These:
- Provide a consistent top navigation bar
- Are used directly from handlers based on auth state
- Are easy to override or extract later without changing handlers

[`templates/components/navigation.templ`](templates/components/navigation.templ:1) is added as a foundation for future extraction if you want:
- A more modular nav system
- Project-specific navigation variants

For now, the starter:
- Keeps `NavigationLoggedIn` / `NavigationLoggedOut` in `layout.templ` (simple, centralized)
- Adds `components/navigation.templ` as a clean place to extend if needed

This satisfies the “simple improvement” philosophy:
- Clear organization
- No hard coupling to complex component trees
- Easy future refactor path

---

## 2️⃣ How to Use This in New Projects

When you fork this starter for a new startup:

1. Update layout SEO + branding:
   - Edit [`templates/layout.templ`](templates/layout.templ:15) `title` suffix, meta defaults, JSON-LD.
2. Update navigation:
   - Adjust `NavigationLoggedOut` / `NavigationLoggedIn` links and labels.
3. Update page content:
   - Edit `home.templ`, `login.templ`, `profile.templ`, `admin_dashboard.templ` only.
4. (Optional) Extract more components:
   - Create `templates/components/*.templ` for repeated UI chunks.

You never touch handlers for basic branding/SEO/visual changes.
All UI / layout concerns stay in `templates/`.

---

## 3️⃣ Why This Matches Simple Improvement #1

- ✅ Single layout entrypoint with SEO + shell
- ✅ Clear separation: layout vs pages vs (optional) components
- ✅ Handlers only compose: `Layout(title, description, navigation, content)`
- ✅ Minimal cognitive load for new projects
- ✅ Easy to customize per startup without rewriting core logic
- ✅ No performance or security downsides (purely structural)

This is the intended baseline for all future projects using this starter.
