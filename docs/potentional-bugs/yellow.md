# Minor Risks

**Migration loading** — The server expects `migrations/001_init.sql` relative to its working directory. Starting the binary outside the repository root requires packaging or an explicit working directory.

**Grid history** — Empty historical periods return null cells by design; the startup seed only covers the latest two hours.
