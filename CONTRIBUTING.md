# Contributing to Journal Guardian

First off, thanks for taking the time to contribute! 🎉

Journal Guardian is an open-source tool written in Go that monitors logs for errors and uses Large Language Models (LLMs) to provide explanations. We value your help in making this tool more robust, efficient, and easier to use.

The following is a set of guidelines for contributing to Journal Guardian. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

---

## 🚀 How Can I Contribute?

### 1. Reporting Bugs
This section guides you through submitting a bug report for Journal Guardian. Following these guidelines helps maintainers and the community understand your report, reproduce the behavior, and find related reports.

* **Check existing issues** to see if the bug has already been reported.
* **Use a clear and descriptive title** for the issue to identify the problem.
* **Describe the exact steps to reproduce the problem** in as much detail as possible.
* **Provide specific examples** to demonstrate the steps. Include logs (sanitize any sensitive info!) and the LLM output if relevant.
* **Describe the behavior you observed** after following the steps and point out what exactly is the problem with that behavior.
* **Explain which behavior you expected to see instead and why.**

### 2. Suggesting Enhancements
This section guides you through submitting an enhancement suggestion, including completely new features (like support for new LLM providers or log formats) and minor improvements.

* **Use a clear and descriptive title** for the issue to identify the suggestion.
* **Provide a step-by-step description of the suggested enhancement** in as much detail as possible.
* **Explain why this enhancement would be useful** to most Journal Guardian users.

### 3. Code Contributions
We love pull requests! Here's a quick guide to getting your code merged.

#### Local Development Setup
1.  **Fork the repository** on GitHub.
2.  **Clone your fork** locally:
    ```bash
    git clone [https://github.com/your-username/journal-guardian.git](https://github.com/your-username/journal-guardian.git)
    cd journal-guardian
    ```
3.  **Install Go dependencies**:
    ```bash
    go mod download
    ```
4.  **Set up Environment Variables**:
    Create a `.env` file based on `.env.example`. You will need API keys for the LLMs you intend to test (e.g., OpenAI, Anthropic, or local LLM endpoints).

#### Specific Areas We Need Help With
* **Log Parsers:** Adding regex or logic to support different log formats (JSON, Syslog, Apache Common, etc.).
* **LLM Integration:** Adding adapters for new LLM providers (e.g., Mistral, Llama via Ollama, etc.).
* **Prompt Engineering:** Improving the system prompts to get better, more concise explanations for errors.
* **Performance:** Optimizing the log tailing and concurrent processing in Go (Goroutines/Channels).

---

## 📝 Styleguides

### Go Styleguide
* We follow standard **[Effective Go](https://go.dev/doc/effective_go)** guidelines.
* Run `go fmt ./...` before committing to ensure standard formatting.
* Run `go vet ./...` to catch common errors.
* **Naming:** Use meaningful variable names. Since this project parses logs, clarify if a variable is a `rawLogLine`, `parsedLogEntry`, or `errorExplanation`.

### LLM & Prompting Guidelines
If you are modifying how we interact with LLMs:
* **Token Economy:** Be mindful of token usage. Do not send entire log files to the LLM; send only the relevant error context.
* **Privacy:** Ensure no PII (Personally Identifiable Information) from logs is inadvertently sent to external LLM APIs.
* **Fallback:** Always implement error handling if the LLM API is unreachable or returns a hallucination/malformed response.

### Commit Messages
* Use the present tense ("Add feature" not "Added feature").
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...").
* Reference issues and pull requests liberally after the first line.

---

## 🧪 Testing

We value stability. Please include tests with your PRs.

* **Unit Tests:** Run `go test ./...` to execute the suite.
* **Mocking LLMs:** Do **not** make real API calls to LLM providers in the test suite. Use the provided Mock interfaces to simulate LLM responses.

---

## ⚖️ License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
