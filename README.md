# Auto Closer

This action closes all issues that have a specific label while keeping the latest one(s) open. Works best for daily, weekly or monthly auto created issues (e.g. [lowply/issue-from-template](https://github.com/lowply/issue-from-template/)). Use with care for the initial run especially when you have a large number of open issues with the label.

## Environment variables

- `AC_LABEL` (_required_): The label that the target issue should have.
- `AC_KEEP` (_optional_): The number of the issues should be kept open. Default value: `1`

## Workflow example

```
name: weekly report
on:
  schedule:
  - cron: "0 0 * * 2"
jobs:
  open:
    name: Open new report issue
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - uses: lowply/issue-from-template@v0.1.0
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        IFT_TEMPLATE_NAME: template.md
  close:
    needs: open
    name: Close old issues
    runs-on: ubuntu-latest
    steps:
    - uses: lowply/auto-closer@v0.0.5
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        AC_LABEL: "report"
        AC_KEEP: 3
```

## Running locally for development

This is designed to be used as a GitHub Action, but you can also just build and run it locally with the following env vars:

```
cd src
export GITHUB_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export GITHUB_REPOSITORY="owner/repository"
export AC_LABEL="label"
go run .
```
