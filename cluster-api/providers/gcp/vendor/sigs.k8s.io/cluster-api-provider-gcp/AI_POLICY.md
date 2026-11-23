# AI Usage Policy

The Cluster API Provider GCP (CAPG) project has rules for AI usage.
The project acknowledges the benefits from thoughtful AI-assisted development,
but contributors must maintain high standards for code quality, security, and collaboration.

## Why do we have the policy
This policy is not to discourage the use of AI and should not be seen as anti-AI.
Instead, it's needed to ensure that the project standards are adhered to and
that contributors fully understand their contributions.

All issues and PRs require time from the human maintainers, reviewers and other contributors.
This policy exists to protect the valuable and often limited time of these humans from poor quality changes.

## Principles

- **Human accountability:** Contributors are responsible for all submitted
  changes, which means they understand the work well enough to explain,
  maintain, and support it.
- **Quality:** AI-assisted work must meet the same quality, testing, security,
  and review standards as any other contribution.
- **Transparency:** Disclose substantive AI assistance so reviewers understand
  how the contribution was produced and validated.

## Rules

- **You must declare AI usage**. When submitting a PR that contains AI assisted work the use of AI this must be stated.
Listing AI tooling as a co-author, co-signing commits using an AI tool, or using the assisted-by, co-developed or similar commit trailer is not allowed.
- **AI must not be used for PR descriptions**. A key indicator to the reviewers and maintainers that you understand
your contribution is a PR description written by yourself with your understanding.
If a PR description looks to be AI generated it will be closed.
Moreover PR descriptions must exactly follow the prescribed .github/PULL_REQUEST_TEMPLATE.md.
- **Write issue descriptions yourself:** Do not use AI to draft or generate
  issue descriptions. Issues should reflect the contributor's own knowledge
  and experience, and follow the applicable template under
  `.github/ISSUE_TEMPLATE/`.
- **Respond to PR review comments yourself:** Do not use AI to draft responses
  to review comments. Reviewers want to engage directly with contributors, not
  generated responses. If you do not engage directly with reviewers, the PR
  will be closed.
- **Own AI-assisted documentation and proposals:** AI may help draft or revise
  documentation, proposals, user stories, or research. Contributors must
  understand, verify, and take responsibility for the final content.

## Best Practices

### ✅ Recommended Uses
- Generating boilerplate code and common patterns
- Creating comprehensive test suites
- Refactoring existing code for clarity
- Generating utility functions and helpers
- Explaining existing code patterns

### ❌ Avoid AI For
- API Version bumps
- CAPI Contract changes
- Complex logic without thorough review
- Security critical authentication/authorization code
- Code you don’t fully understand
- Large architectural changes
