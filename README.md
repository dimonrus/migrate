# migrate
Database migration features

## How to use
``` 
migrations := &migrate.Migration{
    RegistryPath:  "user/project",
    MigrationPath: "user/project/migrations",
    RegistryXPath: "app.Registry",
    DBO:           app.GetDB(),
    Registry:      make(migrate.MigrationRegistry),
}

// Init schema migration
err := migrations.InitMigration("schema")
if err != nil {
    panic(err)
}

// Execute schema migrations
err = migrations.Upgrade("schema")
if err != nil {
    panic(err)
}
```

#### If you find this project useful or want to support the author, you can send tokens to any of these wallets
- Bitcoin: bc1qgx5c3n7q26qv0tngculjz0g78u6mzavy2vg3tf
- Ethereum: 0x62812cb089E0df31347ca32A1610019537bbFe0D
- Dogecoin: DET7fbNzZftp4sGRrBehfVRoi97RiPKajV
