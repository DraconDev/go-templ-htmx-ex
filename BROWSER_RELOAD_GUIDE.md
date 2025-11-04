# Browser Auto-Reload Setup Guide for Go + Templ + HTMX

## ⚠️ Important: HTMX and Browser Reload Considerations

For HTMX applications, traditional browser auto-reload can be **problematic** because:

- **HTMX preserves form state** across partial updates
- **Dynamic content** (modals, pagination, filtered lists) gets lost
- **SPA-like experience** is disrupted with full page reloads
- **Developer workflow** suffers from constant context loss

## Auto Reload Options

### Option 1: Air Proxy with Browser Reload
**Command**: `make air` (now with proxy enabled)

**What happens**:
- Air creates a proxy server on port 4200
- Your app runs on port 8080  
- Browser connects to proxy (http://localhost:4200)
- **Full page reload** when code changes

**Pros**:
- ✅ Automatic browser refresh
- ✅ Simple setup

**Cons**:
- ❌ **Destroys HTMX state** - forms reset, modals close
- ❌ **Disrupts development flow** with constant page reloads
- ❌ **Loses dynamic content** like filtered results, pagination

### Option 2: No Browser Reload (Recommended for HTMX)
**Command**: `make air` (with proxy disabled)

**What happens**:
- Only server reloads when code changes
- Browser shows updated content on **manual refresh** (F5)
- **HTMX state preserved** between manual refreshes

**Pros**:
- ✅ **Perfect for HTMX development**
- ✅ **Preserves form state and dynamic content**
- ✅ **No disruptive page reloads**
- ✅ **Fast server compilation** (~3-4ms)

**Cons**:
- ❌ Requires manual browser refresh

### Option 3: Hybrid Approach (Advanced)
Use browser reload only for **non-interactive pages**:
- Static content pages: Browser reload enabled
- HTMX forms/pages: Browser reload disabled

## Recommended Setup for HTMX Development

### For HTMX Applications:
1. **Disable proxy** (current setup)
2. **Use manual refresh** (F5 or Ctrl+R)
3. **Benefit from fast server reload** (~3ms builds)

### For Traditional Multi-page Applications:
1. **Enable proxy** (Air proxy with browser reload)
2. **Full page reload** works well
3. **No state preservation needed**

## Quick Setup Commands

```bash
# For HTMX development (recommended)
make air

# For traditional development
# Edit .air.toml to enable proxy
# Then: make air
```

## Testing Your Setup

1. **Start the dev server**: `make air`
2. **Open browser**: http://localhost:8080
3. **Make a code change** in any `.go` or `.templ` file
4. **Check terminal**: Should show fast rebuild (~3ms)
5. **Refresh browser**: Press F5 to see changes

## Performance Notes

- **Server reload**: ~3-4ms (extremely fast)
- **Templ generation**: Included in build process
- **File watching**: Instant detection
- **Memory usage**: Minimal with Air

## Best Practices for HTMX Development

1. **Use manual refresh** during active form editing
2. **Test HTMX interactions** before and after code changes
3. **Keep browser console open** to monitor network requests
4. **Use browser dev tools** for debugging HTMX behavior

## Troubleshooting

**If server won't start (port 8080 in use)**:
```bash
pkill -f "microservice-test"
pkill -f "air"
sleep 2
make air
```

**If changes aren't detected**:
- Ensure you're editing `.go` or `.templ` files
- Check that Air is running (`make air`)
- Look for build errors in terminal

**If Templ components aren't updating**:
- Air automatically runs `templ generate` before each build
- Check Templ syntax errors in build output
- Manually run `make generate` if needed
