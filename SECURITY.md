# Security Policy

Thank you for your interest in keeping Journal Guardian secure. As a tool that processes logs and interacts with Large Language Models (LLMs), we take security and data privacy seriously.

## Supported Versions

We currently provide security updates for the following versions of Journal Guardian. If you are using an older version, please upgrade to the latest release to ensure your environment is secure.

| Version | Supported          | Notes                                  |
| :-----: | :----------------: | :------------------------------------- |
| v1.3.x  | :white_check_mark: | Current stable release                 |
| v1.2.x  | :x:                | End of life                            |
| v1.0.x  | :x:                | End of life                            |
| < 1.0   | :x:                | End of life                            |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

If you have discovered a security vulnerability in Journal Guardian, please report it privately.

1.  **Email:** Send a detailed report to **z@emre.xyz**.
2.  **GitHub Private Reporting:** If enabled on this repository, you can go to the **Security** tab -> **Advisories** -> **New draft security advisory**.

Please include:
* The specific version of Journal Guardian you are using.
* The steps to reproduce the vulnerability.
* Any relevant logs or configuration snippets (sanitize secrets first).

We will acknowledge your report within **48 hours** and will provide an estimated timeline for a fix.

## LLM & Data Privacy Specifics

Since Journal Guardian sends log data to external LLM providers (e.g., OpenAI, Anthropic) or local models, please be aware of the following security considerations:

* **PII & Secrets:** We strive to implement sanitization filters, but **you are ultimately responsible** for ensuring your logs do not contain unencrypted secrets (API keys, passwords) or sensitive PII before Journal Guardian processes them.
* **Prompt Injection:** Be aware that malicious log entries could theoretically attempt "Prompt Injection" attacks against the analyzing LLM. While this rarely affects the Go application itself, it may result in misleading explanations from the LLM.

## Security Updates

When a vulnerability is fixed, we will release a patch version (e.g., `v1.3.1`) and publish a security advisory. We encourage all users to update immediately upon release.
