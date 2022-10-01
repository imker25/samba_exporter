# Developer Hints

## VSCode Settings

In case you want develop this software with [VS Code](https://code.visualstudio.com/) you need to add the repositories root folder to the **GOPATH** within the `VS Code Settings` to get golang extension and golang tools work, e. g.:

```json
{
      "go.gopath": "${env:GOPATH}:${workspaceFolder}",
}
```

## Debugging

For debugging the pipe communication you might want to use the `-test-mode` on `samba_statusd` and `samba_exporter`.

**Remark:** Never use `-test-mode` on just one of the two programs.