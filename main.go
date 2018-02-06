package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"
    "fmt"

    "github.com/davecgh/go-spew/spew"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/montana-network/blockchain"
)

type JsonRPC struct {
    Jsonrpc string
    Method string
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        t := time.Now()
        genesisBlock := blockchain.Block{0, t.String(), "", "", 1, "coin", nil}
        spew.Dump(genesisBlock)
        blockchain.Blockchain = append(blockchain.Blockchain, genesisBlock)
    }()

    log.Fatal(run())
}

func run() error {
    mux := makeMuxRouter()
    httpAddr := os.Getenv("MONTANA_NODE_PORT")
    log.Println("Listening on ", os.Getenv("MONTANA_NODE_PORT"))

    s := &http.Server{
        Addr:           ":" + httpAddr,
        Handler:        mux,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    if err := s.ListenAndServe(); err != nil {
        return err
    }

    return nil
}

func makeMuxRouter() http.Handler {
    muxRouter := mux.NewRouter()
    muxRouter.HandleFunc("/rpc", handleRouting).Methods("POST")
    return muxRouter
}

func handleRouting(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)

    var request JsonRPC

    if err := decoder.Decode(&request); err != nil {
        respondWithJSON(w, r, http.StatusBadRequest, r.Body)
        return
    }

    if request.Method == "explorer" {
        respondWithJSON(w, r, http.StatusOK, blockchain.Blockchain)
    }

    fmt.Printf("%v\n", r.Body)
    fmt.Printf("%v\n", request)
    fmt.Printf("%v\n", request.Method + " - wow 2")
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
    response, err := json.MarshalIndent(payload, "", "  ")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("HTTP 500: Internal Server Error"))
        return
    }
    w.WriteHeader(code)
    w.Write(response)
}
