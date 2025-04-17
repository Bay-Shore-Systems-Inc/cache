# 🤝 Contributing Guide

We follow a structured Git workflow and use **Conventional Commits** to ensure all changes are tracked, versioned, and deployed cleanly through our CI/CD pipeline.

---

## 📜 Commit Message Format

Please follow our commit message conventions outlined in:

👉 [COMMIT_CONVENTIONS.md](./COMMIT_CONVENTIONS.md)

This allows us to:
- Generate changelogs automatically
- Apply semantic versioning
- Ensure clean, searchable Git history

---

## ✅ Contribution Checklist

Before opening a PR:

- Ensure your branch is up to date with the target environment (e.g. `main`)
- Write clear, concise commit messages following the convention
- Test your changes locally or in a dev environment
- Ensure secrets or credentials are **never committed**

---

## 🛠️ Tools

- Use `npx commitlint --edit` to validate your commit
- Use `npm run release` (if configured) to generate a changelog
