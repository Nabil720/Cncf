package main

import (
    "log"
    "net/http"
    "os"
    "studentservice/config"
    "studentservice/database"
    "studentservice/handlers"
)

func enableCors(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
    // Step 1: Vault থেকে secrets লোড করো
    log.Println("Initializing Student Service with Vault integration...")
    
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
    port, err := config.GetPort(vaultClient, "student")
    if err != nil {
        log.Printf("Warning: Failed to get port from Vault, using default 5001: %v", err)
        port = 5001
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
    http.HandleFunc("/std/add-student", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        handlers.AddStudent(w, r)
    })

    http.HandleFunc("/std/students", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        handlers.GetStudents(w, r)
    })

    http.HandleFunc("/std/delete-student", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        if r.Method != http.MethodDelete {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handlers.DeleteStudent(w, r)
    })

    http.HandleFunc("/std/update-student", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        if r.Method != http.MethodPut {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handlers.UpdateStudent(w, r)
    })

    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "healthy", "service": "student"}`))
    })

    log.Printf("Student Service running on port %d", port)
    log.Fatal(http.ListenAndServe(":5001", nil))
}