# Contributing to SwiftCal

Thank you for your interest in contributing to SwiftCal! This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- Docker (optional)

### Setup Development Environment

1. Fork the repository
2. Clone your fork:

   ```bash
   git clone https://github.com/yourusername/swiftcal.git
   cd swiftcal
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Set up environment:

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. Run the application:
   ```bash
   make dev
   ```

## Development Workflow

### Branch Naming

- Feature branches: `feature/description`
- Bug fixes: `fix/description`
- Documentation: `docs/description`
- Refactoring: `refactor/description`

### Commit Messages

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

### Pull Request Process

1. Create a feature branch from `main`
2. Make your changes
3. Add tests for new functionality
4. Ensure all tests pass
5. Update documentation if needed
6. Submit a pull request

### Code Style

- Follow Go formatting standards (`gofmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused
- Use error handling consistently

## Testing

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/services

# With coverage
go test -cover ./...
```

### Writing Tests

- Write tests for all new functionality
- Use descriptive test names
- Follow AAA pattern (Arrange, Act, Assert)
- Mock external dependencies

## Documentation

### Code Documentation

- Document exported functions and types
- Use clear and concise comments
- Include examples for complex functions

### API Documentation

- Update API docs when endpoints change
- Include request/response examples
- Document error codes and messages

## Issue Reporting

### Bug Reports

Include:

- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Error messages/logs

### Feature Requests

Include:

- Problem description
- Proposed solution
- Use cases
- Implementation suggestions

## Code Review

### Review Guidelines

- Be constructive and respectful
- Focus on code quality and functionality
- Suggest improvements when possible
- Approve only when satisfied

### Review Checklist

- [ ] Code follows style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No security issues introduced
- [ ] Performance considerations addressed

## Release Process

### Versioning

Follow semantic versioning (MAJOR.MINOR.PATCH):

- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)

### Release Steps

1. Update version in `go.mod`
2. Update CHANGELOG.md
3. Create release tag
4. Build and test release artifacts
5. Publish release notes

## Community Guidelines

- Be respectful and inclusive
- Help others learn and grow
- Share knowledge and best practices
- Report inappropriate behavior

## Getting Help

- Check existing issues and documentation
- Ask questions in discussions
- Join our community chat
- Contact maintainers directly
