# Contributing to Quote Vault

We welcome contributions to Quote Vault! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/quote-vault.git`
3. Create a feature branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit your changes: `git commit -m "feat: add new feature"`
7. Push to your branch: `git push origin feature/your-feature-name`
8. Create a Pull Request

## Development Setup

### Prerequisites
- Go 1.19 or later
- SQLite3
- Git

### Local Development

```bash
# Clone the repository
git clone https://github.com/yourusername/quote-vault.git
cd quote-vault

# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env

# Run the application
go run main.go
```

## Code Style

### Go Code Guidelines
- Follow standard Go conventions
- Use `gofmt` to format your code
- Use meaningful variable and function names
- Add comments for public functions and complex logic
- Keep functions small and focused

### Commit Message Convention

We use conventional commits format:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Build process or auxiliary tool changes

Example:
```
feat: add quote search functionality

Implement search endpoint that allows filtering quotes
by author name or content keywords.
```

## Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./handlers
```

### Writing Tests
- Write unit tests for new functionality
- Use table-driven tests when appropriate
- Mock external dependencies
- Aim for good test coverage

## API Guidelines

### REST Principles
- Use appropriate HTTP methods (GET, POST, PUT, DELETE)
- Use meaningful HTTP status codes
- Follow consistent URL patterns
- Include proper error responses

### Response Format
- Use consistent JSON response structure
- Include pagination metadata for list endpoints
- Provide meaningful error messages

## Pull Request Process

1. Ensure your code follows the style guidelines
2. Add tests for new functionality
3. Update documentation if needed
4. Ensure all tests pass
5. Update the README.md if you change functionality

### PR Description Template
```
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Added/updated unit tests
- [ ] Manual testing completed
- [ ] All tests pass

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
```

## Issue Reporting

### Bug Reports
When reporting bugs, please include:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (Go version, OS)
- Relevant logs or error messages

### Feature Requests
For feature requests, please provide:
- Clear description of the feature
- Use case or motivation
- Proposed implementation approach
- Any related issues or discussions

## Code Review Guidelines

### For Authors
- Keep PRs focused and reasonably sized
- Write clear PR descriptions
- Respond to feedback promptly
- Be open to suggestions

### For Reviewers
- Be constructive and helpful
- Focus on code quality and functionality
- Ask questions if something is unclear
- Approve when satisfied with changes

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release PR
4. Tag release after merge
5. Update documentation if needed

## Getting Help

- Check existing issues and documentation
- Create an issue for questions or problems
- Be specific about your problem or question

Thank you for contributing to Quote Vault!