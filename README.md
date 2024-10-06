<div align="center">

# Flash

[![Go](https://img.shields.io/badge/Go-1.19-blue.svg?style=flat&logo=go)](https://golang.org)
[![Open Source](https://img.shields.io/badge/Open%20Source-%F0%9F%92%9A-brightgreen?style=flat)](https://opensource.org)
[![Azure](https://img.shields.io/badge/Azure-Cloud-blue.svg?style=flat&logo=microsoft-azure)](https://azure.microsoft.com/)
[![Powered by ChatGPT](https://img.shields.io/badge/Powered%20by-ChatGPT-ff69b4.svg?style=flat&logo=openai)](https://openai.com/)

</div>

**Flash** is an AI-powered code vulnerability scanner designed to help developers identify security vulnerabilities in their code. By leveraging AI models like OpenAI and Azure OpenAI, Flash automates the review process for various coding languages and provides detailed reports with potential vulnerabilities, proof of concepts, and recommended fixes. Flash can generate reports in Markdown format, making it easy for developers to integrate security analysis into their workflow.

## Features

- **AI-Powered Code Analysis**: Leverages OpenAI's GPT models to analyze code and detect potential security vulnerabilities.
- **Multi-Platform Support**: Flash works across various platforms and languages, making it a flexible solution for code review.
- **Detailed Vulnerability Reports**: Generates reports with detailed descriptions of identified vulnerabilities, proof of concepts, and recommended fixes.
- **Supports Multiple Languages**: Works with PHP, Python, JavaScript, and more.
- **Markdown Report Generation**: Outputs security analysis in Markdown format for easy integration with GitHub and other platforms.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/secfathy/flash.git
    ```

2. Navigate to the project directory:
    ```bash
    cd flash
    ```

3. Build the application:
    ```bash
    go build
    ```

4. Run the application:
    ```bash
    go run main.go -C <codefile> -O <outputfile.md>
    ```

## Usage

Flash scans code files for vulnerabilities by sending code snippets to AI models, which then return a detailed analysis of the vulnerabilities. The results can be saved as Markdown reports.

### Command-line Options:

- `-C`: Path to the code file to be analyzed.
- `-O`: Optional path to save the Markdown results.
- `-H`: Show help message.

### Example:

```bash
go run main.go -C example.php -O report.md

```

### Environment File Settings

```bash

# Choose which API to use (OpenAI or Azure OpenAI)
# Set USE_OPENAI=true if using OpenAI API
# Set USE_AZURE=true if using Azure OpenAI API
USE_OPENAI=false
USE_AZURE=true

# OpenAI API configuration
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_MODEL=gpt-3.5-turbo

# Azure OpenAI API configuration
AZURE_API_KEY=
AZURE_API_VERSION=
AZURE_API_ENDPOINT=
AZURE_DEPLOYMENT_NAME=

```

## Contribution

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-branch-name`.
3. Make your changes and commit them: `git commit -m 'Add new feature'`.
4. Push to the branch: `git push origin feature-branch-name`.
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For questions or support, feel free to reach out to the repository owner.

---

Developed with ❤️ by [Mohammed Fathy @Secfathy](https://github.com/Secfathy)


