# Best golangci-lint Linters for Real-World Usage

## ✅ Top-Tier Linters (Highly Recommended)

### **Security & Safety (Must-Have)**
- **gosec** - Finds security vulnerabilities (SQL injection, weak crypto, etc.)
- **errcheck** - Unchecked errors ⭐ **(already caught our issues!)**
- **govulncheck** - Checks for known vulnerabilities

### **Code Quality (Essential)**
- **staticcheck** - Comprehensive static analysis ⭐ **(excellent quality)**
- **govet** - Basic but essential checks ⭐ **(already enabled)**
- **unused** - Finds unused code ⭐ **(already enabled)**
- **ineffassign** - Finds dead assignments ⭐ **(already enabled)**

### **Code Style (Quality of Life)**
- **gofmt/gofumpt** - Better formatting than basic gofmt
- **stylecheck** - Modern replacement for golint
- **revive** - Fast, configurable, extensible linter

### **Performance (Optimize)**
- **prealloc** - Suggests slice preallocations
- **perfsprint** - Faster string formatting alternatives

## 🎯 Recommended Configuration

### **Essential Set (Safe for most projects)**
```yaml
linters:
  enable:
    - errcheck
    - staticcheck
    - govet
    - unused
    - ineffassign
    - gosimple
    - revive
  disable:
    - gochecknoglobals  # Too strict for most projects
    - gochecknoinits    # Sometimes inits are needed
```

### **Production Set (For serious projects)**
```yaml
linters:
  enable:
    - errcheck
    - staticcheck
    - govet
    - unused
    - ineffassign
    - gosimple
    - revive
    - gosec
    - gofmt
    - goimports
    - stylecheck
```

### **Security-Focused Set**
```yaml
linters:
  enable:
    - errcheck
    - staticcheck
    - govet
    - gosec
    - govulncheck
    - gofmt
    - goimports
```

## 🔧 Configuration File (.golangci.yml)

```yaml
run:
  timeout: 5m
  tests: true
  skip-dirs:
    - vendor
    - node_modules

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  
  staticcheck:
    checks:
      - "all"
      - "-ST1000"  # Skip package comment checks
  
  revive:
    min-confidence: 0

linters:
  enable:
    - errcheck
    - staticcheck
    - govet
    - unused
    - ineffassign
    - gosimple
    - revive
    - gofmt
    - goimports
  disable:
    - gocyclo
    - funlen
    - gomnd

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
```

## 💡 Pro Tips

1. **Start small** - Enable a few linters first
2. **Use `--fast`** for quick checks during development
3. **Fix errors progressively** - Don't try to fix everything at once
4. **VS Code integration** - Works automatically with save
5. **CI/CD integration** - Use in GitHub Actions for quality gates

## 🚀 Quick Start

```bash
# Test current config
golangci-lint run --fast

# See what would be fixed
golangci-lint run --fast --disable-all --enable=errcheck

# Auto-fix where possible
golangci-lint run --fast --fix