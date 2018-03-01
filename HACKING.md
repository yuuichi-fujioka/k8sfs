# Hacking on K8sFS

## Reporting CPU/Mem Profiles

k8sfs reports cpu/mem profiles when SIGQUIT is handled.

e.g.

```
kill -SIGQUIT <k8sfs PID>
```

cpu.prof and mem.prof will be wrote on the current directory that k8sfs is working.
They are visualized with `go tool pprof cpu.prof`.
