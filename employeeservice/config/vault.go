package config

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/hashicorp/vault/api"
)

type VaultConfig struct {
    Address string
    Token   string
    Client  *api.Client
}

type MongoDBConfig struct {
    URI      string `json:"uri"`
    Username string `json:"username"`
    Password string `json:"password"`
    Database string `json:"database"`
}

type APMConfig struct {
    ServerURL   string `json:"server_url"`
    SecretToken string `json:"secret_token"`
    Environment string `json:"environment"`
}

type ServiceConfig struct {
    MongoDB MongoDBConfig
    APM     APMConfig
    Port    int
}

func InitVaultClient() (*VaultConfig, error) {
    vaultAddr := os.Getenv("VAULT_ADDR")
    if vaultAddr == "" {
        return nil, fmt.Errorf("VAULT_ADDR environment variable is required")
    }

    vaultToken := os.Getenv("VAULT_TOKEN")
    if vaultToken == "" {
        return nil, fmt.Errorf("VAULT_TOKEN environment variable is required")
    }

    config := api.DefaultConfig()
    config.Address = vaultAddr
    config.Timeout = 10 * time.Second

    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create Vault client: %v", err)
    }

    client.SetToken(vaultToken)

    return &VaultConfig{
        Address: vaultAddr,
        Token:   vaultToken,
        Client:  client,
    }, nil
}

func GetSecrets(vc *VaultConfig, serviceName string) (*ServiceConfig, error) {
    ctx := context.Background()

    mongoSecret, err := vc.Client.KVv2("kindergarten").Get(ctx, "mongodb")
    if err != nil {
        return nil, fmt.Errorf("failed to read MongoDB secrets: %v", err)
    }

    apmSecret, err := vc.Client.KVv2("kindergarten").Get(ctx, "apm")
    if err != nil {
        return nil, fmt.Errorf("failed to read APM secrets: %v", err)
    }

    portsSecret, err := vc.Client.KVv2("kindergarten").Get(ctx, "ports")
    if err != nil {
        return nil, fmt.Errorf("failed to read ports secrets: %v", err)
    }

    mongoData, err := json.Marshal(mongoSecret.Data)
    if err != nil {
        return nil, err
    }

    var mongoConfig MongoDBConfig
    if err := json.Unmarshal(mongoData, &mongoConfig); err != nil {
        return nil, err
    }

    apmData, err := json.Marshal(apmSecret.Data)
    if err != nil {
        return nil, err
    }

    var apmConfig APMConfig
    if err := json.Unmarshal(apmData, &apmConfig); err != nil {
        return nil, err
    }

    var port int
    switch serviceName {
    case "student":
        if portVal, ok := portsSecret.Data["student"].(float64); ok {
            port = int(portVal)
        } else {
            port = 5001
        }
    case "teacher":
        if portVal, ok := portsSecret.Data["teacher"].(float64); ok {
            port = int(portVal)
        } else {
            port = 5002
        }
    case "employee":
        if portVal, ok := portsSecret.Data["employee"].(float64); ok {
            port = int(portVal)
        } else {
            port = 5003
        }
    default:
        port = 5000
    }

    return &ServiceConfig{
        MongoDB: mongoConfig,
        APM:     apmConfig,
        Port:    port,
    }, nil
}

func GetMongoDBSecrets(vc *VaultConfig) (*MongoDBConfig, error) {
    secrets, err := GetSecrets(vc, "")
    if err != nil {
        return nil, err
    }
    return &secrets.MongoDB, nil
}

func GetPort(vc *VaultConfig, serviceName string) (int, error) {
    secrets, err := GetSecrets(vc, serviceName)
    if err != nil {
        return 0, err
    }
    return secrets.Port, nil
}