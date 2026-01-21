package main

import (
    "log"
    "net/http"
    "os"
    "teacherservice/config"
    "teacherservice/database"
    "teacherservice/handlers"
)

func enableCors(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
    // Step 1: Vault থেকে secrets লোড করো
    log.Println("Initializing Teacher Service with Vault integration...")
    
    vaultClient, err := config.InitVaultClient()
    if err != nil {
        log.Fatalf("Failed to initialize Vault client: %v", err)
    }

    // Get MongoDB secrets from Vault
    mongoSecrets, err := config.GetMongoDBSecrets(vaultClient)
    if err != nil {
        log.Fatalf("Failed to get MongoDB secrets from Vault: %v", err)
    }

    // Get port from Vault
    port, err := config.GetPort(vaultClient, "teacher")
    if err != nil {
        log.Printf("Warning: Failed to get port from Vault, using default 5002: %v", err)
        port = 5002
    }

    log.Println("Successfully loaded configuration from Vault")
    log.Printf("Service port: %d", port)

    // Step 2: Database connection (এখন Vault থেকে URI পাবে)
    os.Setenv("MONGODB_URI", mongoSecrets.URI)
    os.Setenv("DATABASE_NAME", mongoSecrets.Database)

    log.Printf("Connecting to MongoDB with URI: %s", mongoSecrets.URI)
    
    if err := database.Connect(); err != nil {
        log.Fatal("Database connection failed:", err)
    }

    // Step 3: HTTP routes setup করো
    http.HandleFunc("/tech/add-teacher", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        handlers.AddTeacher(w, r)
    })

    http.HandleFunc("/tech/teachers", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        handlers.GetTeachers(w, r)
    })

    http.HandleFunc("/tech/delete-teacher", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        if r.Method != http.MethodDelete {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handlers.DeleteTeacher(w, r)
    })

    http.HandleFunc("/tech/update-teacher", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        if r.Method != http.MethodPut {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handlers.UpdateTeacher(w, r)
    })

    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "healthy", "service": "teacher"}`))
    })

    log.Printf("Teacher Service running on port %d", port)
    log.Fatal(http.ListenAndServe(":5002", nil))
}