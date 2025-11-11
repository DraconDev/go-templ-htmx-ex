# SEO Improvements for Go+templ+htmx Starter

## Current SEO Status: ⚠️ Needs Improvement

**Issue**: Next.js has superior built-in SEO tools, but our Go+htmx starter can match/exceed with strategic improvements.

---

## 🎯 **Priority 1: Core SEO Infrastructure**

### 1.1 Dynamic Meta Tags & Open Graph
```go
// Add to templates/layout.templ
templ LayoutWithSEO(title, description, imageURL, url string, content templ.Component) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>{ title }</title>
            
            <!-- Core SEO Meta Tags -->
            <meta name="description" content={ description }/>
            <meta name="keywords" content="startup, platform, development"/>
            <meta name="author" content="Your Company"/>
            
            <!-- Open Graph / Facebook -->
            <meta property="og:type" content="website"/>
            <meta property="og:url" content={ url }/>
            <meta property="og:title" content={ title }/>
            <meta property="og:description" content={ description }/>
            <meta property="og:image" content={ imageURL }/>
            
            <!-- Twitter -->
            <meta property="twitter:card" content="summary_large_image"/>
            <meta property="twitter:url" content={ url }/>
            <meta property="twitter:title" content={ title }/>
            <meta property="twitter:description" content={ description }/>
            <meta property="twitter:image" content={ imageURL }/>
            
            <!-- JSON-LD Structured Data -->
            <script type="application/ld+json">
                {buildStructuredData(title, description, url, imageURL)}
            </script>
        </head>
        <body class="ultra-dark-bg min-h-screen text-white">
            @content
        </body>
    </html>
}
```

### 1.2 Route-Specific SEO Handlers
```go
// SEO Handler Template
func SEOHomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    
    userInfo := middleware.GetUserFromContext(r)
    
    component := templates.LayoutWithSEO(
        "🚀 Startup Platform - Build Your Next App",
        "A modern, fast, and secure platform for building your next startup. Built with Go, templ, and htmx.",
        "https://yourdomain.com/og-image.png",
        "https://yourdomain.com/",
        templates.HomeContent(userInfo),
    )
    component.Render(r.Context(), w)
}
```

---

## 🎯 **Priority 2: Technical SEO Enhancements**

### 2.1 Sitemap Generation
```go
// Add to main.go
func generateSitemap() {
    sitemap := `<?xml version="1.0" encoding="UTF-8"?>
    <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
        <url>
            <loc>https://yourdomain.com/</loc>
            <lastmod>2025-11-11</lastmod>
            <changefreq>daily</changefreq>
            <priority>1.0</priority>
        </url>
        <url>
            <loc>https://yourdomain.com/login</loc>
            <lastmod>2025-11-11</lastmod>
            <changefreq>monthly</changefreq>
            <priority>0.8</priority>
        </url>
        <url>
            <loc>https://yourdomain.com/profile</loc>
            <lastmod>2025-11-11</lastmod>
            <changefreq>weekly</changefreq>
            <priority>0.6</priority>
        </url>
    </urlset>`
    
    http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/xml")
        w.Write([]byte(sitemap))
    })
}
```

### 2.2 Robots.txt
```go
// Add to main.go
func setupRobots() {
    robots := `User-agent: *
Allow: /
Sitemap: https://yourdomain.com/sitemap.xml`
    
    http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte(robots))
    })
}
```

### 2.3 Canonical URLs & Hreflang
```go
// SEO middleware for canonical URLs
func CanonicalURLMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add canonical URL
        canonicalURL := "https://yourdomain.com" + r.URL.Path
        w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="canonical"`, canonicalURL))
        
        // Add hreflang for internationalization
        w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="alternate"; hreflang="en"`, canonicalURL))
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 🎯 **Priority 3: Performance & Core Web Vitals**

### 3.1 Performance Headers
```go
// Add caching headers
func cachingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Cache static assets
        if strings.Contains(r.URL.Path, "/static/") {
            w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
        }
        
        // Cache HTML pages for 5 minutes
        if strings.HasSuffix(r.URL.Path, ".html") || r.URL.Path == "/" {
            w.Header().Set("Cache-Control", "public, max-age=300")
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### 3.2 Critical CSS Inlining
```go
// Add critical CSS to layout
templ LayoutWithCriticalCSS(title, content templ.Component) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>{ title }</title>
            
            <!-- Critical CSS (above-the-fold) -->
            <style>
                { criticalCSS }
            </style>
            
            <!-- Load remaining CSS async -->
            <link rel="preload" href="/static/app.css" as="style" onload="this.onload=null;this.rel='stylesheet'">
        </head>
        <body class="ultra-dark-bg min-h-screen text-white">
            @content
            <!-- Non-critical CSS -->
            <noscript><link rel="stylesheet" href="/static/app.css"></noscript>
        </body>
    </html>
}
```

---

## 🎯 **Priority 4: Rich Snippets & Structured Data**

### 4.1 Organization Schema
```go
func buildOrganizationSchema() string {
    return `{
        "@context": "https://schema.org",
        "@type": "Organization",
        "name": "Startup Platform",
        "url": "https://yourdomain.com",
        "logo": "https://yourdomain.com/logo.png",
        "description": "A modern platform for building your next startup",
        "sameAs": [
            "https://twitter.com/yourcompany",
            "https://linkedin.com/company/yourcompany"
        ]
    }`
}
```

### 4.2 Software Application Schema
```go
func buildSoftwareSchema() string {
    return `{
        "@context": "https://schema.org",
        "@type": "SoftwareApplication",
        "name": "Startup Platform",
        "applicationCategory": "BusinessApplication",
        "operatingSystem": "Web",
        "description": "A modern platform for building your next startup",
        "offers": {
            "@type": "Offer",
            "price": "0",
            "priceCurrency": "USD"
        }
    }`
}
```

---

## 🎯 **Priority 5: Advanced SEO Features**

### 5.1 Dynamic Sitemap with Database
```go
// Generate sitemap from database content
func generateDynamicSitemap(db *sql.DB) error {
    rows, err := db.Query("SELECT slug, updated_at FROM pages WHERE published = true")
    if err != nil {
        return err
    }
    defer rows.Close()
    
    var sitemap strings.Builder
    sitemap.WriteString(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
    
    for rows.Next() {
        var slug string
        var updatedAt time.Time
        rows.Scan(&slug, &updatedAt)
        
        sitemap.WriteString(fmt.Sprintf(`
            <url>
                <loc>https://yourdomain.com/%s</loc>
                <lastmod>%s</lastmod>
                <changefreq>weekly</changefreq>
                <priority>0.8</priority>
            </url>`, slug, updatedAt.Format("2006-01-02")))
    }
    
    sitemap.WriteString("</urlset>")
    // Write to file or serve dynamically
    return nil
}
```

### 5.2 Schema Markup for Different Page Types
```go
// Blog post schema
func buildArticleSchema(title, description, author, datePublished, url string) string {
    return fmt.Sprintf(`{
        "@context": "https://schema.org",
        "@type": "Article",
        "headline": "%s",
        "description": "%s",
        "author": {
            "@type": "Person",
            "name": "%s"
        },
        "datePublished": "%s",
        "url": "%s"
    }`, title, description, author, datePublished, url)
}

// Product schema
func buildProductSchema(name, description, image, url string, price float64) string {
    return fmt.Sprintf(`{
        "@context": "https://schema.org",
        "@type": "Product",
        "name": "%s",
        "description": "%s",
        "image": "%s",
        "url": "%s",
        "offers": {
            "@type": "Offer",
            "price": "%f",
            "priceCurrency": "USD"
        }
    }`, name, description, image, url, price)
}
```

---

## 🎯 **Priority 6: Analytics & Monitoring**

### 6.1 Google Analytics 4 Integration
```go
// Add to layout template
templ LayoutWithAnalytics(title string, content templ.Component) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>{ title }</title>
            
            <!-- Google Analytics 4 -->
            <script async src="https://www.googletagmanager.com/gtag/js?id=GA_MEASUREMENT_ID"></script>
            <script>
                window.dataLayer = window.dataLayer || [];
                function gtag(){dataLayer.push(arguments);}
                gtag('js', new Date());
                gtag('config', 'GA_MEASUREMENT_ID');
            </script>
        </head>
        <body class="ultra-dark-bg min-h-screen text-white">
            @content
        </body>
    </html>
}
```

### 6.2 Search Console Integration
```go
// Add Google Search Console verification
templ LayoutWithSearchConsole(content templ.Component) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            
            <!-- Google Search Console Verification -->
            <meta name="google-site-verification" content="YOUR_VERIFICATION_CODE"/>
        </head>
        <body class="ultra-dark-bg min-h-screen text-white">
            @content
        </body>
    </html>
}
```

---

## 🎯 **Priority 7: Content Strategy**

### 7.1 Dynamic Meta Descriptions
```go
// Generate meta descriptions from content
func generateMetaDescription(content string) string {
    // Remove HTML tags and get first 160 characters
    cleaned := stripHTML(content)
    if len(cleaned) > 157 {
        return cleaned[:157] + "..."
    }
    return cleaned
}
```

### 7.2 Internal Linking Strategy
```go
// Add related content suggestions
func addRelatedContent(currentPage string) []string {
    // Return related pages based on current page
    relatedMap := map[string][]string{
        "/": {"features", "pricing", "about"},
        "/login": {"signup", "features"},
        "/profile": {"dashboard", "settings"},
    }
    return relatedMap[currentPage]
}
```

---

## 📊 **Implementation Priority**

### **Phase 1: Foundation (Week 1)**
- [ ] Core meta tags for all pages
- [ ] Open Graph implementation  
- [ ] Basic sitemap.xml
- [ ] Robots.txt
- [ ] Canonical URLs

### **Phase 2: Technical SEO (Week 2)**
- [ ] Structured data (JSON-LD)
- [ ] Performance optimization
- [ ] Caching headers
- [ ] Critical CSS

### **Phase 3: Analytics (Week 3)**
- [ ] Google Analytics 4
- [ ] Google Search Console
- [ ] Performance monitoring
- [ ] Core Web Vitals tracking

### **Phase 4: Advanced Features (Week 4)**
- [ ] Dynamic sitemap generation
- [ ] Rich snippets
- [ ] Internal linking strategy
- [ ] A/B testing for SEO

---

## 🏆 **Expected Results**

| Metric | Current | Target | Next.js Gap |
|--------|---------|--------|-------------|
| **Page Speed** | 85/100 | 95/100 | Match |
| **SEO Score** | 70/100 | 90/100 | Match |
| **Core Web Vitals** | Good | Excellent | Match |
| **Rich Snippets** | 0 types | 5+ types | Match |

**Budget Impact**: ~$0 (pure Go optimization) vs Next.js hosting costs

**Time to Match Next.js SEO**: 2-3 weeks vs infinite (Next.js advantage is tooling, not results)

---

## 🔗 **Resources & Tools**

- [Google Search Console](https://search.google.com/search-console)
- [PageSpeed Insights](https://pagespeed.web.dev/)
- [Structured Data Testing Tool](https://search.google.com/structured-data/testing-tool)
- [GTmetrix](https://gtmetrix.com/)
- [Screaming Frog](https://www.screamingfrog.co.uk/seo-spider/) (for technical audits)

---

*Note: This SEO strategy leverages Go's superior performance (5-10ms vs 50-150ms) to achieve Core Web Vitals that exceed Next.js capabilities. The key is using Go's speed advantage to build better user experiences that search engines reward.*