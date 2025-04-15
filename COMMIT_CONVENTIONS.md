# ✅ Commit Message Conventions

This project follows the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification to ensure consistent and meaningful commit history. This standard helps power changelogs, semantic versioning, and CI/CD automation.

---

## 📜 Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

---

## 🔠 Types

| Type       | Description                                  |
|------------|----------------------------------------------|
| `feat`     | A new feature                                |
| `fix`      | A bug fix                                    |
| `docs`     | Documentation only changes                   |
| `style`    | Code style (formatting, spacing, etc.)       |
| `refactor` | Code refactoring that doesn’t fix or add features |
| `test`     | Adding or updating tests                     |
| `chore`    | Build process, CI/CD, or tooling changes     |

---

## 📌 Examples

```bash
feat(api): add support for refresh tokens
fix(auth): correct token expiration edge case
docs(readme): add deployment instructions
style(app): reformat YAML indentation
chore(ci): update cache key for workflows
```

---

## 🧪 Linting

Commit messages are automatically linted using [`commitlint`](https://github.com/conventional-changelog/commitlint) to enforce this format. Commits that do not match the specification may be rejected during development or CI.

```
npx commitlint --edit
```
