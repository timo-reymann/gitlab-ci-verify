pydoctor-theme
===

Pydoctor Theme I use with all my Python projects.

## Features

- Support dark theme (via prefers-color-scheme media query)
- Add orange navbar
- Refine layout spacing

## Usage with pydoctor

1. Add this repo as submodule
   ```
   git submodule add https://github.com/timo-reymann/pydoctor-theme.git
   ```
2. In order to use it adjust your pydoctor configuration:
    ```ini
    theme        = "base"
    template-dir = "pydoctor-theme"
    ```
