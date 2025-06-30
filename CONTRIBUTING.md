# 🤝 Contributing to gin-mcp

Thank you for your interest in contributing to gin-mcp! We believe that open source is about collaboration, learning, and building something beautiful together. This guide will help you get started.

## 🌟 Our Values

We are guided by the pursuit of **excellence and beauty** in everything we do. This means:

- **Quality over quantity** - Every contribution should be thoughtful and well-crafted
- **Clarity and simplicity** - Code should be self-documenting and easy to understand
- **Inclusivity** - We welcome contributors from all backgrounds and experience levels
- **Continuous improvement** - We're always learning and growing together

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- CGO enabled (for Go plugin support)
- Git

### Development Setup

1. **Fork the repository**
   ```bash
   git clone https://github.com/your-username/gin-mcp.git
   cd gin-mcp
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Build the project**
   ```bash
   make build
   ```

4. **Run tests**
   ```bash
   make test
   ```

## 📋 How to Contribute

### 🐛 Reporting Bugs

Before creating a bug report, please:

1. Check if the issue has already been reported
2. Try to reproduce the issue with the latest version
3. Include as much detail as possible

**Bug Report Template:**
```markdown
## Bug Description
Brief description of the issue

## Steps to Reproduce
1. Step one
2. Step two
3. Step three

## Expected Behavior
What you expected to happen

## Actual Behavior
What actually happened

## Environment
- OS: [e.g., macOS, Linux, Windows]
- Go version: [e.g., 1.21.0]
- gin-mcp version: [e.g., commit hash]

## Additional Information
Any other context, logs, or screenshots
```

### 💡 Suggesting Features

We welcome feature suggestions! Please:

1. Check if the feature has already been requested
2. Explain the problem you're trying to solve
3. Describe your proposed solution
4. Consider the impact on existing functionality

### 🔧 Code Contributions

#### Before You Start

1. **Check existing issues** - Look for issues labeled `good first issue` or `help wanted`
2. **Discuss your approach** - Comment on the issue or create a discussion
3. **Keep it focused** - One feature or fix per pull request

#### Development Workflow

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write clear, self-documenting code
   - Add tests for new functionality
   - Update documentation as needed
   - Follow our coding standards

3. **Test your changes**
   ```bash
   make test
   make test-coverage
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push and create a pull request**
   ```bash
   git push origin feature/your-feature-name
   ```

#### Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

Examples:
```
feat: add support for gRPC model endpoints
fix: resolve timeout issue in Go plugin execution
docs: update API documentation with new endpoints
```

### 📝 Documentation

We value clear, comprehensive documentation. When contributing:

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Include examples and use cases
- Keep documentation in sync with code changes

## 🎨 Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Write meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and concise

### Go Plugin Code (for tools)

- Follow standard Go formatting with `gofmt`
- Use meaningful function and variable names
- Include proper error handling
- Use type hints where appropriate
- Include docstrings for functions
- Handle errors gracefully

### General Principles

- **Readability** - Code should be easy to read and understand
- **Maintainability** - Consider future maintenance and updates
- **Performance** - Write efficient code, but prioritize clarity
- **Security** - Follow security best practices

## 🧪 Testing

### Writing Tests

- Write tests for new functionality
- Ensure good test coverage
- Use descriptive test names
- Test both success and error cases

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./handlers -v
```

## 📦 Pull Request Process

### Before Submitting

1. **Self-review** - Review your own code as if you were reviewing someone else's
2. **Test thoroughly** - Ensure all tests pass and new functionality works
3. **Update documentation** - Keep docs in sync with code changes
4. **Check formatting** - Run `make fmt` to ensure consistent formatting

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Other (please describe)

## Testing
- [ ] Added tests for new functionality
- [ ] All existing tests pass
- [ ] Tested manually

## Documentation
- [ ] Updated README.md
- [ ] Added inline comments
- [ ] Updated API documentation

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] No breaking changes (or documented if necessary)
- [ ] Commit messages follow conventional format
```

### Review Process

1. **Automated checks** - CI/CD pipeline runs tests and linting
2. **Code review** - At least one maintainer reviews the PR
3. **Discussion** - Address any feedback or questions
4. **Merge** - Once approved, your PR will be merged

## 🏷️ Issue Labels

We use labels to organize issues:

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Improvements or additions to documentation
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention is needed
- `question` - Further information is requested
- `wontfix` - This will not be worked on

## 🎯 Getting Help

### Questions and Discussions

- **GitHub Discussions** - For questions, ideas, and general discussion
- **GitHub Issues** - For bug reports and feature requests
- **Code of Conduct** - Please read and follow our community guidelines

### Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://github.com/gin-gonic/gin)
- [Effective Go](https://golang.org/doc/effective_go.html)

## 🙏 Recognition

We believe in recognizing and celebrating contributions:

- Contributors are listed in the README
- Significant contributions are highlighted in release notes
- We maintain a contributors hall of fame

## 📄 License

By contributing to gin-mcp, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to gin-mcp! Your work helps make this project better for everyone. 🌟

*"The beauty of open source is that it belongs to everyone and no one at the same time."* 