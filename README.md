# Pomodoro CLI

![Pomodoro CLI Screenshot](screenshot.png)

Pomodoro CLI is a command-line interface application built in Go that implements the Pomodoro technique. It features a full-screen, 3D-styled digital clock to help you maintain focus during work sessions and breaks, entirely from your terminal.

## Features

*   **Standard Pomodoro Cycle:** Automatically alternates between Work, Short Break, and Long Break sessions.
*   **Full-Screen Distraction-Free UI:** Utilizes the terminal's Alternate Screen Buffer to hide your command history and prompt, providing a clean interface.
*   **Custom 3D Rendering:** Renders large, retro-style digital clock text with a dynamic drop shadow effect.
*   **Customizable Intervals:** Easily configure the duration of work sessions, breaks, and the total number of cycles.
*   **Configuration File Support:** Integrates with Viper, allowing you to define your preferences in a persistent configuration file.
*   **Decoupled Architecture:** Built with a strict separation of concerns, separating the core timer engine from the terminal UI rendering logic.

## Installation

### Prerequisites

Ensure you have Go installed on your system (version 1.18 or higher is recommended).

### Building from Source

1.  Clone the repository to your local machine:
    ```bash
    git clone https://github.com/c4rl0s04/pomodoro-cli.git
    cd pomodoro-cli
    ```

2.  Build the executable using the provided Makefile:
    ```bash
    make build
    ```
    This will generate an executable file named `pomodoro-cli` in the root of the project directory.

## Usage

You can start the timer using the default settings or customize it via command-line flags.

### Basic Usage

To start a standard Pomodoro session (25 minutes work, 5 minutes short break, 15 minutes long break, for 4 cycles):

```bash
./pomodoro-cli start
```

### Customizing Durations

Use flags to override the default durations. All durations are specified in minutes.

*   `--work` or `-w`: Set the duration of the work session.
*   `--short-break` or `-s`: Set the duration of the short break.
*   `--long-break` or `-l`: Set the duration of the long break.
*   `--cycles` or `-c`: Set the number of Pomodoro cycles before the session ends.

**Example:** Start a quick session with 15-minute work intervals and 3-minute breaks:

```bash
./pomodoro-cli start --work 15 --short-break 3
```

### Configuration File

Instead of passing flags every time, you can create a configuration file named `.pomodoro.yaml` in your home directory (`~/.pomodoro.yaml`). The CLI will automatically detect and load these values.

**Example `~/.pomodoro.yaml`:**
```yaml
work: 50
short-break: 10
long-break: 30
cycles: 4
```

## Development

The project is structured into three main packages to ensure maintainability and testability:

*   `cmd/`: Contains the Cobra CLI commands and wiring.
*   `core/pomodoro/`: Contains the pure business logic and timer engine.
*   `ui/`: Contains the pterm rendering logic and terminal manipulation.

### Available Makefile Commands

*   `make build`: Compiles the Go binary.
*   `make run`: Compiles and immediately executes the CLI.
*   `make test`: Runs the unit tests for the core packages.
*   `make lint`: Runs `golangci-lint` to check for stylistic and structural errors.
*   `make clean`: Removes the compiled binary and cleans the Go cache.

## License

This project is open-source and available under the standard MIT License.
