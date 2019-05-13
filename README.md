# Auto Closer

This action closes all issues that have a specific label while keeping the latest one(s) open. Works best for daily, weekly or monthly auto created issues. Use with care for the initial run expecially when you have a large number of open issues with the label.

## Environment variables

- `AC_LABEL` (*required*): The label that the target issue should have.
- `AC_KEEP` (*optional*): The number of the issues should be kept open. Default value: `1`

## Running locally for development

This is designed to be used as a GitHub Action, but you can also just build and run it locally with the following env vars:

```
export GITHUB_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export GITHUB_REPOSITORY="owner/repository"
export GITHUB_WORKSPACE="/path/to/your/local/repository"
```

You can use a different repository by overriding the `GITHUB_REPOSITORY` env var for testing purposes.