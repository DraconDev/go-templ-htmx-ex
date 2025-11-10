# Browser Auto-Reload Setup Guide for Go + Templ + HTMX

## ⚠️ Important: HTMX and Browser Reload Considerations

For HTMX applications, traditional browser auto-reload can be **problematic** because:

- **HTMX preserves form state** across partial updates
- **Dynamic content** (modals, pagination, filtered lists) gets lost
- **SPA-like experience** is disrupted with full page reloads
- **Developer workflow** suffers from constant context loss

## Auto Reload Options

### Option 1: Air Proxy with Browser Reload (CURRENT SETUP)
**Command**: `make air` (proxy enabled)

**What happens**:
- Air creates a proxy server on port 4200
- Your app runs on port 8080  
- Browser connects to proxy (http://localhost:4200) for **automatic reload**
- **Full page reload** when code changes

**Pros**:
- ✅ Automatic browser refresh
- ✅ Simple setup

**Cons**:
- ❌ **Destroys HTMX state** - forms reset, modals close
- ❌ **Disrupts development flow** with constant page reloads
- ❌ **Loses dynamic content** like filtered results, pagination

### Option 2: No Browser Reload (Alternative for HTMX)
**Command**: `make air` (with proxy disabled in .air.toml)

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

## Quick Setup Commands

```bash
# Start development with browser auto-reload (CURRENT)
make air
# Then open: http://localhost:4200

# For HTMX work (manual refresh):
# Edit .air.toml to disable proxy
# Then: make air
# Then open: http://localhost:8080
```

## Testing Your Setup

### With Browser Auto-Reload (Port 4200):
1. **Start the dev server**: `make air`
2. **Open browser**: http://localhost:4200
3. **Make a code change** in any `.go` or `.templ` file
4. **Watch the browser** automatically reload with changes

### Without Browser Auto-Reload (Port 8080):
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

### When Working with Forms/HTMX:
- Use **manual refresh mode** (disable proxy)
- Open: http://localhost:8080
- **Preserves form state** during development

### When Working with Static Content:
- Use **auto-reload mode** (enable proxy) 
- Open: http://localhost:4200
- **Quick iteration** for non-interactive pages

## Switching Between Modes

### Enable Browser Auto-Reload:
```bash
# Edit .air.toml:
[proxy]
  enabled = true
  proxy_port = 4200

# Then:
make air
# Browser: http://localhost:4200
```

### Disable Browser Auto-Reload:
```bash
# Edit .air.toml:
[proxy]
  enabled = false

# Then:
make air
# Browser: http://localhost:8080
```

## Troubleshooting

**If server won't start (port 8080 or 4200 in use)**:
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

**If browser reload isn't working**:
- Ensure you're accessing http://localhost:4200 (not 8080)
- Check browser console for errors
- Verify proxy is enabled in .air.toml
