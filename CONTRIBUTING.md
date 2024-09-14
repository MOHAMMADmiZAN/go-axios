
# Contributing to Go-Axios

We welcome contributions to Go-Axios! Below are guidelines to help you get started.

## How to Contribute

1. Fork the repository by clicking the "Fork" button at the top right.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/go-axios.git
   ```
3. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature-branch
   ```
4. Make your changes in your local repository.
5. Write or update tests as needed.
6. Commit your changes with a descriptive commit message:
   ```bash
   git commit -m "Add feature: <description>"
   ```
7. Push your branch to your fork:
   ```bash
   git push origin feature-branch
   ```
8. Open a pull request from your branch on GitHub.
9. Your pull request will be reviewed, and you may be asked to make additional changes.

## Code Guidelines

- Write clear and concise commit messages.
- Ensure your code follows Go's style conventions. You can use `gofmt` to format your code:
  ```bash
  gofmt -w .
  ```
- Ensure all tests pass before submitting your pull request:
  ```bash
  go test ./...
  ```
- Add new tests for any new features or functionality.

## Reporting Issues

If you encounter any bugs or issues, feel free to open an issue on GitHub. Please provide as much detail as possible, including steps to reproduce the issue.

## License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.
