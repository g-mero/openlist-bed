# OpenList-Bed Project Guidelines

## Project Overview

This is an image hosting service built with:

- **vipsgen** - Go bindings for libvips image processing
- **OpenList/Alist** - File storage backend integration
- **Fiber v3** - Web framework for Go

### Docker Deployment

Uses multi-stage build with Debian + libvips.

## vipsgen API Usage

This project uses **github.com/cshum/vipsgen/vips** for image processing. Key points:

1. **API Design**: The vipsgen API is almost identical to the official libvips C API, with different naming conventions (Go-style naming instead of C-style)
2. **Source Code Reference**: When in doubt about API usage, check the vipsgen source code at `github.com/cshum/vipsgen/vips`
3. **Parameter Documentation**: The source code may not document all parameter meanings. For detailed parameter descriptions, refer to the official libvips C API documentation at https://www.libvips.org/API/current/index.html
4. **Function Mapping**: C functions like `vips_heifsave_buffer` map to Go functions in vipsgen with similar names and signatures

## Code Conventions

1. **Language Usage**: Use **English** for all code documentation, comments, and commit messages (unless explicitly requested otherwise). When responding to user questions, match the user's language (conversational responses only, not code)
2. **File Modifications**: Ask before creating or modifying files unless explicitly instructed
3. **No Unnecessary Documentation**: Do not create README, CHANGELOG, or other documentation files unless explicitly requested
4. **Keep This File Concise**: Avoid unnecessary complexity in project guidelines
